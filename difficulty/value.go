package difficulty

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/x16rs"
	"math/big"
)

const (
	vK int64 = 1024
	vM int64 = vK * vK
	vG int64 = vK * vM
	vT int64 = vK * vG
	vP int64 = vK * vT
	vE int64 = vK * vP
)

func pow(value *big.Int, x int) *big.Int {
	if x == 1 {
		return value
	}
	num := new(big.Int).Set(value)
	for i := 1; i < int(x); i++ {
		num = num.Mul(num, value)
	}
	return num
}

func ConvertPowPowerToShowFormat(value *big.Int) string {
	stepn := new(big.Int).SetUint64(1000)
	if value.Cmp(stepn) <= 0 {
		return value.String() + "H/s"
	}
	exts := "KMGTPEZYBNDCX"
	spxi := 0
	spxn := new(big.Int).SetUint64(1)
	for i := 2; i < len(exts); i++ {
		if value.Cmp(pow(stepn, i)) <= 0 {
			spxi = i - 2
			spxn = pow(stepn, i-1)
			break
		}
	}
	num1000 := new(big.Int).SetUint64(1000)
	numi := new(big.Int).Mul(num1000, value)
	numi = new(big.Int).Div(numi, spxn)
	numf := float64(numi.Int64()) / 1000
	return fmt.Sprintf("%.3f"+string(exts[spxi])+"H/s", numf)
}

func ConvertPowPowerToShowFormat_old(value *big.Int) string {

	base := []int64{vE, vP, vT, vG, vM, vK}
	exts := []string{"E", "P", "T", "G", "M", "K"}

	for i := 0; i < len(base); i++ {
		bsn := big.NewInt(base[i])
		if value.Cmp(bsn) == 1 {
			numi := new(big.Int).Mul(value, big.NewInt(100))
			numi = new(big.Int).Div(numi, bsn)
			numf := float64(numi.Int64())
			return fmt.Sprintf("%.2f"+exts[i]+"H/s", numf/100)
		}
	}
	return value.String() + "H/s"

}

///////////////////////////////////////////

// 计算哈希价值
func CalculateHashWorth(curheight uint64, hash []byte) *big.Int {
	worth := DifficultyHashToBig(antimatterHash(hash))
	repeat := x16rs.HashRepeatForBlockHeight(curheight)
	targetHashWorth := new(big.Int).Mul(worth, new(big.Int).SetUint64(uint64(repeat)))
	return targetHashWorth
}

// 计算难度价值
func CalculateDifficultyWorth(curheight uint64, diffnum uint32) *big.Int {
	diffhx := DifficultyUint32ToHashForAntimatter(diffnum)
	worth := DifficultyHashToBig(antimatterHash(diffhx))
	repeat := x16rs.HashRepeatForBlockHeight(curheight)
	targetHashWorth := new(big.Int).Mul(worth, new(big.Int).SetUint64(uint64(repeat)))
	return targetHashWorth
}

// 计算哈希价值
func CalculateHashWorth_old(hash []byte) *big.Int {
	mulnum := big.NewInt(2)
	worth := big.NewInt(2)
	prezorenum := 0
	wbits := BytesToBits(hash)
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

// 反物质哈希
func antimatterHash(hx []byte) []byte {

	prefixzorenum := 0
	basevalbts := []byte{0, 0, 0, 0}
	cpidx := 0
	for i, v := range hx {
		if v == 0 {
			prefixzorenum++
			continue
		}
		cpidx = i
		break
	}
	copy(basevalbts[1:], hx[cpidx:])
	newbasenum := 16777215 - binary.BigEndian.Uint32(basevalbts) // 反记
	// 新值
	newbasenumbts := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(newbasenumbts, newbasenum)
	// hash
	buf := bytes.NewBuffer(newbasenumbts[1:])
	buf.Write(bytes.Repeat([]byte{255}, prefixzorenum))
	return buf.Bytes()
}

func antimatterHash_old(hx []byte) []byte {

	tar := make([]byte, len(hx))
	a := 0
	for i := len(hx) - 1; i >= 0; i-- {
		tar[a] = 255 - hx[i]
		a += 1
	}

	return tar

}
