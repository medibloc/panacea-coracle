package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

// TODO: Use a better name
const envVarPrefix = "EDG_DATAVAL_"

type Config struct {
	LogLevel                      LogLevel `envconfig:"EDG_DATAVAL_LOG_LEVEL" default:"info"`
	HTTPListenAddr                string   `envconfig:"EDG_DATAVAL_HTTP_LADDR" required:"true"`
	PanaceaGrpcAddress            string   `envconfig:"EDG_DATAVAL_PANACEA_GRPC_ADDR" required:"true"`
	ValidatorMnemonic             string   `envconfig:"EDG_DATAVAL_VALIDATOR_MNEMONIC" required:"true"`
	PublicEndPointUrl             string   `envconfig:"EDG_DATAVAL_PUBLIC_ENDPOINT_URL" required:"true"`
	AWSS3Bucket                   string   `envconfig:"EDG_DATAVAL_AWS_S3_BUCKET" required:"true"`
	AWSS3Region                   string   `envconfig:"EDG_DATAVAL_AWS_S3_REGION" required:"true"`
	AWSS3AccessKeyID              string   `envconfig:"EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID" required:"true"`
	AWSS3SecretAccessKey          string   `envconfig:"EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY" required:"true"`
	ConfigDir                     string   `envconfig:"EDG_DATAVAL_CONFIG_DIR" required:"true"`
	EnclaveAttestationProviderURL string   `envconfig:"EDG_DATAVAL_ENCLAVE_ATTESTATION_PROVIDER_URL" required:"true"`
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
		log.Panic(err)
	}
	return conf
}
