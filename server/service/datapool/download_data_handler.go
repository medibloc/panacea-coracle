package datapool

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	"github.com/medibloc/panacea-data-market-validator/crypto"

	log "github.com/sirupsen/logrus"
)

func (svc *dataPoolService) handleDownloadData(w http.ResponseWriter, r *http.Request) {
	//if err, errStatusCode := validateBasic(r); err != nil {
	//	log.Error(err)
	//	http.Error(w, err.Error(), errStatusCode)
	//	return
	//}
	if r.FormValue("requester_address") == "" {
		log.Error("failed to read query parameter")
		http.Error(w, "failed to read query parameter", http.StatusBadRequest)
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
		certChan := svc.handleRound(poolID, round)
		res = svc.handleCert(certChan, redeemer)
	}

	for data := range res {
		w.Write(data)
		flusher.Flush()
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Connection", "close")
	flusher.Flush()
	return
}

func (svc *dataPoolService) handleRound(poolID, round uint64) <-chan datapooltypes.DataValidationCertificate {
	certs, _ := svc.PanaceaClient.GetDataCertsByRound(poolID, round)

	out := make(chan datapooltypes.DataValidationCertificate, len(certs))

	//go func() {
	for _, n := range certs {
		out <- n
	}
	close(out)
	//}()
	return out
}

func (svc *dataPoolService) handleCert(cert <-chan datapooltypes.DataValidationCertificate, redeemer string) <-chan []byte {
	out := make(chan []byte, len(cert))

	//go func() {
	for n := range cert {
		var path strings.Builder
		path.WriteString(strconv.FormatUint(n.UnsignedCert.PoolId, 10))
		path.WriteString("/")
		path.WriteString(strconv.FormatUint(n.UnsignedCert.Round, 10))

		filename := base64.StdEncoding.EncodeToString(n.UnsignedCert.DataHash)
		//fmt.Print("path : ", path.String(), " | filename : ", filename, "\n")
		fmt.Print(n, "\n")
		// download encrypted data
		cipherData, _ := svc.Store.DownloadFile(path.String(), filename)
		//if err != nil {
		//	return nil, err
		//}

		// decrypt data
		//plainData, _ := crypto.DecryptDataWithAES256(svc.DataEncKey, nil, cipherData)
		//if err != nil {
		//	return nil, err
		//}

		// get pubkey of redeemer
		pubKey, _ := svc.PanaceaClient.GetPubKey(redeemer)
		//if err != nil {
		//	return nil, err
		//}

		// re-encrypt data
		reEncryptedData, _ := crypto.EncryptDataWithSecp256k1(pubKey.Bytes(), cipherData)
		//if err != nil {
		//	return nil, err
		//}

		//return reEncryptedData, nil
		out <- reEncryptedData
	}
	//}()

	return out
}
