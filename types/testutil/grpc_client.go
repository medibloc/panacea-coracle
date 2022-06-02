package testutil

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	datadealtypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type MockGrpcClient struct {
	accountMap      map[string]authtypes.AccountI
	dealMap         map[string]datadealtypes.Deal
	poolMap         map[string]datapooltypes.Pool
	regValidatorMap map[string]datapooltypes.DataValidator
}

func NewMockGrpcClient(
	accounts []authtypes.AccountI,
	deals []datadealtypes.Deal,
	pools []datapooltypes.Pool,
	regValidators []datapooltypes.DataValidator,
) MockGrpcClient {
	cli := MockGrpcClient{
		accountMap:      make(map[string]authtypes.AccountI),
		dealMap:         make(map[string]datadealtypes.Deal),
		poolMap:         make(map[string]datapooltypes.Pool),
		regValidatorMap: make(map[string]datapooltypes.DataValidator),
	}
	if len(accounts) > 0 {
		for _, account := range accounts {
			cli.accountMap[GetAddress(account.GetPubKey())] = account
		}
	}
	if len(deals) > 0 {
		for _, deal := range deals {
			cli.dealMap[strconv.FormatUint(deal.DealId, 10)] = deal
		}
	}
	if len(pools) > 0 {
		for _, pool := range pools {
			cli.poolMap[strconv.FormatUint(pool.PoolId, 10)] = pool
		}
	}
	if len(regValidators) > 0 {
		for _, regVal := range regValidators {
			cli.regValidatorMap[regVal.GetAddress()] = regVal
		}
	}

	return cli
}

func (m MockGrpcClient) GetAccount(panaceaAddr string) (authtypes.AccountI, error) {
	acc, ok := m.accountMap[panaceaAddr]
	if !ok {
		return nil, errors.New("not found")
	}
	return acc, nil
}

func (m MockGrpcClient) GetDeal(id string) (datadealtypes.Deal, error) {
	deal, ok := m.dealMap[id]
	if !ok {
		return datadealtypes.Deal{}, errors.New("not found")
	}
	return deal, nil
}

func (m MockGrpcClient) GetRegisteredDataValidator(address string) (*datapooltypes.DataValidator, error) {
	regVal, ok := m.regValidatorMap[address]
	if !ok {
		return nil, errors.New("not found")
	}
	return &regVal, nil
}

func (m MockGrpcClient) GetPool(id string) (datapooltypes.Pool, error) {
	pool, ok := m.poolMap[id]
	if !ok {
		return datapooltypes.Pool{}, errors.New("not found")
	}
	return pool, nil
}

func (m MockGrpcClient) Close() error {
	log.Info("Call Close()")
	return nil
}

func (m MockGrpcClient) GetPubKey(panaceaAddr string) (types.PubKey, error) {
	addr, err := m.GetAccount(panaceaAddr)
	if err != nil {
		return nil, err
	}
	return addr.GetPubKey(), nil
}
