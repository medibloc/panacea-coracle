package auth

type SignatureAuthentication struct {
	Algorithm string
	KeyId     string
	Nonce     string
	Signature string
}
