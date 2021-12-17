package mint

var (
	SingleBlockMaxSize    uint32 = 1024 * 500 // ≈ 500KB
	SingleBlockMaxTxCount uint32 = 999        // ≈ coinbase + 999

	AdjustTargetDifficultyNumberOfBlocks uint64 = 288 // about one day
	EachBlockRequiredTargetTime          uint64 = 300 // 60 * 5 , about 5 min

	MinTransactionFeePurityOfOneByte uint64 = 10000 * 10000 * 10000 / 200 // = 5000000000:232 = 50铢 max 1844 67440737 09551615
	// one simple trs bytes length = 2+21+10+5+21+10+2+33+64+2 = 170 ~ 200

)
