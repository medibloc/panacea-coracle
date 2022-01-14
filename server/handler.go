package server

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcec"
	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/store"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/validation"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// content type check from header
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Only application/json is supported", http.StatusUnsupportedMediaType)
		return
	}

	// TODO: use r.Body itself (without ReadAll), if possible
	jsonInput, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to read HTTP request body", http.StatusBadRequest)
		return
	}
	log.Debug(string(jsonInput))

	// TODO: get deal information from panacea
	desiredSchemaURI := "https://json.schemastore.org/github-issue-forms.json"

	// TODO: check if data validator is trusted or not

	// validate data (schema check)
	if err := validation.ValidateJSONSchema(jsonInput, desiredSchemaURI); err != nil {
		log.Error(err)
		http.Error(w, "JSON schema validation failed", http.StatusForbidden)
		return
	}

	dealId := mux.Vars(r)[types.DealIdKey]

	// encrypt and store data
	// TODO: get recipient pub key info from blockchain
	tempPrivKey, _ := btcec.NewPrivateKey(btcec.S256())
	tempPubKey := tempPrivKey.PubKey()
	encryptedData, err := crypto.EncryptData(tempPubKey.SerializeCompressed(), jsonInput)
	if err != nil {
		log.Error("failed to encrypt data: ", err)
	}
	log.Debug(encryptedData)

	// make dataHash and upload to S3
	dataHash := hex.EncodeToString(crypto.Hash(jsonInput))
	err = store.UploadFile(dealId, dataHash, encryptedData)
	if err != nil {
		log.Error("failed to store data: ", err)
	}

	// make downloadURL
	dataURL := store.MakeDownloadURL(dealId, dataHash)
	encryptedDataURL, err := crypto.EncryptData(tempPubKey.SerializeCompressed(), []byte(dataURL))
	if err != nil {
		log.Error("failed to make encryptedDataURL: ", err)
	}

	resp := &types.CertificateResponse{}

	resp.Certificate.DealId = dealId
	resp.Certificate.DataHash = dataHash
	resp.Certificate.EncryptedDataURL = hex.EncodeToString(encryptedDataURL)

	// sign certificate
	marshaledResp, err := json.Marshal(resp)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to marshal HTTP JSON response", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, marshaledResp)
}
