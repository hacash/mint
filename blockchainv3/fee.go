package blockchainv3

import (
	"github.com/hacash/core/interfaces"
	"github.com/hacash/mint"
	"sync"
)

/**
 * unit: zhu
 * return: > mint.MinTransactionFeePurity
 */
func (bc *ChainKernel) GetLatestAverageFeePurity() uint32 {
	var fpc = bc.averageFeePurityCounts
	var num = len(fpc)
	if num == 0 {
		return mint.MinTransactionFeePurity
	}
	var ftt = uint64(0)
	for i := 0; i < num; i++ {
		ftt += uint64(fpc[i])
	}
	// ok
	avgf := uint32(ftt / uint64(num))
	if avgf < mint.MinTransactionFeePurity {
		avgf = mint.MinTransactionFeePurity
	}
	return avgf
}

/**
 * deal AverageFeePurity
 */
var handleAverageFeePurityLock sync.Mutex

func (bc *ChainKernel) handleAverageFeePurityByNewBlock(block interfaces.Block) {
	handleAverageFeePurityLock.Lock()
	defer handleAverageFeePurityLock.Unlock()

	//fmt.Println("handleAverageFeePurityByNewBlock", block.GetHeight())
	trslist := block.GetTrsList()
	txn := len(trslist)
	if txn <= 1 {
		return // empty block
	}
	ftt := uint64(0)
	cts := uint64(0)
	for i := 1; i < txn; i++ {
		trs := trslist[i]
		if trs.GetFee().NotEqual(trs.GetFeeOfMinerRealReceived()) {
			continue // not count burn 90% tx
		}
		ftt += uint64(trs.FeePurity())
		cts += 1
	}
	if cts <= 0 {
		return // no tx
	}
	// fee
	avgfeep := ftt / cts
	// ok
	bc.averageFeePurityCounts = append(bc.averageFeePurityCounts, uint32(avgfeep))
	if len(bc.averageFeePurityCounts) > 12 {
		bc.averageFeePurityCounts = bc.averageFeePurityCounts[1:] // max size 12 block about 1 hour
	}

	//

}
