package datapool

import (
	"github.com/medibloc/panacea-oracle/server/middleware/auth"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/medibloc/panacea-oracle/server/service"
)

type dataPoolService struct {
	*service.Service
}

func RegisterHandlers(svc *service.Service, router *mux.Router) {
	s := &dataPoolService{
		Service: svc,
	}

	router.HandleFunc("/v0/data-pool/pools/{poolId}/rounds/{round}/data", s.handleValidateData).Methods(http.MethodPost)
	router.HandleFunc("/v0/data-pool/pools/{poolId}/data", s.handleDownloadData).Methods(http.MethodGet)
}

func RegisterMiddleware(auth *auth.AuthenticationMiddleware) {
	auth.AddURL("/v0/data-pool/pools/{poolId}/data", http.MethodGet)
}
