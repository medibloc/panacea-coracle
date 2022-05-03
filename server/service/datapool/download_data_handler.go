package datapool

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	"github.com/medibloc/panacea-data-market-validator/crypto"

	log "github.com/sirupsen/logrus"
)

func (svc *dataPoolService) handleDownloadData(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("requester_address") == "" {
		log.Error("requester address is required")
		http.Error(w, "requester address is required", http.StatusBadRequest)
		return
	}

	redeemer := r.FormValue("requester_address")

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
	redeemedRound := uint64(3)

	res := make(<-chan []byte)

	// get dataCerts from panacea and re-encrypt all the data
	for round := uint64(1); round <= redeemedRound; round++ {
		res = svc.handleRound(poolID, round, redeemer)
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

	w.WriteHeader(http.StatusOK)
}

func (svc *dataPoolService) handleRound(poolID, round uint64, redeemer string) <-chan []byte {
	certs, _ := svc.PanaceaClient.GetDataCertsByRound(poolID, round)

	out := make(chan []byte, len(certs))

	go func() {
		for _, cert := range certs {
			reEncryptedData, err := svc.handleCert(cert, redeemer)
			if err != nil {
				log.Error("error in handling certificate")
			}
			out <- reEncryptedData
		}
		close(out)
	}()

	return out
}

func (svc *dataPoolService) handleCert(cert datapooltypes.DataValidationCertificate, redeemer string) ([]byte, error) {
	var path strings.Builder
	path.WriteString(strconv.FormatUint(cert.UnsignedCert.PoolId, 10))
	path.WriteString("/")
	path.WriteString(strconv.FormatUint(cert.UnsignedCert.Round, 10))

	filename := base64.StdEncoding.EncodeToString(cert.UnsignedCert.DataHash)

	// download encrypted data
	cipherData, err := svc.Store.DownloadFile(path.String(), filename)
	if err != nil {
		return nil, err
	}

	// decrypt data
	//plainData, _ := crypto.DecryptDataWithAES256(svc.DataEncKey, nil, cipherData)
	//if err != nil {
	//	return nil, err
	//}

	// get pubkey of redeemer
	pubKey, err := svc.PanaceaClient.GetPubKey(redeemer)
	if err != nil {
		return nil, err
	}

	// re-encrypt data
	reEncryptedData, err := crypto.EncryptDataWithSecp256k1(pubKey.Bytes(), cipherData)
	if err != nil {
		return nil, err
	}

	return reEncryptedData, nil
}
