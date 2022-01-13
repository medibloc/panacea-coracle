package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func writeJSONResponse(w http.ResponseWriter, statusCode int, jsonBody []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonBody); err != nil {
		log.Errorf("failed to write HTTP JSON response body: %v", err)
	}
}
