package node

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/aymene01/ledgerNet/pb"
	"github.com/aymene01/ledgerNet/types"
)

type Blockstorer interface {
	Put(*pb.Block) error
	Get(string) (*pb.Block, error)
}

type MemoryBlockStore struct {
	lock   sync.RWMutex
	blocks map[string]*pb.Block
}

func NewMemoryBlockStrore() *MemoryBlockStore {
	return &MemoryBlockStore{
		blocks: make(map[string]*pb.Block),
	}
}

func (s *MemoryBlockStore) Get(hash string) (*pb.Block, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	block, ok := s.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash [%s] does not exist", hash)
	}

	return block, nil
}

func (s *MemoryBlockStore) Put(b *pb.Block) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	hash := hex.EncodeToString(types.HashBlock(b))
	s.blocks[hash] = b
	return nil
}
