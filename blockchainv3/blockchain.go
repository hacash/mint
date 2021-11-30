package blockchainv3

import (
	"fmt"
	"github.com/hacash/chain/blockstorev3"
	"github.com/hacash/chain/chainstatev3"
	"github.com/hacash/core/interfacev3"
	"github.com/hacash/mint/event"
	"sync"
)

// 区块链实例
type BlockChain struct {
	config *BlockChainConfig

	//状态
	stateImmutable *chainstatev3.ChainState
	stateCurrent   *chainstatev3.ChainState

	blockstore *blockstorev3.BlockStore

	// data cache
	prev288BlockTimestamp       map[uint64]uint64
	prev288BlockTimestampLocker *sync.Mutex

	// feed
	validatedBlockInsertFeed *event.Feed
	diamondCreateFeed        *event.Feed

	insertLock *sync.RWMutex
}

func NewBlockChain(cnf *BlockChainConfig) (*BlockChain, error) {

	stocnf := blockstorev3.NewBlockStoreConfig(cnf.cnffile)
	blockstore, e := blockstorev3.NewBlockStore(stocnf)
	if e != nil {
		return nil, e
	}

	scnf := chainstatev3.NewChainStateConfig(cnf.cnffile)
	immutable, e := chainstatev3.NewChainStateImmutable(scnf)
	if e != nil {
		return nil, e
	}

	// 区块储存
	immutable.SetBlockStoreObj(blockstore)

	ins := &BlockChain{
		config:                      cnf,
		stateImmutable:              immutable,
		blockstore:                  blockstore,
		validatedBlockInsertFeed:    &event.Feed{},
		diamondCreateFeed:           &event.Feed{},
		prev288BlockTimestamp:       make(map[uint64]uint64),
		prev288BlockTimestampLocker: &sync.Mutex{},
		insertLock:                  &sync.RWMutex{},
	}

	// 重建不成熟的区块状态，返回最新区块值
	ins.stateCurrent, e = ins.BuildImmatureBlockStates()
	if e != nil {
		return nil, e
	}

	return ins, nil
}

func (b BlockChain) BlockStore() interfacev3.BlockStore {
	return b.blockstore
}

func (bc *BlockChain) Start() error {

	fmt.Println("[BlockChain] Block chain state data dir: \"" + bc.config.Datadir + "\"")

	//bc.ifDoRollback() // set config to do rollback

	// 循环等待下载比特币转移日志
	go bc.blockstore.RunDownLoadBTCMoveLog()

	go bc.loop()

	return nil
}
