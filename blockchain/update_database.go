package blockchain

import (
	"fmt"
	"github.com/hacash/chain/blockstore"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/sys"
)

// 检查升级数据库版本
func CheckAndUpdateBlockchainDatabaseVersion(ini *sys.Inicnf) {
	curversion := sys.BlockChainStateDatabaseCurrentUseVersion
	_, has := ini.MustDataDirCheckVersion(curversion)
	if has {
		return // 当前版本已经存在，正常返回
	}
	// 需要升级，检查历史版本
	olddir := ""
	oldversion := curversion - 1
	for {
		if oldversion < sys.BlockChainStateDatabaseLowestCompatibleVersion {
			// 已经低于最低可兼容版本了，表示区块出现了分叉，必须全部从网络重新同步
			return
		}
		olddir, has = ini.MustDataDirCheckVersion(oldversion)
		if has {
			break
		}
		oldversion--
	}

	// 开始升级
	oldblockdatadir := olddir + "/blockstore"
	cnf1 := blockstore.NewEmptyBlockStoreConfig()
	cnf1.Datadir = oldblockdatadir
	oldblockDB, e0 := blockstore.NewBlockStoreForUpdateDatabaseVersion(cnf1)
	if e0 != nil {
		fmt.Println("Check And Update Blockchain Database Version Error:", e0.Error())
		return // 发生错误，返回
	}
	defer oldblockDB.Close()

	// 建立新状态
	bccnf := NewBlockChainConfig(ini)
	newblockchain, e1 := NewBlockChain(bccnf)
	if e1 != nil {
		fmt.Println("Check And Update Blockchain Database Version, NewBlockChain Error:", e1.Error())
		return // 发生错误，返回
	}
	defer newblockchain.Close()

	// 依次读取区块，并插入新状态
	readblockhei := uint64(1)
	fmt.Print("[Database] Upgrade blockchain database version, Please wait and do not close the program...\n[Database] Checking block height:          0")
	for {
		//fmt.Println("1")
		_, body, e := oldblockDB.ReadBlockBytesByHeight(readblockhei, 0)
		if e != nil {
			fmt.Println("Check And Update Blockchain Database Version, ReadBlockBytesByHeight Error:", e.Error())
			return // 发生错误，返回
		}
		//fmt.Println("2")
		if len(body) == 0 {
			fmt.Printf("\b\b\b\b\b\b\b\b\b\b%10d", readblockhei-1)
			break // 已经读取完毕
		}
		//fmt.Println("3")
		// 解析区块
		blk, _, e2 := blocks.ParseBlock(body, 0)
		if e2 != nil {
			fmt.Println("Check And Update Blockchain Database Version, ParseBlock Error:", e2.Error())
			return // 发生错误，返回
		}
		//fmt.Println("4")
		// 插入区块（升级模式）
		e3 := newblockchain.insertBlockToChainStateAndStoreUnsafe(blk)
		if e3 != nil {
			fmt.Println("Check And Update Blockchain Database Version, InsertBlock Error:", e3.Error())
			return // 发生错误，返回
		}
		//fmt.Println("5")
		// 打印
		if readblockhei%1000 == 0 {
			//fmt.Printf("%d", readblockhei)
			fmt.Printf("\b\b\b\b\b\b\b\b\b\b%10d", readblockhei)
		}
		//fmt.Println("6")
		// next block
		readblockhei++
	}

	//fmt.Println("7", olddir)
	// 删除旧版本
	// defer os.RemoveAll(olddir)

	//fmt.Println("8")
	// 全部区块同步成功
	fmt.Printf(" all finished.\n[Database] version v%d => v%d upgrade successfully!\n", oldversion, curversion)
}