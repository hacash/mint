package blockchain

import (
	"fmt"
	"github.com/hacash/chain/blockstorev2"
	"github.com/hacash/chain/chainstatev2"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
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

	insertLock *sync.RWMutex
}

func NewBlockChain(config *BlockChainConfig) (*BlockChain, error) {

	cscnf := chainstatev2.NewChainStateConfig(config.cnffile)
	// Database rebuild mode
	cscnf.DatabaseVersionRebuildMode = config.DatabaseVersionRebuildMode
	csobject, e1 := chainstatev2.NewChainState(cscnf)
	if e1 != nil {
		fmt.Println("chainstate.NewChainState Error:", e1)
		return nil, e1
	}
	stocnf := blockstorev2.NewBlockStoreConfig(config.cnffile)
	stobject, e2 := blockstorev2.NewBlockStore(stocnf)
	if e2 != nil {
		fmt.Println("blockstore.NewBlockStore Error:", e2)
		return nil, e2
	}
	e3 := csobject.SetBlockStore(stobject) // set chain store
	if e3 != nil {
		fmt.Println("csobject.SetBlockStore Error:", e3)
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
		insertLock:                  &sync.RWMutex{},
	}
	// return
	return blockchain, nil
}

// Replace yourself
func (bc *BlockChain) ReplaceChainstate(new *BlockChain) {
	bc.config = new.config
	bc.chainstate = new.chainstate
}

func (bc *BlockChain) Close() error {
	bc.insertLock.Lock()
	defer bc.insertLock.Unlock()
	if bc.chainstate != nil {
		bc.chainstate.Close()
	}
	return nil
}

func (bc *BlockChain) Start() error {

	fmt.Println("[BlockChain] Block chain state data dir: \"" + bc.config.Datadir + "\"")

	bc.ifDoRollback() // set config to do rollback

	go bc.loop()

	if !bc.config.DisableDownloadBTCMoveLog {
		go bc.downLoadBTCMoveLog()
	}

	return nil
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
func (bc *BlockChain) State() interfacev2.ChainState {
	bc.insertLock.RLock()
	defer bc.insertLock.RUnlock()
	return bc.chainstate
}

func (bc *BlockChain) StateRead() interfaces.ChainStateOperationRead {
	bc.insertLock.RLock()
	defer bc.insertLock.RUnlock()

	return bc.chainstate
}

func (bc *BlockChain) CurrentState() interfaces.ChainState {
	bc.insertLock.RLock()
	defer bc.insertLock.RUnlock()

	return nil
}

func (b *BlockChain) ChainStateIinitializeCall(func(interfaces.ChainStateOperation)) {

}

func (bc *BlockChain) GetChainEngineKernel() interfaces.ChainEngine {
	return bc
}

func (bc *BlockChain) GetRecentArrivedBlocks() []interfaces.Block {
	return []interfaces.Block{}
}

func (bc *BlockChain) GetLatestAverageFeePurity() uint32 {
	panic(any("nevel call this!"))
}

func (bc *BlockChain) SetChainEngineKernel(engine interfaces.ChainEngine) {
}

// Latest blocks (confirmed and immature)
func (bc *BlockChain) LatestBlock() (interfaces.BlockHeadMetaRead, interfaces.BlockHeadMetaRead, error) {
	blk, e := bc.chainstate.ReadLastestBlockHeadMetaForRead()
	return blk, blk, e
}

// Latest block diamonds
func (bc *BlockChain) LatestDiamond() (*stores.DiamondSmelt, error) {
	return bc.chainstate.ReadLastestDiamond()
}
