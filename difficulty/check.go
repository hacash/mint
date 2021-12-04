package difficulty

import (
	"github.com/hacash/core/interfaces"
)

func CheckHashDifficultySatisfyByBlock(blkhash []byte, block interfaces.Block) bool {
	targetdiffhash := Uint32ToHash(block.GetHeight(), block.GetDifficulty())
	return CheckHashDifficultySatisfy(blkhash, targetdiffhash)
}
func CheckHashDifficultySatisfyByDiffnum(blkhash []byte, blkhei uint64, diffnum uint32) bool {
	targetdiffhash := Uint32ToHash(blkhei, diffnum)
	return CheckHashDifficultySatisfy(blkhash, targetdiffhash)
}

func CheckHashDifficultySatisfy(result_hash, target_diffculty_hash []byte) bool {
	if len(result_hash) != 32 || len(target_diffculty_hash) != 32 {
		panic("CheckHashDifficultySatisfy hx1, hx2 size must be 32.")
	}
	for k := 0; k < 32; k++ {
		if result_hash[k] < target_diffculty_hash[k] {
			return true
		} else if result_hash[k] > target_diffculty_hash[k] {
			return false
		}
	}
	return true
}
