package blockchain

import (
	"github.com/hacash/core/sys"
)

type BlockChainConfig struct {
	cnffile          *sys.Inicnf
	Datadir          string
	RollbackToHeight uint64
}

func NewEmptyBlockChainConfig() *BlockChainConfig {
	cnf := &BlockChainConfig{
		RollbackToHeight: 0,
	}
	return cnf
}

//////////////////////////////////////////////////

func NewBlockChainConfig(cnffile *sys.Inicnf) *BlockChainConfig {
	cnf := NewEmptyBlockChainConfig()
	cnf.cnffile = cnffile

	section := cnffile.Section("")
	cnf.RollbackToHeight = section.Key("RollbackToHeight").MustUint64(0)

	cnf.Datadir = cnffile.MustDataDir()

	return cnf

}
