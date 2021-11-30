package blockchainv3

import (
	"github.com/hacash/chain/chainstatev3"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/interfacev3"
)

/**
 * 重建不成熟的区块状态
 */
func (b *BlockChain) BuildImmatureBlockStates() (*chainstatev3.ChainState, error) {
	b.insertLock.Lock()
	defer b.insertLock.Unlock()

	s := b.stateImmutable

	// 读取状态
	latestStatus, e := s.LatestStatusRead()
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
		baseState, e := s.SearchBaseStateByBlockHash(tarblk.GetPrevHash())
		if e != nil {
			return nil, e
		}
		// 插入并获得状态newState
		_, e = b.forkStateWithAppendBlock(baseState, tarblk.(interfacev3.Block))
		if e != nil {
			return nil, e
		}
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

	// 返回
	return curState, nil
}
