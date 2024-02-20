package mint

var (
	SingleBlockMaxSize    uint32 = 1024 * 500 // ≈ 500KB
	SingleBlockMaxTxCount uint32 = 999        // ≈ coinbase + 999

	AdjustTargetDifficultyNumberOfBlocks uint64 = 288 // about one day
	EachBlockRequiredTargetTime          uint64 = 300 // 60 * 5 , about 5 min

	MinTransactionFeePurity uint32 = 1666 // = 10000 / (168 / 32 + 1)
	// one simple trs bytes length = 2+5+21+6+2+(21+10)+2+33+64+2 = 168

)
