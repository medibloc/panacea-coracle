package server

import (
	"net/http"
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
	defer func() {
		if err := ctx.Close(); err != nil {
			log.Panic(err)
		}
	}()

	validateDataHandler, err := NewValidateDataHandler(ctx, conf)
	if err != nil {
		log.Panic(err)
	}
	router := mux.NewRouter()
	router.Handle("/validate-data/{dealId}", validateDataHandler).Methods(http.MethodPost)

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