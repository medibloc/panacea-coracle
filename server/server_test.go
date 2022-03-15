package server_test

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/edgelesssys/ego/attestation"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/server"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestTLS(t *testing.T) {
	mnemonic, err := crypto.GenerateMnemonic()
	require.NoError(t, err)

	conf := config.Config{
		HTTPListenAddr: "localhost:9999",
		PanaceaGrpcAddress: "localhost:9090",
		ValidatorMnemonic: mnemonic,
		AWSS3Bucket: "data-market-test",
		AWSS3Region: "ap-northeast-2",
	}
	go server.Run(&conf)

	time.Sleep(time.Second)

	tokenRes, err := getToken(conf)
	require.NoError(t, err)

	report, err := attestation.VerifyAzureAttestationToken(tokenRes.Token, types.AttestationProviderURL)
	require.NoError(t, err)

	certBytes := report.Data
	cert, err := x509.ParseCertificate(certBytes)
	require.NoError(t, err)

	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	tlsConfig.RootCAs.AddCert(cert)

	requestBody := bytes.NewBufferString("Request validate data")

	res, err := callHttp(tlsConfig, "POST", fmt.Sprintf("https://%s/v1/data-pool/pools/%s/rounds/%s/data", conf.HTTPListenAddr, "1", "1"), requestBody)

	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	fmt.Println(string(body))
}

func getToken(conf config.Config) (types.TokenResponse, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	res, err := callHttp(tlsConfig, "GET", fmt.Sprintf("https://%s/v1/tee/attestation-token", conf.HTTPListenAddr), nil)
	if err != nil {
		return types.TokenResponse{}, err
	}

	defer res.Body.Close()

	if http.StatusCreated != res.StatusCode {
		return types.TokenResponse{}, errors.New("StatusCode is not 201")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return types.TokenResponse{}, err
	}

	tokenRes := types.TokenResponse{}
	err = json.Unmarshal(body, &tokenRes)
	if err != nil {
		return types.TokenResponse{}, err
	}

	return tokenRes, nil
}

func callHttp(tlsConfig *tls.Config, method, url string, body io.Reader) (*http.Response, error) {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}