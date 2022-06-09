package datapool

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	panaceadatapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	"github.com/medibloc/panacea-oracle/codec"
	"github.com/medibloc/panacea-oracle/crypto"
	"github.com/medibloc/panacea-oracle/server/response"
	"github.com/medibloc/panacea-oracle/types"
	datapooltypes "github.com/medibloc/panacea-oracle/types/datapool"
	"github.com/medibloc/panacea-oracle/validation"
	log "github.com/sirupsen/logrus"
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

	poolID := mux.Vars(r)[types.PoolIDKey]

	// get pool info by ID from blockchain
	pool, err := svc.PanaceaClient.GetPool(poolID)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to get pool information", http.StatusInternalServerError)
		return
	}

	// trusted oracle check
	poolParams := pool.PoolParams
	if !validation.Contains(poolParams.TrustedOracles, svc.OracleAccount.GetAddress()) {
		log.Error("not a trusted oracle")
		http.Error(w, "invalid oracle", http.StatusBadRequest)
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

	dataWithAES256, err := crypto.EncryptDataWithAES256(svc.DataEncKey, nil, jsonInput)
	if err != nil {
		log.Error("failed to make encrypted data: ", err)
		http.Error(w, "failed to make encrypted data", http.StatusInternalServerError)
		return
	}

	filename := base64.StdEncoding.EncodeToString(dataHash)

	round := pool.Round

	var path strings.Builder

	path.WriteString(poolID)
	path.WriteString("/")
	path.WriteString(strconv.FormatUint(round, 10))

	err = svc.Store.UploadFile(path.String(), filename, dataWithAES256)
	if err != nil {
		log.Error("failed to store data: ", err)
		http.Error(w, "failed upload to S3", http.StatusInternalServerError)
		return
	}

	// response data
	unsignedCertificate, err := datapooltypes.NewUnsignedDataCert(
		pool,
		dataHash,
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

	resp := &panaceadatapooltypes.DataCert{
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

	log.Infof("data validation completed for pool %s: %s", poolID, base64.StdEncoding.EncodeToString(dataHash))

	response.WriteJSONResponse(w, http.StatusCreated, marshaledResp)
}

func validateBasic(r *http.Request) (error, int) {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("only application/json is supported"), http.StatusUnsupportedMediaType
	} else if r.URL.Query().Get(types.RequesterAddressParamKey) == "" {
		return fmt.Errorf("failed to read query parameter"), http.StatusBadRequest
	}
	return nil, 0
}
