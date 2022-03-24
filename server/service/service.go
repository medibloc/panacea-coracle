package service

import (
	"fmt"

	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/panacea"
	"github.com/medibloc/panacea-data-market-validator/store"
	"github.com/medibloc/panacea-data-market-validator/tee"
)

type Service struct {
	Conf             *config.Config
	ValidatorAccount *panacea.ValidatorAccount
	Store            store.S3Store
	PanaceaClient    *panacea.GrpcClient
	TLSCert          *tee.TLSCertificate
}

func New(conf *config.Config) (*Service, error) {
	validatorAccount, err := panacea.NewValidatorAccount(conf.ValidatorMnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to load validator account: %w", err)
	}

	s3Store, err := store.NewS3Store(conf.AWSS3Bucket, conf.AWSS3Region, conf.AWSS3AccessKeyID, conf.AWSS3SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3Store: %w", err)
	}

	panaceaClient, err := panacea.NewGrpcClient(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create PanaceaGRPCClient: %w", err)
	}

	tlsCert, err := tee.CreateTLSCertificate()
	if err != nil {
		panaceaClient.Close()
		return nil, fmt.Errorf("failed to create TLS certificate: %w", err)
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
