package config

import "fmt"

type Config struct {
	BaseConfig `mapstructure:",squash"`

	HTTP    HTTPConfig    `mapstructure:"http"`
	Panacea PanaceaConfig `mapstructure:"panacea"`
	AWSS3   AWSS3Config   `mapstructure:"aws-s3"`
	Enclave EnclaveConfig `mapstructure:"enclave"`
}

type BaseConfig struct {
	LogLevel              string `mapstructure:"log-level"`
	ValidatorMnemonic     string `mapstructure:"validator-mnemonic"`
	DataEncryptionKeyFile string `mapstructure:"data-encryption-key-file"`
}

type HTTPConfig struct {
	ListenAddr string `mapstructure:"laddr"`
}

type PanaceaConfig struct {
	GRPCAddr string `mapstructure:"grpc-addr"`
}

type AWSS3Config struct {
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key"`
}

type EnclaveConfig struct {
	Enable                  bool   `mapstructure:"enable"`
	AttestationProviderAddr string `mapstructure:"attestation-provider-addr"`
}

func DefaultConfig() *Config {
	return &Config{
		BaseConfig: BaseConfig{
			LogLevel:              "info",
			ValidatorMnemonic:     "",
			DataEncryptionKeyFile: "config/data_encryption_key.sealed",
		},
		HTTP: HTTPConfig{
			ListenAddr: "0.0.0.0:8080",
		},
		Panacea: PanaceaConfig{
			GRPCAddr: "127.0.0.1:9090",
		},
		AWSS3: AWSS3Config{
			Region:          "",
			Bucket:          "",
			AccessKeyID:     "",
			SecretAccessKey: "",
		},
		Enclave: EnclaveConfig{
			Enable:                  true,
			AttestationProviderAddr: "127.0.0.1:9999",
		},
	}
}

func (c *Config) validate() error {
	if c.Enclave.Enable {
		if c.Enclave.AttestationProviderAddr == "" {
			return fmt.Errorf("attestation-provider-addr should be specified if enclave is enabled")
		}
	}

	//TODO: validate other configs

	return nil
}
