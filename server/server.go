package server

import (
	"context"
	"errors"
	"github.com/medibloc/panacea-data-market-validator/server/attestation"
	"github.com/medibloc/panacea-data-market-validator/server/datadeal"
	"github.com/medibloc/panacea-data-market-validator/server/datapool"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	panaceaapp "github.com/medibloc/panacea-core/v2/app"
	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
)

func Run(conf *config.Config) {
	panaceaapp.SetConfig()

	ctx, err := newContext(conf)
	if err != nil {
		log.Panic(err)
	}

	grpcClient, err := NewGrpcClient(ctx.PanaceaConn)
	if err != nil {
		log.Panic(err)
	}

	router := mux.NewRouter()
	router.Handle("/v0/data-deal/validate-data/{dealId}", datadeal.NewValidateDataHandler(grpcClient, conf)).Methods(http.MethodPost)
	router.Handle("/v1/data-pool/pools/{poolId}/rounds/{round}/data", datapool.NewValidateDataHandler(grpcClient, conf)).Methods(http.MethodPost)
	router.Handle("/v1/data-pool/pools/{poolId}/data", datapool.NewDownloadDataHandler(grpcClient)).Methods(http.MethodGet)
	router.Handle("/v1/tee/attestation-token", attestation.NewTokenHandler()).Methods(http.MethodGet)

	server := &http.Server{
		Handler:      router,
		Addr:         conf.HTTPListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	httpServerErrCh := make(chan error, 1)
	go func() {
		log.Infof("👻 Data Validator Server Started 🎃: Serving %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
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
		log.Panicf("error occurs while server shutting down: %v", err)
	}

	log.Info("closing all other resources")
	if err := ctx.Close(); err != nil {
		log.Panicf("error occurs while closing other resources: %v", err)
	}
}
