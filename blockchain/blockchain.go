package blockchain

import (
	"fmt"
	"github.com/hacash/chain/blockstorev2"
	"github.com/hacash/chain/chainstatev2"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/mint/event"
	"os"
	"sync"
)

type BlockChain struct {
	config *BlockChainConfig

	chainstate *chainstatev2.ChainState

	//newBlockArriveQueueCh       chan interfaces.Block
	//newTransactionArriveQueueCh chan interfaces.Transaction

	////////////////////////

	validatedBlockInsertFeed *event.Feed
	diamondCreateFeed        *event.Feed

	////////////////////////

	// data cache
	prev288BlockTimestamp       map[uint64]uint64
	prev288BlockTimestampLocker *sync.Mutex

	////////////////////////

	insertLock *sync.Mutex
}

func NewBlockChain(config *BlockChainConfig) (*BlockChain, error) {

	cscnf := chainstatev2.NewChainStateConfig(config.cnffile)
	// 是否为数据库重建模式
	cscnf.DatabaseVersionRebuildMode = config.DatabaseVersionRebuildMode
	csobject, e1 := chainstatev2.NewChainState(cscnf)
	if e1 != nil {
		fmt.Println("chainstate.NewChainState Error", e1)
		return nil, e1
	}
	stocnf := blockstorev2.NewBlockStoreConfig(config.cnffile)
	stobject, e2 := blockstorev2.NewBlockStore(stocnf)
	if e2 != nil {
		fmt.Println("blockstore.NewBlockStore Error", e2)
		return nil, e2
	}
	e3 := csobject.SetBlockStore(stobject) // set chain store
	if e3 != nil {
		fmt.Println("csobject.SetBlockStore Error", e3)
		return nil, e3
	}
	// new
	blockchain := &BlockChain{
		config:                      config,
		chainstate:                  csobject,
		validatedBlockInsertFeed:    &event.Feed{},
		diamondCreateFeed:           &event.Feed{},
		prev288BlockTimestampLocker: &sync.Mutex{},
		prev288BlockTimestamp:       map[uint64]uint64{},
		insertLock:                  &sync.Mutex{},
	}
	// return
	return blockchain, nil
}

// 替换自己
func (bc *BlockChain) ReplaceSelf(new *BlockChain) {
	bc.config = new.config
	bc.chainstate = new.chainstate
	bc.validatedBlockInsertFeed = new.validatedBlockInsertFeed
	bc.diamondCreateFeed = new.diamondCreateFeed
	bc.prev288BlockTimestamp = new.prev288BlockTimestamp
	bc.prev288BlockTimestampLocker = new.prev288BlockTimestampLocker
	bc.insertLock = new.insertLock
}

func (bc *BlockChain) Close() {
	if bc.chainstate != nil {
		bc.chainstate.Close()
	}
}

func (bc *BlockChain) Start() {

	bc.ifDoRollback() // set config to do rollback

	go bc.loop()

	if !bc.config.DisableDownloadBTCMoveLog {
		go bc.downLoadBTCMoveLog()
	}

}

func (bc *BlockChain) ifDoRollback() {

	if bc.config.RollbackToHeight > 0 {
		tarhei, rollerr := bc.RollbackToBlockHeight(bc.config.RollbackToHeight)
		if rollerr != nil {
			fmt.Println(rollerr.Error())
		} else {
			fmt.Println("Rollback To Block Height", tarhei, "Successfully !")
		}
		os.Exit(0)
	}
}

// interface api
func (bc *BlockChain) State() interfaces.ChainState {
	return bc.chainstate
}
