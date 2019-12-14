package blockchain

import (
	"github.com/hacash/core/interfaces"
	"github.com/hacash/mint"
	"github.com/hacash/mint/difficulty"
	"math/big"
)

func (bc *BlockChain) CalculateNextDiffculty(lastestBlock interfaces.Block) ([]byte, *big.Int, uint32, error) {

	newBlockHeight := lastestBlock.GetHeight() + 1

	// check hash difficulty
	var prev288BlockTimestamp uint64 = 0
	if newBlockHeight%mint.AdjustTargetDifficultyNumberOfBlocks != 0 {
		// read prev288BlockTimestamp value
		t, e := bc.ReadPrev288BlockTimestamp(newBlockHeight)
		if e != nil {
			return nil, nil, 0, e
		}
		prev288BlockTimestamp = t
	}
	res1, res2, res3 := difficulty.CalculateNextTarget(
		lastestBlock.GetDifficulty(),
		newBlockHeight,
		prev288BlockTimestamp,
		lastestBlock.GetTimestamp(),
		mint.EachBlockRequiredTargetTime,
		mint.AdjustTargetDifficultyNumberOfBlocks,
		nil,
	)
	return res1, res2, res3, nil

}

func (bc *BlockChain) CalculateNextTargetDiffculty() ([]byte, *big.Int, uint32, error) {

	lastestBlock, e1 := bc.chainstate.ReadLastestBlockHeadAndMeta()
	if e1 != nil {
		return nil, nil, 0, e1
	}

	return bc.CalculateNextDiffculty(lastestBlock)
}
