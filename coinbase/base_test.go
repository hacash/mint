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

	var ttp float64 = 1

	var tthacount float64 = 0
	var ttdesc float64 = 0

	for {
		// 计算可借出数量
		loanhac, predeshac := CalculationOfInterestDiamondMortgageLoanAmount(ttp)

		// 实际得到
		realgot := loanhac - predeshac

		// 实际年利率 %
		yearrate := predeshac / realgot * 100

		// 统计
		tthacount += realgot
		ttdesc += predeshac

		// 数值打印
		fmt.Printf("抵押总比例: %.2f/100 , HAC可借: %0.2f , 预付息: %0.2f , 实得: %0.2f , 年利率: %0.2f%%\n", ttp, loanhac, predeshac, realgot, yearrate)

		if ttp >= 100 {
			break
		} else {
			ttp += 1
			continue
		}

		// 步进
		if ttp >= 99 {
			break
		} else if ttp >= 10 {
			ttp += 1
			continue
		} else if ttp >= 5 {
			ttp += 0.5
			continue
		} else if ttp >= 3 {
			ttp += 0.2
			continue
		} else if ttp >= 2 {
			ttp += 0.1
			continue
		} else if ttp >= 1 {
			ttp += 0.05
			continue
		}
	}

	fmt.Println("总增发:", tthacount, "预付息:", ttdesc)

	time.Sleep(time.Second)

}
