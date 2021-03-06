package blockchain

import (
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/genesis"
	"github.com/hacash/mint"
)

func (bc *BlockChain) ReadPrev288BlockTimestamp(blockHeight uint64) (uint64, error) {
	bc.prev288BlockTimestampLocker.Lock()
	defer bc.prev288BlockTimestampLocker.Unlock()

	if blockHeight <= mint.AdjustTargetDifficultyNumberOfBlocks {
		return genesis.GetGenesisBlock().GetTimestamp(), nil // genesis block
	}

	blkheight := blockHeight - 1

	prev288height := blkheight / mint.AdjustTargetDifficultyNumberOfBlocks * mint.AdjustTargetDifficultyNumberOfBlocks

	if prev, ok := bc.prev288BlockTimestamp[prev288height]; ok {
		return prev, nil
	}

	//fmt.Println("bc.chainstate.ChainStore  read  prev288height:", prev288height)

	if len(bc.prev288BlockTimestamp) > 200 {
		bc.prev288BlockTimestamp = make(map[uint64]uint64) // clean
	}
	// read
	blockstore := bc.chainstate.BlockStore()
	prev288blockheaddatas, e2 := blockstore.ReadBlockHeadBytesByHeight(prev288height)
	if e2 != nil {
		return 0, e2
	}
	prev288block, _, e3 := blocks.ParseBlockHead(prev288blockheaddatas, 0)
	if e3 != nil {
		return 0, e3
	}
	prev288timestamp := prev288block.GetTimestamp()
	// cache
	bc.prev288BlockTimestamp[prev288height] = prev288timestamp

	//fmt.Println("******************* bc.chainstate.ChainStore ReadPrev288BlockTimestamp  read blockHeight:", blockHeight, "  prev288height:", prev288height, "time", prev288timestamp, "time", time.Unix(int64(prev288timestamp), 0).Format(Time_format_layout))

	// return ok
	return prev288timestamp, nil
}
