package codec_test

import (
	panaceatypes "github.com/medibloc/panacea-core/v2/x/market/types"
	"github.com/medibloc/panacea-data-market-validator/codec"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJsonMarshalAndUnMarshal(t *testing.T) {
	unsignedCertificate, err := types.NewUnsignedDataValidationCertificate(
		"1",
		[]byte("dataHash"),
		[]byte("encryptedDataURL"),
		"requester_address",
		"dataValidatorAddress")

	require.NoError(t, err)

	signature := []byte("signature")

	cert := &panaceatypes.DataValidationCertificate{
		UnsignedCert: &unsignedCertificate,
		Signature:    signature,
	}

	json, err := codec.ProtoMarshalJSON(cert)
	require.NoError(t, err)

	resCert := &panaceatypes.DataValidationCertificate{}

	err = codec.ProtoUnmarshalJSON(json, resCert)
	require.NoError(t, err)
	require.Equal(t, cert, resCert)
}
