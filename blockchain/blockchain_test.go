package blockchain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/blocks"
	"golang.org/x/net/websocket"
	"os"
	"strconv"
	"testing"
)

func Test_t1(t *testing.T) {

	testdir := "/home/shiqiujie/Desktop/Hacash/go/src/github.com/hacash/mint/blockchain/testdata"

	os.RemoveAll(testdir)

	cnf := NewEmptyBlockChainConfig()
	cnf.Datadir = testdir
	blockchain, e1 := NewBlockChain(cnf)
	if e1 != nil {
		fmt.Println(e1)
		return
	}
	blockchain.Start()

	// websocket

	wsConn, e2 := websocket.Dial("ws://127.0.0.1:3338/websocket", "ws", "http://127.0.0.1/")
	if e2 != nil {
		fmt.Println(e1)
		return
	}

	start_block_height := 1

	datasbuf := bytes.NewBuffer([]byte{})
	tagetdataslength := -1

	rdata := make([]byte, 5000)
	for {
		if tagetdataslength == -1 {
			fmt.Println("getblocks  ---  start_block_height", start_block_height)
			wsConn.Write([]byte("getblocks " + strconv.Itoa(start_block_height)))
		}

		rn, e := wsConn.Read(rdata)
		if e != nil {
			fmt.Println(e)
			return
		}
		//fmt.Println("rn", rn)
		data := rdata[0:rn]
		if rn == 9 && bytes.Compare(data, []byte("endblocks")) == 0 {
			fmt.Println("got endblocks.")
			break
		}
		datasbuf.Write(data)
		if datasbuf.Len() < 4 {
			fmt.Println("datasbuf.Len() < 4, continue")
			continue
		}
		if tagetdataslength == -1 {
			tagetdataslength = int(binary.BigEndian.Uint32(data[0:4]))
		}

		if datasbuf.Len() == tagetdataslength+4 {
			datas := datasbuf.Bytes()
			start_block_height, e = newBlocksDataArrive(blockchain, datas[4:])
			fmt.Println("start_block_height", start_block_height)
			if e != nil {
				fmt.Println(e)
				return
			}
			tagetdataslength = -1
			datasbuf = bytes.NewBuffer([]byte{})
		}

	}

	fmt.Println("end of block.")

	//blockchain.TryValidateAppendNewBlockToChainStateAndStore()

}

func newBlocksDataArrive(blockchain *BlockChain, datas []byte) (int, error) {

	start_block_height := 1

	seek := uint32(0)
	for {
		if int(seek)+1 > len(datas) {
			break
		}
		//fmt.Println(seek, datas[seek:seek + 80])
		newblock, sk, e := blocks.ParseBlock(datas, seek)
		if e != nil {
			fmt.Println(e)
			return 0, e
		}
		//fmt.Println(newblock.GetHeight())
		seek = sk
		// do store
		blockchain.InsertBlock(newblock, "test")
		start_block_height = int(newblock.GetHeight()) + 1
	}
	// ok
	return start_block_height, nil
}
