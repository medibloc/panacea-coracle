package tee

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/edgelesssys/ego/ecrypto"
	"github.com/edgelesssys/ego/enclave"
	"github.com/tendermint/tendermint/libs/json"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	CertificateKey      = "certificate"
	CertificateFilename = "certificate_seal"
)

type SealCertificate struct {
	Cert    []byte
	PrivKey *rsa.PrivateKey
}

func GetCertificate(storePath string) ([]byte, *rsa.PrivateKey, error) {
	storeFullPath := filepath.Join(storePath, CertificateFilename)

	if _, err := os.Stat(storeFullPath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
	}

	sealedBody, err := ioutil.ReadFile(storeFullPath)

	if err != nil {
		return nil, nil, err
	}

	bytes, err := ecrypto.Unseal(sealedBody, []byte(CertificateKey))
	if err != nil {
		return nil, nil, err
	}

	cert := &SealCertificate{}
	err = json.Unmarshal(bytes, cert)
	if err != nil {
		return nil, nil, err
	}
	return cert.Cert, cert.PrivKey, nil
}

// CreateCertificate If there is a sealed certificate in the received path, it responds after parsing.
// However, if the certificate does not exist, it responds by creating a new one. It also responds with RSA PrivateKey.
func CreateCertificate(storePath string) ([]byte, *rsa.PrivateKey, error) {
	certBytes, priv, err := createCertificate()
	if err != nil {
		return nil, nil, err
	}

	err = sealAndStore(certBytes, priv, storePath)
	if err != nil {
		return nil, nil, err
	}

	return certBytes, priv, nil
}

// createCertificate Create an x509 certificate and generate an rsa key.
func createCertificate() ([]byte, *rsa.PrivateKey, error) {
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

// sealAndStore Seal the certificate and privKey and store it in a specific path.
//Unsealing is only possible in this enclave.
func sealAndStore(certBytes []byte, priv *rsa.PrivateKey, storePath string) error {
	cert := SealCertificate{
		Cert:    certBytes,
		PrivKey: priv,
	}
	jsonBytes, err := json.Marshal(cert)
	if err != nil {
		return err
	}

	sealedBody, err := ecrypto.SealWithProductKey(jsonBytes, []byte(CertificateKey))
	if err != nil {
		return err
	}

	err = os.MkdirAll(storePath, 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(storePath, CertificateFilename), sealedBody, 0755)
	if err != nil {
		return err
	}
	return nil
}

// CreateAzureAttestationToken If you call this on macOS, you will get the following error.
// SIGSYS: bad system call
// PC=0x407f2d0 m=0 sigcode=0
func CreateAzureAttestationToken(cert []byte, attestationProviderURL string) (string, error) {
	return enclave.CreateAzureAttestationToken(cert, attestationProviderURL)
}
