package panacea

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	txclient "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	cosmostype "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
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
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
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

func (c *GrpcClient) GetChainId() (string, error) {
	client := tmservice.NewServiceClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetLatestBlock(ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return "", err
	}

	return response.GetBlock().GetHeader().ChainID, nil
}

// RegisterDataValidator registers data validator on blockchain
func (c *GrpcClient) RegisterDataValidator(endpoint string, validatorAcc *ValidatorAccount) error {
	interfaceRegistry := c.interfaceRegistry
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txConfig := tx.NewTxConfig(marshaler, []signing.SignMode{signing.SignMode_SIGN_MODE_DIRECT})
	txBuilder := txConfig.NewTxBuilder()

	address := validatorAcc.GetAddress()
	dataValidator := &datapooltypes.DataValidator{
		Address:  address,
		Endpoint: endpoint,
	}

	msgRegisterDataValidator := datapooltypes.NewMsgRegisterDataValidator(dataValidator)

	err := txBuilder.SetMsgs(msgRegisterDataValidator)
	if err != nil {
		return err
	}

	privKey := secp256k1.PrivKey{
		Key: validatorAcc.secp256k1PrivKey.Bytes(),
	}

	account, err := c.GetAccount(address)
	if err != nil {
		return err
	}

	sequence := account.GetSequence()

	//TODO: Fee will be set in Config.toml in near future, now just hard-coded.
	fees := cosmostype.NewCoins(cosmostype.NewInt64Coin("umed", 1000000))
	txBuilder.SetFeeAmount(fees)
	txBuilder.SetGasLimit(200000)

	sigV2 := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}

	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return nil
	}

	//TODO: ChainID will be set in Config.toml in near future, it just hard-coded.
	chainId, err := c.GetChainId()
	if err != nil {
		return err
	}

	accountNumber := account.GetAccountNumber()

	signerData := xauthsigning.SignerData{
		ChainID:       chainId,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}

	sigv2, err := txclient.SignWithPrivKey(signing.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder, &privKey, txConfig, sequence)
	if err != nil {
		return nil
	}

	err = txBuilder.SetSignatures(sigv2)
	if err != nil {
		return nil
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	newTxClient := txtypes.NewServiceClient(c.conn)
	resp, err := newTxClient.BroadcastTx(
		ctx,
		&txtypes.BroadcastTxRequest{
			Mode:    txtypes.BroadcastMode_BROADCAST_MODE_BLOCK,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return nil
	}

	if resp.TxResponse.Code == 0 {
		log.Info("register data validator success")
	} else {
		return fmt.Errorf(resp.TxResponse.RawLog)
	}

	return nil
}
