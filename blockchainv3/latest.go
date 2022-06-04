package blockchainv3

import (
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

// Latest blocks (confirmed and immature)
func (bc *ChainKernel) LatestBlock() (interfaces.BlockHeadMetaRead, interfaces.BlockHeadMetaRead, error) {
	// Disk confirmed
	imm, e := bc.stateImmutable.ImmutableStatusRead()
	if e != nil {
		return nil, nil, e
	}
	return bc.stateCurrent.GetPending().GetPendingBlockHead(), imm.GetImmutableBlockHeadMeta(), nil
}

// Latest block diamonds
func (bc *ChainKernel) LatestDiamond() (*stores.DiamondSmelt, error) {
	status, e := bc.stateCurrent.LatestStatusRead()
	if e != nil {
		return nil, e
	}
	// ok
	return status.ReadLastestDiamond(), nil
}
