package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const DefaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

###############################################################################
###                           Base Configuration                            ###
###############################################################################

log-level = "{{ .BaseConfig.LogLevel }}"
validator-mnemonic = "{{ .BaseConfig.ValidatorMnemonic }}"
data-encryption-key-file = "{{ .BaseConfig.DataEncryptionKeyFile }}"

###############################################################################
###                         HTTP Configuration                              ###
###############################################################################

[http]

laddr = "{{ .HTTP.ListenAddr }}"
endpoint = "{{ .HTTP.Endpoint }}"

###############################################################################
###                         Panacea Configuration                           ###
###############################################################################

[panacea]

grpc-addr = "{{ .Panacea.GRPCAddr }}"

###############################################################################
###                         AWS S3 Configuration                            ###
###############################################################################

[aws-s3]

region = "{{ .AWSS3.Region }}"
bucket = "{{ .AWSS3.Bucket }}"
access-key-id = "{{ .AWSS3.AccessKeyID }}"
secret-access-key = "{{ .AWSS3.SecretAccessKey }}"

###############################################################################
###                         Enclave Configuration                           ###
###############################################################################

[enclave]

enable = {{ .Enclave.Enable }}
attestation-provider-addr = "{{ .Enclave.AttestationProviderAddr }}"

###############################################################################
###                  AuthenticationConfig Configuration                     ###
###############################################################################

[authentication]

expiration = "{{ .Authentication.Expiration }}"
size = {{ .Authentication.Size }}

`

const (
	configFileName = "config"
	configFileExt  = "toml"
)

var configTemplate *template.Template

// init is run automatically when the package is loaded.
func init() {
	tmpl := template.New("configFileTemplate")

	var err error
	if configTemplate, err = tmpl.Parse(DefaultConfigTemplate); err != nil {
		log.Panic(err)
	}
}

func WriteConfigTOML(path string, config *Config) error {
	var buffer bytes.Buffer
	if err := configTemplate.Execute(&buffer, config); err != nil {
		return fmt.Errorf("failed to populate config template: %w", err)
	}

	return ioutil.WriteFile(path, buffer.Bytes(), 0600)
}

func ReadConfigTOML(path string) (*Config, error) {
	fileExt := filepath.Ext(path)

	v := viper.New()
	v.AddConfigPath(filepath.Dir(path))
	v.SetConfigName(strings.TrimSuffix(filepath.Base(path), fileExt))
	v.SetConfigType(fileExt[1:]) // excluding the dot

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var conf Config
	if err := v.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &conf, nil
}
