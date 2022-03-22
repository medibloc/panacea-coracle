package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/edgelesssys/ego/attestation"
	"github.com/medibloc/panacea-data-market-validator/server/tee"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClient(t *testing.T) {
	serverURL := "https://localhost:8081/v0/tee/attestation-token"
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	res := httpGet(tlsConfig, serverURL)

	fmt.Println(string(res))

	tokenRes := tee.TokenResponse{}
	err := json.Unmarshal(res, &tokenRes)
	require.NoError(t, err)

	report, err := attestation.VerifyAzureAttestationToken(tokenRes.Token, types.AttestationProviderURL)
	require.NoError(t, err)

	certBytes := report.Data
	cert, err := x509.ParseCertificate(certBytes)
	require.NoError(t, err)

	tlsConfig = &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
	tlsConfig.RootCAs.AddCert(cert)
	res = httpGet(tlsConfig, serverURL)
	fmt.Println(string(res))
}

func httpGet(tlsConfig *tls.Config, url string) []byte {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}
