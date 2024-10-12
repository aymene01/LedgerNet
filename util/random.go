package util

import (
	"crypto/rand"
	"io"
	randN "math/rand"
	"time"

	"github.com/aymene01/ledgerNet/pb"
)

func RandomHash() []byte {
	hash := make([]byte, 32)
	io.ReadFull(rand.Reader, hash)
	return hash
}

func RandomBlock() *pb.Block {
	header := &pb.Header{
		Version: 1,
		Height: int32(randN.Intn(1000)),
		PrevHash: RandomHash(),
		RootHash: RandomHash(),
		Timestamp: time.Now().UnixNano(),
	}
	return &pb.Block{
		Header: header,
	}
}