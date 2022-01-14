package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const envVarPrefix = "DATAVAL"

type Config struct {
	LogLevel       LogLevel `envconfig:"DATAVAL_LOG_LEVEL" default:"info"`
	HTTPListenAddr string   `envconfig:"DATAVAL_HTTP_LADDR" required:"true"`
}

// LogLevel is a type aliasing for the envconfig custom decoder.
// https://github.com/kelseyhightower/envconfig#custom-decoders
type LogLevel log.Level

func (d *LogLevel) Decode(value string) error {
	lvl, err := log.ParseLevel(value)
	if err != nil {
		return err
	}
	*d = LogLevel(lvl)
	return nil
}

func MustLoad() Config {
	var conf Config
	if err := envconfig.Process(envVarPrefix, &conf); err != nil {
		log.Fatal(err)
	}
	return conf
}
