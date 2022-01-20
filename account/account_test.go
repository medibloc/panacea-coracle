package account_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types"
	panaceaapp "github.com/medibloc/panacea-core/v2/app"
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestAccount(t *testing.T) {
	panaceaapp.SetConfig()

	mnemonic := os.Getenv("OWNER_MNEMONIC")

	acc, err := account.NewValidatorAccount(mnemonic)
	require.NoError(t, err)
	require.Equal(t, "panacea1gtx6lmnjg6ykvv07ruyxamth6yuhgcvmhg3pqz", acc.GetAddress())

	priv2 := secp256k1.PrivKey{Key: acc.GetPrivKey().Bytes()}

	pub, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, priv2.PubKey())
	require.NoError(t, err)
	require.Equal(t, "panaceapub1addwnpepqwa9p79deddu9u3khl728ntfnyj6j37aguaxfht4u68hchvt048kw407c0q", pub)
}
