package attestation

import "net/http"

type TokenHandler struct{}

func NewTokenHandler() http.Handler {
	return TokenHandler{}
}

func (t TokenHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}
