package crypto_test

import (
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/panacea"
	"github.com/stretchr/testify/require"
	"testing"
)

// Create mnemonic and check if PrivKey is extracted normally.
// And data verification test is also performed with the generated PrivKey
func TestGeneratePrivateKeyFromMnemonic(t *testing.T) {
	mnemonic, err := crypto.GenerateMnemonic()
	require.NoError(t, err)
	require.NotEqual(t, "", mnemonic)

	privKey, err := crypto.GeneratePrivateKeyFromMnemonic(mnemonic, panacea.CoinType)
	require.NoError(t, err)

	originData := []byte("Test origin data.")

	sign, err := privKey.Sign(originData)
	require.NoError(t, err)

	require.True(t, privKey.PubKey().VerifySignature(originData, sign))
}
