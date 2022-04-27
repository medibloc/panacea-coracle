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
	// validate basic
	if err, errStatusCode := validateBasic(r); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), errStatusCode)
		return
	}

	redeemer := r.FormValue("requester_address")

	// verify redeemer signature (w/ nonce)

	// verify redeem receipt via panacea

	flusher, _ := w.(http.Flusher)

	poolID := uint64(1)
	redeemedRound := uint64(3)
	// get dataValidationCert from panacea
	for round := uint64(1); round < redeemedRound; round++ {
		certs, err := svc.PanaceaClient.GetDataCertsByRound(poolID, round)
		if err != nil {
			log.Error(err)
			http.Error(w, "failed to get data certificates", http.StatusInternalServerError)
			return
		}

		for _, cert := range certs {
			encryptedCert, err := svc.encryptDataCert(redeemer, cert)
			if err != nil {
				log.Error(err)
				http.Error(w, "failed to handle data certificates", http.StatusInternalServerError)
				return
			}

			w.Write(encryptedCert)
			flusher.Flush()
		}
	}

	//w.Header().Set("Content-Type", "")
	w.WriteHeader(http.StatusOK)
}

func (svc *dataPoolService) encryptDataCert(redeemer string, cert datapooltypes.DataValidationCertificate) ([]byte, error) {
	var path strings.Builder
	path.WriteString(strconv.FormatUint(cert.UnsignedCert.PoolId, 10))
	path.WriteString("/")
	path.WriteString(strconv.FormatUint(cert.UnsignedCert.Round, 10))
	path.WriteString("/")
	path.WriteString(string(cert.UnsignedCert.DataHash))

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

	return reEncryptedData, nil
}
