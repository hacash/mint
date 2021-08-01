package blockchain

import (
	"fmt"
	"os"
	"path"
)

func (bc *BlockChain) RollbackToBlockHeight(targetblockheight uint64) (uint64, error) {

	fmt.Print("[BlockChain] Rollback to block height: ", targetblockheight, "... \n")

	// 依次读取区块，并插入新状态
	fmt.Print("[Database] Replay the block (NOT resynchronized), closing halfway will result in data corruption, Please wait and do not close the program...\n[Database] Checking block height:          0")

	// 关闭旧的
	bc.chainstate.Close()

	// 重命名目录
	olddir := path.Join(path.Dir(bc.config.Datadir), "rbnk")
	e0 := os.RemoveAll(olddir)
	if e0 != nil {
		return 0, e0
	}

	e1 := os.Rename(bc.config.Datadir, olddir)
	if e1 != nil {
		return 0, e1
	}

	// replay
	newbc, e2 := UpdateDatabaseReturnBlockChain(bc.config.cnffile, olddir, targetblockheight)

	e3 := os.RemoveAll(olddir)
	if e3 != nil {
		return 0, e3
	}

	if e2 != nil {
		return 0, e2
	}

	// copy state
	bc.ReplaceSelf(newbc)

	// ok
	return targetblockheight, nil

}

// 旧版
/*
func (bc *BlockChain) RollbackToBlockHeightOld(targetblockheight uint64) (uint64, error) {

	panic("RecoverChainState be deprecated")

	lastest, err := bc.chainstate.ReadLastestBlockHeadAndMeta()
	if err != nil {
		return 0, err
	}
	last_hei := lastest.GetHeight()
	if last_hei <= 1 || last_hei <= targetblockheight {
		return last_hei, nil // end
	}
	fmt.Print("[BlockChain] Rollback to block height: ", targetblockheight, ", lastest height:", last_hei, "... ")
	blockstore := bc.chainstate.BlockStore()
	var rollbackBlock interfaces.Block = nil
	for i := lastest.GetHeight(); i >= 1; i-- {
		// read block
		_, blkdatas, e2 := blockstore.ReadBlockBytesByHeight(i, 0)
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
			}
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

		} else {
			// do rollback
			e4 := block.RecoverChainState(bc.chainstate)
			if e4 != nil {
				return 0, e4
			}
			// delete all trs data ptr
			e5 := blockstore.CancelUniteTransactions(block)
			if e5 != nil {
				return 0, e5
			}
			// do next
		}
	}

	return 0, fmt.Errorf("error break.")

}
*/
