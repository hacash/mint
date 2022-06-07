package blockchainv3

import (
	"github.com/hacash/core/sys"
)

type BlockChainConfig struct {
	cnffile *sys.Inicnf

	Datadir          string
	RollbackToHeight uint64

	// Database rebuild mode
	DatabaseVersionRebuildMode bool
}

func NewEmptyBlockChainConfig() *BlockChainConfig {
	cnf := &BlockChainConfig{
		RollbackToHeight:           0,
		DatabaseVersionRebuildMode: false,
	}
	return cnf
}

//////////////////////////////////////////////////

func NewBlockChainConfig(cnffile *sys.Inicnf) *BlockChainConfig {
	cnf := NewEmptyBlockChainConfig()
	cnf.cnffile = cnffile

	section := cnffile.Section("")
	cnf.RollbackToHeight = section.Key("RollbackToHeight").MustUint64(0)

	cnf.Datadir = cnffile.MustDataDirWithVersion()

	return cnf

}
