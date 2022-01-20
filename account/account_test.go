package account_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	panaceaapp "github.com/medibloc/panacea-core/v2/app"
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// In order for this test to be successful, a valid MNEMONIC must be added to the environmental variable(VALIDATOR_MNEMONIC)
func TestAccount(t *testing.T) {
	panaceaapp.SetConfig()

	mnemonic := os.Getenv(types.VALIDATOR_MNEMONIC)

	acc, err := account.NewValidatorAccount(mnemonic)
	require.NoError(t, err)
	require.Equal(t, "panacea1gtx6lmnjg6ykvv07ruyxamth6yuhgcvmhg3pqz", acc.GetAddress())

	priv2 := secp256k1.PrivKey{Key: acc.GetPrivKey().Bytes()}

	pub, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv2.PubKey())
	require.NoError(t, err)
	require.Equal(t, "panaceapub1addwnpepqwa9p79deddu9u3khl728ntfnyj6j37aguaxfht4u68hchvt048kw407c0q", pub)
}
