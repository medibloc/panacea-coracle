package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
)

func Run(conf *config.Config) {
	SetConfig()

	log.Infof("dial to blockchain: %s", conf.PanaceaGrpcAddress)
	conn, err := grpc.Dial(conf.PanaceaGrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer func() {
		if err := conn.Close(); err != nil {
			log.Panic(err)
		}
	}()

	validateDataHandler, err := NewValidateDataHandler(conf, conn)
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.Handle("/validate-data/{dealId}", validateDataHandler)

	server := &http.Server{
		Handler:      router,
		Addr:         conf.HTTPListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("👻 Data Validator Server Started 🎃: Serving %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
