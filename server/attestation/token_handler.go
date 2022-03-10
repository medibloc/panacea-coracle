package attestation

import (
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/muesli/cache2go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type TokenHandler struct {
	validatorAccount account.ValidatorAccount
	cache            *cache2go.CacheTable
}

func NewTokenHandler(conf *config.Config) http.Handler {
	validatorAccount, err := account.NewValidatorAccount(conf.ValidatorMnemonic)
	if err != nil {
		log.Panic(errors.Wrap(err, "failed to NewValidatorAccount"))
	}

	cache := cache2go.Cache("AttestationToken")

	return TokenHandler{
		validatorAccount: validatorAccount,
		cache:            cache,
	}
}

func (t TokenHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}
