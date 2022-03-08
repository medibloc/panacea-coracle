package server

import (
	"context"
	"fmt"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	panaceaapp "github.com/medibloc/panacea-core/v2/app"
	"github.com/medibloc/panacea-core/v2/app/params"
	markettypes "github.com/medibloc/panacea-core/v2/x/market/types"
	"github.com/medibloc/panacea-data-market-validator/types"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

type GrpcClient struct {
	conn           *grpc.ClientConn
	encodingConfig params.EncodingConfig
}

func NewGrpcClient(conn *grpc.ClientConn) (GrpcClient, error) {
	if conn == nil {
		return GrpcClient{}, types.ErrNoGrpcConnection
	}

	return GrpcClient{
			conn:           conn,
			encodingConfig: panaceaapp.MakeEncodingConfig(),
		},
		nil
}

// GetPubKey gets the public key from blockchain.
func (c GrpcClient) GetPubKey(panaceaAddr string) ([]byte, error) {
	client := authtypes.NewQueryClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Account(ctx, &authtypes.QueryAccountRequest{Address: panaceaAddr})
	if err != nil {
		return nil, fmt.Errorf("failed to get account info via grpc: %w", err)
	}

	var acc authtypes.AccountI
	if err := c.encodingConfig.InterfaceRegistry.UnpackAny(response.GetAccount(), &acc); err != nil {
		return nil, fmt.Errorf("failed to unpack account info: %w", err)
	}
	return acc.GetPubKey().Bytes(), nil
}

// GetDeal gets deal info from blockchain
func (c GrpcClient) GetDeal(id string) (markettypes.Deal, error) {
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
