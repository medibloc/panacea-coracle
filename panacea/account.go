package panacea

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	log "github.com/sirupsen/logrus"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

const (
	CoinType             = 371
	AccountAddressPrefix = "panacea"
)

type ValidatorAccount struct {
	privKey tmcrypto.PrivKey
	pubKey  tmcrypto.PubKey
	hrp     string
}

func NewValidatorAccount(mnemonic string) (ValidatorAccount, error) {
	privKey, err := crypto.GeneratePrivateKeyFromMnemonic(mnemonic, CoinType)

	if err != nil {
		return ValidatorAccount{}, err
	}

	return ValidatorAccount{
		privKey: privKey,
		pubKey:  privKey.PubKey(),
		hrp:     AccountAddressPrefix,
	}, nil
}

func (v ValidatorAccount) GetAddress() string {
	address, err := bech32.ConvertAndEncode(v.hrp, v.pubKey.Address().Bytes())
	if err != nil {
		log.Panic(err)
	}
	return address
}

func (v ValidatorAccount) GetPrivKey() tmcrypto.PrivKey {
	return v.privKey
}

func (v ValidatorAccount) GetPubKey() tmcrypto.PubKey {
	return v.pubKey
}
