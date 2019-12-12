package blockchain

import (
	"github.com/ethereum/go-ethereum/event"
	"github.com/hacash/chain/blockstore"
	"github.com/hacash/chain/chainstate"
	"github.com/hacash/core/interfaces"
	"path"
	"sync"
)

type BlockChain struct {
	config *BlockChainConfig

	chainstate *chainstate.ChainState

	newBlockArriveQueueCh       chan interfaces.Block
	newTransactionArriveQueueCh chan interfaces.Transaction

	////////////////////////

	power  interfaces.PowMaster
	txpool interfaces.TxPool

	////////////////////////

	validatedBlockInsertFeed event.Feed

	////////////////////////

	// data cache
	prev288BlockTimestamp       map[uint64]uint64
	prev288BlockTimestampLocker sync.Mutex
}

func NewBlockChain(config *BlockChainConfig) (*BlockChain, error) {

	cscnf := chainstate.NewEmptyChainStateConfig()
	cscnf.Datadir = path.Join(config.datadir, "chainstate")
	csobject, e1 := chainstate.NewChainState(cscnf)
	if e1 != nil {
		return nil, e1
	}
	stocnf := blockstore.NewEmptyBlockStoreConfig()
	stocnf.Datadir = path.Join(config.datadir, "chainstore")
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
		config:                      config,
		chainstate:                  csobject,
		newBlockArriveQueueCh:       make(chan interfaces.Block, 10),
		newTransactionArriveQueueCh: make(chan interfaces.Transaction, 50),
		prev288BlockTimestamp:       map[uint64]uint64{},
	}
	// return
	return blockchain, nil
}

func (bc *BlockChain) Start() {

	go bc.loop()

}

func (bc *BlockChain) SubscribeValidatedBlockOnInsert(blockCh chan interfaces.Block) {
	bc.validatedBlockInsertFeed.Subscribe(blockCh)
}

// interface api
func (bc *BlockChain) State() interfaces.ChainStateOperation {
	return bc.chainstate
}
