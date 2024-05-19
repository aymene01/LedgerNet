package types

import (
	"crypto/sha256"

	"github.com/aymene01/blocker/crypto"
	"github.com/aymene01/blocker/pb"
	"google.golang.org/protobuf/proto"
)

func HashTransaction(tx *pb.Transaction) []byte {
	b, err := proto.Marshal(tx)

	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func verifyTransaction(tx *pb.Transaction) bool {
	for _, input := range tx.Inputs {
		var (
			sig    = crypto.SignatureFromBytes(input.Signature)
			pubKey = crypto.PublicKeyFromBytes(input.PublicKey)
		)
		// TODO: make sure we don't run into a problem after verification cause we have set
		// signature to nil
		input.Signature = nil

		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}
	}

	return true
}
