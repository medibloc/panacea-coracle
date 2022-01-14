package crypto

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/require"
	"testing"
)

// Success encryption and decryption
func TestEncryptData(t *testing.T) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	require.NoError(t, err)

	pubKey := privKey.PubKey().SerializeCompressed()
	origData := []byte("encrypt origData please")

	cipherText, err := EncryptData(pubKey, origData)
	require.NoError(t, err)

	plainText, err := btcec.Decrypt(privKey, cipherText)

	if !bytes.Equal(origData, plainText) {
		t.Errorf("decrypted data doesn't match original data")
	}
}

// Success encryption but fail decryption
func TestEncryptData_FailDecryption(t *testing.T) {
	privKey1, err1 := btcec.NewPrivateKey(btcec.S256())
	privKey2, err2 := btcec.NewPrivateKey(btcec.S256())
	require.NoError(t, err1)
	require.NoError(t, err2)

	// encrypt to pubKey1
	pubKey := privKey1.PubKey().SerializeCompressed()
	origData := []byte("decryption will be failed")

	cipherText, err := EncryptData(pubKey, origData)
	require.NoError(t, err)

	// try to decrypt using privKey2
	_, err = btcec.Decrypt(privKey2, cipherText)
	require.Error(t, err)
}

// TestHash sha256 and hex conversion test.
func TestHash(t *testing.T) {
	origData := []byte("encrypt origData please")
	hashData := Hash(origData)
	hashDataStr := hex.EncodeToString(hashData)
	require.NotNil(t, hashData)
	require.Equal(t, 32, len(hashData))
	require.Equal(t, 64, len(hashDataStr))

	expectedHashDataStr := "d0b6d2dd9790d430785d8a8cf1cb86dfe9a247664a08c49619267960a9f3a46d"
	require.Equal(t, expectedHashDataStr, hashDataStr)
}