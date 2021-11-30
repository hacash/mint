package blockchainv3

import (
	"fmt"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/interfacev3"
)

const (
	ImmatureBlockMaxLength = 4 // 最多允许四个不成熟的区块
)

/**
 * 发现新区块并尝试插入区块链
 * 返回值：插入新区块的状态，当前使用的最新的状态（切换后的），错误
 */
func (bc *BlockChain) DiscoverNewBlockToInsert(newblock interfacev3.Block, origin string) (interfacev3.ChainState, interfacev3.ChainState, error) {
	bc.insertLock.Lock()
	defer bc.insertLock.Unlock()

	/*
	 */
	oldImmutablePending := bc.stateImmutable.GetPending()

	// 插入新的区块

	// 判断是否可以插入
	prevhash := newblock.GetPrevHash()
	basestate, e := bc.stateImmutable.SearchBaseStateByBlockHash(prevhash)
	if e != nil {
		return nil, nil, e
	}
	if basestate == nil {
		// 未找到上级区块
		return nil, nil, fmt.Errorf("cannot find prev block %s", prevhash.ToHex())
	}

	// 尝试插入
	newstate, e := bc.forkStateWithAppendBlock(basestate, newblock)
	if e != nil {
		return nil, nil, e
	}

	//判断是否需要移动成熟区块
	var isMoveComfirm bool = false
	var doComfirmState interfacev3.ChainState = bc.stateImmutable
	if newblock.GetHeight()-oldImmutablePending.GetPendingBlockHeight() > ImmatureBlockMaxLength {
		// 移动确认
		var comfirmState interfacev3.ChainState = newstate
		for i := 0; i < ImmatureBlockMaxLength; i++ {
			comfirmState = newstate.GetParent()
		}
		if comfirmState != nil {
			doComfirmState = comfirmState // 前移确认
			isMoveComfirm = true
		}
	}

	// 插入成功，更新状态表
	last, e := bc.stateImmutable.LatestStatusRead()
	if e != nil {
		return nil, nil, e
	}
	doComfirmStateImmutable := doComfirmState.(interfacev3.ChainStateImmutable)
	// 更新成熟的区块头以及不成熟的区块哈希表
	immatureBlockHashs, e := doComfirmStateImmutable.SeekImmatureBlockHashs()
	if e != nil {
		return nil, nil, e
	}
	newImmutablePending := doComfirmStateImmutable.GetPending()
	last.SetImmutableBlockHeadMeta(newImmutablePending.GetPendingBlockHead())
	last.SetImmatureBlockHashList(immatureBlockHashs)

	// 保存进磁盘
	if isMoveComfirm {
		_, e := doComfirmState.ImmutableWriteToDisk()
		if e != nil {
			return nil, nil, e // 写入磁盘出错
		}
	}

	// 判断是否切换 current 状态
	var newCurrentState = bc.stateCurrent
	curPdptr := bc.stateCurrent.GetPending()
	if newblock.GetHeight() > curPdptr.GetPendingBlockHeight() {
		// 更新区块指针
		var upPtrState = newstate
		for i := 0; i < ImmatureBlockMaxLength; i++ {
			if upPtrState == nil {
				break
			}
			pd := upPtrState.GetPending()
			e := newstate.BlockStore().UpdateSetBlockHashReferToHeight(
				pd.GetPendingBlockHeight(), pd.GetPendingBlockHash())
			if e != nil {
				return nil, nil, e
			}
			upPtrState = upPtrState.GetParent()
		}
		// 更新最新状态指针
		newCurrentState = newstate
		bc.stateCurrent = newstate

		// send feed
		diamondCreate := newCurrentState.GetPending().GetWaitingSubmitDiamond()
		if diamondCreate != nil {
			// fmt.Println("diamondCreate bc.diamondCreateFeed.Send(diamondCreate), ", diamondCreate, diamondCreate.Diamond, diamondCreate.Number)
			bc.diamondCreateFeed.Send(diamondCreate)
		}

		orimark := newblock.OriginMark()
		if orimark != "" && orimark != "sync" {
			// 发送新区快到达通知
			bc.validatedBlockInsertFeed.Send(newblock.(interfacev2.Block))
		}

	}

	// 成功返回
	return doComfirmState, newCurrentState, nil
}
