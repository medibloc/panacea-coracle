package datapool

import (
	"context"
	"encoding/base64"
	"fmt"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/medibloc/panacea-data-market-validator/crypto"

	"github.com/medibloc/panacea-data-market-validator/types"

	"golang.org/x/sync/errgroup"

	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	log "github.com/sirupsen/logrus"
	"strconv"

	"github.com/gorilla/mux"

	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	"github.com/medibloc/panacea-data-market-validator/types"

	log "github.com/sirupsen/logrus"
)

func (svc *dataPoolService) handleDownloadData(w http.ResponseWriter, r *http.Request) {
	if err, errStatusCode := validateDownloadRequest(r); err != nil {
		log.Errorf("invalid download request: %v", err)
		http.Error(w, err.Error(), errStatusCode)
		return
	}

	redeemer := r.FormValue("requester_address")

	poolID, err := strconv.ParseUint(mux.Vars(r)[types.PoolIDKey], 10, 64)
	if err != nil {
		log.Errorf("invalid pool ID: %v", err)
		http.Error(w, "invalid pool ID", http.StatusBadRequest)
		return
	}

	redeemHistory, err := svc.PanaceaClient.GetDataPassRedeemHistory(redeemer, poolID)
	if err != nil {
		log.Errorf("failed to get redeem receipt: %v", err)
		http.Error(w, "failed to get data pass redeem receipt", http.StatusInternalServerError)
		return
	}

	if len(redeemHistory.DataPassRedeemReceipts) == 0 {
		log.Errorf("redeem receipt not found under %s", redeemHistory.Redeemer)
		http.Error(w, "redeem receipt not found", http.StatusNotFound)
		return
	}

	redeemedRound := getRedeemedRound(redeemHistory.DataPassRedeemReceipts)

	for round := uint64(1); round <= redeemedRound; round++ {
		certs, err := svc.PanaceaClient.GetDataCerts(poolID, round)
		if err != nil {
			log.Errorf("failed to get data certs: %v", err)
			http.Error(w, "failed to get data certs", http.StatusInternalServerError)
			return
		}

		log.Info(certs)
	}

	return
}

func validateDownloadRequest(r *http.Request) (error, int) {
	if r.FormValue("requester_address") == "" {
		return fmt.Errorf("failed to read query parameter"), http.StatusBadRequest
	}

	return nil, 0
}

func getRedeemedRound(receipts []datapooltypes.DataPassRedeemReceipt) uint64 {
	maxRound := receipts[0].Round

	for _, receipt := range receipts {
		if receipt.Round > maxRound {
			maxRound = receipt.Round
		}
	}

	return maxRound
func (svc *dataPoolService) handleDownloadData(w http.ResponseWriter, r *http.Request) {
	if err, errStatusCode := validateDownloadRequest(r); err != nil {
		log.Errorf("invalid download request: %v", err)
		http.Error(w, err.Error(), errStatusCode)
		return
	}

	//redeemer := r.FormValue("requester_address")

	// TODO: verify redeemer signature (w/ nonce)

	// TODO: get redeem receipt from panacea and verify it

	// TODO: get poolID and round from redeem receipt. For now, temp value
	poolIDTemp := uint64(1)
	redeemedRoundTemp := uint64(3)

	// TODO: add file format to filename
	fileFormat := ".json"

	czw := types.NewConcurrentZipWriter(w)
	defer func() {
		if err := czw.Close(); err != nil {
			log.Errorf("error occurred while closing zip writer: %v", err)
			http.Error(w, "failed to download", http.StatusInternalServerError)
		}
	}()

	g, ctx := errgroup.WithContext(context.Background())

	// get data certificates from panacea and return it
	for round := uint64(1); round <= redeemedRoundTemp; round++ {
		certs, _ := svc.PanaceaClient.GetDataCertsByRound(poolIDTemp, round)
		g.Go(func() error {
			for _, cert := range certs {
				select {
				// when ctx done, return and terminate goroutine
				case <-ctx.Done():
					return nil

				default:
					// e.g., pool 1 round 3 data -> 'pool-1-3-{dataHash}'
					filename :=
						"pool-" + strconv.FormatUint(cert.UnsignedCert.PoolId, 10) +
							"-" + strconv.FormatUint(cert.UnsignedCert.Round, 10) +
							"-" + base64.StdEncoding.EncodeToString(cert.UnsignedCert.DataHash) +
							fileFormat

					// download data from storage and decrypt it
					data, err := svc.downloadAndDecryptData(cert)
					if err != nil {
						return fmt.Errorf("error when downloading pool %d, round %d  :%w", cert.UnsignedCert.PoolId, cert.UnsignedCert.Round, err)
					}

					// zip write
					err = czw.ZipWrite(filename, data)
					if err != nil {
						return fmt.Errorf("failed to write data to %s: %w", filename, err)
					}
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Errorf("failed to download: %v", err)
		http.Error(w, "failed to download", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"pool-%d.zip\"", poolIDTemp))

	return
}

// downloadAndDecryptData downloads data by certificate and decrypt it.
func (svc *dataPoolService) downloadAndDecryptData(cert datapooltypes.DataValidationCertificate) ([]byte, error) {
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

func validateDownloadRequest(r *http.Request) (error, int) {
	if r.FormValue("requester_address") == "" {
		return fmt.Errorf("failed to read query parameter"), http.StatusBadRequest
	}

	return nil, 0
}
