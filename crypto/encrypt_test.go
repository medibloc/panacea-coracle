package crypto

import (
	"bytes"
	"encoding/base64"
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

// Test to generate signature and verify
func TestSignData(t *testing.T) {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	require.NoError(t, err)

	origData := []byte("decryption will be failed")

	signatureByte, err := SignData(privKey.Serialize(), origData)
	require.NoError(t, err)

	pubKey := privKey.PubKey()

	signature, err := btcec.ParseSignature(signatureByte, btcec.S256())
	require.NoError(t, err)

	verify := signature.Verify(origData, pubKey)
	require.True(t, verify)
}

// TestHash sha256 and base64 conversion test.
func TestHash(t *testing.T) {
	origData := []byte("encrypt origData please")
	hashData := Hash(origData)
	hashDataStr := base64.StdEncoding.EncodeToString(hashData)

	require.NotNil(t, hashData)
	require.Equal(t, 32, len(hashData))
	require.Equal(t, 44, len(hashDataStr))

	expectedHashDataStr := "0LbS3ZeQ1DB4XYqM8cuG3+miR2ZKCMSWGSZ5YKnzpG0="
	require.Equal(t, expectedHashDataStr, hashDataStr)
}