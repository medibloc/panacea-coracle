package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/medibloc/panacea-data-market-validator/config"

	"github.com/gorilla/mux"
	panaceaapp "github.com/medibloc/panacea-core/v2/app"
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/store"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/validation"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var (
	_ http.Handler = ValidateDataHandler{}
)

type ValidateDataHandler struct {
	validatorAccount account.ValidatorAccount
	grpcClient       GrpcClient
}

// NewValidateDataHandler Create a ValidateData handler.
// Validator_MNEMONIC should be received as an environmental variable.
func NewValidateDataHandler(conf *config.Config) (http.Handler, error) {
	validatorAccount, err := account.NewValidatorAccount(conf.ValidatorMnemonic)
	if err != nil {
		return ValidateDataHandler{}, errors.Wrap(err, "failed to make ValidateDataHandler")
	}

	grpcClient := NewGrpcClient(conf.PanaceaGrpcAddress, panaceaapp.MakeEncodingConfig())

	return ValidateDataHandler{
		validatorAccount: validatorAccount,
		grpcClient:       *grpcClient,
	}, nil
}

func (v ValidateDataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// content type check from header
	if err, errCode := v.validate(r); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), errCode)
		return
	}

	// TODO: use r.Body itself (without ReadAll), if possible
	jsonInput, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to read HTTP request body", http.StatusBadRequest)
		return
	}

	dealId := mux.Vars(r)[types.DealIdKey]

	// New grpc service connected to blockchain
	grpcCli := v.grpcClient

	// get deal info by Id from blockchain
	deal, err := grpcCli.GetDeal(dealId)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to get deal information", http.StatusInternalServerError)
		return
	}

	// get validator account from mnemonic
	valAccount := v.validatorAccount

	// trusted validator check
	if !validation.Contains(deal.TrustedDataValidators, valAccount.GetAddress()) {
		log.Error("not a trusted data-validator")
		http.Error(w, "invalid data validator", http.StatusBadRequest)
		return
	}

	// data schema validation
	for _, uri := range deal.DataSchema {
		if err := validation.ValidateJSONSchema(jsonInput, uri); err != nil {
			log.Error(err)
			http.Error(w, "JSON schema validation failed", http.StatusForbidden)
			return
		}
	}

	// encrypt and store data
	ownerPubKey, err := v.grpcClient.GetPubKey(deal.Owner)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to get public key", http.StatusInternalServerError)
		return
	}

	encryptedData, err := crypto.EncryptData(ownerPubKey.Bytes(), jsonInput)
	if err != nil {
		log.Error("failed to encrypt data: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	log.Debug(encryptedData)

	// make dataHash and upload to s3Store
	dataHash := base64.StdEncoding.EncodeToString(crypto.Hash(jsonInput))

	s3Store, err := store.NewDefaultS3Store()
	if err != nil {
		log.Error("failed to create s3Store: ", err)
		http.Error(w, "failed to create s3Store", http.StatusInternalServerError)
		return
	}
	fileName := s3Store.MakeRandomFilename()
	err = s3Store.UploadFile(dealId, fileName, encryptedData)
	if err != nil {
		log.Error("failed to store data: ", err)
		http.Error(w, "failed upload to S3", http.StatusInternalServerError)
		return
	}

	// make downloadURL
	dataURL := s3Store.MakeDownloadURL(dealId, fileName)
	encryptedDataURL, err := crypto.EncryptData(ownerPubKey.Bytes(), []byte(dataURL))
	if err != nil {
		log.Error("failed to make encryptedDataURL: ", err)
		http.Error(w, "failed to make encryptedDataURL", http.StatusInternalServerError)
		return
	}

	unsignedCertificate, err := types.NewUnsignedDataValidationCertificate(
		dealId,
		dataHash,
		base64.StdEncoding.EncodeToString(encryptedDataURL),
		r.FormValue("requester_address"),
		v.validatorAccount.GetAddress())
	if err != nil {
		log.Error("failed to make unsignedDataValidationCertificate: ", err)
		http.Error(w, "failed to make unsignedDataValidationCertificate", http.StatusInternalServerError)
		return
	}

	serializedCertificate, err := unsignedCertificate.Marshal()
	if err != nil {
		log.Error("failed to make marshal unsignedDataValidationCertificate: ", err)
		http.Error(w, "failed to make marshal unsignedDataValidationCertificate", http.StatusInternalServerError)
		return
	}

	signature, err := crypto.SignData(valAccount.GetPrivKey().Bytes(), serializedCertificate)
	if err != nil {
		log.Error("failed to make signature: ", err)
		http.Error(w, "failed to make signature", http.StatusInternalServerError)
		return
	}

	resp := types.NewDataValidationCertificate(unsignedCertificate, signature)

	// sign certificate
	marshaledResp, err := json.Marshal(resp)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to marshal HTTP JSON response", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, marshaledResp)
}

// validate Verification of parameter
func (v ValidateDataHandler) validate(r *http.Request) (error, int) {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("only application/json is supported"), http.StatusUnsupportedMediaType
	}

	if r.FormValue("requester_address") == "" {
		return fmt.Errorf("failed to read query parameter"), http.StatusBadRequest
	}
	return nil, 0
}