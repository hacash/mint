package blockchain

import (
	"github.com/hacash/core/inicnf"
	"github.com/hacash/core/sys"
)

type BlockChainConfig struct {
	datadir string
}

func NewBlockChainConfig(datadir string) *BlockChainConfig {
	cnf := &BlockChainConfig{
		datadir: datadir,
	}
	return cnf
}

//////////////////////////////////////////////////

func NewBlockChainByIniCnf(cnffile *inicnf.File) (*BlockChain, error) {

	data_dir := sys.CnfMustDataDir(cnffile.Section("").Key("data_dir").String())

	cnf := NewBlockChainConfig(data_dir)
	blockchain, e1 := NewBlockChain(cnf)
	if e1 != nil {
		return nil, e1
	}

	return blockchain, nil
}
