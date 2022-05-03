package datapool

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/server/service"
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
	router.HandleFunc("sample/pool/{poolId}/round/{round}", s.handleSampleData).Methods(http.MethodGet)
}
