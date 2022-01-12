package crypto

import (
	"bytes"
	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/require"
	"testing"
)

// Success encryption and decryption
func TestEncryptData(t *testing.T) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		t.Fatal("failed to generate key: ", err)
	}

	pubKey := privKey.PubKey().SerializeCompressed()
	origData := []byte("encrypt origData please")

	cipherText, err := EncryptData(pubKey, origData)
	if err != nil {
		t.Fatal("failed to encrypt origData: ", err)
	}

	plainText, err := btcec.Decrypt(privKey, cipherText)

	if !bytes.Equal(origData, plainText) {
		t.Errorf("decrypted data doesn't match original data")
	}
}

// Success encryption but fail decryption
func TestEncryptData_FailDecryption(t *testing.T) {
	privKey1, err1 := btcec.NewPrivateKey(btcec.S256())
	privKey2, err2 := btcec.NewPrivateKey(btcec.S256())
	if err1 != nil || err2 != nil {
		t.Fatal("failed to generate key: ", err1, err2)
	}

	// encrypt to pubKey1
	pubKey := privKey1.PubKey().SerializeCompressed()
	origData := []byte("decryption will be failed")

	cipherText, err := EncryptData(pubKey, origData)
	if err != nil {
		t.Fatal("failed to encrypt origData: ", err)
	}

	// try to decrypt using privKey2
	_, err = btcec.Decrypt(privKey2, cipherText)
	require.Error(t, err)
}
