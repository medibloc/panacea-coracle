package datadeal

import (
	"fmt"

	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gorilla/mux"
	markettypes "github.com/medibloc/panacea-core/v2/x/market/types"
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/medibloc/panacea-data-market-validator/codec"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/server/response"
	"github.com/medibloc/panacea-data-market-validator/store"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/validation"
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

// NewValidateDataHandler creates a ValidateData handler.
func NewValidateDataHandler(validatorAccount account.ValidatorAccount, grpcClient grpcClient, conf *config.Config) http.Handler {
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

	// get deal info by Id from blockchain
	deal, err := v.grpcClient.GetDeal(dealId)
	if err != nil {
		log.Error("failed to get deal information: ", err)
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
		log.Error("failed to get public key: ", err)
		http.Error(w, "failed to get public key", http.StatusInternalServerError)
		return
	}

	encryptedData, err := crypto.EncryptData(ownerPubKey, jsonInput)
	if err != nil {
		log.Error("failed to encrypt data: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	log.Debug(encryptedData)

	// make dataHash and upload to s3Store
	dataHash := crypto.Hash(jsonInput)

	fileName := v.store.MakeRandomFilename()
	err = v.store.UploadFile(dealId, fileName, encryptedData)
	if err != nil {
		log.Error("failed to store data: ", err)
		http.Error(w, "failed upload to S3", http.StatusInternalServerError)
		return
	}

	// make downloadURL
	dataURL := v.store.MakeDownloadURL(dealId, fileName)
	encryptedDataURL, err := crypto.EncryptData(ownerPubKey, []byte(dataURL))
	if err != nil {
		log.Error("failed to make encryptedDataURL: ", err)
		http.Error(w, "failed to make encryptedDataURL", http.StatusInternalServerError)
		return
	}

	unsignedCertificate, err := types.NewUnsignedDataValidationCertificate(
		dealId,
		dataHash,
		encryptedDataURL,
		r.URL.Query().Get("requester_address"),
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

	signature, err := valAccount.GetSecp256PrivKey().Sign(serializedCertificate)
	if err != nil {
		log.Error("failed to make signature: ", err)
		http.Error(w, "failed to make signature", http.StatusInternalServerError)
		return
	}

	resp := &markettypes.DataValidationCertificate{
		UnsignedCert: &unsignedCertificate,
		Signature:    signature,
	}

	// sign certificate
	marshaledResp, err := codec.ProtoMarshalJSON(resp)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to marshal HTTP JSON response", http.StatusInternalServerError)
		return
	}

	response.WriteJSONResponse(w, http.StatusCreated, marshaledResp)
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
