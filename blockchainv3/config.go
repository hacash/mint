package blockchainv3

import (
	"github.com/hacash/core/sys"
)

type BlockChainConfig struct {
	cnffile *sys.Inicnf

	Datadir          string
	RollbackToHeight uint64
	// btc move
	DownloadBTCMoveLogUrl     string
	DisableDownloadBTCMoveLog bool // 不下载日志

	// 数据库重建模式
	DatabaseVersionRebuildMode bool
}

func NewEmptyBlockChainConfig() *BlockChainConfig {
	cnf := &BlockChainConfig{
		RollbackToHeight:           0,
		DownloadBTCMoveLogUrl:      "",
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

	sec2 := cnffile.Section("btcmovecheck")
	if sec2.Key("enable").MustBool(false) {
		cnf.DownloadBTCMoveLogUrl = sec2.Key("logs_url").MustString("")
	}
	// 不下载日志
	cnf.DisableDownloadBTCMoveLog = sec2.Key("disable_download").MustBool(false)

	cnf.Datadir = cnffile.MustDataDirWithVersion()

	return cnf

}
