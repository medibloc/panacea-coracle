package codec_test

import (
	datadealtypes "github.com/medibloc/panacea-data-market-validator/types/datadeal"
	"testing"

	panaceatypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	"github.com/medibloc/panacea-data-market-validator/codec"
	"github.com/stretchr/testify/require"
)

func TestJsonMarshalAndUnMarshal(t *testing.T) {
	unsignedCertificate, err := datadealtypes.NewUnsignedDataValidationCertificate(
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
