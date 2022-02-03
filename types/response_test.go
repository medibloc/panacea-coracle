package types_test

import (
	"encoding/base64"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDataValidationCertificateResponse(t *testing.T) {
	dealId := "1"
	dataHash := []byte("dataHash")
	encryptedDataUrl := []byte("encryptedDataUrl")
	requesterAddress := "requesterAddress"
	dataValidatorAddress := "dataValidatorAddress"

	unsignedCert, err := types.NewUnsignedDataValidationCertificate(dealId, dataHash, encryptedDataUrl, requesterAddress, dataValidatorAddress)
	require.NoError(t, err)

	signature := []byte("signature")
	dataValCertResp := types.NewDataValidationCertificateResponse(unsignedCert, signature)

	unsignedCertResp := dataValCertResp.UnsignedCert
	require.Equal(t, unsignedCert.DealId, unsignedCertResp.DealId)
	require.Equal(t, base64.StdEncoding.EncodeToString(unsignedCert.DataHash), unsignedCertResp.DataHashBase64)
	require.Equal(t, base64.StdEncoding.EncodeToString(unsignedCert.EncryptedDataUrl), unsignedCertResp.EncryptedDataUrlBase64)
	require.Equal(t, unsignedCert.RequesterAddress, unsignedCertResp.RequesterAddress)
	require.Equal(t, unsignedCert.DataValidatorAddress, unsignedCertResp.DataValidatorAddress)
	require.Equal(t, base64.StdEncoding.EncodeToString(signature), dataValCertResp.SignatureBase64)
}