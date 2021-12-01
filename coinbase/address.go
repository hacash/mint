package coinbase

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

//
func UpdateCoinbaseAddress(tx interfaces.Transaction, address fields.Address) {
	tx.SetAddress(address)
}

//
func UpdateBlockCoinbaseAddress(block interfaces.Block, address fields.Address) {
	UpdateCoinbaseAddress(block.GetTrsList()[0], address)
}
