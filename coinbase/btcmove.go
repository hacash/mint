package coinbase

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

/**
 * 货币发行算法：BTC 转移增发 Hac
 */

/*

LV:  1     BTC:       1,       1     HAC: 1048576,  1048576
LV:  2     BTC:       2,       3     HAC:  524288,  2097152
LV:  3     BTC:       4,       7     HAC:  262144,  3145728
LV:  4     BTC:       8,      15     HAC:  131072,  4194304
LV:  5     BTC:      16,      31     HAC:   65536,  5242880
LV:  6     BTC:      32,      63     HAC:   32768,  6291456
LV:  7     BTC:      64,     127     HAC:   16384,  7340032
LV:  8     BTC:     128,     255     HAC:    8192,  8388608
LV:  9     BTC:     256,     511     HAC:    4096,  9437184
LV: 10     BTC:     512,    1023     HAC:    2048, 10485760
LV: 11     BTC:    1024,    2047     HAC:    1024, 11534336
LV: 12     BTC:    2048,    4095     HAC:     512, 12582912
LV: 13     BTC:    4096,    8191     HAC:     256, 13631488
LV: 14     BTC:    8192,   16383     HAC:     128, 14680064
LV: 15     BTC:   16384,   32767     HAC:      64, 15728640
LV: 16     BTC:   32768,   65535     HAC:      32, 16777216
LV: 17     BTC:   65536,  131071     HAC:      16, 17825792
LV: 18     BTC:  131072,  262143     HAC:       8, 18874368
LV: 19     BTC:  262144,  524287     HAC:       4, 19922944
LV: 20     BTC:  524288, 1048575     HAC:       2, 20971520
LV: 21     BTC: 1048576, 2097151     HAC:       1, 22020096

共转移 2097151 枚 BTC，增发 22020096 枚HAC

*/

func powf2(n int) int64 {
	res := math.Pow(2.0, float64(n))
	return int64(res)
}

// 第几枚BTC增发HAC数量（单位：枚）
func MoveBtcCoinRewardNumber(btcidx int64) int64 {
	var lvn = 21
	if btcidx == 1 {
		return powf2(lvn - 1)
	}
	if btcidx > powf2(lvn)-1 {
		return 1 // 最后始终增发一枚
	}
	var tarlv int
	for i := 0; i < lvn; i++ {
		l := powf2(i) - 1
		r := powf2(i+1) - 1
		if btcidx > l && btcidx <= r {
			tarlv = i + 1
			break
		}
	}
	return powf2(lvn - tarlv)
}

func PrintMoveBtcCoinRewardNumberTable() {

	const maxi = 21
	var (
		sumnum     = 1
		coin       = 1048576
		total_num  = 0
		total_coin = 0
	)
	for i := maxi; i > 0; i-- {
		total_coin += coin * sumnum
		total_num += sumnum
		fmt.Println("LV: "+padding(maxi-i+1, 2, " "),
			"    BTC: "+padding(sumnum, 7, " ")+", "+padding(total_num, 7, " "),
			"    HAC: "+padding(coin, 7, " ")+", "+padding(total_coin, 8, " "),
		)
		sumnum = sumnum * 2
		coin = coin / 2
	}

	fmt.Printf("共转移 %d 枚 BTC，增发 %d 枚 HAC \n", total_num, total_coin)

}

func padding(num, w int, prx string) string {
	str := strings.Repeat(prx, w) + strconv.Itoa(num)
	return string([]byte(str)[len(str)-w:])
}

/*

 */
