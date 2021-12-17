package blockchain

import (
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/mint"
	"github.com/hacash/mint/coinbase"
	"github.com/hacash/mint/difficulty"
)

func (bc *BlockChain) CreateNextBlockByValidateTxs(txlist []interfaces.Transaction) (interfaces.Block, []interfaces.Transaction, uint32, error) {

	lastest, e1 := bc.chainstate.ReadLastestBlockHeadAndMeta()
	if e1 != nil {
		return nil, nil, 0, e1
	}
	// create
	nextblock := blocks.NewEmptyBlockVersion1(lastest)
	if nextblock.GetHeight() < mint.AdjustTargetDifficultyNumberOfBlocks {
		nextblock.Difficulty = fields.VarUint4(difficulty.LowestDifficultyCompact)
	} else {
		// change diffculty
		_, _, bits, err := bc.CalculateNextDiffculty(lastest)
		//fmt.Println("CalculateNextDiffculty - - - - - ", lastest.GetHeight()+1, " - - - ", hex.EncodeToString(tarhx), hex.EncodeToString(difficulty.Uint32ToHash(lastest.GetHeight(), bits)))
		if err != nil {
			return nil, nil, 0, err
		}
		nextblock.Difficulty = fields.VarUint4(bits)
	}
	// coinbase tx
	nextblock.AddTransaction(coinbase.CreateCoinbaseTx(nextblock.GetHeight()))
	// state run
	blockTempState, e2 := bc.chainstate.NewSubBranchTemporaryChainState()
	if e2 != nil {
		return nil, nil, 0, e2
	}
	blockTempState.SetPendingBlockHeight(nextblock.GetHeight())
	defer blockTempState.DestoryTemporary()
	// append tx
	removeTxs := make([]interfaces.Transaction, 0)
	totaltxs := uint32(0)
	totaltxssize := uint32(0)

	for _, tx := range txlist {
		// 检查tx是否存在
		txinchain, e0 := bc.chainstate.CheckTxHash(tx.Hash())
		if e0 != nil || txinchain {
			removeTxs = append(removeTxs, tx) // remove it , its already in chain
			continue
		}
		if totaltxs >= mint.SingleBlockMaxTxCount || totaltxssize >= mint.SingleBlockMaxSize {
			break // overflow block max size or max num
		}
		txTempState, e1 := blockTempState.NewSubBranchTemporaryChainState()
		if e1 != nil {
			return nil, nil, 0, e1
		}
		err := tx.(interfacev2.Transaction).WriteinChainState(txTempState)
		if err != nil {
			//fmt.Println("********************  create block error  ***********************")
			//fmt.Println(err)
			removeTxs = append(removeTxs, tx) // remove it
			continue
		}
		// add
		nextblock.AddTransaction(tx.(interfacev2.Transaction))
		// 统计
		totaltxs += 1
		totaltxssize += tx.Size()
		// 合并状态
		e2 := blockTempState.MergeCoverWriteChainState(txTempState)
		if e2 != nil {
			txTempState.DestoryTemporary()
			return nil, nil, 0, e2
		}
		// clear
		txTempState.DestoryTemporary()
		// next
	}

	//fmt.Println("CreateNextBlockByValidateTxs:", totaltxs)

	// ok return
	return nextblock, removeTxs, totaltxssize, nil
}
