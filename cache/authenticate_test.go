package cache_test

import (
	"github.com/medibloc/panacea-data-market-validator/cache"
	"github.com/medibloc/panacea-data-market-validator/config"
	types "github.com/medibloc/panacea-data-market-validator/types/auth"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func makeTestAuthenticationCache() *cache.AuthenticationCache {
	conf := config.DefaultConfig()

	return cache.NewAuthenticationCache(conf)
}

// TestSetAndGetCache performs a test of add and remove of the cache.
func TestSetAndGetCache(t *testing.T) {
	c := makeTestAuthenticationCache()

	authentication := types.SignatureAuthentication{
		Algorithm: "es256k1-sha256",
		KeyId:     "panacea1xxx",
		Nonce:     "123",
		Signature: "signature",
	}
	err := c.Set(&authentication)
	require.NoError(t, err)

	authenticationResult := c.Get(authentication.KeyId, authentication.Nonce)
	require.NotNil(t, authenticationResult)
	require.Equal(t, authentication.Algorithm, authenticationResult.Algorithm)
	require.Equal(t, authentication.KeyId, authenticationResult.KeyId)
	require.Equal(t, authentication.Nonce, authenticationResult.Nonce)
	require.Equal(t, authentication.Signature, authenticationResult.Signature)

	time.Sleep(12 * time.Second)

	authenticationResult = c.Get(authentication.KeyId, authentication.Nonce)
	// This value must be followed by nil. Because the cache time has expired.
	require.Nil(t, authenticationResult)
}

// TestAddMoreThanCacheSize tests the behavior when more than the maximum size of the cache is added.
func TestAddMoreThanCacheSize(t *testing.T) {
	c := makeTestAuthenticationCache()

	for i := 0; i < 100000; i++ {
		authentication := types.SignatureAuthentication{
			Algorithm: "es256k1-sha256",
			KeyId:     "panacea1xxx",
			Nonce:     strconv.Itoa(i),
			Signature: "signature",
		}
		err := c.Set(&authentication)
		require.NoError(t, err)
	}

	// The first stored c is deleted. (LRU: Least Recently Used)
	for i := 0; i < 50000; i++ {
		require.Nil(t, c.Get("panacea1xxx", strconv.Itoa(i)))
	}

	for i := 50000; i < 100000; i++ {
		require.NotNil(t, c.Get("panacea1xxx", strconv.Itoa(i)))
	}
}
