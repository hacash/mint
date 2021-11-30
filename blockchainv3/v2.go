package blockchainv3

import (
	"fmt"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/genesis"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/interfacev3"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
	"github.com/hacash/mint"
	"github.com/hacash/mint/coinbase"
	"github.com/hacash/mint/difficulty"
	"github.com/hacash/x16rs"
	"strings"
	"time"
)

func (bc *BlockChain) InsertBlock(newblk interfacev2.Block, origin string) error {
	blkv3 := newblk.(interfacev3.Block)
	_, _, e := bc.DiscoverNewBlockToInsert(blkv3, origin)
	return e
}

func (bc *BlockChain) StateImmutable() interfacev3.ChainStateImmutable {
	return bc.stateImmutable
}

func (b BlockChain) StateRead() interfaces.ChainStateOperationRead {
	return b.stateCurrent
}

func (bc *BlockChain) ValidateTransactionForTxPool(newtx interfacev2.Transaction) error {
	newtxhash := newtx.Hash()
	txhxhex := newtxhash.ToHex()
	exist, e0 := bc.StateRead().CheckTxHash(newtxhash)
	//fmt.Println(exist, exist_tx_bytes)
	if e0 != nil {
		return e0
	}
	if exist {
		return fmt.Errorf("tx %s is exist in blockchain.", txhxhex)
	}
	// check
	if newtx.GetTimestamp() > uint64(time.Now().Unix()) {
		return fmt.Errorf("tx %s timestamp cannot more than now.", txhxhex)
	}
	// fee purity
	if newtx.FeePurity() < mint.MinTransactionFeePurityOfOneByte {
		return fmt.Errorf("tx %s handling fee is too low for miners to accept.", txhxhex)
	}
	// sign
	ok, e1 := newtx.VerifyAllNeedSigns()
	if !ok || e1 != nil {
		return fmt.Errorf("tx %s verify signature error", txhxhex)
	}
	// try run
	lastestBlock, e1 := bc.StateRead().ReadLastestBlockHeadMetaForRead()
	if e1 != nil {
		return e1
	}
	// create temp state
	newTxState, e2 := bc.stateCurrent.ForkNextBlock(lastestBlock.GetHeight()+1, nil, nil)
	if e2 != nil {
		return e2
	}
	newTxState.SetInMemTxPool(true) // 标记是矿池状态
	defer newTxState.Destory()      // clean data
	// validate
	//newTxState.SetPendingBlockHeight(lastestBlock.GetHeight() + 1)
	runerr := newtx.(interfacev3.Transaction).WriteInChainState(newTxState)
	if runerr != nil {
		return runerr
	}
	// have diamond in tx
	// diamond, _ := newTxState.GetPendingSubmitStoreDiamond()
	// ok pass check !
	return nil
}

func (b *BlockChain) ValidateDiamondCreateAction(action interfacev2.Action) error {

	act, ok := action.(*actions.Action_4_DiamondCreate)
	if !ok {
		return fmt.Errorf("its not Action_4_DiamondCreate Action.")
	}

	// 开发者模式，不做检查
	if sys.TestDebugLocalDevelopmentMark {
		return nil // 开发者模式不检查返回成功
	}

	last, err := b.StateRead().ReadLastestDiamond()
	if err != nil {
		return err
	}
	if last == nil { // is first
		genesisblk := genesis.GetGenesisBlock()
		last = &stores.DiamondSmelt{
			Number:           0,
			ContainBlockHash: genesisblk.Hash(),
		}
	}
	if uint32(act.Number) != uint32(last.Number)+1 {
		return fmt.Errorf("Diamond number error.")
	}
	if last.ContainBlockHash.Equal(act.PrevHash) != true {
		return fmt.Errorf("Diamond prev block hash error.")
	}
	hashave, e := b.StateRead().Diamond(act.Diamond)
	if e != nil {
		return e
	}
	if hashave != nil {
		return fmt.Errorf("Diamond <%s> already exist.", act.Diamond)
	}
	// 检查钻石挖矿计算
	sha3hash, diamond_resbytes, diamond_str := x16rs.Diamond(uint32(act.Number), act.PrevHash, act.Nonce, act.Address, act.GetRealCustomMessage())
	diamondstrval, isdia := x16rs.IsDiamondHashResultString(diamond_str)
	if !isdia {
		return fmt.Errorf("String <%s> is not diamond.", diamond_str)
	}
	if strings.Compare(diamondstrval, string(act.Diamond)) != 0 {
		return fmt.Errorf("Diamond need <%s> but got <%s>", act.Diamond, diamondstrval)
	}
	// 检查钻石难度值
	difok := x16rs.CheckDiamondDifficulty(uint32(act.Number), sha3hash, diamond_resbytes)
	if !difok {
		return fmt.Errorf("Diamond difficulty not meet the requirements.")
	}
	// check ok
	return nil
}

func (bc *BlockChain) CreateNextBlockByValidateTxs(txlist []interfacev2.Transaction) (interfacev2.Block, []interfacev2.Transaction, uint32, error) {

	lastest, e1 := bc.StateRead().ReadLastestBlockHeadMetaForRead()
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
	blockTempState, e2 := bc.stateCurrent.ForkNextBlockObj(nextblock.GetHeight(), nil, nil)
	if e2 != nil {
		return nil, nil, 0, e2
	}
	//blockTempState.SetPendingBlockHeight(nextblock.GetHeight())
	defer blockTempState.Destory()
	// append tx
	removeTxs := make([]interfacev2.Transaction, 0)
	totaltxs := uint32(0)
	totaltxssize := uint32(0)

	for _, tx := range txlist {
		// 检查tx是否存在
		txinchain, e0 := bc.StateRead().CheckTxHash(tx.Hash())
		if e0 != nil || txinchain {
			removeTxs = append(removeTxs, tx) // remove it , its already in chain
			continue
		}
		if totaltxs >= 2000 || totaltxssize >= mint.SingleBlockMaxSize {
			break // overflow block max size or max num
		}
		txTempState, e1 := blockTempState.ForkSubChildObj()
		if e1 != nil {
			return nil, nil, 0, e1
		}
		err := tx.(interfacev3.Transaction).WriteInChainState(txTempState)
		if err != nil {
			//fmt.Println("********************  create block error  ***********************")
			//fmt.Println(err)
			removeTxs = append(removeTxs, tx) // remove it
			continue
		}
		// add
		nextblock.AddTransaction(tx)
		// 统计
		totaltxs += 1
		totaltxssize += tx.Size()
		// 合并状态
		e2 := blockTempState.TraversalCopyByObj(txTempState)
		if e2 != nil {
			txTempState.Destory()
			return nil, nil, 0, e2
		}
		// clear
		txTempState.Destory()
		// next
	}

	//fmt.Println("CreateNextBlockByValidateTxs:", totaltxs)

	// ok return
	return nextblock, removeTxs, totaltxssize, nil

	//return nil, nil, 0, nil
}

func (bc *BlockChain) SubscribeValidatedBlockOnInsert(blockCh chan interfacev2.Block) {
	bc.validatedBlockInsertFeed.Subscribe(blockCh)
}

func (bc *BlockChain) SubscribeDiamondOnCreate(diamondCh chan *stores.DiamondSmelt) {
	bc.diamondCreateFeed.Subscribe(diamondCh)
}

func (bc *BlockChain) RollbackToBlockHeight(uint64) (uint64, error) {
	return 0, nil
}