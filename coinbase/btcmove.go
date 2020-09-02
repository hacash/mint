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

LV:  1     BTC:       1,       1     HAC: 1048576,  1048576     LOCK: 1024w, 19.69y, 1024
LV:  2     BTC:       2,       3     HAC:  524288,  2097152     LOCK:  512w, 9.846y, 1024
LV:  3     BTC:       4,       7     HAC:  262144,  3145728     LOCK:  256w, 4.923y, 1024
LV:  4     BTC:       8,      15     HAC:  131072,  4194304     LOCK:  128w, 2.461y, 1024
LV:  5     BTC:      16,      31     HAC:   65536,  5242880     LOCK:   64w, 1.230y, 1024
LV:  6     BTC:      32,      63     HAC:   32768,  6291456     LOCK:   32w, 0.615y, 1024
LV:  7     BTC:      64,     127     HAC:   16384,  7340032     LOCK:   16w, 0.307y, 1024
LV:  8     BTC:     128,     255     HAC:    8192,  8388608     LOCK:    8w, 0.153y, 1024
LV:  9     BTC:     256,     511     HAC:    4096,  9437184     LOCK:    4w, 0.076y, 1024
LV: 10     BTC:     512,    1023     HAC:    2048, 10485760     LOCK:    2w, 0.038y, 1024
LV: 11     BTC:    1024,    2047     HAC:    1024, 11534336     LOCK:    1w, 0.019y, 1024
LV: 12     BTC:    2048,    4095     HAC:     512, 12582912     LOCK:    0w,     0y,  512
LV: 13     BTC:    4096,    8191     HAC:     256, 13631488     LOCK:    0w,     0y,  256
LV: 14     BTC:    8192,   16383     HAC:     128, 14680064     LOCK:    0w,     0y,  128
LV: 15     BTC:   16384,   32767     HAC:      64, 15728640     LOCK:    0w,     0y,   64
LV: 16     BTC:   32768,   65535     HAC:      32, 16777216     LOCK:    0w,     0y,   32
LV: 17     BTC:   65536,  131071     HAC:      16, 17825792     LOCK:    0w,     0y,   16
LV: 18     BTC:  131072,  262143     HAC:       8, 18874368     LOCK:    0w,     0y,    8
LV: 19     BTC:  262144,  524287     HAC:       4, 19922944     LOCK:    0w,     0y,    4
LV: 20     BTC:  524288, 1048575     HAC:       2, 20971520     LOCK:    0w,     0y,    2
LV: 21     BTC: 1048576, 2097151     HAC:       1, 22020096     LOCK:    0w,     0y,    1

共转移2097151枚BTC，增发22020096枚HAC


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
