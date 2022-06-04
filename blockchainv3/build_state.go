package blockchainv3

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/chain/chainstatev3"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

/**
 * 重建不成熟的区块状态
 */
func (b *ChainKernel) BuildImmatureBlockStates() (*chainstatev3.ChainState, error) {
	b.insertLock.Lock()
	defer b.insertLock.Unlock()

	s := b.stateImmutable

	// Read status
	latestStatus, e := s.ImmutableStatusRead()
	if e != nil {
		return nil, e
	}
	ithxs := latestStatus.GetImmatureBlockHashList()
	if len(ithxs) == 0 {
		// There is no state to rebuild
		return s, nil
	}
	// Start reading block data
	store := s.BlockStore()
	fmt.Printf("[BlockChain] Build %d immature block states: ", len(ithxs))
	for _, hx := range ithxs {
		if len(hx) != fields.HashSize {
			fmt.Printf("BuildImmatureBlockStates error: len(hx) != fields.HashSize\n")
			continue
		}
		tarblkbts, e := store.ReadBlockBytesByHash(hx)
		if e != nil {
			return nil, e
		}
		// Parsing block
		tarblk, _, e := blocks.ParseBlock(tarblkbts, 0)
		if e != nil {
			return nil, e
		}
		// Search inserted parent status
		baseState, e := s.SearchBaseStateByBlockHashObj(tarblk.GetPrevHash())
		if e != nil {
			return nil, e
		}
		if baseState == nil {
			fmt.Printf("BuildImmatureBlockStates error: cannot find base state for block %d\n", tarblk.GetHeight())
			continue // cannot find base state
		}
		// Insert and get state newstate
		_, e = b.forkStateWithAppendBlock(baseState, tarblk.(interfaces.Block))
		if e != nil {
			return nil, e
		}
		fmt.Printf("%d ", tarblk.GetHeight())
	}
	// After the block is created, find the latest header current state
	currenthash := latestStatus.GetLatestBlockHash()
	if currenthash == nil {
		return s, nil // The latest status is the current status
	}
	// search for
	curState, e := s.SearchBaseStateByBlockHashObj(currenthash)
	if e != nil {
		return nil, e
	}

	curhx := curState.GetPendingBlockHash()
	if len(curhx) > 8 {
		curhx = curhx[len(curhx)-8:]
	}
	fmt.Printf("finished. current block: %d hash: ....%s\n", curState.GetPendingBlockHeight(), hex.EncodeToString(curhx))

	// return
	return curState, nil
}
