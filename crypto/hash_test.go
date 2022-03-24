package crypto_test

import (
	"encoding/base64"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestHash sha256 and base64 conversion test.
func TestHash(t *testing.T) {
	origData := []byte("encrypt origData please")
	hashData := crypto.Hash(origData)
	hashDataStr := base64.StdEncoding.EncodeToString(hashData)

	require.NotNil(t, hashData)
	require.Equal(t, 32, len(hashData))
	require.Equal(t, 44, len(hashDataStr))

	expectedHashDataStr := "0LbS3ZeQ1DB4XYqM8cuG3+miR2ZKCMSWGSZ5YKnzpG0="
	require.Equal(t, expectedHashDataStr, hashDataStr)
}
