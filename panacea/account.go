package panacea

import "github.com/medibloc/panacea-data-market-validator/account"

const (
	CoinType             = 371
	AccountAddressPrefix = "panacea"
)

func NewValidatorAccount(mnemonic string) (account.ValidatorAccount, error) {
	return account.NewValidatorAccount(mnemonic, AccountAddressPrefix, CoinType)
}
