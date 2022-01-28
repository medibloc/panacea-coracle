package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func gracefulShutdown(handler http.Handler) http.Handler {
	log.Info("graceful shutdown middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		select {
		case <-ctx.Done():
			log.Info("Graceful handler exit")
			w.WriteHeader(http.StatusInternalServerError)
			return

		default:
			handler.ServeHTTP(w, r)
		}
	})
}
