package tee

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/edgelesssys/ego/enclave"
	"math/big"
	"time"
)

// CreateTLSCertificate creates an x509 certificate and generate an rsa key.
func CreateTLSCertificate() ([]byte, *rsa.PrivateKey, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: "DataValidator"},
		NotAfter:     time.Now().AddDate(1, 0, 0),
		DNSNames:     []string{"localhost"},
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	return certBytes, priv, nil
}

// CreateAzureAttestationToken sends the x509 certificate and remote report to Azure to verify
// that it is working in a trusted environment and returns a JWT token.
func CreateAzureAttestationToken(cert []byte, attestationProviderURL string) (string, error) {
	return enclave.CreateAzureAttestationToken(cert, attestationProviderURL)
}
