package datadeal

import (
	"encoding/base64"
	"fmt"
	datadealtypes "github.com/medibloc/panacea-oracle/types/datadeal"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	markettypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	"github.com/medibloc/panacea-oracle/codec"
	"github.com/medibloc/panacea-oracle/crypto"
	"github.com/medibloc/panacea-oracle/server/response"
	"github.com/medibloc/panacea-oracle/types"
	"github.com/medibloc/panacea-oracle/validation"
	log "github.com/sirupsen/logrus"
)

func (svc *dataDealService) handleValidateData(w http.ResponseWriter, r *http.Request) {
	// content type check from header
	if err, errorStatusCode := validateBasic(r); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), errorStatusCode)
		return
	}

	// TODO: use r.Body itself (without ReadAll), if possible
	jsonInput, err := ioutil.ReadAll(r.Body)
	if err != nil || len(jsonInput) == 0 {
		log.Error(err)
		http.Error(w, "failed to read HTTP request body", http.StatusBadRequest)
		return
	}

	dealID := mux.Vars(r)[types.DealIDKey]

	// get deal info by ID from blockchain
	deal, err := svc.PanaceaClient.GetDeal(dealID)
	if err != nil {
		log.Error("failed to get deal information: ", err)
		http.Error(w, "failed to get deal information", http.StatusInternalServerError)
		return
	}

	// trusted oracle check
	if !validation.Contains(deal.TrustedOracles, svc.OracleAccount.GetAddress()) {
		log.Error("not a trusted oracle")
		http.Error(w, "invalid oracle", http.StatusBadRequest)
		return
	}

	// data schema validation
	if err := validation.ValidateJSONSchemata(jsonInput, deal.DataSchema); err != nil {
		log.Error(err)
		http.Error(w, "JSON schema validation failed", http.StatusForbidden)
		return
	}

	// encrypt and store data
	ownerPubKey, err := svc.PanaceaClient.GetPubKey(deal.Owner)
	if err != nil {
		log.Error("failed to get public key: ", err)
		http.Error(w, "failed to get public key", http.StatusInternalServerError)
		return
	}

	encryptedData, err := crypto.EncryptDataWithSecp256k1(ownerPubKey.Bytes(), jsonInput)
	if err != nil {
		log.Error("failed to encrypt data: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	log.Debug(encryptedData)

	// make dataHash and upload to s3Store
	dataHash := crypto.Hash(jsonInput)

	fileName := svc.Store.MakeRandomFilename()

	err = svc.Store.UploadFile(dealID, fileName, encryptedData)
	if err != nil {
		log.Error("failed to store data: ", err)
		http.Error(w, "failed upload to S3", http.StatusInternalServerError)
		return
	}

	// make downloadURL
	dataURL := svc.Store.MakeDownloadURL(dealID, fileName)
	encryptedDataURL, err := crypto.EncryptDataWithSecp256k1(ownerPubKey.Bytes(), []byte(dataURL))
	if err != nil {
		log.Error("failed to make encryptedDataURL: ", err)
		http.Error(w, "failed to make encryptedDataURL", http.StatusInternalServerError)
		return
	}

	unsignedCertificate, err := datadealtypes.NewUnsignedDataCert(
		dealID,
		dataHash,
		encryptedDataURL,
		r.URL.Query().Get(types.RequesterAddressParamKey),
		svc.OracleAccount.GetAddress())
	if err != nil {
		log.Error("failed to make unsignedDataCert: ", err)
		http.Error(w, "failed to make unsignedDataCert", http.StatusInternalServerError)
		return
	}

	serializedCertificate, err := unsignedCertificate.Marshal()
	if err != nil {
		log.Error("failed to make marshal unsignedDataCert: ", err)
		http.Error(w, "failed to make marshal unsignedDataCert", http.StatusInternalServerError)
		return
	}

	signature, err := svc.OracleAccount.GetSecp256k1PrivKey().Sign(serializedCertificate)
	if err != nil {
		log.Error("failed to make signature: ", err)
		http.Error(w, "failed to make signature", http.StatusInternalServerError)
		return
	}

	resp := &markettypes.DataCert{
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

	log.Infof("data validation completed for deal %s: %s", dealID, base64.StdEncoding.EncodeToString(dataHash))
	response.WriteJSONResponse(w, http.StatusCreated, marshaledResp)
}

func validateBasic(r *http.Request) (error, int) {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("only application/json is supported"), http.StatusUnsupportedMediaType
	}

	if r.URL.Query().Get(types.RequesterAddressParamKey) == "" {
		return fmt.Errorf("failed to read query parameter"), http.StatusBadRequest
	}
	return nil, 0
}
