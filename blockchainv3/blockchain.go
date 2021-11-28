package blockchainv3

import (
	"github.com/hacash/core/interfacev3"
	"sync"
)

// 区块链实例
type BlockChain struct {
	config *BlockChainConfig

	//状态
	stateImmutable *interfacev3.ChainStateImmutable
	stateCurrent   *interfacev3.ChainState

	insertLock sync.RWMutex
}

func NewBlockChain(cnf *BlockChainConfig) (*BlockChain, error) {

	ins := &BlockChain{
		config:     cnf,
		insertLock: sync.RWMutex{},
	}
	return ins, nil
}
