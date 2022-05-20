package auth

import (
	"fmt"
	"strings"
)

type SignatureAuthentication struct {
	Algorithm string
	KeyId     string
	Nonce     string
	Signature string
}

func (m *SignatureAuthentication) ToWWWAuthenticateValue() string {
	s := []string{
		fmt.Sprintf("algorithm=%s", m.Algorithm),
		fmt.Sprintf("keyId=%s", m.Algorithm),
		fmt.Sprintf("nonce=%s", m.Algorithm),
	}
	return strings.Join(s, ", ")
}
