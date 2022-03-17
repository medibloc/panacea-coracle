package account

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

type ValidatorAccount struct {
	privKey tmcrypto.PrivKey
	pubKey  tmcrypto.PubKey
	hrp     string
}

func NewValidatorAccount(mnemonic, hrp string, coinType uint32) (ValidatorAccount, error) {
	privKey, err := crypto.GeneratePrivateKeyFromMnemonic(mnemonic, coinType)

	if err != nil {
		return ValidatorAccount{}, err
	}

	return ValidatorAccount{
		privKey: privKey,
		pubKey:  privKey.PubKey(),
		hrp:     hrp,
	}, nil
}

func (v ValidatorAccount) GetAddress() string {
	address, err := bech32.ConvertAndEncode(v.hrp, v.pubKey.Address().Bytes())
	if err != nil {
		panic(err)
	}
	return address
}

func (v ValidatorAccount) GetPrivKey() tmcrypto.PrivKey {
	return v.privKey
}

func (v ValidatorAccount) GetPubKey() tmcrypto.PubKey {
	return v.pubKey
}
