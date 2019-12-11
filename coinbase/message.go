package coinbase

import (
	"encoding/binary"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
)

//
func ParseMinerPoolCoinbaseMessage(msgwords string, minernum uint32) [16]byte {
	var msg [16]byte
	copy(msg[0:11], []byte(msgwords)[0:11]) // minerpoolcn
	binary.BigEndian.PutUint32(msg[12:16], minernum)
	return msg
}

//
func UpdateCoinbaseMessageForMinerPool(tx transactions.Transaction_0_Coinbase, minernum uint32) {
	newmsg := ParseMinerPoolCoinbaseMessage(string(tx.Message), minernum)
	tx.Message = fields.TrimString16(string(newmsg[:]))
}
