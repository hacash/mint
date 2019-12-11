package blockchain

type BlockChainConfig struct {
	datadir string
}

func NewBlockChainConfig(datadir string) *BlockChainConfig {
	cnf := &BlockChainConfig{
		datadir: datadir,
	}
	return cnf
}
