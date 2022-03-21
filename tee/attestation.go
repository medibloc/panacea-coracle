package tee

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/edgelesssys/ego/enclave"
	"github.com/medibloc/panacea-data-market-validator/types"
	"math/big"
)

func CreateCertificate(priv *ecdsa.PrivateKey) ([]byte, error) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: "localhost"},
		DNSNames:     []string{"localhost"},
	}

	return x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
}

// CreateAzureAttestationToken If you call this on macOS, you will get the following error.
// SIGSYS: bad system call
// PC=0x407f2d0 m=0 sigcode=0
func CreateAzureAttestationToken(cert []byte) (string, error) {
	return enclave.CreateAzureAttestationToken(cert, types.AttestationProviderURL)
}