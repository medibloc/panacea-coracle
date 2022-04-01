package datapool

import (
	"fmt"
	"github.com/gorilla/mux"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	"github.com/medibloc/panacea-data-market-validator/codec"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/server/response"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/validation"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func (svc *dataPoolService) handleValidateData(w http.ResponseWriter, r *http.Request) {
	if err, errorStatusCode := validateBasic(r); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), errorStatusCode)
		return
	}

	// TODO: use r.Body itself (without ReadAll), if possible
	jsonInput, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to read HTTP request body", http.StatusBadRequest)
		return
	}

	// get pool info by ID from blockchain
	pool, err := svc.PanaceaClient.GetPool(mux.Vars(r)[types.PoolIDKey])
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to get pool information", http.StatusInternalServerError)
		return
	}

	// trusted validator check
	poolParams := pool.PoolParams
	if !validation.Contains(poolParams.TrustedDataValidators, svc.ValidatorAccount.GetAddress()) {
		log.Error("not a trusted data-validator")
		http.Error(w, "invalid data validator", http.StatusBadRequest)
		return
	}

	// TODO data schemata validation
	if err := validation.ValidateJSONSchemata(jsonInput, poolParams.DataSchema); err != nil {
		log.Error(err)
		http.Error(w, "JSON schema validation failed", http.StatusForbidden)
		return
	}

	// make dataHash
	dataHash := crypto.Hash(jsonInput)

	// TODO encrypt and store data

	// response data
	unsignedCertificate, err := types.NewUnsignedDataValidationCertificateOfDataPool(
		pool,
		dataHash,
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

	resp := &datapooltypes.DataValidationCertificate{
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

func validateBasic(r *http.Request) (error, int) {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("only application/json is supported"), http.StatusUnsupportedMediaType
	} else if r.FormValue("requester_address") == "" {
		return fmt.Errorf("failed to read query parameter"), http.StatusBadRequest
	}
	return nil, 0
}
