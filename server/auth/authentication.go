package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"net/http"
	"strings"
)

const prefix = "Signature"

type SignatureAuthentication struct {
	Algorithm string
	Realm string
	Qop string
	Nonce string
	Opaque string
	Signature string
}

func (m SignatureAuthentication) BasicValidate() error {
	if m.Algorithm != "es256k1-sha256" {
		return errors.New(fmt.Sprintf("is not supported value. (algorithm: %s)", m.Algorithm))
	} else if m.Realm != "PanaceaAccount" {
		return errors.New(fmt.Sprintf("is not supported value. (realm: %s)", m.Realm))
	} else if m.Qop != "auth" {
		return errors.New(fmt.Sprintf("is not supported value. (qop: %s)", m.Qop))
	} else if m.Nonce == "" {
		return errors.New("'nonce' cannot be empty")
	} else if m.Opaque == "" {
		return errors.New("'opaque' cannot be empty")
	} else if m.Signature == "" {
		return errors.New("'signature' cannot be empty")
	}
	return nil
}

func ParseSignatureAuthentication(r *http.Request) (*SignatureAuthentication, error) {
	auth := r.Header.Get("Authorization")

	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return nil, errors.New("not supported authentication type")
	}

	signatureAuth := &SignatureAuthentication{}

	for _, headers := range strings.Split(auth[len(prefix):], ",") {
		header := strings.Split(headers, "=")
		switch header[0] {
		case "algorithm":
			signatureAuth.Algorithm = header[1]
		case "realm":
			signatureAuth.Realm = header[1]
		case "qop":
			signatureAuth.Qop = header[1]
		case "nonce":
			signatureAuth.Nonce = header[1]
		case "opaque":
			signatureAuth.Opaque = header[1]
		case "signature":
			signatureAuth.Signature = header[1]
		}
	}

	return signatureAuth, nil
}

func MakeNonce() (string, error) {
	token := make([]byte, 24)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(crypto.Hash(token)), nil
}