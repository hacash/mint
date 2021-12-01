package blockchainv3

import (
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

// 最新的区块（已确认的，和未成熟的）
func (bc *ChainKernel) LatestBlock() (interfaces.BlockHeadMetaRead, interfaces.BlockHeadMetaRead, error) {
	// 已经磁盘确认的
	imm, e := bc.stateImmutable.ImmutableStatusRead()
	if e != nil {
		return nil, nil, e
	}
	return bc.stateCurrent.GetPending().GetPendingBlockHead(), imm.GetImmutableBlockHeadMeta(), nil
}

// 最新的区块钻石
func (bc *ChainKernel) LatestDiamond() (*stores.DiamondSmelt, error) {
	status, e := bc.stateCurrent.LatestStatusRead()
	if e != nil {
		return nil, e
	}
	// ok
	return status.ReadLastestDiamond(), nil
}
