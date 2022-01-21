package server

import (
	"context"
	"fmt"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/medibloc/panacea-core/v2/app/params"
	markettypes "github.com/medibloc/panacea-core/v2/x/market/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"time"
)

type GrpcClient struct {
	addr           string
	encodingConfig params.EncodingConfig
}

func NewGrpcClient(grpcAddr string, encodingConfig params.EncodingConfig) *GrpcClient {
	return &GrpcClient{
		addr:           grpcAddr,
		encodingConfig: encodingConfig,
	}
}

// GetPubKey gets the public key from blockchain.
func (cli GrpcClient) GetPubKey(panaceaAddr string) ([]byte, error) {
	log.Infof("Dial to %s", cli.addr)
	conn, err := grpc.Dial(cli.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial grpc: %w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("failed to close grpc connection %v", err)
		}
	}()

	client := authtypes.NewQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Account(ctx, &authtypes.QueryAccountRequest{Address: panaceaAddr})
	if err != nil {
		return nil, fmt.Errorf("failed to get account info via grpc: %w", err)
	}

	var acc authtypes.AccountI
	if err := cli.encodingConfig.InterfaceRegistry.UnpackAny(response.GetAccount(), &acc); err != nil {
		return nil, fmt.Errorf("failed to unpack account info: %w", err)
	}
	return acc.GetPubKey().Bytes(), nil
}

func (cli GrpcClient) GetDeal(id string) (markettypes.Deal, error) {
	dealId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return markettypes.Deal{}, fmt.Errorf("failed to parse deal id: %w", err)
	}

	log.Infof("Dial to %s", cli.addr)
	conn, err := grpc.Dial(cli.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return markettypes.Deal{}, fmt.Errorf("failed to dial grpc: %w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("failed to close grpc connection %v", err)
		}
	}()

	client := markettypes.NewQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Deal(ctx, &markettypes.QueryDealRequest{DealId: dealId})
	if err != nil {
		return markettypes.Deal{}, fmt.Errorf("failed to get deal info: %w", err)
	}

	return *response.GetDeal(), nil
}
