package difficulty

import (
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/genesis"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/mint"
	"math/big"
	"sync"
)

var (
	// data cache
	prev288BlockTimestampMaps   = make(map[uint64]uint64)
	prev288BlockTimestampLocker = sync.Mutex{}
)

func CalculateNextDiffculty(store interfaces.BlockStoreRead, lastestBlock interfaces.BlockHeadMetaRead) ([]byte, *big.Int, uint32, error) {

	newBlockHeight := lastestBlock.GetHeight() + 1

	// check hash difficulty
	var prev288BlockTimestamp uint64 = 0
	if newBlockHeight%mint.AdjustTargetDifficultyNumberOfBlocks == 0 {
		// read prev288BlockTimestamp value
		t, e := ReadPrev288BlockTimestamp(store, newBlockHeight)
		if e != nil {
			return nil, nil, 0, e
		}
		prev288BlockTimestamp = t
	}

	//var change string = ""

	res1, res2, res3 := CalculateNextTarget(
		lastestBlock.GetDifficulty(),
		newBlockHeight,
		prev288BlockTimestamp,
		lastestBlock.GetTimestamp(),
		mint.EachBlockRequiredTargetTime,
		mint.AdjustTargetDifficultyNumberOfBlocks,
		nil, // &change,
	)

	return res1, res2, res3, nil

}

func ReadPrev288BlockTimestamp(store interfaces.BlockStoreRead, blockHeight uint64) (uint64, error) {
	prev288BlockTimestampLocker.Lock()
	defer prev288BlockTimestampLocker.Unlock()

	if blockHeight <= mint.AdjustTargetDifficultyNumberOfBlocks {
		return genesis.GetGenesisBlock().GetTimestamp(), nil // genesis block
	}

	blkheight := blockHeight - 1

	prev288height := blkheight / mint.AdjustTargetDifficultyNumberOfBlocks * mint.AdjustTargetDifficultyNumberOfBlocks

	if prev, ok := prev288BlockTimestampMaps[prev288height]; ok {
		return prev, nil
	}

	//fmt.Println("bc.chainstate.ChainStore  read  prev288height:", prev288height)

	if len(prev288BlockTimestampMaps) > 200 {
		prev288BlockTimestampMaps = make(map[uint64]uint64) // clean
	}
	// read
	_, prev288blockheaddatas, e2 := store.ReadBlockBytesByHeight(prev288height)
	if e2 != nil {
		return 0, e2
	}
	prev288block, _, e3 := blocks.ParseBlockHead(prev288blockheaddatas, 0)
	if e3 != nil {
		return 0, e3
	}
	prev288timestamp := prev288block.GetTimestamp()
	// cache
	prev288BlockTimestampMaps[prev288height] = prev288timestamp

	//fmt.Println("******************* bc.chainstate.ChainStore ReadPrev288BlockTimestamp  read blockHeight:", blockHeight, "  prev288height:", prev288height, "time", prev288timestamp, "time", time.Unix(int64(prev288timestamp), 0).Format(Time_format_layout))

	// return ok
	return prev288timestamp, nil
}
