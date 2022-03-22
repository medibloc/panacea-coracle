package tee

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/edgelesssys/ego/ecrypto"
	"github.com/edgelesssys/ego/enclave"
	"github.com/medibloc/panacea-data-market-validator/types"
	log "github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/json"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
)

const (
	CertificateKey      = "certificate"
	CertificateFilename = "certificate_seal"
)

type SealCertificate struct {
	Cert    []byte
	PrivKey *rsa.PrivateKey
}

// CreateCertificate If there is a sealed certificate in the received path, it responds after parsing.
// However, if the certificate does not exist, it responds by creating a new one. It also responds with RSA PrivateKey.
func CreateCertificate(storePath string) ([]byte, *rsa.PrivateKey, error) {
	fileFullPath := filepath.Join(storePath, CertificateFilename)
	if exists(fileFullPath) {
		log.Info("A sealed certificate exists. Is doing read the certificate.")
		return readFileAndGetCertificate(fileFullPath)
	}

	log.Info("There is no certificate. Generate a new certificate.")

	return createCertificate(storePath)
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func readFileAndGetCertificate(path string) ([]byte, *rsa.PrivateKey, error) {
	sealedBody, err := ioutil.ReadFile(path)

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

func createCertificate(storePath string) ([]byte, *rsa.PrivateKey, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: "DataValidator"},
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	err = sealAndStore(certBytes, priv, storePath)
	if err != nil {
		return nil, nil, err
	}

	return certBytes, priv, nil
}

func sealAndStore(certBytes []byte, priv *rsa.PrivateKey, fileFullPath string) error {
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

	err = os.MkdirAll(fileFullPath, 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileFullPath, sealedBody, 0755)
	if err != nil {
		return err
	}
	return nil
}

// CreateAzureAttestationToken If you call this on macOS, you will get the following error.
// SIGSYS: bad system call
// PC=0x407f2d0 m=0 sigcode=0
func CreateAzureAttestationToken(cert []byte) (string, error) {
	return enclave.CreateAzureAttestationToken(cert, types.AttestationProviderURL)
}
