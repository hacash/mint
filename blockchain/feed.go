package blockchain

import (
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

func (bc *BlockChain) SubscribeValidatedBlockOnInsert(blockCh chan interfaces.Block) {
	bc.validatedBlockInsertFeed.Subscribe(blockCh)
}

func (bc *BlockChain) SubscribeDiamondOnCreate(diamondCh chan *stores.DiamondSmelt) {
	bc.diamondCreateFeed.Subscribe(diamondCh)
}
