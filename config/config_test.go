package config_test

import (
	"testing"

	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMustLoad(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("DATAVAL_AWS_S3_REGION", "ap-northeast-2")

	c := config.MustLoad()
	require.Equal(t, log.DebugLevel, log.Level(c.LogLevel))
	require.Equal(t, "0.0.0.0:8181", c.HTTPListenAddr)
	require.Equal(t, "0.0.0.0:9191", c.PanaceaGrpcAddress)
	require.Equal(t, "Your MNEMONIC", c.ValidatorMnemonic)
	require.Equal(t, "my-s3-bucket", c.AWSS3Bucket)
	require.Equal(t, "ap-northeast-2", c.AWSS3Region)
}

func TestMustLoad_MissingRequiredEnv_HttpAddr(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("DATAVAL_AWS_S3_REGION", "ap-northeast-2")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_GrpcAddr(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("DATAVAL_AWS_S3_REGION", "ap-northeast-2")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_ValidatorMnemonic(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("DATAVAL_AWS_S3_REGION", "ap-northeast-2")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_AWSS3Bucket(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("DATAVAL_AWS_S3_REGION", "ap-northeast-2")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_AWSS3Region(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_InvalidLogLevel(t *testing.T) {
	t.Setenv("DATAVAL_LOG_LEVEL", "hello")
	t.Setenv("DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")

	require.Panics(t, func() {
		config.MustLoad()
	})
}
