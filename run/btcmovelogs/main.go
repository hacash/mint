package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
	"github.com/hacash/mint/coinbase"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

/**
 * BTC 单向转移至 Hacash 主网，验证日志接口
 */

/**

// test

export GOPATH=/media/yangjie/500GB/hacash/go
cd mint/run/btcmovelogs
go run main.go

passwd 123456
gentx btcmove 1 1001 1596702752 0 1 1048576 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 8deb5180a3388fee4991674c62705041616980e76288a8888b65530e41ccf90d 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 HAC4:244
gentx release_lockbls 000000000000000000000000000000000001 HAC1024:248 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 HAC1:248

go build -ldflags '-w -s' -o   hacash_btc_move_log_2021_04_25_01  mint/run/btcmovelogs/main.go

*/

var cacheDatas = make([]*stores.SatoshiGenesis, 0)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	seekAlllogFiles()

	listenport := "3366"
	if len(os.Args) > 1 {
		listenport = os.Args[1]
	}

	// http server listen
	go func() {
		servmux := http.NewServeMux()
		servmux.HandleFunc("/btcmovelogs", dealQuery)
		fmt.Println("ListenAndServe:", "0.0.0.0:"+listenport)
		http.ListenAndServe("0.0.0.0:"+listenport, servmux)
	}()

	s := <-c
	fmt.Println("Got signal:", s)

}

func exitError(err, textfile string, curtrsno, fl int64) {
	fmt.Println("[Error]", err, "\n[Error] position file:", textfile, ", trsno:", curtrsno, ", line:", fl)
	os.Exit(0)
}

func seekAlllogFiles() {

	var curtrsno int64 = 1
	var prevGenesis *stores.SatoshiGenesis = nil
	var filenum int = 0

	for i := 1; ; i++ {
		textfile := fmt.Sprintf("./btcmovelogs%d.txt", i)
		// read file
		file, err := os.Open(textfile)
		if err != nil {
			break
		}
		fmt.Printf("load log file %s\n", textfile)

		scanner := bufio.NewScanner(file)
		var fl int64 = 1
		for scanner.Scan() {
			line := scanner.Text()
			//fmt.Println(line)
			genis := parseGenesis(line, textfile, curtrsno, fl)
			//fmt.Println(genis)
			if genis == nil {
				continue
			}
			// verification
			ckok := checkGenesis(prevGenesis, genis, textfile, curtrsno, fl)
			if !ckok {
				return
			}
			// Add to cache
			cacheDatas = append(cacheDatas, genis)
			// next row
			prevGenesis = genis
			curtrsno++
			fl++
		}
		filenum++
		// Next file
		file.Close()
	}
	//fmt.Println(cacheDatas)
	fmt.Printf("all %d files total %d lines load logs ok.\n", filenum, len(cacheDatas))
}

func parseGenesis(line, textfile string, curtrsno, fl int64) *stores.SatoshiGenesis {

	fixline := strings.Replace(line, " ", "", -1)
	if len(fixline) == 0 {
		return nil // Ignore blank lines
	}
	dts := strings.Split(fixline, ",")
	if len(dts) != 8 {
		exrr := fmt.Sprintf("data format error \"%s\" ", line)
		exitError(exrr, textfile, curtrsno, fl)
		return nil
	}
	// Parse data
	nums := make([]int64, 6)
	for i := 0; i < 6; i++ {
		n, e := strconv.ParseInt(dts[i], 10, 0)
		if e != nil {
			exrr := fmt.Sprintf("data format error \"%s\" ", line)
			exitError(exrr, textfile, curtrsno, fl)
			return nil
		}
		nums[i] = n
	}
	// Check address and txhx
	addr, ae := fields.CheckReadableAddress(dts[6])
	if ae != nil {
		exrr := fmt.Sprintf("address format error \"%s\" ", dts[6])
		exitError(exrr, textfile, curtrsno, fl)
		return nil
	}
	trshx, te := hex.DecodeString(dts[7])
	if te != nil {
		exrr := fmt.Sprintf("tx hash format error \"%s\" ", dts[7])
		exitError(exrr, textfile, curtrsno, fl)
		return nil
	}
	if len(trshx) != 32 {
		exrr := fmt.Sprintf("tx hash length error \"%s\" ", dts[7])
		exitError(exrr, textfile, curtrsno, fl)
		return nil
	}
	// return
	return &stores.SatoshiGenesis{
		fields.VarUint4(nums[0]),
		fields.VarUint4(nums[1]),
		fields.VarUint5(nums[2]),
		fields.VarUint4(nums[3]),
		fields.VarUint4(nums[4]),
		fields.VarUint4(nums[5]),
		*addr,
		trshx,
	}
}

func checkGenesis(prevGenesis, genis *stores.SatoshiGenesis, textfile string, curtrsno, fl int64) bool {

	//fmt.Println(int64(genis.TransferNo), curtrsno)

	// Verify trsno
	if int64(genis.TransferNo) != curtrsno {
		exrr := fmt.Sprintf("TransferNo need %d but got %d", curtrsno, genis.TransferNo)
		exitError(exrr, textfile, curtrsno, fl)
		return false
	}

	if prevGenesis == nil {
		if genis.TransferNo != 1 || genis.BitcoinEffectiveGenesis != 0 {
			exrr := fmt.Sprintf("first line data error TransferNo:%d, BitcoinEffectiveGenesis:%d",
				genis.TransferNo, genis.BitcoinEffectiveGenesis)
			exitError(exrr, textfile, curtrsno, fl)
			return false
		}
	} else {
		// Verify block height
		if genis.BitcoinBlockHeight < prevGenesis.BitcoinBlockHeight {
			exrr := fmt.Sprintf("BitcoinBlockHeight need no less than %d but got %d",
				prevGenesis.BitcoinBlockHeight, genis.BitcoinBlockHeight)
			exitError(exrr, textfile, curtrsno, fl)
			return false
		}
		// Validate timestamp
		if genis.BitcoinBlockTimestamp < prevGenesis.BitcoinBlockTimestamp {
			exrr := fmt.Sprintf("BitcoinBlockTimestamp need no less than %d but got %d",
				prevGenesis.BitcoinBlockTimestamp, genis.BitcoinBlockTimestamp)
			exitError(exrr, textfile, curtrsno, fl)
			return false
		}
		// Verify the number of BTCs that have been transferred
		effbtc := prevGenesis.BitcoinEffectiveGenesis + prevGenesis.BitcoinQuantity
		if genis.BitcoinEffectiveGenesis != effbtc {
			exrr := fmt.Sprintf("BitcoinEffectiveGenesis need %d but got %d",
				effbtc, genis.BitcoinEffectiveGenesis)
			exitError(exrr, textfile, curtrsno, fl)
			return false
		}
	}

	// Verify bitcoin quantity
	mvbtc := int64(genis.BitcoinQuantity)
	if mvbtc < 1 || mvbtc > 10000 {
		exrr := fmt.Sprintf("BitcoinQuantity need between %s but got %d", "1 ~ 1000", genis.BitcoinQuantity)
		exitError(exrr, textfile, curtrsno, fl)
		return false
	}
	// Verify the number of additional HACs
	var ttHacNum int64 = 0
	for i := genis.BitcoinEffectiveGenesis + 1; i <= genis.BitcoinEffectiveGenesis+genis.BitcoinQuantity; i++ {
		ttHacNum += coinbase.MoveBtcCoinRewardNumber(int64(i))
	}
	if int64(genis.AdditionalTotalHacAmount) != ttHacNum {
		exrr := fmt.Sprintf("AdditionalTotalHacAmount need %d but got %d", ttHacNum, genis.AdditionalTotalHacAmount)
		exitError(exrr, textfile, curtrsno, fl)
		return false
	}

	// Check successful
	return true
}

func dealQuery(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	trsnostr := request.Form.Get("trsno")
	limitstr := request.Form.Get("limit")
	var trsno int64 = 0
	if n, ok := strconv.ParseInt(trsnostr, 10, 0); ok == nil {
		trsno = n
	}
	var limit int64 = 1
	if n, ok := strconv.ParseInt(limitstr, 10, 0); ok == nil {
		limit = n
	}
	if trsno == 0 {
		//w.Write([]byte("not find"))
		return
	}
	// obtain
	seekIdx := (trsno - 1)
	// read
	// Print
	allretstr := []string{}
	for i := seekIdx; i < int64(len(cacheDatas)) && i < seekIdx+limit; i++ {
		genesis := cacheDatas[i]
		resstr := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%s,%s",
			genesis.TransferNo,
			genesis.BitcoinBlockHeight,
			genesis.BitcoinBlockTimestamp,
			genesis.BitcoinEffectiveGenesis,
			genesis.BitcoinQuantity,
			genesis.AdditionalTotalHacAmount,
			genesis.OriginAddress.ToReadable(),
			hex.EncodeToString(genesis.BitcoinTransferHash),
		)
		//
		allretstr = append(allretstr, resstr)
	}

	// Output results
	w.Write([]byte(strings.Join(allretstr, "|")))
	return

}
