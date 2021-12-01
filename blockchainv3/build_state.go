package blockchainv3

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/chain/chainstatev3"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/interfaces"
)

/**
 * 重建不成熟的区块状态
 */
func (b *ChainKernel) BuildImmatureBlockStates() (*chainstatev3.ChainState, error) {
	b.insertLock.Lock()
	defer b.insertLock.Unlock()

	s := b.stateImmutable

	// 读取状态
	latestStatus, e := s.ImmutableStatusRead()
	if e != nil {
		return nil, e
	}
	ithxs := latestStatus.GetImmatureBlockHashList()
	if len(ithxs) == 0 {
		// 不存在需要重建的状态
		return s, nil
	}
	// 开始读取区块数据
	store := s.BlockStore()
	fmt.Printf("[BlockChain] Build %d immature block states: ", len(ithxs))
	for _, hx := range ithxs {
		tarblkbts, e := store.ReadBlockBytesByHash(hx)
		if e != nil {
			return nil, e
		}
		// 解析区块
		tarblk, _, e := blocks.ParseBlock(tarblkbts, 0)
		if e != nil {
			return nil, e
		}
		// 搜寻插入的父级状态
		baseState, e := s.SearchBaseStateByBlockHashObj(tarblk.GetPrevHash())
		if e != nil {
			return nil, e
		}
		// 插入并获得状态newState
		_, e = b.forkStateWithAppendBlock(baseState, tarblk.(interfaces.Block))
		if e != nil {
			return nil, e
		}
		fmt.Printf("%d ", tarblk.GetHeight())
	}
	// 区块建立完毕，找出最新头部 current state
	currenthash := latestStatus.GetLatestBlockHash()
	if currenthash == nil {
		return s, nil // 最新状态就是当前状态
	}
	// 搜寻
	curState, e := s.SearchBaseStateByBlockHashObj(currenthash)
	if e != nil {
		return nil, e
	}

	curhx := curState.GetPendingBlockHash()
	if len(curhx) > 8 {
		curhx = curhx[len(curhx)-8:]
	}
	fmt.Printf("finished. current block: %d hash: ....%s\n", curState.GetPendingBlockHeight(), hex.EncodeToString(curhx))

	// 返回
	return curState, nil
}
