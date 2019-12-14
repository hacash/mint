package coinbase

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"strings"
)

func CreateCoinbaseTx(blockheight uint64) *transactions.Transaction_0_Coinbase {

	coinbasetx := transactions.NewTransaction_0_Coinbase()
	coinbasetx.Reward = *BlockCoinBaseReward(blockheight)
	rwdaddr, _ := fields.CheckReadableAddress("1AVRuFXNFi3rdMrPH4hdqSgFrEBnWisWaS")
	coinbasetx.Address = *rwdaddr
	coinbasetx.Message = fields.TrimString16(strings.Repeat(" ", 16))

	return coinbasetx
}
