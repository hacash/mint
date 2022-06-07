package blockchain

import (
	"fmt"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/mint"
	"time"
)

func (bc *BlockChain) ValidateTransactionForTxPool(newtx interfaces.Transaction) error {
	newtxhash := newtx.Hash()
	txhxhex := newtxhash.ToHex()
	exist, e0 := bc.chainstate.CheckTxHash(newtxhash)
	//fmt.Println(exist, exist_tx_bytes)
	if e0 != nil {
		return e0
	}
	if exist {
		return fmt.Errorf("tx %s is exist in blockchain.", txhxhex)
	}
	// check
	if newtx.GetTimestamp() > uint64(time.Now().Unix()) {
		return fmt.Errorf("tx %s timestamp cannot more than now.", txhxhex)
	}
	// fee purity
	if newtx.FeePurity() < mint.MinTransactionFeePurityOfOneByte {
		return fmt.Errorf("tx %s handling fee is too low for miners to accept.", txhxhex)
	}
	// sign
	ok, e1 := newtx.VerifyAllNeedSigns()
	if !ok || e1 != nil {
		return fmt.Errorf("tx %s verify signature error", txhxhex)
	}
	// try run
	lastestBlock, e1 := bc.chainstate.ReadLastestBlockHeadAndMeta()
	if e1 != nil {
		return e1
	}
	// create temp state
	newTxState, e2 := bc.chainstate.NewSubBranchTemporaryChainState()
	if e2 != nil {
		return e2
	}
	newTxState.SetInTxPool(true)        // Mark is pool status
	defer newTxState.DestoryTemporary() // clean data
	// validate
	newTxState.SetPendingBlockHeight(lastestBlock.GetHeight() + 1)
	runerr := newtx.(interfacev2.Transaction).WriteinChainState(newTxState)
	if runerr != nil {
		return runerr
	}
	// have diamond in tx
	// diamond, _ := newTxState.GetPendingSubmitStoreDiamond()
	// ok pass check !
	return nil
}
