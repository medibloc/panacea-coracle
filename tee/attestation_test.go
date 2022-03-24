package tee

import (
	"crypto/x509"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTLSCertificate(T *testing.T) {
	certBytes, priv, err := CreateTLSCertificate()
	require.NoError(T, err)

	cert, err := x509.ParseCertificate(certBytes)
	require.NoError(T, err)

	require.Equal(T, "DataValidator", cert.Subject.CommonName)
	require.Equal(T, priv.Public(), cert.PublicKey)
	require.Equal(T, x509.RSA, cert.PublicKeyAlgorithm)
}
