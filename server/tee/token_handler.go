package tee

import (
	"github.com/medibloc/panacea-data-market-validator/server/response"
	"github.com/medibloc/panacea-data-market-validator/tee"
	log "github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/json"
	"net/http"
)

type TokenResponse struct {
	Token string
}

func (svc *teeService) handleToken(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	token, err := tee.CreateAzureAttestationToken(svc.Cert, svc.AttestationProviderURL)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to create azure attestation token", http.StatusInternalServerError)
	}

	jsonBody, err := json.Marshal(TokenResponse{Token: token})
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
	}

	response.WriteJSONResponse(w, http.StatusCreated, jsonBody)
}
