package blockchainv3

import (
	"fmt"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

func (bc *ChainKernel) SetConfirmTxIndexer(indexer interfaces.ConfirmTxIndexer) {
	bc.txindexer = indexer
}

func (bc *ChainKernel) execTxIndexer(blkhx fields.Hash) {
	//fmt.Printf("---- execTxIndexer 1 ---- %s\n", blkhx.ToHex())
	if bc.txindexer == nil {
		return
	}
	//fmt.Println("---- execTxIndexer 2 ----")
	// load block
	var block interfaces.Block = nil
	blkptr, ok := bc.recentArrivedBlocks.Load(string(blkhx))
	if ok {
		block = blkptr.(interfaces.Block)
	}
	//fmt.Println("---- execTxIndexer 3 ----")
	// read from dist
	if block == nil {
		blkbts, e := bc.blockstore.ReadBlockBytesByHash(blkhx)
		if e != nil {
			fmt.Println(e)
			return // error
		}
		block, _, e = blocks.ParseBlock(blkbts, 0)
		if e != nil {
			fmt.Println(e)
			return // error
		}
	}
	// loop tx
	txs := block.GetTrsList()
	//fmt.Printf("---- execTxIndexer 4 txs: %d ----\n", len(txs))
	if len(txs) <= 1 {
		return // block empty
	}
	for i := 1; i < len(txs); i++ {
		var tx = txs[i]
		resmark := bc.txindexer.ScanTx(block, tx)
		if resmark == 1 {
			// next block
			break
		}
		// next tx
		continue
	}

	// next block

}
