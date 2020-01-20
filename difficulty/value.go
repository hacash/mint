package difficulty

import (
	"fmt"
	"github.com/hacash/blockmint/miner/difficulty"
	"math/big"
)

const (
	vK int64 = 1000
	vM int64 = 1000 * vK
	vG int64 = 1000 * vM
	vT int64 = 1000 * vG
	vP int64 = 1000 * vT
	vE int64 = 1000 * vP
)

func ConvertPowPowerToShowFormat( value *big.Int ) string {

	base := []int64{vE,vP,vT,vG,vM,vK}
	exts := []string{"E","P","T","G","M","K"}

	for i:=0; i<len(base); i++ {
		bsn := big.NewInt(base[i])
		if value.Cmp( bsn ) == 1 {
			numi := new(big.Int).Mul( value, big.NewInt(100) )
			numi = new(big.Int).Div( numi, bsn )
			numf := float64(numi.Int64())
			return fmt.Sprintf("%.2f"+exts[i]+"H/s", numf/100)
		}
	}
	return value.String() + "H/s"

}


///////////////////////////////////////////

// 计算哈希价值
func CalculateHashWorth(hash []byte) *big.Int {
	mulnum := big.NewInt(2)
	worth := big.NewInt(2)
	prezorenum := 0
	wbits := difficulty.BytesToBits(hash)
	for i, v := range wbits {
		if v != 0 {
			prezorenum = i
			break
		}
	}
	//
	for i := 0; i < prezorenum; i++ {
		worth = worth.Mul(worth, mulnum)
	}
	return worth
}
