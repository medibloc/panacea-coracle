package config

import (
	"fmt"
	"time"
)

type Config struct {
	BaseConfig `mapstructure:",squash"`

	HTTP         HTTPConfig         `mapstructure:"http"`
	Panacea      PanaceaConfig      `mapstructure:"panacea"`
	AWSS3        AWSS3Config        `mapstructure:"aws-s3"`
	Enclave      EnclaveConfig      `mapstructure:"enclave"`
	Authenticate AuthenticateConfig `mapsutrcture:"authenticate"`
}

type BaseConfig struct {
	LogLevel              string `mapstructure:"log-level"`
	ValidatorMnemonic     string `mapstructure:"validator-mnemonic"`
	DataEncryptionKeyFile string `mapstructure:"data-encryption-key-file"`
}

type HTTPConfig struct {
	ListenAddr string `mapstructure:"laddr"`
	Endpoint   string `mapstructure:"endpoint"`
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

type AuthenticateConfig struct {
	Expiration time.Duration `mapstruct:"expiration"`
	Size       int           `mapstruct:"size"`
}

func DefaultConfig() *Config {
	return &Config{
		BaseConfig: BaseConfig{
			LogLevel:              "info",
			ValidatorMnemonic:     "",
			DataEncryptionKeyFile: ".dataval/config/data_encryption_key.sealed",
		},
		HTTP: HTTPConfig{
			ListenAddr: "0.0.0.0:8080",
			Endpoint:   "",
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
		Authenticate: AuthenticateConfig{
			Expiration: 5 * time.Second,
			Size:       50000,
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
