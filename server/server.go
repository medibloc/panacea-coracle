package server

import (
	"context"
	"github.com/medibloc/panacea-data-market-validator/types"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	validateDataHandler, err := NewValidateDataHandler(conf)
	if err != nil {
		log.Panic(err)
	}

	router := mux.NewRouter()
	router.Handle("/validate-data/{dealId}", validateDataHandler).Methods(http.MethodPost)
	router.Use(gracefulShutdown)

	server := &http.Server{
		Handler:      router,
		Addr:         conf.HTTPListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Infof("👻 Data Validator Server Started 🎃: Serving %s", server.Addr)
		return server.ListenAndServe()
	})

	g.Go(func() error {
		// When os signal is detected, graceful shutdown starts
		// gRPC connection is closed first
		<-gCtx.Done()

		log.Info("grpc connection is closing")

		conn := gCtx.Value(types.CtxGrpcConnKey)
		if conn == nil {
			return types.ErrNoGrpcConnection
		}

		if err := conn.(*grpc.ClientConn).Close(); err != nil {
			return err
		}

		log.Info("server is closing")
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		log.Errorf("exit reason : %s \n", err)
	}
}
