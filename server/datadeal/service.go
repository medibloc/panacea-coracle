package datadeal

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/server/service"
)

type dataDealService struct {
	*service.Service
}

func RegisterHandlers(svc *service.Service, router *mux.Router) {
	s := &dataDealService{
		Service: svc,
	}

	router.HandleFunc("/v0/data-deal/deals/{dealId}/data", s.handleValidateData).Methods(http.MethodPost)
}
