package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	resp := &types.CertificateResponse{}

	dealId := mux.Vars(r)[types.DealIdKey]
	resp.Certificate.DealId = dealId

	// file format check
	data, err := utils.ReadFormFile(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, e := w.Write([]byte(err.Error())); e != nil {
			log.Error("Response write failed: ", e)
		}
		log.Error(err)
		return
	}

	fmt.Println(data)

	// get deal information from panacea

	// check if data validator is trusted or not

	// validate data (schema check)

	// encrypt and store data

	// sign certificate

	marshaledData, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, e := w.Write([]byte(err.Error())); e != nil {
			log.Error("Response write failed: ", e)
		}
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, e := w.Write(marshaledData); e != nil {
		log.Error("Response write failed: ", e)
	}
}
