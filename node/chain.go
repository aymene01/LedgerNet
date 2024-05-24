package node

import (
	"encoding/hex"

	"github.com/aymene01/ledgerNet/pb"
)

type Chain struct {
	blockStore Blockstorer
}

func NewChain(bs Blockstorer) *Chain {
	return &Chain{
		blockStore: bs,
	}
}

func (c *Chain) AddBlock(b *pb.Block) error {
	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHeight(height int) (*pb.Block, error) {
	return nil, nil
}

func (c *Chain) GetBlockByHash(hash []byte) (*pb.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}
