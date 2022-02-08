package e2e

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"
	"github.com/medibloc/panacea-data-market-validator/types"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestValidateData(t *testing.T) {
	buyerMnemonic := os.Getenv("E2E_DATA_BUYER_MNEMONIC")
	require.NotEmpty(t, buyerMnemonic)
	datavalHTTPAddr := os.Getenv("E2E_DATAVAL_HTTP_ADDR")
	require.NotEmpty(t, datavalHTTPAddr)

	dealID := 1
	requesterAddr := "panacea1c7yh0ql0rhvyqm4vuwgaqu0jypafnwqdc6x60e"
	data := `{
		"name": "This is a name",
		"description": "This is a description"
	}`

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/validate-data/%d?requester_address=%s", datavalHTTPAddr, dealID, requesterAddr),
		strings.NewReader(data),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("cert: %v", string(body))

	var cert types.DataValidationCertificateResponse
	err = json.Unmarshal(body, &cert)
	require.NoError(t, err)

	encDataURL, err := base64.StdEncoding.DecodeString(cert.UnsignedCert.EncryptedDataUrlBase64)
	require.NoError(t, err)

	privKey := getPrivateKey(t, buyerMnemonic)
	dataURL := string(decrypt(t, privKey, encDataURL))
	t.Logf("dataURL: %v", dataURL)

	downloadedData := downloadFile(t, dataURL)
	decryptedData := decrypt(t, privKey, downloadedData)

	require.Equal(t, data, string(decryptedData))
}

const (
	accountNum = 0
	coinType   = 371
	addressIdx = 0
)

func getPrivateKey(t *testing.T, mnemonic string) []byte {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	require.NoError(t, err)

	hdPath := hd.NewFundraiserParams(accountNum, coinType, addressIdx).String()
	masterPriv, chainCode := hd.ComputeMastersFromSeed(seed)

	privKey, err := hd.DerivePrivateKeyForPath(masterPriv, chainCode, hdPath)
	require.NoError(t, err)
	return privKey
}

func decrypt(t *testing.T, privKeyBz []byte, data []byte) []byte {
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBz[:])
	decrypted, err := btcec.Decrypt(privKey, data)
	require.NoError(t, err)
	return decrypted
}

func downloadFile(t *testing.T, url string) []byte {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	return data
}
