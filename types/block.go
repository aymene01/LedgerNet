package types

import (
	"crypto/sha256"

	"github.com/aymene01/ledgerNet/crypto"
	"github.com/aymene01/ledgerNet/pb"
	"google.golang.org/protobuf/proto"
)

func SignBlock(pk *crypto.PrivateKey, b *pb.Block) *crypto.Signature {
	return pk.Sign(HashBlock(b))
}

func HashBlock(block *pb.Block) []byte {
	b, err := proto.Marshal(block)

	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}
