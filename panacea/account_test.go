package panacea_test

import (
	"github.com/medibloc/panacea-oracle/crypto"
	"github.com/medibloc/panacea-oracle/panacea"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// This test creates MNEMONIC and checks whether address and publicKey are created normally
func TestAccount(t *testing.T) {
	mnemonic, err := crypto.GenerateMnemonic()

	acc, err := panacea.NewOracleAccount(mnemonic)
	require.NoError(t, err)
	require.Equal(t, 46, len(acc.GetAddress()))
	require.True(t, strings.HasPrefix(acc.GetAddress(), "panacea1"))
}
