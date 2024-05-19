package crypto

import (
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()

	// Check if the private key is not nil
	if privKey == nil {
		t.Errorf("GeneratePrivateKey() returned nil")
	}

	// Check if the private key length is as expected
	privKeyBytes := privKey.Bytes()
	if len(privKeyBytes) != privKeyLen {
		t.Errorf("Expected private key length to be %d, but got %d", privKeyLen, len(privKeyBytes))
	}

	// Check if the public key derived from the private key is valid
	pubKey := privKey.Public()
	if pubKey == nil {
		t.Errorf("PrivateKey.Public() returned nil")
	}

	pubKeyBytes := pubKey.Bytes()
	if len(pubKeyBytes) != pubKeyLen {
		t.Errorf("Expected public key length to be %d, but got %d", pubKeyLen, len(pubKeyBytes))
	}
}

func TestBytes(t *testing.T) {
	privKey := GeneratePrivateKey()
	privKeyBytes := privKey.Bytes()

	// Check if the returned bytes are not nil
	if privKeyBytes == nil {
		t.Errorf("PrivateKey.Bytes() returned nil")
	}

	// Check if the length of the returned bytes is as expected
	if len(privKeyBytes) != privKeyLen {
		t.Errorf("Expected private key bytes length to be %d, but got %d", privKeyLen, len(privKeyBytes))
	}

	pubKey := privKey.Public()
	pubKeyBytes := pubKey.Bytes()

	// Check if the returned bytes are not nil
	if pubKeyBytes == nil {
		t.Errorf("PublicKey.Bytes() returned nil")
	}

	// Check if the length of the returned bytes is as expected
	if len(pubKeyBytes) != pubKeyLen {
		t.Errorf("Expected public key bytes length to be %d, but got %d", pubKeyLen, len(pubKeyBytes))
	}
}

func TestSignAndVerify(t *testing.T) {
	// Generate a private key
	privKey := GeneratePrivateKey()

	// Derive the public key from the private key
	pubKey := privKey.Public()

	// Create a message to sign
	msg := []byte("Hello, world!")

	// Sign the message with the private key
	sig := privKey.Sign(msg)

	// Verify the signature with the public key
	if !sig.Verify(pubKey, msg) {
		t.Error("Signature verification failed")
	}

	// Modify the message and verify the signature again (it should fail)
	msg[0] = 'h'
	if sig.Verify(pubKey, msg) {
		t.Error("Signature verification should have failed, but it succeeded")
	}

	// Modify the signature and verify it again (it should fail)
	sigBytes := sig.Bytes()
	sigBytes[0] = 'x'
	sig.value = sigBytes
	if sig.Verify(pubKey, msg) {
		t.Error("Signature verification should have failed, but it succeeded")
	}
}

func TestPublic(t *testing.T) {
	privKey := GeneratePrivateKey()

	// Check if the public key derived from the private key is not nil
	pubKey := privKey.Public()
	if pubKey == nil {
		t.Errorf("PrivateKey.Public() returned nil")
	}

	// Check if the length of the public key is as expected
	pubKeyBytes := pubKey.Bytes()
	if len(pubKeyBytes) != pubKeyLen {
		t.Errorf("Expected public key length to be %d, but got %d", pubKeyLen, len(pubKeyBytes))
	}

	// Check if the public key is the same as the last 32 bytes of the private key
	privKeyBytes := privKey.Bytes()
	if !compareBytes(pubKeyBytes, privKeyBytes[32:]) {
		t.Error("Public key does not match the last 32 bytes of the private key")
	}
}

func compareBytes(b1, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}

	for i := range b1 {
		if b1[i] != b2[i] {
			return false
		}
	}

	return true
}
