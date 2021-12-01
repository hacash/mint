package coinbase

import (
	"encoding/binary"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

//
func ParseMinerPoolCoinbaseMessage(msgwords string, minernum uint32) [16]byte {
	var msg [16]byte
	copy(msg[:], msgwords) // minerpoolcom
	binary.BigEndian.PutUint32(msg[12:16], minernum)
	return msg
}

//
func UpdateCoinbaseMessageForMiner(tx interfaces.Transaction, minernum uint32) {
	newmsg := ParseMinerPoolCoinbaseMessage(string(tx.GetMessage()), minernum)
	tx.SetMessage(fields.TrimString16(string(newmsg[:])))
}

//
func UpdateBlockCoinbaseMessageForMiner(block interfaces.Block, minernum uint32) {
	UpdateCoinbaseMessageForMiner(block.GetTrsList()[0], minernum)
}

//
func UpdateCoinbaseMessage(tx interfaces.Transaction, msgstr string) {
	tx.SetMessage(fields.TrimString16(string(msgstr[:])))
}

//
func UpdateBlockCoinbaseMessage(block interfaces.Block, msgstr string) {
	UpdateCoinbaseMessage(block.GetTrsList()[0], msgstr)
}
