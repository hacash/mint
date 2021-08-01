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

	validatedBlockInsertFeed event.Feed
	diamondCreateFeed        event.Feed

	////////////////////////

	// data cache
	prev288BlockTimestamp       map[uint64]uint64
	prev288BlockTimestampLocker sync.Mutex

	////////////////////////

	insertLock sync.Mutex
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
		config:     config,
		chainstate: csobject,
		// newBlockArriveQueueCh:       make(chan interfaces.Block, 10),
		// newTransactionArriveQueueCh: make(chan interfaces.Transaction, 50),
		prev288BlockTimestamp: map[uint64]uint64{},
	}
	// return
	return blockchain, nil
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
