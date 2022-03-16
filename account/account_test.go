package account_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	panaceaapp "github.com/medibloc/panacea-core/v2/app"
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// This test creates MNEMONIC and checks whether address and publicKey are created normally
func TestAccount(t *testing.T) {
	panaceaapp.SetConfig()

	mnemonic, err := crypto.GenerateMnemonic()
	require.NoError(t, err)

	acc, err := account.NewValidatorAccount(mnemonic)
	require.NoError(t, err)
	require.Equal(t, 46, len(acc.GetAddress()))
	require.True(t, strings.HasPrefix(acc.GetAddress(), "panacea1"))

	priv := secp256k1.PrivKey{Key: acc.GetPrivKey().Bytes()}

	pub, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, 78, len(pub))
	require.True(t, strings.HasPrefix(pub, "panaceapub1"))
}
