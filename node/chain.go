package node

import (
	"encoding/hex"
	"fmt"

	"github.com/aymene01/ledgerNet/pb"
	"github.com/aymene01/ledgerNet/types"
)

type HeaderList struct {
	headers []*pb.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{
		headers: []*pb.Header{},
	}
}

func (list *HeaderList) Get(index int) *pb.Header {
	if index > list.Height() {
		panic("index too hight")
	}

	return list.headers[index]
}

func (list *HeaderList) Add(h *pb.Header){
	list.headers = append(list.headers, h)
}

func (list *HeaderList) Height() int {
	return len(list.headers) - 1
}

type Chain struct {
	blockStore Blockstorer
	headers    *HeaderList
}

func NewChain(bs Blockstorer) *Chain {
	return &Chain{
		blockStore: bs,
		headers: NewHeaderList(),
	}
}

func (c *Chain) Height() int {
	return c.headers.Height()
}

func (c *Chain) AddBlock(b *pb.Block) error {
	// add header to the list of headers
	c.headers.Add(b.Header)
	// validation
	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHeight(height int) (*pb.Block, error) {
	if c.Height() < height {
		return nil, fmt.Errorf("given heigh (%d) too hight - height (%d)", height, c.Height())
	}

	header := c.headers.Get(height)
	hash := types.HashHeader(header)
	return c.GetBlockByHash(hash)
}

func (c *Chain) GetBlockByHash(hash []byte) (*pb.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}
