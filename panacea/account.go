package panacea

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/btcsuite/btcd/btcec"
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
	secp256PrivKey tmcrypto.PrivKey
	secp256PubKey  tmcrypto.PubKey
	ecdsaPrivKey   *ecdsa.PrivateKey
	ecdsaPubKey  *ecdsa.PublicKey
	hrp     string
}

func NewValidatorAccount(mnemonic string) (*ValidatorAccount, error) {
	privKey, err := crypto.GeneratePrivateKeyFromMnemonic(mnemonic, CoinType)

	if err != nil {
		return &ValidatorAccount{}, err
	}

	ecdsaPrivKey, ecdsaPubKey := btcec.PrivKeyFromBytes(elliptic.P256(), privKey.Bytes())

	return &ValidatorAccount{
		secp256PrivKey: privKey,
		secp256PubKey:  privKey.PubKey(),
		ecdsaPrivKey:   ecdsaPrivKey.ToECDSA(),
		ecdsaPubKey:    ecdsaPubKey.ToECDSA(),
		hrp:            AccountAddressPrefix,
	}, nil
}

func (v ValidatorAccount) GetAddress() string {
	address, err := bech32.ConvertAndEncode(v.hrp, v.secp256PubKey.Address().Bytes())
	if err != nil {
		log.Panic(err)
	}
	return address
}


func (v ValidatorAccount) GetSecp256PrivKey() tmcrypto.PrivKey {
	return v.secp256PrivKey
}

func (v ValidatorAccount) GetSecp256PubKey() tmcrypto.PubKey {
	return v.secp256PubKey
}

func (v ValidatorAccount) GetEcdsaPrivKey() *ecdsa.PrivateKey {
	return v.ecdsaPrivKey
}

func (v ValidatorAccount) GetEcdsaPubKey() *ecdsa.PublicKey {
	return v.ecdsaPubKey
}
