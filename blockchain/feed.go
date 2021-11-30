package blockchain

import (
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
)

func (bc *BlockChain) SubscribeValidatedBlockOnInsert(blockCh chan interfacev2.Block) {
	bc.validatedBlockInsertFeed.Subscribe(blockCh)
}

func (bc *BlockChain) SubscribeDiamondOnCreate(diamondCh chan *stores.DiamondSmelt) {
	bc.diamondCreateFeed.Subscribe(diamondCh)
}
