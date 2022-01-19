package server

import (
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-core/v2/app/params"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/store"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/validation"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

func handleRequest(grpcAddr string, encodingConfig params.EncodingConfig) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		dealId := mux.Vars(r)[types.DealIdKey]

		// New grpc service connected to blockchain
		grpcSvc := NewGrpcService(grpcAddr, encodingConfig)

		// get deal info by Id from blockchain
		deal, err := grpcSvc.GetDeal(dealId)
		if err != nil {
			log.Error(err)
			http.Error(w, "failed to get deal information", http.StatusInternalServerError)
			return
		}

		// trusted validator check
		if !validation.Contains(deal.TrustedDataValidators, types.DataValidatorAddress) {
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

		// get public key of deal owner from blockchain
		pubKeyBytes, err := grpcSvc.GetPubKey(deal.Owner)
		if err != nil {
			log.Error(err)
			http.Error(w, "failed to get public key", http.StatusInternalServerError)
			return
		}

		dealId := mux.Vars(r)[types.DealIdKey]

		// encrypt and store data
		// TODO: get recipient pub key info from blockchain
		encryptedData, err := crypto.EncryptData(pubKeyBytes, jsonInput)
		if err != nil {
			log.Error(err)
			http.Error(w, "failed to encrypt data", http.StatusInternalServerError)
		}
		log.Debug(encryptedData)

		// make dataHash and upload to s3Store
		dataHash := hex.EncodeToString(crypto.Hash(jsonInput))

		s3Store, err := store.NewDefaultS3Store()
		if err != nil {
			log.Error("failed to create s3Store: ", err)
		}

		fileName := s3Store.MakeRandomFilename()
		err = s3Store.UploadFile(dealId, fileName, encryptedData)
		if err != nil {
			log.Error("failed to store data: ", err)
		}

		// make downloadURL
		dataURL := s3Store.MakeDownloadURL(dealId, fileName)
		encryptedDataURL, err := crypto.EncryptData(pubKeyBytes, []byte(dataURL))
		if err != nil {
			log.Error("failed to make encryptedDataURL: ", err)
		}

		resp := &types.CertificateResponse{}

		resp.Certificate.DealId = dealId
		resp.Certificate.DataHash = dataHash
		resp.Certificate.EncryptedDataURL = hex.EncodeToString(encryptedDataURL)

		marshaledResp, err := json.Marshal(resp)
		if err != nil {
			log.Error(err)
			http.Error(w, "failed to marshal HTTP JSON response", http.StatusInternalServerError)
			return
		}

		writeJSONResponse(w, http.StatusCreated, marshaledResp)
	}
}
