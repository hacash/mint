package blockchain

import (
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/mint"
	"github.com/hacash/mint/coinbase"
	"github.com/hacash/mint/difficulty"
)

func (bc *BlockChain) CreateNextBlockByValidateTxs(txlist []interfaces.Transaction) (interfaces.Block, uint32, error) {

	lastest, e1 := bc.chainstate.ReadLastestBlockHeadAndMeta()
	if e1 != nil {
		return nil, 0, e1
	}
	// create
	nextblock := blocks.NewEmptyBlock_v1(lastest)
	if nextblock.GetHeight() < mint.AdjustTargetDifficultyNumberOfBlocks {
		nextblock.Difficulty = fields.VarInt4(difficulty.LowestDifficultyCompact)
	} else if nextblock.GetHeight()%mint.AdjustTargetDifficultyNumberOfBlocks == 0 {
		// change diffculty
		_, _, bits, err := bc.CalculateNextDiffculty(lastest)
		if err != nil {
			return nil, 0, err
		}
		nextblock.Difficulty = fields.VarInt4(bits)
	}
	// coinbase tx
	nextblock.AddTransaction(coinbase.CreateCoinbaseTx(nextblock.GetHeight()))
	// append tx
	totaltxs := uint32(0)
	totaltxssize := uint32(0)
	for _, tx := range txlist {

		totaltxs += 1
		totaltxssize += tx.Size()
		if totaltxs > 2000 || totaltxssize > mint.SingleBlockMaxSize {
			break // overflow block max size or max num
		}
		// add
		nextblock.AddTransaction(tx)
	}
	// change mkrl root
	nextblock.SetMrklRoot(blocks.CalculateMrklRoot(nextblock.GetTransactions()))

	// ok return
	return nextblock, totaltxssize, nil
}
