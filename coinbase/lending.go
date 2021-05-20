package coinbase

// 钻石抵押借贷： 计算可借数量
// totalLendingPercentage 已经借出的总百分比
// 返回可借数量和预付利息
func CalculationOfInterestDiamondMortgageLoanAmount(totalLendingPercentage float64) (float64, float64) {

	ttp := totalLendingPercentage

	// 可借出数量
	loanhac := ((float64(100)/ttp)-1)*10 + 1
	loanculty := float64(20) // 借出难度
	if ttp <= loanculty+1 {
		loanhac -= (loanhac * 0.75) * ((loanculty - (ttp - 1)) / loanculty)
	}
	// 预先付息数量
	predeshac := float64(1)
	rateculty := float64(14) // 前期利息难度
	if ttp <= rateculty+1 {
		predeshac += loanhac * (0.14 * ((rateculty - (ttp - 1)) / rateculty))
	}

	// 返回可借数量和预付利息
	return loanhac, predeshac
}
