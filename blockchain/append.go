package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/transactions"
	"github.com/hacash/mint"
	"github.com/hacash/mint/coinbase"
	"github.com/hacash/mint/difficulty"
	"time"
)

const (
	Time_format_layout = "01/02 15:04:05"
)

// interface api
func (bc *BlockChain) InsertBlock(newblock interfaces.Block, origin string) error {
	if origin != "" {
		newblock.SetOriginMark(origin)
	}
	return bc.tryValidateAppendNewBlockToChainStateAndStore(newblock)
}

// append block
func (bc *BlockChain) tryValidateAppendNewBlockToChainStateAndStore(newblock interfaces.Block) error {
	bc.insertLock.Lock()
	defer bc.insertLock.Unlock()
	// insert
	prevblock, e1 := bc.chainstate.ReadLastestBlockHeadAndMeta()
	if e1 != nil {
		return e1
	}
	newBlockHeight := newblock.GetHeight()
	newBlockTimestamp := newblock.GetTimestamp()
	newBlockHash := newblock.HashFresh()
	newBlockHashHexStr := newBlockHash.ToHex()
	errmsgprifix := fmt.Sprintf("Warning: try insert append new block height:%d, hx:%s to chain, ", newBlockHeight, newBlockHashHexStr)
	// check max size in p2p node message on get one
	// check height
	if newBlockHeight != prevblock.GetHeight()+1 {
		return fmt.Errorf(errmsgprifix+"Need block height %d but got %d.", prevblock.GetHeight()+1, newBlockHeight)
	}
	// check prev hash
	if bytes.Compare(newblock.GetPrevHash(), prevblock.HashFresh()) != 0 {
		prevhash := prevblock.HashFresh()
		newblkprevhash := newblock.GetPrevHash()
		return fmt.Errorf(errmsgprifix+"Need block prev hash %s but got %s.", prevhash.ToHex(), newblkprevhash.ToHex())
	}
	// check time now
	if int64(newBlockTimestamp) > int64(time.Now().Unix()) {
		createtime := time.Unix(int64(newBlockTimestamp), 0).Format(Time_format_layout)
		nowtime := time.Now().Format(Time_format_layout)
		return fmt.Errorf(errmsgprifix+"Block create timestamp cannot equal or more than now %s but got %s.", nowtime, createtime)
	}
	// check time prev
	if int64(newBlockTimestamp) <= int64(prevblock.GetTimestamp()) {
		prevtime := time.Unix(int64(prevblock.GetTimestamp()), 0).Format(Time_format_layout)
		currtime := time.Unix(int64(newBlockTimestamp), 0).Format(Time_format_layout)
		return fmt.Errorf(errmsgprifix+"Block create timestamp cannot equal or less than prev %s but got %s.", prevtime, currtime)
	}
	// check tx count
	if uint32(len(newblock.GetTransactions())) != newblock.GetTransactionCount() {
		return fmt.Errorf(errmsgprifix+"Transaction count wrong, accept %d, but got %d.",
			len(newblock.GetTransactions()),
			newblock.GetTransactionCount())
	}
	// check mkrl root
	newblktxs := newblock.GetTransactions()
	newblockRealMkrlRoot := blocks.CalculateMrklRoot(newblktxs)
	newblkmkrlroot := newblock.GetMrklRoot()
	if bytes.Compare(newblockRealMkrlRoot, newblkmkrlroot) != 0 {
		err := fmt.Errorf(errmsgprifix+"Need block mkrl root %s but got %s.", newblockRealMkrlRoot.ToHex(), newblkmkrlroot.ToHex())
		//fmt.Println(err); os.Exit(0)
		fmt.Println(err)
		for i, v := range newblktxs {
			fmt.Println("tx", i, v.Hash())
		}
		fmt.Println("- - - - - - - - - - - - - mkrl error block body hex - - - - - - - - - - - - -")
		testprintblkdts, _ := newblock.Serialize()
		fmt.Println(hex.EncodeToString(testprintblkdts))
		fmt.Println("- - - - - - - - - - - - - mkrl error block body hex end - - - - - - - - - - -")
		return err
	}
	//fmt.Println("mkrl:", newblockRealMkrlRoot.ToHex(), newblkmkrlroot.ToHex())
	// check coinbase tx
	if len(newblktxs) < 1 {
		return fmt.Errorf(errmsgprifix + "Block not included any transactions.")
	}
	var newblockCoinbaseReward *fields.Amount
	if cb1, ok := newblktxs[0].(*transactions.Transaction_0_Coinbase); ok {
		newblockCoinbaseReward = &cb1.Reward
	} else {
		return fmt.Errorf(errmsgprifix + "Not find coinbase tx in transactions at first.")
	}
	// check coinbase reward
	shouldrewards := coinbase.BlockCoinBaseReward(newBlockHeight)
	if newblockCoinbaseReward.NotEqual(shouldrewards) {
		return fmt.Errorf(errmsgprifix+"Block coinbase reward need %s got %s.", shouldrewards, newblockCoinbaseReward.ToFinString())
	}
	// check hash difficulty
	targetDiffHash, _, _, e5 := bc.CalculateNextDiffculty(prevblock)
	if e5 != nil {
		return e5
	}
	if difficulty.CheckHashDifficultySatisfy(newBlockHash, targetDiffHash) == false {
		return fmt.Errorf(errmsgprifix+"Maximum accepted hash diffculty is %s but got %s.", hex.EncodeToString(targetDiffHash), newBlockHashHexStr)
	}
	// 检查验证全部交易签名
	sigok, e6 := newblock.VerifyNeedSigns()
	if e6 != nil {
		return e6
	}
	if sigok != true {
		return fmt.Errorf(errmsgprifix + "Block signature verify faild.")
	}
	// 判断包含交易是否已经存在 和 区块大小 和 交易时间戳
	timenow := uint64(time.Now().Unix())
	totaltxsize := uint32(0)
	//blockstore := bc.chainstate.BlockStore()
	for i := 1; i < len(newblktxs); i++ { // ignore coinbase tx
		if newblktxs[i].GetTimestamp() > timenow {
			return fmt.Errorf(errmsgprifix+"Tx timestamps %d is not more than now %d.", newblktxs[i].GetTimestamp(), timenow)
		}
		txhashnofee := newblktxs[i].Hash()
		ok, e := bc.chainstate.CheckTxHash(txhashnofee)
		if e != nil {
			return e
		}
		if ok == true {
			return fmt.Errorf(errmsgprifix+"Tx %s is exist.", txhashnofee.ToHex())
		}
		totaltxsize += newblktxs[i].Size()
	}
	if totaltxsize > mint.SingleBlockMaxSize {
		return fmt.Errorf(errmsgprifix+"Txs total size %d is overflow max size %d.", totaltxsize, mint.SingleBlockMaxSize)
	}
	// 执行验证区块的每一笔交易
	newBlockChainState, e7 := bc.chainstate.NewSubBranchTemporaryChainState()
	if e7 != nil {
		return e7
	}
	newBlockChainState.SetPendingBlockHeight(newBlockHeight) // set pending
	newBlockChainState.SetPendingBlockHash(newBlockHash)     // set pending
	defer newBlockChainState.DestoryTemporary()
	// setup debug
	if newblock.GetHeight() == 1 {
		setupDebugChainState(newBlockChainState) // first state setup
	}
	err2 := newblock.WriteinChainState(newBlockChainState)
	if err2 != nil {
		return err2
	}
	// 储存状态数据
	err3 := bc.chainstate.MergeCoverWriteChainState(newBlockChainState)
	if err3 != nil {
		return err3
	}
	diamondCreate, err5 := newBlockChainState.GetPendingSubmitStoreDiamond()
	if err5 != nil {
		return err5
	}
	err4 := bc.chainstate.SubmitDataStoreWriteToInvariableDisk(newblock)
	if err4 != nil {
		return err4
	}
	// ok
	// send feed
	if diamondCreate != nil {
		// fmt.Println("diamondCreate bc.diamondCreateFeed.Send(diamondCreate), ", diamondCreate, diamondCreate.Diamond, diamondCreate.Number)
		bc.diamondCreateFeed.Send(diamondCreate)
	}

	orimark := newblock.OriginMark()
	if orimark != "" && orimark != "sync" {
		// 发送新区快通知
		bc.validatedBlockInsertFeed.Send(interfaces.Block(newblock))
	}

	// return
	return nil
}

// 不安全的升级数据库
func (bc *BlockChain) insertBlockToChainStateAndStoreUnsafe(newblock interfaces.Block) error {
	// 状态
	newBlockChainState, e7 := bc.chainstate.NewSubBranchTemporaryChainState()
	if e7 != nil {
		return e7
	}
	newBlockChainState.SetPendingBlockHeight(newblock.GetHeight()) // set pending
	newBlockChainState.SetPendingBlockHash(newblock.Hash())        // set pending
	defer newBlockChainState.DestoryTemporary()
	// setup debug
	if newblock.GetHeight() == 1 {
		setupDebugChainState(newBlockChainState) // first state setup
	}
	err2 := newblock.WriteinChainState(newBlockChainState)
	if err2 != nil {
		return err2
	}
	// 储存状态数据
	err3 := bc.chainstate.MergeCoverWriteChainState(newBlockChainState)
	if err3 != nil {
		return err3
	}
	err4 := bc.chainstate.SubmitDataStoreWriteToInvariableDisk(newblock)
	if err4 != nil {
		return err4
	}
	// ok
	return nil
}

// first debug amount
func setupDebugChainState(chainstate interfaces.ChainStateOperation) {

	addr1, _ := fields.CheckReadableAddress("12vi7DEZjh6KrK5PVmmqSgvuJPCsZMmpfi")
	addr2, _ := fields.CheckReadableAddress("1LsQLqkd8FQDh3R7ZhxC5fndNf92WfhM19")
	addr3, _ := fields.CheckReadableAddress("1NUgKsTgM6vQ5nxFHGz1C4METaYTPgiihh")
	amt1, _ := fields.NewAmountFromFinString("ㄜ1:244")
	amt2, _ := fields.NewAmountFromFinString("ㄜ12:244")
	chainstate.BalanceSet(*addr1, stores.NewBalanceWithAmount(amt2))
	chainstate.BalanceSet(*addr2, stores.NewBalanceWithAmount(amt1))
	chainstate.BalanceSet(*addr3, stores.NewBalanceWithAmount(amt1))

}
