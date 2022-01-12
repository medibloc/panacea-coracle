package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/utils"
	"github.com/medibloc/panacea-data-market-validator/validation"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	resp := &types.CertificateResponse{}

	dealId := mux.Vars(r)[types.DealIdKey]
	resp.Certificate.DealId = dealId

	// read data from request
	data, err := utils.ReadData(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, e := w.Write([]byte(err.Error())); e != nil {
			log.Error("response write failed: ", e)
		}
		log.Error("read data failed: ", err)
		return
	}

	fmt.Println(data)

	// TODO: get deal information from panacea
	desiredSchemaURI := "https://json.schemastore.org/github-issue-forms.json"

	// TODO: check if data validator is trusted or not

	// validate data (schema check)
	if err := validation.ValidateJSONSchema(data, desiredSchemaURI); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusForbidden)
		if _, e := w.Write([]byte(err.Error())); e != nil {
			log.Error("response write failed: ", e)
		}
		return
	}

	// TODO: encrypt and store data

	// TODO: sign certificate

	marshaledResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, e := w.Write([]byte(err.Error())); e != nil {
			log.Error("response write failed: ", e)
		}
		log.Error("response marshal failed: ", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, e := w.Write(marshaledResp); e != nil {
		log.Error("response write failed: ", e)
	}
}
