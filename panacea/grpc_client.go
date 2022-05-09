package panacea

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/types"

	sdk "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/std"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	datadealtypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
	datavalcodec "github.com/medibloc/panacea-data-market-validator/codec"
	"github.com/medibloc/panacea-data-market-validator/config"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	dataValPrivKey = secp256k1.GenPrivKey()
	dataValPubKey  = dataValPrivKey.PubKey()
	dataVal1       = cosmossdk.AccAddress(dataValPubKey.Address())

	requesterPrivKey = secp256k1.GenPrivKey()
	requesterPubKey  = requesterPrivKey.PubKey()
	requesterAddr    = cosmossdk.AccAddress(requesterPubKey.Address())
)

type GrpcClient struct {
	conn              *grpc.ClientConn
	interfaceRegistry sdk.InterfaceRegistry
}

func NewGrpcClient(conf *config.Config) (*GrpcClient, error) {
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

func (c GrpcClient) GetDataCertsByRound(poolID, round uint64) ([]datapooltypes.DataValidationCertificate, error) {
	certs, err := MakeTestDataCerts(poolID, round)
	if err != nil {
		return nil, err
	}

	//if round == uint64(2) {
	//	return nil, errors.New("failed to get data certs")
	//}

	return certs, nil
}

// MakeTestDataCerts returns list of 3 data certs
func MakeTestDataCerts(poolID, round uint64) ([]datapooltypes.DataValidationCertificate, error) {
	var res []datapooltypes.DataValidationCertificate

	for i := uint64(1); i < 10; i++ {

		//dataHash := crypto.Hash(data)
		unsignedCert := &datapooltypes.UnsignedDataValidationCertificate{
			PoolId:        poolID,
			Round:         round,
			DataHash:      []byte("data-" + strconv.FormatUint(poolID, 10) + "-" + strconv.FormatUint(round, 10) + "-" + strconv.FormatUint(i, 10)),
			DataValidator: dataVal1.String(),
			Requester:     requesterAddr.String(),
		}

		json, err := datavalcodec.ProtoMarshalJSON(unsignedCert)
		if err != nil {
			return nil, err
		}

		sign, err := dataValPrivKey.Sign(json)
		if err != nil {
			return nil, err
		}

		cert := datapooltypes.DataValidationCertificate{
			UnsignedCert: unsignedCert,
			Signature:    sign,
		}

		res = append(res, cert)
	}

	return res, nil
}
