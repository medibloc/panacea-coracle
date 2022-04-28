package datapool

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"fmt"

	"github.com/tendermint/tendermint/libs/json"

	"github.com/medibloc/panacea-data-market-validator/crypto"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/panacea"
	"github.com/medibloc/panacea-data-market-validator/types"
	log "github.com/sirupsen/logrus"
)

func (svc *dataPoolService) setSample(w http.ResponseWriter, r *http.Request) {
	poolID := mux.Vars(r)[types.PoolIDKey]
	poolIDUint, _ := strconv.ParseUint(poolID, 10, 64)
	round, _ := strconv.ParseUint(mux.Vars(r)["round"], 10, 64)
	
	fmt.Print("pool id : ", poolID, "\n")
	fmt.Print("round : ", round, "\n")

	var path strings.Builder

	path.WriteString(poolID)
	path.WriteString("/")
	path.WriteString(strconv.FormatUint(round, 10))

	certs, err := panacea.MakeTestDataCerts(poolIDUint, round)
	if err != nil {
		log.Error("failed to make sample certs")
		return
	}
	for _, cert := range certs {
		marshaledCert, err := json.Marshal(cert)
		if err != nil {
			log.Error("failed marshal cert")
			return
		}
		dataWithAES256, err := crypto.EncryptDataWithAES256(svc.DataEncKey, nil, marshaledCert)
		if err != nil {
			log.Error("failed to encrypt data")
			return
		}
		filename := base64.StdEncoding.EncodeToString(cert.UnsignedCert.DataHash)
		fmt.Print("path : ", path.String(), " / filename : ", filename, "\n")
		svc.Store.UploadFile(path.String(), filename, dataWithAES256)
	}

	w.WriteHeader(http.StatusCreated)
}

