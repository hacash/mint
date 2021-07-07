package coinbase

import (
	"fmt"
	"testing"
	"time"
)

func Test_t1(t *testing.T) {

	//PrintMoveBtcCoinRewardNumberTable()

	for i := 1; i < 20; i++ {
		n := i + 16380
		fmt.Println(n, MoveBtcCoinRewardNumber(int64(n)))
	}

}

// 比特币抵押借贷数据
func Test_t2(t *testing.T) {

	var ttp float64 = 0
	var step float64 = 0.5
	var tthacount float64 = 0
	var ttdesc float64 = 0

	fmt.Printf("|抵押总比例|HAC可借数|预付利息|实际借得|年利率|\n")
	fmt.Printf("|---|---|---|---|---|\n")

	for {
		// 计算可借出数量
		loanhac, predeshac := CalculationOfInterestBitcoinMortgageLoanAmount(ttp)

		// 实际得到
		realgot := loanhac - predeshac

		// 实际年利率 %
		yearrate := predeshac / realgot * 100

		// 统计
		tthacount += realgot
		ttdesc += predeshac

		// 数值打印
		fmt.Printf("|%.2f%% | %0.2f | %0.2f | %0.2f | %0.2f%%|\n", ttp, loanhac, predeshac, realgot, yearrate)
		if ttp < 11 {
			ttp += 0.5
		} else {
			ttp += 1
		}

		if ttp >= 99.8 {
			break
		} else {
			// ttp += step
			continue
		}

	}

	fmt.Println("\n总增发:", tthacount*step, "预付息:", ttdesc*step)

	time.Sleep(time.Second)

}
