package panacea

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	types3 "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	tx2 "github.com/cosmos/cosmos-sdk/x/auth/tx"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

func TestTx(t *testing.T) {
	conf := config.MustLoad()

	account, err := NewValidatorAccount(conf.ValidatorMnemonic)
	require.NoError(t, err)

	privKey := secp256k1.PrivKey{
		Key: account.secp256k1PrivKey.Bytes(),
	}


	client, err := NewGrpcClient(&conf)
	require.NoError(t, err)

	interfaceRegistry := client.interfaceRegistry
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txConfig := tx2.NewTxConfig(marshaler, []signing.SignMode{signing.SignMode_SIGN_MODE_DIRECT})
	txBuilder := txConfig.NewTxBuilder()

	msgSend := types2.MsgSend{
		FromAddress: account.GetAddress(),
		ToAddress: account.GetAddress(),
		Amount: types3.NewCoins(types3.NewInt64Coin("umed", 1000000000)),
	}
	err = txBuilder.SetMsgs(&msgSend)
	fees := types3.NewCoins(types3.NewInt64Coin("umed", 1000000))

	require.NoError(t, err)
	txBuilder.SetGasLimit(200000)
	txBuilder.SetFeeAmount(fees)
	txBuilder.SetMemo("Memo")

	var sigsV2 []signing.SignatureV2
	sigV2 := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: 1,
	}
	sigsV2 = append(sigsV2, sigV2)

	err = txBuilder.SetSignatures(sigsV2...)
	require.NoError(t, err)

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	require.NoError(t, err)

	txJSON := string(txBytes)
	fmt.Println(txJSON)


	txClient := tx.NewServiceClient(client.conn)
	grpcRes, err := txClient.BroadcastTx(
		context.Background(),
		&tx.BroadcastTxRequest{
			Mode: tx.BroadcastMode_BROADCAST_MODE_BLOCK,
			TxBytes: txBytes,
		},
	)

	require.NoError(t, err)
	fmt.Println(grpcRes.TxResponse)
}
