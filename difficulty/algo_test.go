package difficulty

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"
)

//

// Calculate hash value
func CalculateHashWorthV2_old(hash []byte) *big.Int {
	mulstep := big.NewInt(256)
	worth := big.NewInt(1)
	for _, v := range hash {
		if v == 0 {
			worth = worth.Mul(worth, mulstep)
			continue
		}
		mulnum := big.NewInt(256 - int64(v))
		worth = worth.Mul(worth, mulnum)
		break
	}
	return worth
}

func Test_t148375092345(t *testing.T) {
	restr := ""
	upnum := 999
	for i := float64(0); i < 256; i++ {
		nx := int(math.Exp2(8 - math.Log2(i)))
		if upnum != nx {
			upnum = nx
			restr += fmt.Sprintf("{%d, %d},\n", nx, int(i)-1)
		}
	}
	fmt.Println(restr)
}

func Test_t1(t *testing.T) {

	basehx := bytes.Repeat([]byte{255}, 32)

	for i := 255; i >= 0; i-- {

		copy(basehx, []byte{uint8(i)})
		//fmt.Println( basehx, BytesToBits(basehx) )
		fmt.Println(basehx[0:4], BytesToBits(basehx)[0:16], CalculateHashWorthForTest(basehx))

	}
}

func Test_t2(t *testing.T) {

	hxstrs := []string{
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"00feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"00000000fffffeffffffffffffffffffffffffffffffffffffffffffffffffff",
		"00000000fdffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"00000000fd60ee37af3db8dcad3fa8cae3beb9a54d78906fe3c859a8efc16a93",
		"000000000030ee37af3db8dcad3fa8cae3beb9a54d78906fe3c859a8efc16a93",
		"000000000000e037af3db8dcad3fa8cae3beb9a54d78906fe3c859a8efc16a93",
		"000000000000b037af3db8dcad3fa8cae3beb9a54d78906fe3c859a8efc16a93",
		"0000000000006037af3db8dcad3fa8cae3beb9a54d78906fe3c859a8efc16a93",
	}

	for _, v := range hxstrs {
		hx, _ := hex.DecodeString(v)
		worth := DifficultyHashToBig(antimatterHash(hx))
		fmt.Println(v, antimatterHash(hx), worth, ConvertPowPowerToShowFormat(worth))
		//fmt.Print("\n-------------------\n\n")
		//fmt.Println( hx, BytesToBits(hx), CalculateHashWorthForTest(hx), CalculateHashWorthV2(hx) )
		//fmt.Println(v, antimatterHash_old2(hx), DifficultyHashToBig(antimatterHash_old2(hx)).String())
		//fmt.Println(v, antimatterHash(hx),DifficultyHashToBig(antimatterHash(hx)).String())
	}

}

func Test_t3(t *testing.T) {

	hxs := [][]byte{
		{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	for _, hx := range hxs {
		fmt.Println(hex.EncodeToString(hx), DifficultyHashToBig(antimatterHash(hx)).String())
	}

}

func Test_t4(t *testing.T) {

	hxs := [][]byte{
		HexDecodeString("0000000005d28eb241b004a2da93e2248d9b820d2e794408651aa8d55c819a8d"),
		HexDecodeString("000000001748ae6594aeb201652232f9a7df3c148106f548506f259579966235"),
		HexDecodeString("000000000afd7f82d86251c57b00ddfc106f84f63c716662ef28a13dc51ce983"),
		HexDecodeString("0000000002b5c4fe535c9f2a7ef1edca24dca3cb26fa04874d48a90b1b42efda"),
		HexDecodeString("00000000174aa267220b821af1d2faaee899eab2f70a40ea69c4819c8e6efb57"),
		HexDecodeString("0000000015e7e5627aa9b7064e9a4b4991dcc5c0d4a21dc7db8bcb1b5b97d1bd"),
		HexDecodeString("000000000c0e3a40c518a625878ed55b15f1bd6ed60d903eeaa26ae9ee032255"),
		HexDecodeString("000000001839e4abc416eba7843c48dfc372b157c9376463561bb972e94945d0"),
	}

	for _, hx := range hxs {
		fmt.Println(hex.EncodeToString(hx), CalculateHashWorthForTest(hx))
	}

}

func HexDecodeString(hexstr string) []byte {
	bts, e := hex.DecodeString(hexstr)
	if e != nil {
		panic(e)
	}
	return bts
}

func FAN(hx []byte) []byte {

	tar := make([]byte, len(hx))
	a := 0
	for i := len(hx) - 1; i >= 0; i-- {
		tar[a] = 255 - hx[i]
		a += 1
	}

	return tar

}
