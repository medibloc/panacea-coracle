package datadeal

import (
	"fmt"

	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	markettypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	"github.com/medibloc/panacea-data-market-validator/codec"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/server/response"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/validation"
	log "github.com/sirupsen/logrus"
)

func (svc *dataDealService) handleValidateData(w http.ResponseWriter, r *http.Request) {
	// content type check from header
	if err, errCode := validateHeaders(r); err != nil {
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
	deal, err := svc.PanaceaClient.GetDeal(dealId)
	if err != nil {
		log.Error("failed to get deal information: ", err)
		http.Error(w, "failed to get deal information", http.StatusInternalServerError)
		return
	}

	// trusted validator check
	if !validation.Contains(deal.TrustedDataValidators, svc.ValidatorAccount.GetAddress()) {
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
	ownerPubKey, err := svc.PanaceaClient.GetPubKey(deal.Owner)
	if err != nil {
		log.Error("failed to get public key: ", err)
		http.Error(w, "failed to get public key", http.StatusInternalServerError)
		return
	}

	encryptedData, err := crypto.EncryptDataWithSecp256k1(ownerPubKey, jsonInput)
	if err != nil {
		log.Error("failed to encrypt data: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	log.Debug(encryptedData)

	// make dataHash and upload to s3Store
	dataHash := crypto.Hash(jsonInput)

	fileName := svc.Store.MakeRandomFilename()
	err = svc.Store.UploadFile(dealId, fileName, encryptedData)
	if err != nil {
		log.Error("failed to store data: ", err)
		http.Error(w, "failed upload to S3", http.StatusInternalServerError)
		return
	}

	// make downloadURL
	dataURL := svc.Store.MakeDownloadURL(dealId, fileName)
	encryptedDataURL, err := crypto.EncryptDataWithSecp256k1(ownerPubKey, []byte(dataURL))
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
		svc.ValidatorAccount.GetAddress())
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

	signature, err := svc.ValidatorAccount.GetSecp256k1PrivKey().Sign(serializedCertificate)
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

func validateHeaders(r *http.Request) (error, int) {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("only application/json is supported"), http.StatusUnsupportedMediaType
	}

	if r.FormValue("requester_address") == "" {
		return fmt.Errorf("failed to read query parameter"), http.StatusBadRequest
	}
	return nil, 0
}
