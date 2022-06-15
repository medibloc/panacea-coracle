package panacea

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	datadealtypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	oracletypes "github.com/medibloc/panacea-core/v2/x/oracle/types"
	"github.com/medibloc/panacea-oracle/config"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type GrpcClientI interface {
	Close() error

	GetPubKey(panaceaAddr string) (types.PubKey, error)

	GetAccount(panaceaAddr string) (authtypes.AccountI, error)

	GetDeal(id string) (datadealtypes.Deal, error)

	GetPool(id string) (datapooltypes.Pool, error)

	GetDataPassRedeemHistory(redeemer string, poolID uint64) (datapooltypes.DataPassRedeemHistory, error)

	GetDataCerts(poolID, round uint64) ([]datapooltypes.DataCert, error)

	GetRegisteredOracle(address string) (*oracletypes.Oracle, error)
}

var _ GrpcClientI = (*GrpcClient)(nil)

const pageLimit = 30

type GrpcClient struct {
	conn              *grpc.ClientConn
	interfaceRegistry sdk.InterfaceRegistry
}

func NewGrpcClient(conf *config.Config) (GrpcClientI, error) {
	log.Infof("dialing to Panacea gRPC endpoint: %s", conf.Panacea.GRPCAddr)

	var conn *grpc.ClientConn
	var err error

	u, err := url.Parse(conf.Panacea.GRPCAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse grpc addr: %w", err)
	}

	prefixLen := len(u.Scheme + "://")
	addrBody := conf.Panacea.GRPCAddr[prefixLen:]

	var creds credentials.TransportCredentials

	if u.Scheme == "tcp" || u.Scheme == "http" {
		creds = insecure.NewCredentials()
	} else if u.Scheme == "https" {
		creds = credentials.NewClientTLSFromCert(nil, "")
	} else if u.Scheme == "" {
		return nil, fmt.Errorf("empty panacea grpc address")
	} else {
		return nil, fmt.Errorf("invalid panacea grpc addr: %s", conf.Panacea.GRPCAddr)
	}

	conn, err = grpc.Dial(addrBody, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Panacea: %w", err)
	}

	return &GrpcClient{
		conn:              conn,
		interfaceRegistry: makeInterfaceRegistry(),
	}, nil
}

// makeInterfaceRegistry
func makeInterfaceRegistry() sdk.InterfaceRegistry {
	interfaceRegistry := sdk.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	datadealtypes.RegisterInterfaces(interfaceRegistry)
	datapooltypes.RegisterInterfaces(interfaceRegistry)
	oracletypes.RegisterInterfaces(interfaceRegistry)
	return interfaceRegistry
}

func (c *GrpcClient) Close() error {
	log.Info("closing Panacea gRPC connection")
	return c.conn.Close()
}

// GetPubKey gets the public key from blockchain.
func (c *GrpcClient) GetPubKey(panaceaAddr string) (types.PubKey, error) {
	acc, err := c.GetAccount(panaceaAddr)
	if err != nil {
		return nil, err
	}

	return acc.GetPubKey(), nil
}

func (c *GrpcClient) GetAccount(panaceaAddr string) (authtypes.AccountI, error) {
	client := authtypes.NewQueryClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Account(ctx, &authtypes.QueryAccountRequest{Address: panaceaAddr})
	if err != nil {
		return nil, fmt.Errorf("failed to get account info via grpc: %w", err)
	}

	var acc authtypes.AccountI
	if err := c.interfaceRegistry.UnpackAny(response.GetAccount(), &acc); err != nil {
		return nil, fmt.Errorf("failed to unpack account info: %w", err)
	}
	return acc, nil
}

// GetDeal gets deal info from blockchain
func (c *GrpcClient) GetDeal(id string) (datadealtypes.Deal, error) {
	dealId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return datadealtypes.Deal{}, fmt.Errorf("failed to parse deal id: %w", err)
	}

	client := datadealtypes.NewQueryClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Deal(ctx, &datadealtypes.QueryDealRequest{DealId: dealId})
	if err != nil {
		return datadealtypes.Deal{}, fmt.Errorf("failed to get deal info: %w", err)
	}

	return *response.GetDeal(), nil
}

// GetRegisteredOracle gets registered oracle
func (c *GrpcClient) GetRegisteredOracle(address string) (*oracletypes.Oracle, error) {
	client := oracletypes.NewQueryClient(c.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.Oracle(ctx, &oracletypes.QueryOracleRequest{Address: address})
	if err != nil {
		return nil, fmt.Errorf("failed to get oracle info: %w", err)
	}

	return res.GetOracle(), nil
}

func (c *GrpcClient) GetPool(id string) (datapooltypes.Pool, error) {
	poolId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return datapooltypes.Pool{}, fmt.Errorf("failed to parse pool id: %w", err)
	}

	client := datapooltypes.NewQueryClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Pool(ctx, &datapooltypes.QueryPoolRequest{
		PoolId: poolId,
	})
	if err != nil {
		return datapooltypes.Pool{}, fmt.Errorf("failed to get pool info: %w", err)
	}

	return *response.GetPool(), nil

}

func (c *GrpcClient) GetDataPassRedeemHistory(redeemer string, poolID uint64) (datapooltypes.DataPassRedeemHistory, error) {
	client := datapooltypes.NewQueryClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.DataPassRedeemHistory(ctx, &datapooltypes.QueryDataPassRedeemHistoryRequest{
		Redeemer: redeemer,
		PoolId:   poolID,
	})
	if err != nil {
		return datapooltypes.DataPassRedeemHistory{}, err
	}

	return response.GetDataPassRedeemHistories(), nil
}

func (c *GrpcClient) GetDataCerts(poolID, round uint64) ([]datapooltypes.DataCert, error) {
	client := datapooltypes.NewQueryClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var certs []datapooltypes.DataCert

	pageReq := &datapooltypes.QueryDataCertsRequest{
		PoolId: poolID,
		Round:  round,
		Pagination: &query.PageRequest{
			Key:   nil,
			Limit: pageLimit,
		},
	}

	for {
		response, err := client.DataCerts(ctx, pageReq)
		if err != nil {
			return nil, err
		}

		certs = append(certs, response.GetDataCerts()...)

		if response.Pagination.NextKey == nil {
			break
		}

		pageReq.Pagination.Key = response.Pagination.NextKey
	}

	return certs, nil
}
