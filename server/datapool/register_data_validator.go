package datapool

import (
	"encoding/json"
	"github.com/medibloc/panacea-data-market-validator/server/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type DataValidator struct {
	Address  string `json:"address"`
	Endpoint string `json:"endpoint"`
}

func (svc *dataPoolService) handleRegisterDataValidator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dataValidator DataValidator

	err := json.NewDecoder(r.Body).Decode(&dataValidator)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to decode data validator", http.StatusInternalServerError)
		return
	}

	err = svc.PanaceaClient.RegisterDataValidator(dataValidator.Address, dataValidator.Endpoint)
	if err != nil {
		log.Error("failed to register data validator: ", err)
		http.Error(w, "failed to register data validator", http.StatusInternalServerError)
		return
	}

	marshaledRes, err := json.Marshal(dataValidator)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to marshal HTTP JSON response", http.StatusInternalServerError)
		return
	}

	response.WriteJSONResponse(w, http.StatusOK, marshaledRes)
}
