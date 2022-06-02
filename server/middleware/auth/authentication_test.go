package auth_test

import (
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gorilla/context"
	"github.com/medibloc/panacea-data-market-validator/cache"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/panacea"
	"github.com/medibloc/panacea-data-market-validator/server/middleware/auth"
	"github.com/medibloc/panacea-data-market-validator/server/service"
	"github.com/medibloc/panacea-data-market-validator/server/service/datapool"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/types/testutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	requesterPrivKey = secp256k1.GenPrivKey()
	requesterPubKey  = requesterPrivKey.PubKey()
	requesterAddress = testutil.GetAddress(requesterPubKey)
)

func TestParseSignatureAuthorizationNotSupportType(t *testing.T) {
	middleware := auth.NewMiddleware(makeMockSvc())
	datapool.RegisterMiddleware(middleware)
	req := httptest.NewRequest(http.MethodGet, "/v0/data-pool/pools/1/data", nil)
	recorder := httptest.NewRecorder()

	middleware.Middleware(makeCheckRequesterHandler()).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	require.Equal(t, "not supported auth type\n", recorder.Body.String())
}

func TestBasicValidateAlgorithmWrong(t *testing.T) {
	middleware := auth.NewMiddleware(makeMockSvc())
	datapool.RegisterMiddleware(middleware)

	authHeader := makeAuthorizationHeader("es256k1", requesterAddress, "", "")
	req := httptest.NewRequest(http.MethodGet, "/v0/data-pool/pools/1/data", nil)
	req.Header.Add("Authorization", authHeader)
	recorder := httptest.NewRecorder()

	middleware.Middleware(makeCheckRequesterHandler()).ServeHTTP(recorder, req)

	res := recorder.Result()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, "is not supported value. (Algorithm: es256k1)\n", recorder.Body.String())
}

func TestBasicValidateKeyIdEmpty(t *testing.T) {
	middleware := auth.NewMiddleware(makeMockSvc())
	datapool.RegisterMiddleware(middleware)

	authHeader := makeAuthorizationHeader(auth.EsSha256, "", "", "")
	req := httptest.NewRequest(http.MethodGet, "/v0/data-pool/pools/1/data", nil)
	req.Header.Add("Authorization", authHeader)
	recorder := httptest.NewRecorder()

	middleware.Middleware(makeCheckRequesterHandler()).ServeHTTP(recorder, req)

	res := recorder.Result()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Equal(t, "'KeyId' cannot be empty\n", recorder.Body.String())
}

func TestParseSignatureAuthenticationParts(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v0/tee/attestation-token", nil)
	req.Header.Add("Authorization", "Signature\n    keyId=\"panacea1xxxx\",\n    algorithm=\"es256k1-sha256\",\n    nonce=\"\",\n    signature=\"\"")

	parts, err := auth.ParseSignatureAuthorizationParts(req.Header.Get("Authorization"))
	require.NoError(t, err)

	require.Equal(t, "panacea1xxxx", parts[types.AuthKeyIDHeaderKey])
	require.Equal(t, "es256k1-sha256", parts[types.AuthAlgorithmHeaderKey])
	require.Equal(t, "", parts[types.AuthNonceHeaderKey])
	require.Equal(t, "", parts[types.AuthSignatureHeaderKey])
}

func makeNormalHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Success execute normalHandler")
	}
}

func makeCheckRequesterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Success execute checkRequestHandler")
		requester := context.Get(r, types.RequesterAddressKey)
		if requesterAddress != requester {
			panic(fmt.Sprintf("Do not matched requester(%s) and requesterAddress(%s)", requester, requesterAddress))
		}
		log.Info("requesterAddress: ", requester)
	}
}

func makeMockSvc() *service.Service {
	conf := config.DefaultConfig()
	grpcClient := makeMockGrpcClient()
	return &service.Service{
		Conf:          conf,
		Cache:         cache.NewAuthenticationCache(conf),
		PanaceaClient: grpcClient,
	}
}

func makeMockGrpcClient() panacea.GrpcClientI {
	accounts := []authtypes.AccountI{
		testutil.NewBaseAccount(requesterPubKey, 0, 0),
	}

	return testutil.NewMockGrpcClient(
		accounts,
		nil,
		nil,
		nil,
	)
}

func makeAuthorizationHeader(algorithm, keyId, nonce, signature string) string {
	return fmt.Sprintf(
		"Signature keyId=\"%s\", "+
			"algorithm=\"%s\", "+
			"nonce=\"%s\", "+
			"signature=\"%s\"",
		keyId,
		algorithm,
		nonce,
		signature,
	)
}

func TestEntireAuthenticationProcess(t *testing.T) {
	svc := makeMockSvc()
	middleware := auth.NewMiddleware(svc)
	datapool.RegisterMiddleware(middleware)

	authHeader := makeAuthorizationHeader(auth.EsSha256, requesterAddress, "", "")
	req := httptest.NewRequest(http.MethodGet, "/v0/data-pool/pools/1/data", nil)
	req.Header.Add("Authorization", authHeader)
	recorder := httptest.NewRecorder()

	middleware.Middleware(makeCheckRequesterHandler()).ServeHTTP(recorder, req)

	res := recorder.Result()
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	authorizationParts, err := auth.ParseSignatureAuthorizationParts(res.Header.Get("WWW-Authenticate"))
	algorithm := authorizationParts[types.AuthAlgorithmHeaderKey]
	keyID := authorizationParts[types.AuthKeyIDHeaderKey]
	nonce := authorizationParts[types.AuthNonceHeaderKey]

	require.Equal(t, auth.EsSha256, algorithm)
	require.Equal(t, requesterAddress, keyID)
	require.NotEqual(t, "", nonce)
	require.NoError(t, err)
	signature, err := requesterPrivKey.Sign([]byte(nonce))
	require.NoError(t, err)

	// Check is exist authentication in cache
	inCachedAuthentication := svc.Cache.Get(keyID, nonce)
	require.NotNil(t, inCachedAuthentication)

	authHeader = makeAuthorizationHeader(
		auth.EsSha256,
		requesterAddress,
		nonce,
		base64.StdEncoding.EncodeToString(signature),
	)
	req.Header.Set("Authorization", authHeader)
	recorder = httptest.NewRecorder()

	middleware.Middleware(makeCheckRequesterHandler()).ServeHTTP(recorder, req)

	res = recorder.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "", res.Header.Get("WWW-Authenticate"))

	// Check is not exist authentication in cache
	inCachedAuthentication = svc.Cache.Get(keyID, nonce)
	require.Nil(t, inCachedAuthentication)
}

func TestNotIncludeAuthenticationURL(t *testing.T) {
	middleware := auth.NewMiddleware(makeMockSvc())
	datapool.RegisterMiddleware(middleware)

	authHeader := makeAuthorizationHeader(auth.EsSha256, requesterAddress, "", "")
	req := httptest.NewRequest(http.MethodPost, "/v0/data-pool/pools/1/rounds/{round}/data", nil)
	req.Header.Add("Authorization", authHeader)
	recorder := httptest.NewRecorder()

	middleware.Middleware(makeNormalHandler()).ServeHTTP(recorder, req)

	res := recorder.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "", res.Header.Get("WWW-Authenticate"))
}

func TestSameURLDifferentMethod(t *testing.T) {
	middleware := auth.NewMiddleware(makeMockSvc())
	datapool.RegisterMiddleware(middleware)

	authHeader := makeAuthorizationHeader(auth.EsSha256, requesterAddress, "", "")
	req := httptest.NewRequest(http.MethodPost, "/v0/data-pool/pools/1/data", nil)
	req.Header.Add("Authorization", authHeader)
	recorder := httptest.NewRecorder()

	middleware.Middleware(makeNormalHandler()).ServeHTTP(recorder, req)

	res := recorder.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "", res.Header.Get("WWW-Authenticate"))
}
