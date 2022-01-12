package server

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/utils"
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

	// get deal information from panacea

	// check if data validator is trusted or not

	// validate data (schema check)

	// encrypt and store data
	// TODO: get recipient pub key info from blockchain
	tempPrivKey, _ := btcec.NewPrivateKey(btcec.S256())
	encryptedData, err := crypto.EncryptData(tempPrivKey.PubKey().SerializeCompressed(), data)
	if err != nil {
		log.Error("failed to encrypt data: ", err)
	}
	fmt.Println(encryptedData)

	// sign certificate

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
