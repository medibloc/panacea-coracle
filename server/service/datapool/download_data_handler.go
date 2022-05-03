package datapool

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/medibloc/panacea-data-market-validator/crypto"

	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	log "github.com/sirupsen/logrus"
)

func (svc *dataPoolService) handleDownloadData(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("requester_address") == "" {
		log.Error("requester address is required")
		http.Error(w, "requester address is required", http.StatusBadRequest)
		return
	}

	//redeemer := r.FormValue("requester_address")

	// TODO: verify redeemer signature (w/ nonce)

	// TODO: get redeem receipt from panacea and verify it

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Panic("expected http.ResponseWriter to be an http.Flusher")
		http.Error(w, "internal error in data download", http.StatusInternalServerError)
		return
	}

	// TODO: get poolID and round from redeem receipt. For now, temp value
	poolID := uint64(1)
	redeemedRound := uint64(2)

	res := make(<-chan []byte)

	// get dataCerts from panacea and re-encrypt all the data
	for round := uint64(1); round <= redeemedRound; round++ {
		res = svc.handleRound(poolID, round)
	}

	for data := range res {
		_, err := w.Write(data)
		if err != nil {
			log.Error(err)
			http.Error(w, "internal error in data download", http.StatusInternalServerError)
			return
		}
		flusher.Flush()
	}

	//w.WriteHeader(http.StatusOK)
}

func (svc *dataPoolService) handleRound(poolID, round uint64) <-chan []byte {
	certs, _ := svc.PanaceaClient.GetDataCertsByRound(poolID, round)

	out := make(chan []byte, len(certs))

	go func() {
		for _, cert := range certs {
			data, err := svc.handleCert(cert)
			if err != nil {
				log.Error(err.Error())
			}
			out <- data
		}
		close(out)
	}()

	return out
}

func (svc *dataPoolService) handleCert(cert datapooltypes.DataValidationCertificate) ([]byte, error) {
	fmt.Print(cert)
	var path strings.Builder
	path.WriteString(strconv.FormatUint(cert.UnsignedCert.PoolId, 10))
	path.WriteString("/")
	path.WriteString(strconv.FormatUint(cert.UnsignedCert.Round, 10))

	filename := base64.StdEncoding.EncodeToString([]byte("data-" + strconv.FormatUint(cert.UnsignedCert.PoolId, 10) + "-" + strconv.FormatUint(cert.UnsignedCert.Round, 10)))

	fmt.Print("round : ", cert.UnsignedCert.Round, " | filename : ", filename, "\n")

	// download encrypted data
	cipherData, err := svc.Store.DownloadFile(path.String(), filename)
	if err != nil {
		return nil, err
	}

	// decrypt data
	data, err := crypto.DecryptDataWithAES256(svc.DataEncKey, nil, cipherData)
	if err != nil {
		return nil, err
	}

	return data, nil
}
