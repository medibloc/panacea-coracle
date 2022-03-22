package tee

import (
	"crypto/x509"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestIsExistsCertificate(T *testing.T) {
	storePath := "./test/1/2"
	storeFullPath := filepath.Join(storePath, CertificateFilename)
	err := os.MkdirAll(storePath, 0755)
	require.NoError(T, err)

	defer func() {
		err := os.RemoveAll("./test")
		require.NoError(T, err)
	} ()

	err = ioutil.WriteFile(storeFullPath, []byte("test"), 0755)
	require.NoError(T, err)
}

func TestCreateCertificate(T *testing.T) {
	storePath := "./test/1/2"
	certBytes, priv, err := createCertificate(storePath)
	require.NoError(T, err)

	cert, err := x509.ParseCertificate(certBytes)
	require.NoError(T, err)

	require.Equal(T, "DataValidator", cert.Subject.CommonName)
	require.Equal(T, priv.Public(), cert.PublicKey)
	require.Equal(T, x509.RSA, cert.PublicKeyAlgorithm)
}