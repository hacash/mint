package blockchain

import (
	"github.com/hacash/core/sys"
)

type BlockChainConfig struct {
	datadir string
}

func NewEmptyBlockChainConfig() *BlockChainConfig {
	cnf := &BlockChainConfig{}
	return cnf
}

//////////////////////////////////////////////////

func NewBlockChainConfig(cnffile *sys.Inicnf) *BlockChainConfig {
	cnf := NewEmptyBlockChainConfig()

	cnf.datadir = cnffile.MustDataDir()

	return cnf

}
