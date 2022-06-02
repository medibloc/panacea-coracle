package tee

import (
	"crypto/rsa"
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTLSCertificate(t *testing.T) {
	tlsCert, err := CreateTLSCertificate()
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	require.NoError(t, err)

	rsaPrivKey, ok := tlsCert.PrivateKey.(*rsa.PrivateKey)
	require.True(t, ok)

	require.Equal(t, "Oracle", cert.Subject.CommonName)
	require.Equal(t, rsaPrivKey.Public(), cert.PublicKey)
	require.Equal(t, x509.RSA, cert.PublicKeyAlgorithm)
}
