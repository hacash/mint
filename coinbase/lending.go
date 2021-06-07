package coinbase

import corecb "github.com/hacash/core/coinbase"

// 比特币抵押借贷： 计算可借数量
// totalLendingPercentage 已经借出的总百分比
// 返回可借数量和预付利息
func CalculationOfInterestBitcoinMortgageLoanAmount(totalLendingPercentage float64) (float64, float64) {
	return corecb.CalculationOfInterestBitcoinMortgageLoanAmount(totalLendingPercentage)
}
