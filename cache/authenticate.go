package cache

import (
	"fmt"
	"github.com/bluele/gcache"
	"github.com/medibloc/panacea-data-market-validator/config"
	types "github.com/medibloc/panacea-data-market-validator/types/auth"
)

type AuthenticationCache struct {
	Cache gcache.Cache
}

func NewAuthenticationCache(conf *config.Config) *AuthenticationCache {
	cache := gcache.
		New(conf.Authenticate.Size).
		LRU().
		Expiration(conf.Authenticate.Expiration).
		Build()
	return &AuthenticationCache{
		Cache: cache,
	}
}

func (m AuthenticationCache) Set(auth *types.SignatureAuthentication) error {
	key := makeKey(auth.KeyId, auth.Nonce)
	err := m.Cache.Set(key, auth)
	if err != nil {
		return err
	}
	return nil
}

func (m AuthenticationCache) Get(keyId, nonce string) *types.SignatureAuthentication {
	key := makeKey(keyId, nonce)
	value, err := m.Cache.Get(key)
	if err != nil {
		return nil
	}

	sig, ok := value.(*types.SignatureAuthentication)

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
	return fmt.Sprintf("%s|%s", keyId, nonce)
}
