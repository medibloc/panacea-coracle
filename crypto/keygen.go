package crypto

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultAccount      = 0
	defaultAddressIndex = 0
)

// GenerateMnemonic create a new mnemonic
func GenerateMnemonic() (string, error) {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// GeneratePrivateKeyFromMnemonic when a valid mnemonic is inputted, it returns a PrivKey and error is nil.
// If the mnemonic is not valid, an error is not nil.
func GeneratePrivateKeyFromMnemonic(mnemonic string) (secp256k1.PrivKey, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}


	hdPath := hd.NewFundraiserParams(defaultAccount, sdk.GetConfig().GetCoinType(), defaultAddressIndex).String()
	master, ch := hd.ComputeMastersFromSeed(bip39.NewSeed(mnemonic, ""))

	return hd.DerivePrivateKeyForPath(master, ch, hdPath)
}