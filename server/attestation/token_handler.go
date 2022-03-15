package attestation

import (
	"github.com/edgelesssys/ego/enclave"
	"github.com/medibloc/panacea-data-market-validator/account"
	"github.com/medibloc/panacea-data-market-validator/server/response"
	"github.com/medibloc/panacea-data-market-validator/types"
	log "github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/json"
	"net/http"
)

type TokenHandler struct {
	validatorAccount account.ValidatorAccount
	cert             []byte
}

func NewTokenHandler(validatorAccount account.ValidatorAccount, cert []byte) http.Handler {
	return TokenHandler{
		validatorAccount: validatorAccount,
		cert:             cert,
	}
}

func (t TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := enclave.CreateAzureAttestationToken(t.cert, types.AttestationProviderURL)
	if err != nil {
		log.Error(err)
		http.Error(w, "failed to create azure attestation token.", http.StatusInternalServerError)
	}

	jsonData, err := json.Marshal(types.TokenResponse{
		Token: token,
	})

	if err != nil {
		log.Error(err)
		http.Error(w, "failed to ", http.StatusInternalServerError)
	}

	response.WriteJSONResponse(w, http.StatusCreated, jsonData)
}
