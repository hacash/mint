package difficulty

import "github.com/hacash/core/interfaces"

func CheckHashDifficultySatisfyByBlock(blkhash []byte, block interfaces.Block) bool {
	targetdiffhash := Uint32ToHash(block.GetHeight(), block.GetDifficulty())
	return CheckHashDifficultySatisfy(blkhash, targetdiffhash)
}
func CheckHashDifficultySatisfyByDiffnum(blkhash []byte, blkhei uint64, diffnum uint32) bool {
	targetdiffhash := Uint32ToHash(blkhei, diffnum)
	return CheckHashDifficultySatisfy(blkhash, targetdiffhash)
}

func CheckHashDifficultySatisfy(hx1, hx2 []byte) bool {
	if len(hx1) != 32 || len(hx2) != 32 {
		panic("CheckHashDifficultySatisfy hx1, hx2 size must be 32.")
	}
	for k := 0; k < 32; k++ {
		if hx1[k] < hx2[k] {
			return true
		} else if hx1[k] > hx2[k] {
			return false
		}
	}
	return true
}
