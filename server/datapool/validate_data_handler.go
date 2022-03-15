package datapool

import (
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/store"
	"github.com/pkg/errors"
	"io/ioutil"
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

func NewValidateDataHandler(validatorAccount account.ValidatorAccount, grpcClient grpcClient, conf *config.Config) ValidateDataHandler {
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

func (v ValidateDataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Error(err)
		http.Error(w, "failed to read HTTP request body", http.StatusBadRequest)
		return
	}

	log.Info(string(requestBody))
	w.Write([]byte("Called validate data"))

}
