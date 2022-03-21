package tee_test

import (
	"crypto/ecdsa"
	"crypto/x509"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/panacea"
	"github.com/medibloc/panacea-data-market-validator/tee"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestCreateCertificate(T *testing.T) {
	mnemonic, err := crypto.GenerateMnemonic()
	require.NoError(T, err)

	acc, err := panacea.NewValidatorAccount(mnemonic)

	certBytes, err := tee.CreateCertificate(acc.GetEcdsaPrivKey())
	require.NoError(T, err)

	cert, err := x509.ParseCertificate(certBytes)
	require.NoError(T, err)

	err = cert.VerifyHostname("localhost")
	require.NoError(T, err)

	require.Equal(T, reflect.TypeOf(&ecdsa.PublicKey{}), reflect.TypeOf(cert.PublicKey))
	pubKey := cert.PublicKey.(*ecdsa.PublicKey)

	require.Equal(T, acc.GetEcdsaPubKey(), pubKey)
	require.Equal(T, x509.ECDSA, cert.PublicKeyAlgorithm)
	require.Equal(T, "localhost", cert.Subject.CommonName)
}