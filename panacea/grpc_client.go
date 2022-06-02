package panacea

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	datadealtypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientI interface {
	GetPubKey(panaceaAddr string) (types.PubKey, error)

	GetAccount(panaceaAddr string) (authtypes.AccountI, error)

	GetDeal(id string) (datadealtypes.Deal, error)

	GetRegisteredDataValidator(address string) (*datapooltypes.DataValidator, error)

	GetPool(id string) (datapooltypes.Pool, error)

	Close() error
}

var _ GrpcClientI = (*GrpcClient)(nil)

type GrpcClient struct {
	conn              *grpc.ClientConn
	interfaceRegistry sdk.InterfaceRegistry
}

func NewGrpcClient(conf *config.Config) (GrpcClientI, error) {
	log.Infof("dialing to Panacea gRPC endpoint: %s", conf.Panacea.GRPCAddr)
	conn, err := grpc.Dial(conf.Panacea.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

// GetRegisteredDataValidator gets registered data validator
func (c *GrpcClient) GetRegisteredDataValidator(address string) (*datapooltypes.DataValidator, error) {
	client := datapooltypes.NewQueryClient(c.conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.DataValidator(ctx, &datapooltypes.QueryDataValidatorRequest{Address: address})
	if err != nil {
		return nil, fmt.Errorf("failed to get data validator info: %w", err)
	}

	return res.GetDataValidator(), nil
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
