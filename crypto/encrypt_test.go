package crypto

import (
	"bytes"
	"crypto/rand"
	"fmt"
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

	cipherText, err := EncryptDataWithSecp256k1(pubKey, origData)
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

	cipherText, err := EncryptDataWithSecp256k1(pubKey, origData)
	require.NoError(t, err)

	// try to decrypt using privKey2
	_, err = btcec.Decrypt(privKey2, cipherText)
	require.Error(t, err)
}

func TestEncryptDataWithAES256(t *testing.T) {
	fmt.Println(1)
	secretKey, err := randomBytes(32)
	require.NoError(t, err)
	additional := Hash([]byte(fmt.Sprintf("data-pool-%v", 1)))

	data, err := randomBytes(100000)
	require.NoError(t, err)

	cipherText, err := EncryptDataWithAES256(secretKey, additional, data)
	require.NoError(t, err)
	require.NotEqual(t, data, cipherText)

	decryptText, err := DecryptDataWithAES256(secretKey, additional, cipherText)
	require.NoError(t, err)

	require.Equal(t, data, decryptText)
}

func randomBytes(size int) ([]byte, error) {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}
