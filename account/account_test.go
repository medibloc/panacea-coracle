package account_test

import (
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/panacea"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// This test creates MNEMONIC and checks whether address and publicKey are created normally
func TestAccount(t *testing.T) {
	mnemonic, err := crypto.GenerateMnemonic()

	acc, err := account.NewValidatorAccount(mnemonic, panacea.AccountAddressPrefix, panacea.CoinType)
	require.NoError(t, err)
	require.Equal(t, 46, len(acc.GetAddress()))
	require.True(t, strings.HasPrefix(acc.GetAddress(), "panacea1"))
}
