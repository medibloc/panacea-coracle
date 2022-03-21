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

func CreateCertificate(storePath string) (*SealCertificate, error) {
	fileFullPath := filepath.Join(storePath, CertificateFilename)
	if exists(fileFullPath) {
		log.Info("A sealed certificate exists. Is doing read the certificate.")
		return readFileAndGetCertificate(fileFullPath)
	}

	log.Info("There is no certificate. Generate a new certificate.")

	sealCertificate, err := createSealCertificate()
	if err != nil {
		return nil, err
	}

	err = sealAndStore(sealCertificate, fileFullPath)
	if err != nil {
		return nil, err
	}

	return sealCertificate, nil
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func readFileAndGetCertificate(path string) (*SealCertificate, error) {
	sealedBody, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	bytes, err := ecrypto.Unseal(sealedBody, []byte(CertificateKey))
	if err != nil {
		return nil, err
	}

	cert := &SealCertificate{}
	err = json.Unmarshal(bytes, cert)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func createSealCertificate() (*SealCertificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: "DataValidator"},
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	return &SealCertificate{
		Cert:    certBytes,
		PrivKey: priv,
	}, nil
}

func sealAndStore(certificateWithPrivKey *SealCertificate, fileFullPath string) error {
	jsonBytes, err := json.Marshal(certificateWithPrivKey)
	if err != nil {
		return err
	}

	sealedBody, err := ecrypto.SealWithProductKey(jsonBytes, []byte(CertificateKey))
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
