package cache_test

import (
	"github.com/medibloc/panacea-data-market-validator/cache"
	"github.com/medibloc/panacea-data-market-validator/config"
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

	authentication := make(map[string]string)
	authentication["algorithm"] = "es256k1-sha256"
	authentication["keyId"] = "panacea1xxx"
	authentication["nonce"] = "123"
	authentication["signature"] = "signature"
	err := c.Set(authentication["keyId"], authentication["nonce"], authentication)
	require.NoError(t, err)

	authenticationResult := c.Get(authentication["keyId"], authentication["nonce"])
	require.NotNil(t, authenticationResult)
	require.Equal(t, authentication["algorithm"], authenticationResult["algorithm"])
	require.Equal(t, authentication["keyId"], authenticationResult["keyId"])
	require.Equal(t, authentication["nonce"], authenticationResult["nonce"])
	require.Equal(t, authentication["signature"], authenticationResult["signature"])

	time.Sleep(12 * time.Second)

	authenticationResult = c.Get(authentication["keyId"], authentication["nonce"])
	// This value must be followed by nil. Because the cache time has expired.
	require.Nil(t, authenticationResult)
}

// TestAddMoreThanCacheSize tests the behavior when more than the maximum size of the cache is added.
func TestAddMoreThanCacheSize(t *testing.T) {
	c := makeTestAuthenticationCache()

	for i := 0; i < 100000; i++ {
		authentication := make(map[string]string)
		authentication["algorithm"] = "es256k1-sha256"
		authentication["keyId"] = "panacea1xxx"
		authentication["nonce"] = strconv.Itoa(i)
		authentication["signature"] = "signature"
		err := c.Set(authentication["keyId"], authentication["nonce"], authentication)
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

func TestAddAndRemove(t *testing.T) {
	c := makeTestAuthenticationCache()

	authentication := make(map[string]string)
	authentication["algorithm"] = "es256k1-sha256"
	authentication["keyId"] = "panacea1xxx"
	authentication["nonce"] = "123"
	authentication["signature"] = "signature"
	err := c.Set(authentication["keyId"], authentication["nonce"], authentication)
	require.NoError(t, err)

	ok := c.Remove(authentication["keyId"], authentication["nonce"])
	require.True(t, ok)

	authenticationResult := c.Get(authentication["keyId"], authentication["nonce"])
	require.Nil(t, authenticationResult)
}
