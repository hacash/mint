package coinbase

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfacev2"
)

//
func UpdateCoinbaseAddress(tx interfacev2.Transaction, address fields.Address) {
	tx.SetAddress(address)
}

//
func UpdateBlockCoinbaseAddress(block interfacev2.Block, address fields.Address) {
	UpdateCoinbaseAddress(block.GetTransactions()[0], address)
}
