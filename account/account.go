package account

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/btcsuite/btcd/btcec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

type ValidatorAccount struct {
	privKey      tmcrypto.PrivKey
	pubKey       tmcrypto.PubKey
	ecdsaPrivKey *ecdsa.PrivateKey
	ecdsaPubKey  *ecdsa.PublicKey
}

func NewValidatorAccount(mnemonic string) (ValidatorAccount, error) {
	privKey, err := crypto.GeneratePrivateKeyFromMnemonic(mnemonic)

	if err != nil {
		return ValidatorAccount{}, err
	}

	ecdsaPrivKey, ecdsaPubKey := btcec.PrivKeyFromBytes(elliptic.P256(), privKey.Bytes())

	return ValidatorAccount{
		privKey:      privKey,
		pubKey:       privKey.PubKey(),
		ecdsaPrivKey: (*ecdsa.PrivateKey)(ecdsaPrivKey),
		ecdsaPubKey:  (*ecdsa.PublicKey)(ecdsaPubKey),
	}, nil
}

func (v ValidatorAccount) GetAddress() string {
	return sdk.AccAddress(v.pubKey.Address().Bytes()).String()
}

func (v ValidatorAccount) GetSecp256PrivKey() tmcrypto.PrivKey {
	return v.privKey
}

func (v ValidatorAccount) GetSecp256PubKey() tmcrypto.PubKey {
	return v.pubKey
}

func (v ValidatorAccount) GetEcdsaPrivKey() *ecdsa.PrivateKey {
	return v.ecdsaPrivKey
}

func (v ValidatorAccount) GetEcdsaPubKey() *ecdsa.PublicKey {
	return v.ecdsaPubKey
}
