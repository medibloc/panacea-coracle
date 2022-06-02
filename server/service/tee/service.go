package tee

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-oracle/server/service"
)

type teeService struct {
	*service.Service
}

func RegisterHandlers(svc *service.Service, router *mux.Router) {
	s := &teeService{
		Service: svc,
	}

	router.HandleFunc("/v0/tee/attestation-token", s.handleToken).Methods(http.MethodGet)
}
