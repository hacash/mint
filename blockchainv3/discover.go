package blockchainv3

import (
	"fmt"
	"github.com/hacash/core/interfaces"
	"runtime"
)

/**
 * 发现新区块并尝试插入区块链
 * 返回值：插入新区块的状态，当前使用的最新的状态（切换后的），错误
 */
func (bc *ChainKernel) DiscoverNewBlockToInsert(newblock interfaces.Block, origin string) (interfaces.ChainState, interfaces.ChainState, error) {
	bc.insertLock.Lock()
	defer bc.insertLock.Unlock()

	if newblock.GetHeight()%321 == 0 {
		runtime.GC() // Start garbage collection every 321 blocks
	}

	/*
	 */
	blockOriginIsSync := origin == "sync" || origin == "" || origin == "update"

	oldImmutablePending := bc.stateImmutable.GetPending()

	// Insert new block

	// Judge whether it can be inserted
	prevhash := newblock.GetPrevHash()
	basestate, e := bc.stateImmutable.SearchBaseStateByBlockHashObj(prevhash)
	if e != nil {
		return nil, nil, e
	}
	if basestate == nil {
		// Parent block not found
		return nil, nil, fmt.Errorf("cannot find prev block %s", prevhash.ToHex())
	}

	// Try inserting
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
		// Move confirmation
		var comfirmState = newstate
		for i := 0; i < ImmatureBlockMaxLength; i++ {
			comfirmState = comfirmState.GetParentObj()
		}
		doComfirmState = comfirmState // Forward confirmation
		isMoveComfirm = true
	}

	// Insert succeeded, update status table
	immutableStatus, e := bc.stateImmutable.ImmutableStatusRead()
	if e != nil {
		return nil, nil, e
	}
	doComfirmStateImmutable := doComfirmState
	// Update mature block headers and immature block hash tables
	immatureBlockHashs, e := doComfirmStateImmutable.SeekImmatureBlockHashs()
	if e != nil {
		return nil, nil, e
	}
	newImmutablePending := doComfirmStateImmutable.GetPending()
	immutableStatus.SetImmutableBlockHeadMeta(newImmutablePending.GetPendingBlockHead())
	immutableStatus.SetImmatureBlockHashList(immatureBlockHashs)

	// Save blocks to disk
	e = doComfirmState.BlockStore().SaveBlock(newblock)
	if e != nil {
		return nil, nil, e
	}

	diamondCreate := newstate.GetPending().GetWaitingSubmitDiamond()

	curPdptr := bc.stateCurrent.GetPending()
	isChangeCurrentState := newblock.GetHeight() > curPdptr.GetPendingBlockHeight()
	isChangeCurrentForkHead := false // 是否切换了分叉头
	if isChangeCurrentState {
		immutableStatus.SetLatestBlockHash(newblock.Hash())
		if false == newblock.GetPrevHash().Equal(curPdptr.GetPendingBlockHash()) {
			isChangeCurrentForkHead = true // The prev hash is inconsistent and the fork head is switched
			//fmt.Println("prev hash 不一致，切换了分叉头")
		} else {
			//fmt.Println("prev hash 一致，只更新一个编号")
		}
	}

	// Save status
	e = doComfirmState.ImmutableStatusSet(immutableStatus)
	if e != nil {
		return nil, nil, e
	}

	// Save status to disk
	if isMoveComfirm {
		newComfirmImmutableState, e := doComfirmState.ImmutableWriteToDiskObj()
		if e != nil {
			return nil, nil, e // Error writing to disk
		}
		// Release old state
		bc.stateImmutable.Destory() // Memory
		bc.stateImmutable = newComfirmImmutableState
	}

	// Judge whether to switch current status
	var newCurrentStateForReturn = bc.stateCurrent
	if isChangeCurrentState {
		// Update block pointer
		var upPtrState = newstate
		upNUmPtrMax := 1 // 当为同步区块或没有改变fork指针时，只需要更新最后一个区块的指针
		if isChangeCurrentForkHead {
			// Switch the fork head and change the pointing height of five blocks in the history of the state path
			upNUmPtrMax = ImmatureBlockMaxLength + 1
		}
		for i := 0; i < upNUmPtrMax; i++ {
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

		// Update latest status pointer
		newCurrentStateForReturn = newstate
		bc.stateCurrent = newstate

		// feed
		if isChangeCurrentForkHead || false == blockOriginIsSync {
			// Send new area fast arrival notice
			bc.validatedBlockInsertFeed.Send(newblock)
			// send feed
			if diamondCreate != nil {
				// Newly confirmed diamond
				bc.diamondCreateFeed.Send(diamondCreate)
			}
		}
	}

	// Successful return
	return doComfirmState, newCurrentStateForReturn, nil
}
