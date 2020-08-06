package coinbase

import (
	"fmt"
	"testing"
)

func Test_t1(t *testing.T) {

	//PrintMoveBtcCoinRewardNumberTable()

	for i := 1; i < 20; i++ {
		n := i + 16380
		fmt.Println(n, MoveBtcCoinRewardNumber(int64(n)))
	}

}
