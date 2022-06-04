package blockchain

import (
	"fmt"
	"github.com/hacash/chain/blockstorev2"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/sys"
	"sync"
)

func UpdateDatabaseReturnBlockChain(ini *sys.Inicnf, olddatadir string, maxtarhei uint64, isclosenew bool) (*BlockChain, error) {

	// Start upgrade
	oldblockdatadir := olddatadir + "/blockstore"
	cnf1 := blockstorev2.NewEmptyBlockStoreConfig()
	cnf1.Datadir = oldblockdatadir
	oldblockDB, e0 := blockstorev2.NewBlockStoreForUpdateDatabaseVersion(cnf1)
	if e0 != nil {
		// Error occurred, return
		return nil, fmt.Errorf("Check And Update Blockchain Database Version Error: %s", e0.Error())
	}
	defer oldblockDB.Close()

	// Create new status
	bccnf := NewBlockChainConfig(ini)
	bccnf.DatabaseVersionRebuildMode = true // Database rebuild mode
	newblockchain, e1 := NewBlockChain(bccnf)
	if e1 != nil {
		return nil, fmt.Errorf("Check And Update Blockchain Database Version, NewBlockChain Error: %s", e1.Error())
		// Error occurred, return
	}
	// Mode recovery
	defer func() {
		bccnf.DatabaseVersionRebuildMode = false                   // Mode recovery
		newblockchain.State().SetDatabaseVersionRebuildMode(false) // Mode recovery
		// External decision whether to close
		if isclosenew {
			newblockchain.Close()
		}
	}()

	// Parallel read and write
	updateDataCh := make(chan []byte, 50)
	updateBlockCh := make(chan interfacev2.Block, 50)
	finishWait := sync.WaitGroup{}
	finishWait.Add(3)

	// Read data
	go func() {
		readblockhei := uint64(0)
		for {
			readblockhei++
			//fmt.Println("1")
			_, body, e := oldblockDB.ReadBlockBytesLengthByHeight(readblockhei, 0)
			if e != nil {
				fmt.Println("Check And Update Blockchain Database Version, ReadBlockBytesLengthByHeight Error:", e.Error())
				break // Error occurred, return
			}
			if len(body) == 0 {
				break // End all
			}
			// Write data
			updateDataCh <- body
			// Judge maximum synchronization
			if maxtarhei > 0 && maxtarhei <= readblockhei {
				break // Complete all
			}
		}
		// Read complete
		updateDataCh <- nil
		finishWait.Done()
	}()

	// Parsing block
	go func() {
		for {
			body := <-updateDataCh
			if body == nil {
				break // complete
			}
			//fmt.Println("3")
			// Parsing block
			blk, _, e2 := blocks.ParseBlock(body, 0)
			if e2 != nil {
				fmt.Println("Check And Update Blockchain Database Version, ParseBlock Error:", e2.Error())
				break // Error occurred, return
			}
			// Write data
			updateBlockCh <- blk.(interfacev2.Block)
		}
		// Read complete
		updateBlockCh <- nil
		finishWait.Done()
	}()

	// Write block data
	go func() {
		readblockhei := uint64(1)
		for {
			blk := <-updateBlockCh
			if blk == nil {
				fmt.Printf("\b\b\b\b\b\b\b\b\b\b%10d", readblockhei)
				break // complete
			}

			//fmt.Println("4")
			// Insert block (upgrade mode)
			e3 := newblockchain.insertBlockToChainStateAndStoreUnsafe(blk)
			if e3 != nil {
				fmt.Println("Check And Update Blockchain Database Version, InsertBlock Error:", e3.Error())
				break // Error occurred, return
			}
			//fmt.Println("5")
			// Print
			if readblockhei%1000 == 0 {
				//fmt.Printf("%d", readblockhei)
				fmt.Printf("\b\b\b\b\b\b\b\b\b\b%10d", readblockhei)
			}
			//fmt.Println("6")
			// next block
			readblockhei++
		}
		// Insert end
		finishWait.Done()
	}()

	finishWait.Wait()

	return newblockchain, nil
}

// Check upgrade database version
func CheckAndUpdateBlockchainDatabaseVersion(ini *sys.Inicnf) error {
	curversion, compatible := ini.GetDatabaseVersion()
	_, has := ini.MustDataDirCheckVersion(curversion)
	if has {
		return nil // The current version already exists. Return normally
	}
	// Upgrade required, check historical version
	olddir := ""
	oldversion := curversion - 1
	for {
		if oldversion < compatible {
			// It is lower than the minimum compatible version, which indicates that the blocks are forked and must be resynchronized from the network
			return nil
		}
		olddir, has = ini.MustDataDirCheckVersion(oldversion)
		if has {
			break
		}
		oldversion--
	}

	// Read blocks in sequence and insert new status
	fmt.Printf("[Database] Upgrade blockchain database version v%d to v%d, block data is NOT resynchronized, Please wait and do not close the program...\n[Database] Checking block height:          0", oldversion, curversion)

	_, e := UpdateDatabaseReturnBlockChain(ini, olddir, 0, true)
	if e != nil {
		err := fmt.Errorf("Check And Update Blockchain Database Version, NewBlockChain Error: %s\n", e.Error())
		fmt.Println(err.Error())
		// Error occurred, return
		return err
	}

	//fmt.Println("7", olddir)
	// Delete old version
	// defer os.RemoveAll(olddir)

	//fmt.Println("8")
	// All blocks synchronized successfully
	fmt.Printf(" all finished.\n[Database] version v%d => v%d upgrade successfully!\n", oldversion, curversion)

	return nil
}
