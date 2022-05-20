package datapool

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/medibloc/panacea-data-market-validator/types"

	log "github.com/sirupsen/logrus"
)

func (svc *dataPoolService) handleDownloadData(w http.ResponseWriter, r *http.Request) {
	if err, errStatusCode := validateDownloadRequest(r); err != nil {
		log.Errorf("invalid download request: %v", err)
		http.Error(w, err.Error(), errStatusCode)
		return
	}

	poolID, err := strconv.ParseUint(mux.Vars(r)[types.PoolIDKey], 10, 64)
	if err != nil {
		log.Errorf("invalid pool ID: %v", err)
		http.Error(w, "invalid pool ID", http.StatusBadRequest)
		return
	}

	round, err := strconv.ParseUint(mux.Vars(r)[types.RoundKey], 10, 64)
	if err != nil {
		log.Errorf("invalid round: %v", err)
		http.Error(w, "invalid round", http.StatusBadRequest)
		return
	}

	dataPassID, err := strconv.ParseUint(mux.Vars(r)[types.DataPassIDKey], 10, 64)
	if err != nil {
		log.Errorf("invalid data pass ID: %v", err)
		http.Error(w, "invalid data pass ID", http.StatusBadRequest)
		return
	}

	redeemReceipt, err := svc.PanaceaClient.GetDataPassRedeemReceipt(poolID, round, dataPassID)
	if err != nil {
		log.Errorf("failed to get redeem receipt: %v", err)
		http.Error(w, "failed to get data pass redeem receipt", http.StatusInternalServerError)
		return
	}

	redeemer := r.FormValue("requester_address")

	if redeemReceipt.Redeemer != redeemer {
		log.Errorf("redeemer is not matched: requested redeemer is \"%s\", but actual redeemer is \"%s\"", redeemer, redeemReceipt.Redeemer)
		http.Error(w, "redeemer is not matched", http.StatusBadRequest)
		return
	}

	return
}

func validateDownloadRequest(r *http.Request) (error, int) {
	if r.FormValue("requester_address") == "" {
		return fmt.Errorf("failed to read query parameter"), http.StatusBadRequest
	}

	return nil, 0
}
