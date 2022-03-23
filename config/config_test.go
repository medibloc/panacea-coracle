package config_test

import (
	"testing"

	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMustLoad(t *testing.T) {
	t.Setenv("EDG_DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("EDG_DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("EDG_DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("EDG_DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("EDG_DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("EDG_DATAVAL_AWS_S3_REGION", "ap-northeast-2")
	t.Setenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID", "my-access-key")
	t.Setenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY", "my-secret-access-key")
	t.Setenv("EDG_DATAVAL_CONFIG_DIR", "my-config-dir")
	t.Setenv("EDG_DATAVAL_ATTESTATION_PROVIDER_URL", "my-attestation-provider-url")

	c := config.MustLoad()
	require.Equal(t, log.DebugLevel, log.Level(c.LogLevel))
	require.Equal(t, "0.0.0.0:8181", c.HTTPListenAddr)
	require.Equal(t, "0.0.0.0:9191", c.PanaceaGrpcAddress)
	require.Equal(t, "Your MNEMONIC", c.ValidatorMnemonic)
	require.Equal(t, "my-s3-bucket", c.AWSS3Bucket)
	require.Equal(t, "ap-northeast-2", c.AWSS3Region)
	require.Equal(t, "my-access-key", c.AWSS3AccessKeyID)
	require.Equal(t, "my-secret-access-key", c.AWSS3SecretAccessKey)
	require.Equal(t, "my-config-dir", c.ConfigDir)
	require.Equal(t, "my-attestation-provider-url", c.EnclaveAttestationProviderURL)
}

func TestMustLoad_MissingRequiredEnv_HttpAddr(t *testing.T) {
	t.Setenv("EDG_DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("EDG_DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("EDG_DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("EDG_DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("EDG_DATAVAL_AWS_S3_REGION", "ap-northeast-2")
	t.Setenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID", "my-access-key")
	t.Setenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY", "my-secret-access-key")
	t.Setenv("EDG_DATAVAL_CONFIG_DIR", "my-config-dir")
	t.Setenv("EDG_DATAVAL_ATTESTATION_PROVIDER_URL", "my-attestation-provider-url")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_GrpcAddr(t *testing.T) {
	t.Setenv("EDG_DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("EDG_DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("EDG_DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("EDG_DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("EDG_DATAVAL_AWS_S3_REGION", "ap-northeast-2")
	t.Setenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID", "my-access-key")
	t.Setenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY", "my-secret-access-key")
	t.Setenv("EDG_DATAVAL_CONFIG_DIR", "my-config-dir")
	t.Setenv("EDG_DATAVAL_ATTESTATION_PROVIDER_URL", "my-attestation-provider-url")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_ValidatorMnemonic(t *testing.T) {
	t.Setenv("EDG_DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("EDG_DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("EDG_DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("EDG_DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("EDG_DATAVAL_AWS_S3_REGION", "ap-northeast-2")
	t.Setenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID", "my-access-key")
	t.Setenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY", "my-secret-access-key")
	t.Setenv("EDG_DATAVAL_CONFIG_DIR", "my-config-dir")
	t.Setenv("EDG_DATAVAL_ATTESTATION_PROVIDER_URL", "my-attestation-provider-url")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_AWSS3Bucket(t *testing.T) {
	t.Setenv("EDG_DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("EDG_DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("EDG_DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("EDG_DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("EDG_DATAVAL_AWS_S3_REGION", "ap-northeast-2")
	t.Setenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID", "my-access-key")
	t.Setenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY", "my-secret-access-key")
	t.Setenv("EDG_DATAVAL_CONFIG_DIR", "my-config-dir")
	t.Setenv("EDG_DATAVAL_ATTESTATION_PROVIDER_URL", "my-attestation-provider-url")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_AWSS3Region(t *testing.T) {
	t.Setenv("EDG_DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("EDG_DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("EDG_DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("EDG_DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("EDG_DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID", "my-access-key")
	t.Setenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY", "my-secret-access-key")
	t.Setenv("EDG_DATAVAL_CONFIG_DIR", "my-config-dir")
	t.Setenv("EDG_DATAVAL_ATTESTATION_PROVIDER_URL", "my-attestation-provider-url")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_CertificateStorePath(t *testing.T) {
	t.Setenv("EDG_DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("EDG_DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("EDG_DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("EDG_DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("EDG_DATAVAL_AWS_S3_REGION", "ap-northeast-2")
	t.Setenv("EDG_DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID", "my-access-key")
	t.Setenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY", "my-secret-access-key")
	t.Setenv("EDG_DATAVAL_ATTESTATION_PROVIDER_URL", "my-attestation-provider-url")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_MissingRequiredEnv_AttestationProviderURL(t *testing.T) {
	t.Setenv("EDG_DATAVAL_LOG_LEVEL", "debug")
	t.Setenv("EDG_DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("EDG_DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("EDG_DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("EDG_DATAVAL_AWS_S3_REGION", "ap-northeast-2")
	t.Setenv("EDG_DATAVAL_AWS_S3_BUCKET", "my-s3-bucket")
	t.Setenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID", "my-access-key")
	t.Setenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY", "my-secret-access-key")
	t.Setenv("EDG_DATAVAL_CONFIG_DIR", "my-config-dir")

	require.Panics(t, func() {
		config.MustLoad()
	})
}

func TestMustLoad_InvalidLogLevel(t *testing.T) {
	t.Setenv("EDG_DATAVAL_LOG_LEVEL", "hello")
	t.Setenv("EDG_DATAVAL_HTTP_LADDR", "0.0.0.0:8181")
	t.Setenv("EDG_DATAVAL_PANACEA_GRPC_ADDR", "0.0.0.0:9191")
	t.Setenv("EDG_DATAVAL_VALIDATOR_MNEMONIC", "Your MNEMONIC")
	t.Setenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID", "my-access-key")
	t.Setenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY", "my-secret-access-key")
	t.Setenv("EDG_DATAVAL_CONFIG_DIR", "my-config-dir")
	t.Setenv("EDG_DATAVAL_ATTESTATION_PROVIDER_URL", "my-attestation-provider-url")

	require.Panics(t, func() {
		config.MustLoad()
	})
}
