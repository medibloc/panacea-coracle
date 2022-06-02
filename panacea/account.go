package panacea

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/medibloc/panacea-oracle/crypto"
	log "github.com/sirupsen/logrus"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

const (
	CoinType             = 371
	AccountAddressPrefix = "panacea"
)

type OracleAccount struct {
	secp256k1PrivKey tmcrypto.PrivKey
	secp256k1PubKey  tmcrypto.PubKey
	hrp              string
}

func NewOracleAccount(mnemonic string) (*OracleAccount, error) {
	privKey, err := crypto.GeneratePrivateKeyFromMnemonic(mnemonic, CoinType)

	if err != nil {
		return &OracleAccount{}, err
	}

	return &OracleAccount{
		secp256k1PrivKey: privKey,
		secp256k1PubKey:  privKey.PubKey(),
		hrp:              AccountAddressPrefix,
	}, nil
}

func (v OracleAccount) GetAddress() string {
	address, err := bech32.ConvertAndEncode(v.hrp, v.secp256k1PubKey.Address().Bytes())
	if err != nil {
		log.Panic(err)
	}
	return address
}

func (v OracleAccount) AccAddressFromBech32() sdk.AccAddress {
	return v.secp256k1PubKey.Bytes()
}

func (v OracleAccount) GetSecp256k1PrivKey() tmcrypto.PrivKey {
	return v.secp256k1PrivKey
}

func (v OracleAccount) GetSecp256k1PubKey() tmcrypto.PubKey {
	return v.secp256k1PubKey
}
