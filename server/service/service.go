package service

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/edgelesssys/ego/ecrypto"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/panacea"
	"github.com/medibloc/panacea-data-market-validator/store"
	"github.com/medibloc/panacea-data-market-validator/tee"
	tos "github.com/tendermint/tendermint/libs/os"
)

type Service struct {
	Conf             *config.Config
	ValidatorAccount *panacea.ValidatorAccount
	Store            store.Storage
	PanaceaClient    *panacea.GrpcClient
	TLSCert          *tls.Certificate
	DataEncKey       []byte
}

func New(conf *config.Config) (*Service, error) {
	validatorAccount, err := panacea.NewValidatorAccount(conf.ValidatorMnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to load validator account: %w", err)
	}

	s3Store, err := store.NewS3Store(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWSS3Storage: %w", err)
	}

	panaceaClient, err := panacea.NewGrpcClient(conf)
	if err != nil {

		return nil, fmt.Errorf("failed to create PanaceaGRPCClient: %w", err)
	}

	_, err = panaceaClient.GetRegisteredDataValidator(validatorAccount.GetAddress())
	if err != nil {
		if strings.HasSuffix(err.Error(), datapooltypes.ErrDataValidatorNotFound.Error()) {
			return nil, fmt.Errorf("this data validator is not registered in Panacea yet")
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

	key, err := generateDataEncryptionKeyFile(conf.DataEncryptionKeyFile, err)
	if err != nil {
		panaceaClient.Close()
		return nil, err
	}

	return &Service{
		Conf:             conf,
		ValidatorAccount: validatorAccount,
		Store:            s3Store,
		PanaceaClient:    panaceaClient,
		TLSCert:          tlsCert,
		DataEncKey:       key,
	}, nil
}

func generateDataEncryptionKeyFile(dataEncryptionKeyFile string, err error) ([]byte, error) {
	var key []byte

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

		key, err = ecrypto.Unseal(file, nil)
		if err != nil {
			return nil, err
		}
	} else {
		key, err = crypto.GenerateRandomKey(32)
		if err != nil {
			return nil, err
		}

		sealed, err := ecrypto.SealWithProductKey(key, nil)
		if err != nil {
			return nil, err
		}

		var sealedSavedDir strings.Builder

		dir, file := filepath.Split(dataEncryptionKeyFile)

		// ex) .dataval/config/data_encryption_file.sealed
		// sealedSavedDir = $HOME/.dataval/config/, file = data_encryption_file.sealed
		sealedSavedDir.WriteString(userHomeDir)
		sealedSavedDir.WriteString("/")
		sealedSavedDir.WriteString(dir)

		err = tos.EnsureDir(sealedSavedDir.String(), 0755)
		if err != nil {
			return nil, err
		}

		// sealedSavedDir = $HOME/.dataval/config/data_encryption_file.sealed
		sealedSavedDir.WriteString("/")
		sealedSavedDir.WriteString(file)

		err = tos.WriteFile(sealedSavedDir.String(), sealed, 0755)
		if err != nil {
			return nil, err
		}
	}
	return key, nil
}

func (svc *Service) Close() {
	svc.PanaceaClient.Close()
}
