package datapool

import (
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/store"
	"github.com/pkg/errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	_ http.Handler = ValidateDataHandler{}
)

type ValidateDataHandler struct {
	validatorAccount account.ValidatorAccount
	store            store.S3Store
	grpcClient       grpcClient
}

func NewValidateDataHandler(grpcClient grpcClient, conf *config.Config) ValidateDataHandler {
	validatorAccount, err := account.NewValidatorAccount(conf.ValidatorMnemonic)
	if err != nil {
		log.Panic(errors.Wrap(err, "failed to NewValidatorAccount"))
	}

	s3Store, err := store.NewS3Store(conf.AWSS3Bucket, conf.AWSS3Region)
	if err != nil {
		log.Panic(errors.Wrap(err, "failed to create S3Store"))
	}

	return ValidateDataHandler{
		validatorAccount: validatorAccount,
		store:            s3Store,
		grpcClient:       grpcClient,
	}
}

func (v ValidateDataHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}
