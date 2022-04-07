package service

import (
	"crypto/tls"
	"fmt"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/panacea"
	"github.com/medibloc/panacea-data-market-validator/store"
	"github.com/medibloc/panacea-data-market-validator/tee"
	"strings"
)

type Service struct {
	Conf             *config.Config
	ValidatorAccount *panacea.ValidatorAccount
	Store            store.S3Store
	PanaceaClient    *panacea.GrpcClient
	TLSCert          *tls.Certificate
}

func New(conf *config.Config) (*Service, error) {
	validatorAccount, err := panacea.NewValidatorAccount(conf.ValidatorMnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to load validator account: %w", err)
	}

	s3Store, err := store.NewS3Store(conf.AWSS3.Bucket, conf.AWSS3.Region, conf.AWSS3.AccessKeyID, conf.AWSS3.SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3Store: %w", err)
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

	return &Service{
		Conf:             conf,
		ValidatorAccount: validatorAccount,
		Store:            s3Store,
		PanaceaClient:    panaceaClient,
		TLSCert:          tlsCert,
	}, nil
}

func (svc *Service) Close() {
	svc.PanaceaClient.Close()
}
