package blockchain

import (
	"fmt"
	"github.com/hacash/chain/blockstore"
	"github.com/hacash/chain/chainstate"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/mint/event"
	"os"
	"sync"
)

type BlockChain struct {
	config *BlockChainConfig

	chainstate *chainstate.ChainState

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

	cscnf := chainstate.NewChainStateConfig(config.cnffile)
	csobject, e1 := chainstate.NewChainState(cscnf)
	if e1 != nil {
		return nil, e1
	}
	stocnf := blockstore.NewBlockStoreConfig(config.cnffile)
	stobject, e2 := blockstore.NewBlockStore(stocnf)
	if e2 != nil {
		return nil, e2
	}
	e3 := csobject.SetBlockStore(stobject) // set chain store
	if e3 != nil {
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
