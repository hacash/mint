package blockchain

import (
	"fmt"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/interfaces"
)

func (bc *BlockChain) RollbackToBlockHeight(targetblockheight uint64) (uint64, error) {
	lastest, err := bc.chainstate.ReadLastestBlockHeadAndMeta()
	if err != nil {
		return 0, err
	}
	last_hei := lastest.GetHeight()
	if last_hei <= 1 || last_hei <= targetblockheight {
		return last_hei, nil // end
	}
	fmt.Print("[BlockChain] Rollback to block height:", targetblockheight, "lastest height:", last_hei, "... ")
	var rollbackBlock interfaces.Block = nil
	for i := lastest.GetHeight(); i >= 1; i-- {
		// read block
		_, blkdatas, e2 := bc.chainstate.BlockStore().ReadBlockBytesByHeight(i, 0)
		if e2 != nil {
			return 0, e2
		}
		block, _, e3 := blocks.ParseBlock(blkdatas, 0)
		if e3 != nil {
			return 0, e3
		}
		rollbackBlock = block
		// check
		if i == targetblockheight {
			// set status
			e5 := bc.chainstate.SetLastestBlockHeadAndMeta(rollbackBlock)
			if e5 != nil {
				return 0, e5
			} else {
				// save status
				rollerr1 := bc.chainstate.IncompleteSaveLastestBlockHeadAndMeta()
				if rollerr1 != nil {
					return 0, rollerr1
				}
				rollerr2 := bc.chainstate.IncompleteSaveLastestDiamond()
				if rollerr2 != nil {
					return 0, rollerr2
				}
				fmt.Print("successfully !\n")
				return i, nil // ok finish
			}
		} else {
			// do rollback
			e4 := block.RecoverChainState(bc.chainstate)
			if e4 != nil {
				return 0, e4
			}
			// do next
		}
	}

	return 0, fmt.Errorf("error break.")

}
