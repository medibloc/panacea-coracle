package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	authmiddleware "github.com/medibloc/panacea-data-market-validator/server/middleware/auth"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/medibloc/panacea-data-market-validator/server/service"
	"github.com/medibloc/panacea-data-market-validator/server/service/datadeal"
	"github.com/medibloc/panacea-data-market-validator/server/service/datapool"
	"github.com/medibloc/panacea-data-market-validator/server/service/tee"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
)

func Run(conf *config.Config) error {
	svc, err := service.New(conf)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	defer svc.Close()

	router := mux.NewRouter()
	datadeal.RegisterHandlers(svc, router)
	datapool.RegisterHandlers(svc, router)
	tee.RegisterHandlers(svc, router)

	middleware := authmiddleware.NewMiddleware(svc)
	router.Use(middleware.Middleware)
	datapool.RegisterMiddleware(middleware)

	server := &http.Server{
		Handler:      router,
		Addr:         conf.HTTP.ListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	httpServerErrCh := make(chan error, 1)
	go func() {
		log.Infof("👻 Data Validator Server Started 🎃: Serving %s", server.Addr)
		if err := listenAndServe(server, svc.TLSCert); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				httpServerErrCh <- err
			} else {
				close(httpServerErrCh)
			}
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	select {
	case err := <-httpServerErrCh:
		if err != nil {
			log.Errorf("http server was closed with an error: %v", err)
		}
	case <-signalCh:
		log.Info("signal detected")
	}
	log.Info("starting the graceful shutdown")

	log.Info("terminating HTTP server")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Shutdown(ctxTimeout); err != nil {
		return fmt.Errorf("error occurs while server shutting down: %w", err)
	}

	return nil
}

func listenAndServe(server *http.Server, tlsCert *tls.Certificate) error {
	if tlsCert != nil {
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{*tlsCert},
		}
		return server.ListenAndServeTLS("", "")
	} else {
		return server.ListenAndServe()
	}
}
