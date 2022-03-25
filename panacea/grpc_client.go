package panacea

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	markettypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	conn              *grpc.ClientConn
	interfaceRegistry sdk.InterfaceRegistry
}

func NewGrpcClient(conf *config.Config) (*GrpcClient, error) {
	log.Infof("dialing to Panacea gRPC endpoint: %s", conf.PanaceaGrpcAddress)
	conn, err := grpc.Dial(conf.PanaceaGrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	markettypes.RegisterInterfaces(interfaceRegistry)
	return interfaceRegistry
}

func (c *GrpcClient) Close() {
	log.Info("closing Panacea gRPC connection")
	c.conn.Close()
}

// GetPubKey gets the public key from blockchain.
func (c *GrpcClient) GetPubKey(panaceaAddr string) ([]byte, error) {
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
	return acc.GetPubKey().Bytes(), nil
}

// GetDeal gets deal info from blockchain
func (c *GrpcClient) GetDeal(id string) (markettypes.Deal, error) {
	dealId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return markettypes.Deal{}, fmt.Errorf("failed to parse deal id: %w", err)
	}

	client := markettypes.NewQueryClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Deal(ctx, &markettypes.QueryDealRequest{DealId: dealId})
	if err != nil {
		return markettypes.Deal{}, fmt.Errorf("failed to get deal info: %w", err)
	}

	return *response.GetDeal(), nil
}

// RegisterDataValidator registers data validator on blockchain
func (c *GrpcClient) RegisterDataValidator(address, endpoint string) error {
	client := datapooltypes.NewMsgClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dataValidator := datapooltypes.DataValidator{
		Address:  address,
		Endpoint: endpoint,
	}

	_, err := client.RegisterDataValidator(ctx, &datapooltypes.MsgRegisterDataValidator{ValidatorDetail: &dataValidator})
	if err != nil {
		return err
	}

	return nil
}
