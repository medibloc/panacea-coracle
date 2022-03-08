package datapool

import (
	"net/http"
)

var (
	_ http.Handler = DownloadDataHandler{}
)

type DownloadDataHandler struct {
	grpcClient grpcClient
}

func NewDownloadDataHandler(grpcClient grpcClient) http.Handler {
	return DownloadDataHandler{
		grpcClient: grpcClient,
	}
}

func (d DownloadDataHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}
