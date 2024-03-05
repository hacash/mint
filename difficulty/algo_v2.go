package difficulty

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"
)

var (
	LowestDifficultyCompact       = uint32(4294967294) // Preset difficulty value before first adjustment of difficulty
	usedVersionV2AboveBlockHeight = uint64(288 * 160)  // Start using the new algorithm after the 160th difficulty cycle
)

// Encapsulated version external interface V1 + V2

func CalculateNextTarget(
	lastestBits uint32,
	currentHeight uint64,
	prev288BlockTimestamp uint64,
	lastestTimestamp uint64,
	eachblocktime uint64,
	changeblocknum uint64,
	printInfo *string,
) ([]byte, *big.Int, uint32) {
	if lastestBits == 0 {
		lastestBits = LowestDifficultyCompact // deal genesis block
	}
	// Use new version
	if currentHeight >= usedVersionV2AboveBlockHeight {
		if currentHeight >= uint64(288*450) { // 129600
			// In the past, the code actually only queried the total time of the first 287 blocks. After a certain height is added 300 seconds, it becomes 288 blocks
			lastestTimestamp += eachblocktime
		}
		return DifficultyCalculateNextTarget_v2(lastestBits, currentHeight, prev288BlockTimestamp, lastestTimestamp, eachblocktime, changeblocknum, printInfo)
	}
	// Oldest version
	b1, u1 := CalculateNextTargetDifficulty_v1(lastestBits, currentHeight, prev288BlockTimestamp, lastestTimestamp, eachblocktime, changeblocknum, printInfo)
	//return BigToHash256_v1(b1), b1, u1
	//if bytes.Compare(BigToHash256_v1(b1), Uint32ToHash256_v1(u1)) != 0 {
	//	fmt.Println("CalculateNextTargetDifficulty_v1: ", currentHeight, hex.EncodeToString(BigToHash256_v1(b1)), hex.EncodeToString(Uint32ToHash256_v1(u1)))
	//}
	return Uint32ToHash256_v1(u1), b1, u1
}

func Uint32ToBig(currentHeight uint64, diff_num uint32) *big.Int {
	// Use new version
	if currentHeight >= usedVersionV2AboveBlockHeight {
		return DifficultyUint32ToBig(diff_num)
	}
	// Oldest version
	return Uint32ToBig_v1(diff_num)
}

func Uint32ToHash(currentHeight uint64, diff_num uint32) []byte {
	// Use new version
	if currentHeight >= usedVersionV2AboveBlockHeight {
		return DifficultyUint32ToHash(diff_num)
	}
	// Oldest version
	return Uint32ToHash256_v1(diff_num)
}

func HashToBig(currentHeight uint64, hash []byte) *big.Int {
	// Use new version
	if currentHeight >= usedVersionV2AboveBlockHeight {
		return DifficultyHashToBig(hash)
	}
	// Oldest version
	return HashToBig_v1(hash)
}

func HashToUint32(currentHeight uint64, hash []byte) uint32 {
	// Use new version
	if currentHeight >= usedVersionV2AboveBlockHeight {
		return DifficultyHashToUint32(hash)
	}
	// Oldest version
	return Hash256ToUint32_v1(hash)
}

func BigToHash(currentHeight uint64, bignum *big.Int) []byte {
	// Use new version
	if currentHeight >= usedVersionV2AboveBlockHeight {
		return DifficultyBigToHash(bignum)
	}
	// Oldest version
	return BigToHash256_v1(bignum)
}

//////////////////////////////////////////////////////////////////////////////////////////

// Calculate the block difficulty of the next stage
func DifficultyCalculateNextTarget_v2(
	currentBits uint32,
	currentHeight uint64,
	prevTimestamp uint64,
	lastTimestamp uint64,
	eachblocktime uint64,
	changeblocknum uint64,
	printInfo *string,
) ([]byte, *big.Int, uint32) {

	powTargetTimespan := time.Second * time.Duration(eachblocktime*changeblocknum)
	// If the height of the new block is not an integer multiple of 288, it does not need to be updated. It is still the bits of the last block
	if currentHeight%changeblocknum != 0 {
		currentBig := DifficultyUint32ToBig(currentBits)
		return DifficultyBigToHash(currentBig), currentBig, currentBits
	}
	prev2016blockTimestamp := time.Unix(int64(prevTimestamp), 0)
	lastBlockTimestamp := time.Unix(int64(lastTimestamp), 0)
	// Calculate the block out time of 288 blocks
	actualTimespan := lastBlockTimestamp.Sub(prev2016blockTimestamp)
	if actualTimespan < powTargetTimespan/4 {
		// If it is less than 1/4 day, it calculated as 1/4 day
		actualTimespan = powTargetTimespan / 4
	} else if actualTimespan > powTargetTimespan*4 {
		// If it exceeds 4 days, it shall be calculated as 4 days
		actualTimespan = powTargetTimespan * 4
	}

	lastTarget := DifficultyUint32ToBig(currentBits)
	// formula: target = lastTarget * actualTime / expectTime
	newTarget := lastTarget.Mul(lastTarget, big.NewInt(int64(actualTimespan.Seconds())))
	newTarget = newTarget.Div(newTarget, big.NewInt(int64(powTargetTimespan.Seconds())))

	nextBits := DifficultyBigToUint32(newTarget)
	nextHash := DifficultyBigToHash(newTarget)

	// print data
	if printInfo != nil {
		actual_t, target_t := uint64(actualTimespan.Seconds()), uint64(powTargetTimespan.Seconds())
		nhs := strings.TrimRight(hex.EncodeToString(nextHash), "0")
		printStr := fmt.Sprintf("==== new ==== difficulty calculate next target at height %d ==== %ds/%ds ==== %ds/%ds ==== %d -> %d ==== "+nhs+" ====",
			currentHeight,
			actual_t/changeblocknum,
			target_t/changeblocknum,
			actual_t,
			target_t,
			currentBits,
			nextBits)
		*printInfo = printStr
	}

	return nextHash, newTarget, nextBits
}

func DifficultyUint32ToBig(diff_num uint32) *big.Int {
	hashbyte := DifficultyUint32ToHash(diff_num)
	return DifficultyHashToBig(hashbyte)
}

func DifficultyHashToBig(hashbyte []byte) *big.Int {
	cur_big := new(big.Int).SetBytes(bytes.TrimLeft(hashbyte, string([]byte{0})))
	return cur_big
}

func DifficultyUint32ToHash(diff_num uint32) []byte {
	return DifficultyUint32ToHashEx(diff_num, 0)
}

func DifficultyUint32ToHashForAntimatter(diff_num uint32) []byte {
	return DifficultyUint32ToHashEx(diff_num, 1)
}

func DifficultyUint32ToHashEx(diff_num uint32, filltail uint8) []byte {
	diff_byte := make([]byte, 4)
	binary.BigEndian.PutUint32(diff_byte, diff_num)

	// reduction
	originally_bits_1 := bytes.Repeat([]byte{0}, 255-int(diff_byte[0]))
	//fmt.Println("originally_bits_1:", len(originally_bits_1), originally_bits_1)
	originally_bits_2 := BytesToBits([]byte{diff_byte[1], diff_byte[2], diff_byte[3]})
	//fmt.Println("originally_bits_2:", len(originally_bits_2), originally_bits_2)
	originally_yushu := 256 - len(originally_bits_1) - len(originally_bits_2)
	originally_bits_3 := []byte{}
	if originally_yushu > 0 {
		originally_bits_3 = bytes.Repeat([]byte{filltail}, originally_yushu)
	}
	originally_bits_bufs := bytes.NewBuffer(originally_bits_1)
	originally_bits_bufs.Write(originally_bits_2)
	originally_bits_bufs.Write(originally_bits_3)
	originally_bits := originally_bits_bufs.Bytes()
	//fmt.Println("originally_bits:", len(originally_bits), originally_bits)
	originally_byte := BitsToBytes(originally_bits)[0:32]
	//fmt.Println("originally_byte:", len(originally_byte), originally_byte)
	return originally_byte
}

func DifficultyBigToHash(diff_big *big.Int) []byte {
	bigbytes := diff_big.Bytes()
	if len(bigbytes) > 32 {
		bigbytes = bytes.Repeat([]byte{255}, 32) // Maximum value when exceeding
	}
	buf := bytes.NewBuffer(bytes.Repeat([]byte{0}, 32-len(bigbytes)))
	buf.Write(bigbytes)
	return buf.Bytes()
}

func DifficultyBigToUint32(diff_big *big.Int) uint32 {
	bighash := DifficultyBigToHash(diff_big)
	return DifficultyHashToUint32(bighash)
}

func DifficultyHashToUint32(hash_byte []byte) uint32 {

	//hash_byte, _ := hex.DecodeString(hash)
	//
	//fmt.Println("\n--------------", hash, "-------------")
	//fmt.Println("           byte:", len(hash_byte), hash_byte)

	// HASH256 to UINT32
	//fmt.Println(hash_byte)
	hash_bits := BytesToBits(hash_byte)
	//fmt.Println(len(hash_bits), hash_bits)
	headzero := 0
	for _, v := range hash_bits {
		if v != 0 {
			break
		} else {
			headzero++
		}
	}
	hash_bits = append(hash_bits, bytes.Repeat([]byte{1}, 3*8+12)...)
	//fmt.Println(len(hash_bits), hash_bits)
	//fmt.Println(headzero, headzero+3*8)
	hash_bytes := BitsToBytes(hash_bits[headzero : headzero+3*8])
	//fmt.Println(len(hash_bits_2), hash_bits_2)
	//
	diff_byte := make([]byte, 4)
	diff_byte[0] = 255 - uint8(headzero)
	diff_byte[1] = hash_bytes[0]
	diff_byte[2] = hash_bytes[1]
	diff_byte[3] = hash_bytes[2]

	diff_number := binary.BigEndian.Uint32(diff_byte)
	//fmt.Println("diff_number:", diff_number)

	return diff_number
}
