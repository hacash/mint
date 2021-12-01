package blockchainv3

import (
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"sync"
)

// 区块链实例
type BlockChain struct {
	config *BlockChainConfig

	chainEngine interfaces.ChainEngineKernel

	mux *sync.RWMutex
}

func NewBlockChain(cnf *BlockChainConfig) (*BlockChain, error) {

	engcnf := NewChainKernelConfig(cnf.cnffile)
	engine, e := NewChainKernel(engcnf)
	if e != nil {
		return nil, e
	}

	// 首次初始化状态
	engine.ChainStateIinitializeCall(setupHacashChainState)

	ins := &BlockChain{
		config:      cnf,
		chainEngine: engine,
		mux:         &sync.RWMutex{},
	}

	return ins, nil
}

func (bc *BlockChain) GetChainEngineKernel() interfaces.ChainEngineKernel {
	bc.mux.RLock()
	defer bc.mux.RUnlock()

	return bc.chainEngine
}

func (bc *BlockChain) SetChainEngineKernel(engine interfaces.ChainEngineKernel) {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	bc.chainEngine = engine
}

func (bc *BlockChain) Start() error {

	fmt.Println("[BlockChain] Block chain state data dir: \"" + bc.config.Datadir + "\"")

	//bc.ifDoRollback() // set config to do rollback

	e := bc.chainEngine.Start()
	if e != nil {
		return e
	}

	// 循环等待下载比特币转移日志
	go bc.chainEngine.CurrentState().BlockStore().RunDownLoadBTCMoveLog()

	go bc.loop()

	return nil
}

// first debug amount
func setupHacashChainState(chainstate interfaces.ChainStateOperation) {
	addr1, _ := fields.CheckReadableAddress("12vi7DEZjh6KrK5PVmmqSgvuJPCsZMmpfi")
	addr2, _ := fields.CheckReadableAddress("1LsQLqkd8FQDh3R7ZhxC5fndNf92WfhM19")
	addr3, _ := fields.CheckReadableAddress("1NUgKsTgM6vQ5nxFHGz1C4METaYTPgiihh")
	amt1, _ := fields.NewAmountFromFinString("ㄜ1:244")
	amt2, _ := fields.NewAmountFromFinString("ㄜ12:244")
	chainstate.BalanceSet(*addr1, stores.NewBalanceWithAmount(amt2))
	chainstate.BalanceSet(*addr2, stores.NewBalanceWithAmount(amt1))
	chainstate.BalanceSet(*addr3, stores.NewBalanceWithAmount(amt1))
}
