package blockchainv3

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfacev3"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/transactions"
	"github.com/hacash/mint"
	"github.com/hacash/mint/coinbase"
	"github.com/hacash/mint/difficulty"
	"time"
)

const (
	block_time_format_layout = "01/02 15:04:05"
)

// 新建状态去插入区块
func (bc *BlockChain) forkStateWithAppendBlock(baseState interfacev3.ChainState, newblock interfacev3.Block) (interfacev3.ChainState, error) {
	// 检查区块高度和 prev 哈希等等
	prevPending := baseState.GetPending()
	prevblock := prevPending.GetPendingBlockHead()
	if prevblock == nil {
		return nil, fmt.Errorf("baseState.pendingBlock cannot be nil.")
	}
	prevblockHeight := prevPending.GetPendingBlockHeight()
	prevblockHash := prevPending.GetPendingBlockHash()
	// 开始写入
	newBlockHeight := newblock.GetHeight()
	newBlockTimestamp := newblock.GetTimestamp()
	newBlockHash := newblock.HashFresh()
	newBlockHashHexStr := newBlockHash.ToHex()
	errmsgprifix := fmt.Sprintf("Warning: try insert append new block height:%d, hx:%s to chain, ", newBlockHeight, newBlockHashHexStr)
	// check max size in p2p node message on get one
	// check height
	if newBlockHeight != prevblockHeight+1 {
		return nil, fmt.Errorf(errmsgprifix+"Need block height %d but got %d.", prevblockHeight+1, newBlockHeight)
	}
	// check prev hash
	if bytes.Compare(newblock.GetPrevHash(), prevblockHash) != 0 {
		newblkprevhash := newblock.GetPrevHash()
		return nil, fmt.Errorf(errmsgprifix+"Need block prev hash %s but got %s.", prevblockHash.ToHex(), newblkprevhash.ToHex())
	}
	// check time now
	if int64(newBlockTimestamp) > int64(time.Now().Unix()) {
		createtime := time.Unix(int64(newBlockTimestamp), 0).Format(block_time_format_layout)
		nowtime := time.Now().Format(block_time_format_layout)
		return nil, fmt.Errorf(errmsgprifix+"Block create timestamp cannot equal or more than now %s but got %s.", nowtime, createtime)
	}
	// check time prev
	if int64(newBlockTimestamp) <= int64(prevblock.GetTimestamp()) {
		prevtime := time.Unix(int64(prevblock.GetTimestamp()), 0).Format(block_time_format_layout)
		currtime := time.Unix(int64(newBlockTimestamp), 0).Format(block_time_format_layout)
		return nil, fmt.Errorf(errmsgprifix+"Block create timestamp cannot equal or less than prev %s but got %s.", prevtime, currtime)
	}
	// check tx count
	newblktxs := newblock.GetTrsList()
	if uint32(len(newblktxs)) != newblock.GetTransactionCount() {
		return nil, fmt.Errorf(errmsgprifix+"Transaction count wrong, accept %d, but got %d.",
			len(newblktxs),
			newblock.GetTransactionCount())
	}
	// check mkrl root
	var txallhxs = make([]fields.Hash, len(newblktxs))
	for i, v := range newblktxs {
		txallhxs[i] = v.HashWithFee()
	}
	newblockRealMkrlRoot := blocks.CalculateMrklRootByHashWithFee(txallhxs)
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
		return nil, err
	}
	//fmt.Println("mkrl:", newblockRealMkrlRoot.ToHex(), newblkmkrlroot.ToHex())
	// check coinbase tx
	if len(newblktxs) < 1 {
		return nil, fmt.Errorf(errmsgprifix + "Block not included any transactions.")
	}
	var newblockCoinbaseReward *fields.Amount
	if cb1, ok := newblktxs[0].(*transactions.Transaction_0_Coinbase); ok {
		newblockCoinbaseReward = &cb1.Reward
	} else {
		return nil, fmt.Errorf(errmsgprifix + "Not find coinbase tx in transactions at first.")
	}
	// check coinbase reward
	shouldrewards := coinbase.BlockCoinBaseReward(newBlockHeight)
	if newblockCoinbaseReward.NotEqual(shouldrewards) {
		return nil, fmt.Errorf(errmsgprifix+"Block coinbase reward need %s got %s.", shouldrewards, newblockCoinbaseReward.ToFinString())
	}
	// check hash difficulty
	targetDiffHash, _, _, e5 := bc.CalculateNextDiffculty(prevblock)
	if e5 != nil {
		return nil, e5
	}
	if difficulty.CheckHashDifficultySatisfy(newBlockHash, targetDiffHash) == false {
		return nil, fmt.Errorf(errmsgprifix+"Maximum accepted hash diffculty is %s but got %s.", hex.EncodeToString(targetDiffHash), newBlockHashHexStr)
	}
	// 检查验证全部交易签名
	sigok, e6 := newblock.VerifyNeedSigns()
	if e6 != nil {
		return nil, e6
	}
	if sigok != true {
		return nil, fmt.Errorf(errmsgprifix + "Block signature verify faild.")
	}
	// 判断包含交易是否已经存在 和 区块大小 和 交易时间戳
	timenow := uint64(time.Now().Unix())
	totaltxsize := uint32(0)
	//blockstore := bc.chainstate.BlockStore()
	for i := 1; i < len(newblktxs); i++ { // ignore coinbase tx
		if newblktxs[i].GetTimestamp() > timenow {
			return nil, fmt.Errorf(errmsgprifix+"Tx timestamps %d is not more than now %d.", newblktxs[i].GetTimestamp(), timenow)
		}
		txhashnofee := newblktxs[i].Hash()
		ok, e := baseState.CheckTxHash(txhashnofee)
		if e != nil {
			return nil, e
		}
		if ok == true {
			return nil, fmt.Errorf(errmsgprifix+"Tx %s is exist.", txhashnofee.ToHex())
		}
		totaltxsize += newblktxs[i].Size()
	}
	if totaltxsize > mint.SingleBlockMaxSize {
		return nil, fmt.Errorf(errmsgprifix+"Txs total size %d is overflow max size %d.", totaltxsize, mint.SingleBlockMaxSize)
	}
	// 执行验证区块的每一笔交易
	// fork state
	newBlockChainState, e := baseState.ForkNextBlock(newblock.GetHeight(), newblock.Hash(), newblock)
	if e != nil {
		return nil, e
	}
	//newBlockChainState.SetPendingBlockHeight(newBlockHeight) // set pending
	//newBlockChainState.SetPendingBlockHash(newBlockHash)     // set pending
	//defer newBlockChainState.Destory()
	// setup debug
	if newblock.GetHeight() == 1 {
		setupDebugChainState(newBlockChainState) // first state setup
	}
	// 写入区块状态
	err2 := newblock.WriteInChainState(newBlockChainState)
	if err2 != nil {
		return nil, err2
	}

	// 储存状态数据
	/*
		err3 := bc.chainstate.MergeCoverWriteChainState(newBlockChainState)
		if err3 != nil {
			return nil, err3
		}
		diamondCreate, err5 := newBlockChainState.GetPendingSubmitStoreDiamond()
		if err5 != nil {
			return nil, err5
		}
		err4 := bc.chainstate.SubmitDataStoreWriteToInvariableDisk(newblock)
		if err4 != nil {
			return nil, err4
		}
	*/

	/*
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
	*/

	// return ok
	return newBlockChainState, nil
}

// first debug amount
func setupDebugChainState(chainstate interfacev3.ChainStateOperation) {
	addr1, _ := fields.CheckReadableAddress("12vi7DEZjh6KrK5PVmmqSgvuJPCsZMmpfi")
	addr2, _ := fields.CheckReadableAddress("1LsQLqkd8FQDh3R7ZhxC5fndNf92WfhM19")
	addr3, _ := fields.CheckReadableAddress("1NUgKsTgM6vQ5nxFHGz1C4METaYTPgiihh")
	amt1, _ := fields.NewAmountFromFinString("ㄜ1:244")
	amt2, _ := fields.NewAmountFromFinString("ㄜ12:244")
	chainstate.BalanceSet(*addr1, stores.NewBalanceWithAmount(amt2))
	chainstate.BalanceSet(*addr2, stores.NewBalanceWithAmount(amt1))
	chainstate.BalanceSet(*addr3, stores.NewBalanceWithAmount(amt1))
}
