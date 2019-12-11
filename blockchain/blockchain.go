package blockchain

import (
	"github.com/hacash/chain/chainstate"
	"github.com/hacash/chain/chainstore"
	"github.com/hacash/core/interfaces"
	"path"
	"sync"
)

type BlockChain struct {
	config *BlockChainConfig

	chainstate *chainstate.ChainState

	newBlockArriveQueueCh       chan interfaces.Block
	newTransactionArriveQueueCh chan interfaces.Transaction

	power interfaces.PowMaster

	txpool interfaces.TxPool

	////////////////////////

	// data cache
	prev288BlockTimestamp       map[uint64]uint64
	prev288BlockTimestampLocker sync.Mutex
}

func NewBlockChain(config *BlockChainConfig) (*BlockChain, error) {

	cscnf := chainstate.NewChainStateConfig(path.Join(config.datadir, "chainstate"))
	csobject, e1 := chainstate.NewChainState(cscnf)
	if e1 != nil {
		return nil, e1
	}
	stocnf := chainstore.NewChainStoreConfig(path.Join(config.datadir, "chainstore"))
	stobject, e2 := chainstore.NewChainStore(stocnf)
	if e2 != nil {
		return nil, e2
	}
	e3 := csobject.SetChainStore(stobject) // set chain store
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
