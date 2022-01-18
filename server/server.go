package server

import (
	panaceaapp "github.com/medibloc/panacea-core/v2/app"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
)

func Run(conf *config.Config) {
	// encodingConfig for decoding google.protobuf.Any type in grpc response
	encodingConfig := panaceaapp.MakeEncodingConfig()

	router := mux.NewRouter()
	router.HandleFunc("/validate-data/{dealId}", handleRequest(conf.GrpcAddress, encodingConfig)).Methods(http.MethodPost)

	server := &http.Server{
		Handler:      router,
		Addr:         conf.HTTPListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("👻 Data Validator Server Started 🎃: Serving %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
