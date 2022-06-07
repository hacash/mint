package coinbase

import corecb "github.com/hacash/core/coinbase"

// Bitcoin mortgage loan: calculate the quantity that can be borrowed
// Totallendingpercentage total percentage lent
// Return lendable quantity and prepaid interest
func CalculationOfInterestBitcoinMortgageLoanAmount(totalLendingPercentage float64) (float64, float64) {
	return corecb.CalculationOfInterestBitcoinMortgageLoanAmount(totalLendingPercentage)
}
