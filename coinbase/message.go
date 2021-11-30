package coinbase

import (
	"encoding/binary"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfacev2"
)

//
func ParseMinerPoolCoinbaseMessage(msgwords string, minernum uint32) [16]byte {
	var msg [16]byte
	copy(msg[:], msgwords) // minerpoolcom
	binary.BigEndian.PutUint32(msg[12:16], minernum)
	return msg
}

//
func UpdateCoinbaseMessageForMiner(tx interfacev2.Transaction, minernum uint32) {
	newmsg := ParseMinerPoolCoinbaseMessage(string(tx.GetMessage()), minernum)
	tx.SetMessage(fields.TrimString16(string(newmsg[:])))
}

//
func UpdateBlockCoinbaseMessageForMiner(block interfacev2.Block, minernum uint32) {
	UpdateCoinbaseMessageForMiner(block.GetTransactions()[0], minernum)
}

//
func UpdateCoinbaseMessage(tx interfacev2.Transaction, msgstr string) {
	tx.SetMessage(fields.TrimString16(string(msgstr[:])))
}

//
func UpdateBlockCoinbaseMessage(block interfacev2.Block, msgstr string) {
	UpdateCoinbaseMessage(block.GetTransactions()[0], msgstr)
}
