package tee

import (
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTLSCertificate(t *testing.T) {
	tlsCert, err := CreateTLSCertificate()
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(tlsCert.Cert)
	require.NoError(t, err)

	require.Equal(t, "DataValidator", cert.Subject.CommonName)
	require.Equal(t, tlsCert.PrivKey.Public(), cert.PublicKey)
	require.Equal(t, x509.RSA, cert.PublicKeyAlgorithm)
}
