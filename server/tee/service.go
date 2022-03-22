package tee

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/server/service"
)

type teeService struct {
	*service.Service
	Cert                   []byte
	AttestationProviderURL string
}

func RegisterHandlers(svc *service.Service, cert []byte, attestationProviderURL string, router *mux.Router) {
	s := &teeService{
		Service:                svc,
		Cert:                   cert,
		AttestationProviderURL: attestationProviderURL,
	}

	router.HandleFunc("/v0/tee/attestation-token", s.handleToken).Methods(http.MethodGet)
}
