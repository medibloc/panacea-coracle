package datapool

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/crypto"

	log "github.com/sirupsen/logrus"
)

func (svc *dataPoolService) handleSampleData(w http.ResponseWriter, r *http.Request) {
	poolID := mux.Vars(r)["poolId"]
	round := mux.Vars(r)["round"]

	data := []byte("{ " +
		"\"name\": \"This is a name\", " +
		"\"description\": \"pool - " + poolID + " | round : " + round + "\"" +
		`"body": [{ "type": "markdown", "attributes": { "value": "val1" } }]
	}`)

	dataHash := crypto.Hash(data)

	dataWithAES256, err := crypto.EncryptDataWithAES256(svc.DataEncKey, nil, data)
	if err != nil {
		log.Error("failed to make encrypted data: ", err)
		http.Error(w, "failed to make encrypted data", http.StatusInternalServerError)
		return
	}

	filename := base64.StdEncoding.EncodeToString(dataHash)

	var path strings.Builder

	path.WriteString(poolID)
	path.WriteString("/")
	path.WriteString(round)

	err = svc.Store.UploadFile(path.String(), filename, dataWithAES256)
	if err != nil {
		log.Error("failed to store data: ", err)
		http.Error(w, "failed upload to S3", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
