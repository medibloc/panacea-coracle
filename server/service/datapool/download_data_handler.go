package datapool

import (
	"github.com/gorilla/context"
	"github.com/medibloc/panacea-data-market-validator/types"
	"net/http"
)

func (svc *dataPoolService) handleDownloadData(writer http.ResponseWriter, request *http.Request) {
	requesterAddress := context.Get(request, types.RequesterAddressKey)
	if requesterAddress != request.FormValue("requester_address") {
		http.Error(writer, "Do not matched signer and requester", http.StatusBadRequest)
		return
	}
	panic("implement me")
}
