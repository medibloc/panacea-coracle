package datapool

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/types/datapool"

	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	log "github.com/sirupsen/logrus"
)

func (svc *dataPoolService) handleDownloadData(w http.ResponseWriter, r *http.Request) {
	if err, errStatusCode := validateStreamRequest(r); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), errStatusCode)
		return
	}

	//redeemer := r.FormValue("requester_address")

	// TODO: verify redeemer signature (w/ nonce)

	// TODO: get redeem receipt from panacea and verify it

	// TODO: get poolID and round from redeem receipt. For now, temp value
	poolIDTemp := uint64(1)
	redeemedRoundTemp := uint64(3)

	filename := "pool-" + strconv.FormatUint(poolIDTemp, 10) + "-data"

	zw := zip.NewWriter(w)

	zipWriter, err := zw.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	merger := datapool.NewMerger()
	errPipeline := make(chan error, 1)
	defer close(errPipeline)

	// get data certificates from panacea and return it
	for round := uint64(1); round <= redeemedRoundTemp; round++ {
		// add output channel to merger
		merger.Add(svc.setDataPipeline(w, errPipeline, poolIDTemp, round))
	}

	for data := range merger.Merge(errPipeline) {
		if _, err = zipWriter.Write(data); err != nil {
			log.Error(err)
			http.Error(w, "internal error in data download", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", filename))

	if err := zw.Close(); err != nil {
		log.Error("error occurred while closing zip writer", err)
		return
	}

	return
}

// setDataPipeline sets pipeline for data
func (svc *dataPoolService) setDataPipeline(w http.ResponseWriter, errPipeline chan error, poolID, round uint64) <-chan []byte {
	certs, err := svc.PanaceaClient.GetDataCertsByRound(poolID, round)
	if err != nil {
		log.Error("failed to get data certificates from panacea", err)
		http.Error(w, "internal error in data download", http.StatusInternalServerError)
		errPipeline <- err
	}

	out := make(chan []byte, len(certs))

	go func() {
		defer close(out)

		for _, cert := range certs {
			select {
			case <-errPipeline:
				return
			default:
				data, err := svc.handleCert(cert)
				if err != nil {
					log.Error("error in handling certificate", err)
					http.Error(w, "internal error in data download", http.StatusInternalServerError)
					errPipeline <- err
					return
				}
				out <- data
			}

		}
	}()

	return out
}

// handleCert handles data certificate by downloading data and decrypt it.
func (svc *dataPoolService) handleCert(cert datapooltypes.DataValidationCertificate) ([]byte, error) {
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
	data, err := crypto.DecryptDataWithAES256(svc.DataEncKey, nil, cipherData)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func validateStreamRequest(r *http.Request) (error, int) {
	if r.Header.Get("Content-Type") != "application/octet-stream" {
		return fmt.Errorf("only application/octet-stream is supported"), http.StatusUnsupportedMediaType
	} else if r.FormValue("requester_address") == "" {
		return fmt.Errorf("failed to read query parameter"), http.StatusBadRequest
	}
	return nil, 0
}
