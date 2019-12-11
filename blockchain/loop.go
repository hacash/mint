package blockchain

import (
	"fmt"
	"os"
)

func (bc *BlockChain) loop() {

	for {
		select {

		case newblk := <-bc.newBlockArriveQueueCh:
			//fmt.Println(newblk)
			err := bc.TryValidateAppendNewBlockToChainStateAndStore(newblk)
			if err != nil {
				fmt.Println("try Append New Block To Chain Error:", err)
				os.Exit(0) // test
			} else {
				//fmt.Println("successfully insert block height:", newblk.GetHeight(), "hash:", hex.EncodeToString(newblk.Hash()))
				//os.Exit(0) // test
				//time.Sleep(time.Second)
			}

		case newtx := <-bc.newTransactionArriveQueueCh:
			fmt.Println(newtx)

		}
	}

}
