package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

func gracefulShutdown(wg *sync.WaitGroup) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			select {
			case <-ctx.Done():
				log.Info("Graceful handler exit")
				w.WriteHeader(http.StatusInternalServerError)
				return

			default:
				wg.Add(1)
				defer wg.Done()
				handler.ServeHTTP(w, r)
			}
		})
	}
}
