package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-data-market-validator/server/service"
	"github.com/medibloc/panacea-data-market-validator/types"
	"github.com/medibloc/panacea-data-market-validator/validation"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strings"
)

const (
	prefixType = "Signature"

	Algorithm = "algorithm"
	KeyId     = "keyId"
	Nonce     = "nonce"
	Signature = "signature"
)

var authorizationHeaders = []string{Algorithm, KeyId, Nonce, Signature}
var authenticateHeaders = []string{Algorithm, KeyId, Nonce}

type AuthenticationMiddleware struct {
	service *service.Service
	url     map[string][]string
}

func NewMiddleware(svc *service.Service) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		service: svc,
		url:     make(map[string][]string),
	}
}

func (amw *AuthenticationMiddleware) AddURL(path string, methods ...string) {
	// Router is only used to convert path to regex.
	pathRegex, err := mux.NewRouter().Path(path).Methods(http.MethodGet).GetPathRegexp()
	if err != nil {
		panic(err)
	}
	m, ok := amw.url[pathRegex]
	if !ok {
		amw.url[pathRegex] = methods
	} else {
		amw.url[pathRegex] = append(m, methods...)
	}
}

func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if it is not an authentication URL, the authentication check is not performed.
		if !amw.isAuthenticationURL(r) {
			next.ServeHTTP(w, r)
			return
		}

		sigAuthParts, err := ParseSignatureAuthorizationParts(r.Header.Get("Authorization"))
		//signatureAuthentication, err := parseSignatureAuthentication(r)
		if err != nil {
			// result code to 400. Bad request
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = basicValidate(sigAuthParts)
		if err != nil {
			// result code to 400. Bad request
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if sigAuthParts[Nonce] == "" {
			err = amw.generateAuthenticationAndSetHeader(w, sigAuthParts)
			if err != nil {
				log.Error("failed to generate authentication", err)
				http.Error(w, "failed to generate authentication", http.StatusInternalServerError)
			}
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		cache := amw.service.Cache
		auth := cache.Get(sigAuthParts[KeyId], sigAuthParts[Nonce])
		if auth == nil {
			// expired
			err = amw.generateAuthenticationAndSetHeader(w, sigAuthParts)
			if err != nil {
				log.Error("failed to generate authentication", err)
				http.Error(w, "failed to generate authentication", http.StatusInternalServerError)
			}
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		nonce := sigAuthParts[Nonce]
		signature, err := base64.StdEncoding.DecodeString(sigAuthParts[Signature])
		if err != nil {
			http.Error(w, "failed to decode signature", http.StatusBadRequest)
			return
		}

		requesterAddress := sigAuthParts[KeyId]
		pubKey, err := amw.service.PanaceaClient.GetPubKey(requesterAddress)
		if err != nil {
			log.Error(fmt.Sprintf("failed to get the account's public key account: %s", requesterAddress), err)
			http.Error(w, "failed to get the account's public key", http.StatusInternalServerError)
			return
		}

		if !pubKey.VerifySignature([]byte(nonce), signature) {
			http.Error(w, "failed to verification signature", http.StatusBadRequest)
			return
		}

		context.Set(r, types.RequesterAddressKey, requesterAddress)
		next.ServeHTTP(w, r)

		err = amw.generateAuthenticationAndSetHeader(w, sigAuthParts)
		if err != nil {
			log.Error("failed to generate authentication", err)
			http.Error(w, "failed to generate authentication", http.StatusInternalServerError)
		}
	})
}

func (amw *AuthenticationMiddleware) isAuthenticationURL(r *http.Request) bool {
	for path, methods := range amw.url {
		ok, err := regexp.MatchString(path, r.URL.Path)
		if err == nil && ok {
			if validation.Contains(methods, r.Method) {
				return true
			}
		}
	}
	return false
}

// ParseSignatureAuthorizationParts parses Authorization value in Header according to `Signature` type.
func ParseSignatureAuthorizationParts(auth string) (map[string]string, error) {
	if len(auth) < len(prefixType) || !strings.EqualFold(auth[:len(prefixType)], prefixType) {
		return nil, errors.New("not supported auth type")
	}

	headers := strings.Split(auth[len(prefixType):], ",")
	parts := make(map[string]string, len(authorizationHeaders))
	for _, header := range headers {
		for _, w := range authorizationHeaders {
			if strings.Contains(header, w) {
				parts[w] = strings.Split(header, `"`)[1]
			}
		}
	}

	return parts, nil
}

func basicValidate(sigAuthParts map[string]string) error {
	if sigAuthParts[Algorithm] != "es256k1-sha256" {
		return errors.New(fmt.Sprintf("is not supported value. (Algorithm: %s)", sigAuthParts[Algorithm]))
	} else if sigAuthParts[KeyId] == "" {
		return errors.New("'KeyId' cannot be empty")
	}
	return nil
}

func (amw *AuthenticationMiddleware) generateAuthenticationAndSetHeader(w http.ResponseWriter, sigAuthParts map[string]string) error {
	err := amw.generateAuthentication(sigAuthParts)
	if err != nil {
		return err
	}
	auth := makeAuthenticationHeader(sigAuthParts)

	w.Header().Set("WWW-Authenticate", auth)

	return nil
}

func (amw *AuthenticationMiddleware) generateAuthentication(sigAuthParts map[string]string) error {
	err := setNewNonce(sigAuthParts)
	if err != nil {
		return err
	}

	cache := amw.service.Cache
	err = cache.Set(sigAuthParts[KeyId], sigAuthParts[Nonce], sigAuthParts)
	if err != nil {
		return err
	}

	return nil
}

func setNewNonce(sigAuthParts map[string]string) error {
	randomBytes := make([]byte, 24)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return err
	}

	sigAuthParts[Nonce] = base64.StdEncoding.EncodeToString(randomBytes)

	return nil
}

func makeAuthenticationHeader(sigAuthParts map[string]string) string {
	var authenticate []string
	for _, h := range authenticateHeaders {
		authenticate = append(authenticate, fmt.Sprintf("%s=\"%s\"", h, sigAuthParts[h]))
	}

	return prefixType + " " + strings.Join(authenticate, ", ")
}
