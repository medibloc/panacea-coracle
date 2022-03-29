package panacea

import (
	"context"
	"fmt"
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

// MakeInterfaceRegistry
func makeInterfaceRegistry() sdk.InterfaceRegistry {
	interfaceRegistry := sdk.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	markettypes.RegisterInterfaces(interfaceRegistry)
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
	return acc.GetPubKey(), nil
}

func (c *GrpcClient) GetAccountNumber(panaceaAddr string) (uint64, error) {
	client := authtypes.NewQueryClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Account(ctx, &authtypes.QueryAccountRequest{Address: panaceaAddr})
	if err != nil {
		return 0, fmt.Errorf("failed to get account info via grpc: %w", err)
	}

	var acc authtypes.AccountI
	if err := c.interfaceRegistry.UnpackAny(response.GetAccount(), &acc); err != nil {
		return 0, fmt.Errorf("failed to unpack account info: %w", err)
	}
	return acc.GetAccountNumber(), nil
}

func (c *GrpcClient) GetSequence(panaceaAddr string) (uint64, error) {
	client := authtypes.NewQueryClient(c.conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Account(ctx, &authtypes.QueryAccountRequest{Address: panaceaAddr})
	if err != nil {
		return 0, fmt.Errorf("failed to get account info via grpc: %w", err)
	}

	var acc authtypes.AccountI
	if err := c.interfaceRegistry.UnpackAny(response.GetAccount(), &acc); err != nil {
		return 0, fmt.Errorf("failed to unpack account info: %w", err)
	}
	return acc.GetSequence(), nil
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
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	interfaceRegistry := makeInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txConfig := tx.NewTxConfig(marshaler, []signing.SignMode{signing.SignMode_SIGN_MODE_DIRECT})
	txBuilder := txConfig.NewTxBuilder()

	dataValidator := &datapooltypes.DataValidator{
		Address:  address,
		Endpoint: endpoint,
	}

	msgRegisterDataValidator := datapooltypes.NewMsgRegisterDataValidator(dataValidator)

	err := txBuilder.SetMsgs(msgRegisterDataValidator)
	if err != nil {
		return err
	}

	conf := config.MustLoad()
	account, err := NewValidatorAccount(conf.ValidatorMnemonic)
	if err != nil {
		return err
	}

	privKey := secp256k1.PrivKey{
		Key: account.secp256k1PrivKey.Bytes(),
	}

	sequence, err := c.GetSequence(address)
	if err != nil {
		return err
	}

	//TODO: Fee will be set in Config.toml in near future, now just hard-coded.
	fees := cosmostype.NewCoins(cosmostype.NewInt64Coin("umed", 1000000))
	txBuilder.SetFeeAmount(fees)
	txBuilder.SetGasLimit(200000)

	var sigsV2 []signing.SignatureV2
	sigV2 := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	sigsV2 = append(sigsV2, sigV2)

	err = txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil
	}

	accountNumber, err := c.GetAccountNumber(address)
	if err != nil {
		return err
	}

	sigsV2 = []signing.SignatureV2{}
	signerData := xauthsigning.SignerData{
		ChainID:       "panacea-3",
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}

	sigv2, err := txclient.SignWithPrivKey(signing.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder, &privKey, txConfig, sequence)
	if err != nil {
		return nil
	}
	sigsV2 = append(sigsV2, sigv2)

	err = txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	newTxClient := txtypes.NewServiceClient(c.conn)
	resp, err := newTxClient.BroadcastTx(
		context.Background(),
		&txtypes.BroadcastTxRequest{
			Mode:    txtypes.BroadcastMode_BROADCAST_MODE_BLOCK,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return nil
	}

	if resp.TxResponse.Code == 0 {
		log.Info("transaction successfully broadcast")
	} else {
		return fmt.Errorf(resp.TxResponse.RawLog)
	}

	return nil
}
