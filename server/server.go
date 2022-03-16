package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/medibloc/panacea-data-market-validator/server/datadeal"
	"github.com/medibloc/panacea-data-market-validator/server/datapool"
	"github.com/medibloc/panacea-data-market-validator/server/service"
	"github.com/medibloc/panacea-data-market-validator/server/tee"

	"github.com/gorilla/mux"
	panaceaapp "github.com/medibloc/panacea-core/v2/app"
	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
)

func Run(conf *config.Config) {
	panaceaapp.SetConfig()

	svc, err := service.New(conf)
	if err != nil {
		log.Panicf("failed to create service: %w", err)
	}
	defer svc.Close()

	router := mux.NewRouter()
	datadeal.RegisterHandlers(svc, router)
	datapool.RegisterHandlers(svc, router)
	tee.RegisterHandlers(svc, router)

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
		log.Errorf("error occurs while server shutting down: %v", err)
	}
}
