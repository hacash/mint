package blockchain

import (
	"github.com/hacash/core/sys"
)

type BlockChainConfig struct {
	datadir          string
	rollbackToHeight uint64
}

func NewEmptyBlockChainConfig() *BlockChainConfig {
	cnf := &BlockChainConfig{
		rollbackToHeight: 0,
	}
	return cnf
}

//////////////////////////////////////////////////

func NewBlockChainConfig(cnffile *sys.Inicnf) *BlockChainConfig {
	cnf := NewEmptyBlockChainConfig()

	section := cnffile.Section("")
	cnf.rollbackToHeight = section.Key("rollbackToHeight").MustUint64(0)

	cnf.datadir = cnffile.MustDataDir()

	return cnf

}
