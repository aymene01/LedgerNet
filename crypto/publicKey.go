package crypto

import "crypto/ed25519"

type PublicKey struct {
	key ed25519.PublicKey
}

func (p *PublicKey) Address() Address {
	return Address{
		value: p.key[len(p.key)-addressLen:],
	}
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}
