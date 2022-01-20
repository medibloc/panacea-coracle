package account

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	panaceacrypto "github.com/medibloc/panacea-core/v2/x/did/client/crypto"
	"github.com/tendermint/tendermint/crypto"
)

type ValidatorAccount struct {
	privKey crypto.PrivKey
	pubKey  crypto.PubKey
}

func NewValidatorAccount(mnemonic string) (ValidatorAccount, error) {
	if mnemonic == "" {
		return ValidatorAccount{}, fmt.Errorf("failed to get MNEMONIC ")
	}

	privKey, err := panaceacrypto.GenSecp256k1PrivKey(mnemonic, "")
	if err != nil {
		return ValidatorAccount{}, err
	}

	return ValidatorAccount{
		privKey: privKey,
		pubKey: privKey.PubKey(),
	}, nil
}

func (v ValidatorAccount) GetAddress() string {
	return types.AccAddress(v.pubKey.Address().Bytes()).String()
}

func (v ValidatorAccount) GetPrivKey() crypto.PrivKey {
	return v.privKey
}

func (v ValidatorAccount) GetPubKey() crypto.PubKey {
	return v.pubKey
}