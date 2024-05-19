package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
)

func NewPrivateKeyFromString(s string) *PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(b)
}

func NewPrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != seedLen {
		panic("invalid seed length, must be 32")
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}
