package cache

import (
	"github.com/bluele/gcache"
	"github.com/medibloc/panacea-data-market-validator/config"
	"strings"
)

type AuthenticationCache struct {
	Cache gcache.Cache
}

func NewAuthenticationCache(conf *config.Config) *AuthenticationCache {
	cache := gcache.
		New(conf.Authentication.Size).
		LRU().
		Expiration(conf.Authentication.Expiration).
		Build()
	return &AuthenticationCache{
		Cache: cache,
	}
}

func (m AuthenticationCache) Set(keyId, nonce string, sigAuthParts map[string]string) error {
	key := makeKey(keyId, nonce)
	err := m.Cache.Set(key, sigAuthParts)
	if err != nil {
		return err
	}
	return nil
}

func (m AuthenticationCache) Get(keyId, nonce string) map[string]string {
	key := makeKey(keyId, nonce)
	value, err := m.Cache.Get(key)
	if err != nil {
		return nil
	}

	sig, ok := value.(map[string]string)

	if !ok {
		return nil
	}

	return sig
}

func (m AuthenticationCache) Remove(keyId, nonce string) bool {
	key := makeKey(keyId, nonce)
	return m.Cache.Remove(key)
}

func makeKey(keyId, nonce string) string {
	return strings.Join([]string{keyId, nonce}, "|")
}
