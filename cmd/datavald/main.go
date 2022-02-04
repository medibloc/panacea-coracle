package main

import (
	"time"

	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf := config.MustLoad()

	log.SetLevel(log.Level(conf.LogLevel))
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

	server.Run(&conf)
}
