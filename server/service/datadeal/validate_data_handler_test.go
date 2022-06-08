package datadeal

import (
	"bytes"
	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gorilla/mux"
	datadealtypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	markettypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	"github.com/medibloc/panacea-oracle/cache"
	"github.com/medibloc/panacea-oracle/codec"
	"github.com/medibloc/panacea-oracle/config"
	"github.com/medibloc/panacea-oracle/crypto"
	"github.com/medibloc/panacea-oracle/panacea"
	"github.com/medibloc/panacea-oracle/server/service"
	"github.com/medibloc/panacea-oracle/store"
	"github.com/medibloc/panacea-oracle/types"
	"github.com/medibloc/panacea-oracle/types/testutil"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	requesterPrivKey = secp256k1.GenPrivKey()
	requesterPubKey  = requesterPrivKey.PubKey()
	requesterAddress = testutil.GetAddress(requesterPubKey)

	buyerPrivKey = secp256k1.GenPrivKey()
	buyerPubKey  = buyerPrivKey.PubKey()
	buyerAddress = testutil.GetAddress(buyerPubKey)

	oracleMnemonic, _ = crypto.GenerateMnemonic()
	oracleAccount, _  = panacea.NewOracleAccount(oracleMnemonic)

	dealPrivKey = secp256k1.GenPrivKey()
	dealPubKey  = dealPrivKey.PubKey()
	dealAddress = testutil.GetAddress(dealPubKey)

	defaultData = []byte(`{
		"name": "This is a name",
		"description": "This is a description, man",
		"body": [{ "type": "markdown", "attributes": { "value": "val1" } }]
	}`)
)

func makeMockStore() store.Storage {
	mockStore := testutil.NewMockStore()
	// set store file. e.g) store.Upload(path, fileName, data)

	return mockStore
}

func makeMockSvc() *service.Service {
	conf := config.DefaultConfig()
	conf.OracleMnemonic = oracleMnemonic
	return &service.Service{
		Conf:          conf,
		Cache:         cache.NewAuthenticationCache(conf),
		PanaceaClient: makeMockGrpcClient(),
		Store:         makeMockStore(),
		OracleAccount: oracleAccount,
	}
}

func makeMockGrpcClient() panacea.GrpcClientI {
	accounts := []authtypes.AccountI{
		testutil.NewBaseAccount(requesterPubKey, 0, 0),
		testutil.NewBaseAccount(buyerPubKey, 1, 0),
	}

	budget := sdk.NewCoin("umed", sdk.NewInt(1000000000))
	deals := []datadealtypes.Deal{
		{
			DealId:         1,
			DealAddress:    dealAddress,
			DataSchema:     []string{"https://json.schemastore.org/github-issue-forms.json"},
			Budget:         &budget,
			TrustedOracles: []string{oracleAccount.GetAddress()},
			MaxNumData:     10,
			CurNumData:     1,
			Owner:          buyerAddress,
			Status:         "PENDING",
		},
	}
	//deals := make([]datadealtypes.Deal, 0)
	return testutil.NewMockGrpcClient(
		accounts,
		deals,
		nil,
		nil,
		nil,
		nil,
	)
}

func makeDataDealService() dataDealService {
	return dataDealService{
		makeMockSvc(),
	}
}

func setDefaultURLVars(req *http.Request) *http.Request {
	pathMap := make(map[string]string)
	pathMap[types.DealIDKey] = "1"
	return mux.SetURLVars(req, pathMap)
}

func TestHandlerValidateData(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v0/data-deal/deals/1/data?requester_address="+requesterAddress, bytes.NewReader(defaultData))
	req = setDefaultURLVars(req)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	dealService := makeDataDealService()

	dealService.handleValidateData(recorder, req)

	res := recorder.Result()

	require.Equal(t, http.StatusCreated, res.StatusCode)
	resBody := recorder.Body.Bytes()

	cert := markettypes.DataCert{}
	err := codec.ProtoUnmarshalJSON(resBody, &cert)
	require.NoError(t, err)
	unsignedCert := cert.UnsignedCert
	require.Equal(t, uint64(1), unsignedCert.DealId)
	// check if dataUrl is encrypted with buyer's pubKey
	buyerPrivKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), buyerPrivKey.Bytes())
	url, err := btcec.Decrypt(buyerPrivKey, unsignedCert.EncryptedDataUrl)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(string(url), "1/"))
	require.Equal(t, oracleAccount.GetAddress(), unsignedCert.OracleAddress)
	require.Equal(t, crypto.Hash(defaultData), unsignedCert.DataHash)
	require.Equal(t, requesterAddress, unsignedCert.RequesterAddress)
	// verify that oracle signed it
	serializedCertificate, err := unsignedCert.Marshal()
	require.NoError(t, err)
	require.True(t, oracleAccount.GetSecp256k1PubKey().VerifySignature(serializedCertificate, cert.Signature))
}

func TestHandlerValidateDataInvalidContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v0/data-deal/deals/1/data?requester_address="+requesterAddress, bytes.NewReader(defaultData))
	req = setDefaultURLVars(req)

	req.Header.Set("Content-Type", "application/text")

	recorder := httptest.NewRecorder()

	dealService := makeDataDealService()

	dealService.handleValidateData(recorder, req)

	require.Equal(t, http.StatusUnsupportedMediaType, recorder.Result().StatusCode)
	require.Equal(t, "only application/json is supported\n", recorder.Body.String())
}

func TestHandlerValidateDataInvalidParameter(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v0/data-deal/deals/1/data", bytes.NewReader(defaultData))
	req = setDefaultURLVars(req)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	dealService := makeDataDealService()

	dealService.handleValidateData(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	require.Equal(t, "failed to read query parameter\n", recorder.Body.String())
}

func TestHandlerValidateDataEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v0/data-deal/deals/1/data?requester_address="+requesterAddress, nil)
	req = setDefaultURLVars(req)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	dealService := makeDataDealService()

	dealService.handleValidateData(recorder, req)

	res := recorder.Result()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, "failed to read HTTP request body\n", recorder.Body.String())
}

func TestHandlerValidateDataInvalidBodySchema(t *testing.T) {
	data := []byte(`{
		"name": "This is a name",
		"description": "This is a description, man"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/v0/data-deal/deals/1/data?requester_address="+requesterAddress, bytes.NewReader(data))
	req = setDefaultURLVars(req)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	dealService := makeDataDealService()

	dealService.handleValidateData(recorder, req)

	res := recorder.Result()

	require.Equal(t, http.StatusForbidden, res.StatusCode)
	require.Equal(t, "JSON schema validation failed\n", recorder.Body.String())
}
