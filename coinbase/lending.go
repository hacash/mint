package coinbase

// 钻石抵押借贷： 计算可借数量
// totalLendingPercentage 已经借出的总百分比
// 返回可借数量和预付利息
func CalculationOfInterestDiamondMortgageLoanAmount(totalLendingPercentage float64) (float64, float64) {

	ttp := totalLendingPercentage

	if ttp < 1 {
		ttp = 1 // 最低为1
	}

	// 可借出数量
	var loanhac float64 = 0
	loanhac = ((float64(100)/ttp)-1)*10 + 1

	// 接触数量调整
	/*
		if totalLendingPercentage < baserate {
			loanhac += 100 *  ((baserate - totalLendingPercentage) / baserate)
		}
	*/

	loanculty := float64(10) // 借出难度
	loanper := float64(0.78002)
	if ttp <= loanculty+1 {
		loanhac -= (loanhac * loanper) * ((loanculty - (ttp - 1)) / loanculty)
	}

	// 预先付息数量
	predeshac := float64(1)
	rateculty := float64(16) // 前期利息难度
	if ttp <= rateculty+1 {
		predeshac += loanhac * (loanper / 10 * ((rateculty - (ttp - 1)) / rateculty))
	}

	// 返回可借数量和预付利息
	return loanhac, predeshac
}
