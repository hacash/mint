package blockchainv3

import (
	"fmt"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
)

/**
 * 发现新区块并尝试插入区块链
 * 返回值：插入新区块的状态，当前使用的最新的状态（切换后的），错误
 */
func (bc *ChainKernel) DiscoverNewBlockToInsert(newblock interfaces.Block, origin string) (interfaces.ChainState, interfaces.ChainState, error) {
	bc.insertLock.Lock()
	defer bc.insertLock.Unlock()

	/*
	 */
	blockOriginIsSync := origin == "sync" || origin == ""

	oldImmutablePending := bc.stateImmutable.GetPending()

	// 插入新的区块

	// 判断是否可以插入
	prevhash := newblock.GetPrevHash()
	basestate, e := bc.stateImmutable.SearchBaseStateByBlockHashObj(prevhash)
	if e != nil {
		return nil, nil, e
	}
	if basestate == nil {
		// 未找到上级区块
		return nil, nil, fmt.Errorf("cannot find prev block %s", prevhash.ToHex())
	}

	// 尝试插入
	//fmt.Printf("insert block - %d ... ", newblock.GetHeight())
	newstate, e := bc.forkStateWithAppendBlock(basestate, newblock)
	if e != nil {
		return nil, nil, e
	}
	//fmt.Println("ok")

	//判断是否需要移动成熟区块
	var isMoveComfirm bool = false
	var doComfirmState = bc.stateImmutable
	if newblock.GetHeight()-oldImmutablePending.GetPendingBlockHeight() > ImmatureBlockMaxLength {
		// 移动确认
		var comfirmState = newstate
		for i := 0; i < ImmatureBlockMaxLength; i++ {
			comfirmState = comfirmState.GetParentObj()
		}
		doComfirmState = comfirmState // 前移确认
		isMoveComfirm = true
	}

	// 插入成功，更新状态表
	immutableStatus, e := bc.stateImmutable.ImmutableStatusRead()
	if e != nil {
		return nil, nil, e
	}
	doComfirmStateImmutable := doComfirmState
	// 更新成熟的区块头以及不成熟的区块哈希表
	immatureBlockHashs, e := doComfirmStateImmutable.SeekImmatureBlockHashs()
	if e != nil {
		return nil, nil, e
	}
	newImmutablePending := doComfirmStateImmutable.GetPending()
	immutableStatus.SetImmutableBlockHeadMeta(newImmutablePending.GetPendingBlockHead())
	immutableStatus.SetImmatureBlockHashList(immatureBlockHashs)

	// 区块保存进磁盘
	e = doComfirmState.BlockStore().SaveBlock(newblock)
	if e != nil {
		return nil, nil, e
	}

	diamondCreate := newstate.GetPending().GetWaitingSubmitDiamond()

	curPdptr := bc.stateCurrent.GetPending()
	isChangeCurrentState := newblock.GetHeight() > curPdptr.GetPendingBlockHeight()
	if isChangeCurrentState {
		immutableStatus.SetLatestBlockHash(newblock.Hash())
	}

	// 保存状态
	e = doComfirmState.ImmutableStatusSet(immutableStatus)
	if e != nil {
		return nil, nil, e
	}

	// 状态保存进磁盘
	if isMoveComfirm {
		newComfirmImmutableState, e := doComfirmState.ImmutableWriteToDiskObj()
		if e != nil {
			return nil, nil, e // 写入磁盘出错
		}
		bc.stateImmutable = newComfirmImmutableState
	}

	// 判断是否切换 current 状态
	var newCurrentStateForReturn = bc.stateCurrent
	if isChangeCurrentState {
		// 更新区块指针
		var upPtrState = newstate
		for i := 0; i < ImmatureBlockMaxLength; i++ {
			if upPtrState == nil {
				break
			}
			pd := upPtrState.GetPending()
			//fmt.Println("UpdateSetBlockHashReferToHeight: ", pd.GetPendingBlockHeight(), pd.GetPendingBlockHash().ToHex())
			e := newstate.BlockStore().UpdateSetBlockHashReferToHeight(
				pd.GetPendingBlockHeight(), pd.GetPendingBlockHash())
			if e != nil {
				return nil, nil, e
			}
			upPtrState = upPtrState.GetParentObj()
		}

		// 更新最新状态指针
		newCurrentStateForReturn = newstate
		bc.stateCurrent = newstate

		// feed
		if false == blockOriginIsSync {
			// send feed
			if diamondCreate != nil {
				// 新确认了钻石
				bc.diamondCreateFeed.Send(diamondCreate)
			}
			// 发送新区快到达通知
			bc.validatedBlockInsertFeed.Send(newblock.(interfacev2.Block))
		}
	}

	// 成功返回
	return doComfirmState, newCurrentStateForReturn, nil
}
