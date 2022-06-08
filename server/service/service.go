package service

import (
	"crypto/tls"
	"fmt"
	"github.com/edgelesssys/ego/ecrypto"
	oracletypes "github.com/medibloc/panacea-core/v2/x/oracle/types"

	"github.com/medibloc/panacea-oracle/cache"
	"github.com/medibloc/panacea-oracle/config"
	"github.com/medibloc/panacea-oracle/crypto"
	"github.com/medibloc/panacea-oracle/panacea"
	"github.com/medibloc/panacea-oracle/store"
	"github.com/medibloc/panacea-oracle/tee"
	tos "github.com/tendermint/tendermint/libs/os"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	Conf          *config.Config
	OracleAccount *panacea.OracleAccount
	Store         store.Storage
	PanaceaClient panacea.GrpcClientI
	TLSCert       *tls.Certificate
	DataEncKey    []byte
	Cache         *cache.AuthenticationCache
}

func New(conf *config.Config) (*Service, error) {
	oracleAccount, err := panacea.NewOracleAccount(conf.OracleMnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to load oracle account: %w", err)
	}

	s3Store, err := store.NewS3Store(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWSS3Storage: %w", err)
	}

	panaceaClient, err := panacea.NewGrpcClient(conf)
	if err != nil {

		return nil, fmt.Errorf("failed to create PanaceaGRPCClient: %w", err)
	}

	_, err = panaceaClient.GetRegisteredOracle(oracleAccount.GetAddress())
	if err != nil {
		if strings.HasSuffix(err.Error(), oracletypes.ErrOracleNotFound.Error()) {
			return nil, fmt.Errorf("this oracle is not registered in Panacea yet")
		}
		return nil, err
	}

	var tlsCert *tls.Certificate
	if conf.Enclave.Enable {
		tlsCert, err = tee.CreateTLSCertificate()
		if err != nil {
			panaceaClient.Close()
			return nil, fmt.Errorf("failed to create TLS certificate: %w", err)
		}
	}

	key, err := generateDataEncryptionKeyFile(conf, err)
	if err != nil {
		panaceaClient.Close()
		return nil, err
	}

	authenticateCache := cache.NewAuthenticationCache(conf)

	return &Service{
		Conf:          conf,
		OracleAccount: oracleAccount,
		Store:         s3Store,
		PanaceaClient: panaceaClient,
		TLSCert:       tlsCert,
		DataEncKey:    key,
		Cache:         authenticateCache,
	}, nil
}

func generateDataEncryptionKeyFile(conf *config.Config, err error) ([]byte, error) {
	var key []byte
	dataEncryptionKeyFile := conf.DataEncryptionKeyFile
	isEnclave := conf.Enclave.Enable
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fileFullPath := filepath.Join(userHomeDir, dataEncryptionKeyFile)

	if tos.FileExists(fileFullPath) {
		file, err := tos.ReadFile(fileFullPath)
		if err != nil {
			return nil, err
		}

		key, err = unSeal(isEnclave, file)
		if err != nil {
			return nil, err
		}
	} else {
		key, err = crypto.GenerateRandomKey(32)
		if err != nil {
			return nil, err
		}

		sealed, err := seal(isEnclave, key)
		if err != nil {
			return nil, err
		}

		var sealedSavedDir strings.Builder

		dir, file := filepath.Split(dataEncryptionKeyFile)

		// ex) .oracle/config/data_encryption_file.sealed
		// sealedSavedDir = $HOME/.oracle/config/, file = data_encryption_file.sealed
		sealedSavedDir.WriteString(userHomeDir)
		sealedSavedDir.WriteString("/")
		sealedSavedDir.WriteString(dir)

		err = tos.EnsureDir(sealedSavedDir.String(), 0755)
		if err != nil {
			return nil, err
		}

		// sealedSavedDir = $HOME/.oracle/config/data_encryption_file.sealed
		sealedSavedDir.WriteString("/")
		sealedSavedDir.WriteString(file)

		err = tos.WriteFile(sealedSavedDir.String(), sealed, 0755)
		if err != nil {
			return nil, err
		}
	}
	return key, nil
}

func unSeal(isEnclave bool, key []byte) ([]byte, error) {
	if isEnclave {
		return ecrypto.Unseal(key, nil)
	} else {
		return key, nil
	}
}

func seal(isEnclave bool, key []byte) ([]byte, error) {
	if isEnclave {
		return ecrypto.SealWithProductKey(key, nil)
	} else {
		return key, nil
	}
}

func (svc *Service) Close() {
	svc.PanaceaClient.Close()
}
