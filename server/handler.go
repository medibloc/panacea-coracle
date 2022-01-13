package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

	// TODO: encrypt and store data

	// TODO: sign certificate

	resp := &types.CertificateResponse{}

	dealId := mux.Vars(r)[types.DealIdKey]
	resp.Certificate.DealId = dealId

	marshaledResp, err := json.Marshal(resp)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to marshal HTTP JSON response", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, marshaledResp)
}
