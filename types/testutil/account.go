package testutil

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/medibloc/panacea-oracle/panacea"
)

func GetAddress(pubKey cryptotypes.PubKey) string {
	address, err := bech32.ConvertAndEncode(panacea.AccountAddressPrefix, pubKey.Address())
	if err != nil {
		panic(err)
	}
	return address
}

func NewBaseAccount(pubKey cryptotypes.PubKey, accountNumber, sequence uint64) *authtypes.BaseAccount {
	acc := authtypes.BaseAccount{
		Address:       GetAddress(pubKey),
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
	err := acc.SetPubKey(pubKey)
	if err != nil {
		panic(err)
	}

	return &acc
}
