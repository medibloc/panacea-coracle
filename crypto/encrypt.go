package crypto

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec"
)

// EncryptData encrypts data using recipient public key (ECIES)
// Secp256k1 is the only supported elliptic curve for encryption.
func EncryptData(pubKeyByte []byte, data []byte) ([]byte, error) {
	// parse public key
	pubKey, err := btcec.ParsePubKey(pubKeyByte, btcec.S256())
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key bytes: %w", err)
	}

	// encrypt data
	encryptedData, err := btcec.Encrypt(pubKey, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	return encryptedData, nil
}
