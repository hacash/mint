package difficulty

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

// 转换为 哈希 为算力显示
func ConvertDifficultyToRateShow(diffnum uint32, usetimesec int64) string {
	hxworth := CalculateDifficultyWorth(diffnum)
	hashrate := new(big.Int).Div(hxworth, big.NewInt(usetimesec))
	hashrateshow := ConvertPowPowerToShowFormat(hashrate)
	return hashrateshow
}

// 转换为 哈希 为算力显示
func ConvertHashToRateShow(hx []byte, usetimesec int64) string {
	hxworth := CalculateHashWorth(hx)
	hashrate := new(big.Int).Div(hxworth, big.NewInt(usetimesec))
	hashrateshow := ConvertPowPowerToShowFormat(hashrate)
	return hashrateshow
}

func ConvertPowPowerToShowFormat(value *big.Int) string {
	divn := new(big.Float).SetUint64(1000)
	exts := "KMGTPEZYBNDCX"
	basn := new(big.Float).SetInt(value)
	ext := ""
	if basn.Cmp(divn) == 1 {
		for i := 0; i < len(exts); i++ {
			resv := new(big.Float).Quo(basn, divn)
			ext = string(exts[i])
			basn = resv
			if resv.Cmp(divn) == -1 {
				break
			}
		}
	}
	return fmt.Sprintf("%.3f%sH/s", basn, ext)
}

func ConvertPowPowerToShowFormat_old2(value *big.Int) string {
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
func antiByte(bt uint8) uint8 {
	var antiByteSplitNums = [][]uint8{
		{255, 0},
		{128, 1},
		{85, 2},
		{64, 3},
		{51, 4},
		{42, 5},
		{36, 6},
		{32, 7},
		{28, 8},
		{25, 9},
		{23, 10},
		{21, 11},
		{19, 12},
		{18, 13},
		{17, 14},
		{16, 15},
		{15, 16},
		{14, 17},
		{13, 18},
		{12, 19},
		{11, 21},
		{10, 23},
		{9, 25},
		{8, 28},
		{7, 32},
		{6, 36},
		{5, 42},
		{4, 51},
		{3, 64},
		{2, 85},
		{1, 128},
		{0, 255},
	}
	for _, v := range antiByteSplitNums {
		if bt >= v[0] {
			return v[1]
		}
	}
	return 0
}

// 计算哈希价值

func CalculateHashWorth(hash []byte) *big.Int {
	bigbytes := []byte{0, 0, 0}
	zore := 0
	for i := 0; i < 29; i++ {
		if hash[i] > 0 {
			//fmt.Println(hash, i)
			bigbytes[0] = antiByte(hash[i])
			bigbytes[1] = antiByte(hash[i+1])
			bigbytes[2] = antiByte(hash[i+2])
			break
		}
		zore++
	}
	if zore > 2 {
		zore -= 2
	}
	bigbytes = append(bigbytes, bytes.Repeat([]byte{0}, zore)...)
	return new(big.Int).SetBytes(bigbytes)
}
func CalculateHashWorth_old_2022_02_08(hash []byte) *big.Int {
	worth := DifficultyHashToBig(antimatterHash(hash))
	return worth
	//repeat := x16rs.HashRepeatForBlockHeight(curheight)
	//targetHashWorth := new(big.Int).Mul(worth, new(big.Int).SetUint64(uint64(repeat)))
	//return targetHashWorth
}

// 计算难度价值

func CalculateDifficultyWorth(diffnum uint32) *big.Int {
	diffhx := DifficultyUint32ToHashForAntimatter(diffnum)
	return CalculateHashWorth(diffhx)
}
func CalculateDifficultyWorth_old_2022_02_08(diffnum uint32) *big.Int {
	diffhx := DifficultyUint32ToHashForAntimatter(diffnum)
	worth := DifficultyHashToBig(antimatterHash(diffhx))
	return worth
	//repeat := x16rs.HashRepeatForBlockHeight(curheight)
	//targetHashWorth := new(big.Int).Mul(worth, new(big.Int).SetUint64(uint64(repeat)))
	//return targetHashWorth
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
func antimatterHash_old(hx []byte) []byte {

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

func antimatterHash_old2(hx []byte) []byte {

	tar := make([]byte, 0)
	for i := 0; i < len(hx); i++ {
		if hx[i] == 0 {
			tar = append(tar, 255)
		} else if hx[i] == 255 {
			break // 结束
		} else {
			prx := []byte{255 - hx[i]}
			tar = append(prx, tar...)
			break
		}
	}
	return tar

}

func antimatterHash(hx []byte) []byte {

	size := len(hx)
	zorenum := 1
	basenum := make([]byte, 0)
	// 第一步：拆解
	for i := 0; i < size; i++ {
		if hx[i] == 0 {
			zorenum++
		} else {
			for a := 0; a < 3; a++ {
				if i+a < size {
					basenum = append(basenum, 255-hx[i+a])
				} else {
					basenum = append(basenum, 0)
				}
			}
			break
		}
	}
	sbxzore := bytes.Repeat([]byte{0}, zorenum)
	//fmt.Println(basenum, zorenum, sbxzore)
	// 第二步 合并
	for x := 0; x < len(sbxzore) && x < len(basenum); x++ {
		sbxzore[x] = basenum[x]
	}
	if zorenum > 0 {
		if sbxzore[0] == 0 {
			sbxzore[0] = 1
		}
	}
	// 返回
	return sbxzore

}
