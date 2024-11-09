package coinbase

import (
	"fmt"
	"testing"
	"time"
)

func Test_t0(t *testing.T) {
	for i := uint64(1); i < 1000000000; i++ {
		var blkh = i
		var rwdn = BlockCoinBaseRewardNumber(blkh)
		if blkh%10000 == 0 {
			fmt.Println(blkh, rwdn)
		}
	}
}

func Test_t1(t *testing.T) {

	//PrintMoveBtcCoinRewardNumberTable()

	for i := 1; i < 20; i++ {
		n := i + 16380
		fmt.Println(n, MoveBtcCoinRewardNumber(int64(n)))
	}

}

// Bitcoin mortgage loan data
func Test_t2(t *testing.T) {

	var ttp float64 = 0
	var step float64 = 0.5
	var tthacount float64 = 0
	var ttdesc float64 = 0

	fmt.Printf("|抵押总比例|HAC可借数|预付利息|实际借得|年利率|\n")
	fmt.Printf("|---|---|---|---|---|\n")

	for {
		// Calculate lendable quantity
		loanhac, predeshac := CalculationOfInterestBitcoinMortgageLoanAmount(ttp)

		// Actually obtained
		realgot := loanhac - predeshac

		// Effective annual interest rate%
		yearrate := predeshac / realgot * 100

		// Statistics
		tthacount += realgot
		ttdesc += predeshac

		// Numeric printing
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
