package testutil

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	datadealtypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type MockGrpcClient struct {
	accountMap         map[string]authtypes.AccountI
	dealMap            map[string]datadealtypes.Deal
	poolMap            map[string]datapooltypes.Pool
	regOracleMap       map[string]datapooltypes.Oracle
	dpRedeemHistoryMap map[string]datapooltypes.DataPassRedeemHistory
	dataCertsMap       map[string][]datapooltypes.DataCert
}

func NewMockGrpcClient(
	accounts []authtypes.AccountI,
	deals []datadealtypes.Deal,
	pools []datapooltypes.Pool,
	regOracles []datapooltypes.Oracle,
	dpRedeemHistories []datapooltypes.DataPassRedeemHistory,
	dataCerts []datapooltypes.DataCert,
) MockGrpcClient {
	cli := MockGrpcClient{
		accountMap:         make(map[string]authtypes.AccountI),
		dealMap:            make(map[string]datadealtypes.Deal),
		poolMap:            make(map[string]datapooltypes.Pool),
		regOracleMap:       make(map[string]datapooltypes.Oracle),
		dpRedeemHistoryMap: make(map[string]datapooltypes.DataPassRedeemHistory),
		dataCertsMap:       make(map[string][]datapooltypes.DataCert),
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
	if len(regOracles) > 0 {
		for _, regVal := range regOracles {
			cli.regOracleMap[regVal.GetAddress()] = regVal
		}
	}

	if len(dpRedeemHistories) > 0 {
		for _, dpRedeemHistory := range dpRedeemHistories {
			redeemer := dpRedeemHistory.Redeemer
			poolID := strconv.FormatUint(dpRedeemHistory.PoolId, 10)
			cli.dpRedeemHistoryMap[combinesKey(redeemer, poolID)] = dpRedeemHistory
		}
	}

	if len(dataCerts) > 0 {
		for _, cert := range dataCerts {
			poolID := strconv.FormatUint(cert.UnsignedCert.PoolId, 10)
			round := strconv.FormatUint(cert.UnsignedCert.Round, 10)
			certs, ok := cli.dataCertsMap[combinesKey(poolID, round)]
			if !ok {
				certs = make([]datapooltypes.DataCert, 0)
			}
			cli.dataCertsMap[combinesKey(poolID, round)] = append(certs, cert)
		}
	}

	return cli
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

func (m MockGrpcClient) GetRegisteredOracle(address string) (*datapooltypes.Oracle, error) {
	regVal, ok := m.regOracleMap[address]
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

func (m MockGrpcClient) GetDataPassRedeemHistory(redeemer string, poolID uint64) (datapooltypes.DataPassRedeemHistory, error) {
	dpRedeemHistory, ok := m.dpRedeemHistoryMap[combinesKey(redeemer, strconv.FormatUint(poolID, 10))]
	if !ok {
		return datapooltypes.DataPassRedeemHistory{}, errors.New("not found")
	}
	return dpRedeemHistory, nil
}

func (m MockGrpcClient) GetDataCerts(poolID, round uint64) ([]datapooltypes.DataCert, error) {
	dataCerts, ok := m.dataCertsMap[combinesKey(strconv.FormatUint(poolID, 10), strconv.FormatUint(round, 10))]
	if !ok {
		return []datapooltypes.DataCert{}, errors.New("not found")
	}
	return dataCerts, nil
}

func combinesKey(key ...string) string {
	return strings.Join(key, "|")
}
