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
	if err, errStatusCode := validateBasic(r); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), errStatusCode)
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

	// get dataCerts from panacea and re-encrypt all the data
	for round := uint64(1); round <= redeemedRound; round++ {
		certs, err := svc.PanaceaClient.GetDataCertsByRound(poolID, round)
		if err != nil {
			log.Error(err)
			http.Error(w, "failed to get data certificates", http.StatusInternalServerError)
			return
		}

		for _, cert := range certs {
			encryptedCert, err := svc.getAndEncryptDataCert(redeemer, cert)
			if err != nil {
				log.Error(err)
				http.Error(w, "failed to handle data certificates", http.StatusInternalServerError)
				return
			}

			_, err = w.Write(encryptedCert)
			if err != nil {
				log.Error(err)
				http.Error(w, "failed to write data", http.StatusInternalServerError)
				return
			}

			flusher.Flush()
		}
	}

	w.WriteHeader(http.StatusOK)
	flusher.Flush()
}

func (svc *dataPoolService) getAndEncryptDataCert(redeemer string, cert datapooltypes.DataValidationCertificate) ([]byte, error) {
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
	plainData, err := crypto.DecryptDataWithAES256(svc.DataEncKey, nil, cipherData)
	if err != nil {
		return nil, err
	}

	// get pubkey of redeemer
	pubKey, err := svc.PanaceaClient.GetPubKey(redeemer)
	if err != nil {
		return nil, err
	}

	// re-encrypt data
	reEncryptedData, err := crypto.EncryptDataWithSecp256k1(pubKey.Bytes(), plainData)
	if err != nil {
		return nil, err
	}

	return reEncryptedData, nil
}
