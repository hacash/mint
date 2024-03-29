package blockchain

import (
	"fmt"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/genesis"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
	"github.com/hacash/x16rs"
	"strings"
)

func (b *BlockChain) ValidateDiamondCreateAction(action interfaces.Action) error {
	act, ok := action.(*actions.Action_4_DiamondCreate)
	if !ok {
		return fmt.Errorf("its not Action_4_DiamondCreate Action.")
	}

	// Developer mode, no check
	if sys.TestDebugLocalDevelopmentMark {
		return nil // Developer mode does not check and returns success
	}

	last, err := b.StateRead().ReadLastestDiamond()
	if err != nil {
		return err
	}
	if last == nil { // is first
		genesisblk := genesis.GetGenesisBlock()
		last = &stores.DiamondSmelt{
			Number:           0,
			ContainBlockHash: genesisblk.Hash(),
		}
	}
	if uint32(act.Number) != uint32(last.Number)+1 {
		return fmt.Errorf("Diamond number error.")
	}
	if last.ContainBlockHash.Equal(act.PrevHash) != true {
		return fmt.Errorf("Diamond prev block hash error.")
	}
	hashave, e := b.StateRead().Diamond(act.Diamond)
	if e != nil {
		return e
	}
	if hashave != nil {
		return fmt.Errorf("Diamond <%s> already exist.", act.Diamond)
	}
	// Check diamond mining calculation
	sha3hash, diamond_resbytes, diamond_str := x16rs.Diamond(uint32(act.Number), act.PrevHash, act.Nonce, act.Address, act.GetRealCustomMessage())
	diamondstrval, isdia := x16rs.IsDiamondHashResultString(diamond_str)
	if !isdia {
		return fmt.Errorf("String <%s> is not diamond.", diamond_str)
	}
	if strings.Compare(diamondstrval, string(act.Diamond)) != 0 {
		return fmt.Errorf("Diamond need <%s> but got <%s>", act.Diamond, diamondstrval)
	}
	// Check diamond difficulty value
	difok := x16rs.CheckDiamondDifficulty(uint32(act.Number), sha3hash, diamond_resbytes)
	if !difok {
		return fmt.Errorf("Diamond difficulty not meet the requirements.")
	}
	// check ok
	return nil
}
