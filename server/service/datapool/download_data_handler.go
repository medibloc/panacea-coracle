package datapool

import (
	"fmt"
	"net/http"
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

	fmt.Print(poolID, redeemedRound)

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
}
