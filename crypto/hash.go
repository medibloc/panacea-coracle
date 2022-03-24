package crypto

import (
	"crypto/rand"
	"crypto/sha256"
)

// Hash with SHA256.
func Hash(data []byte) []byte {
	hash := sha256.New()

	hash.Write(data)

	return hash.Sum(nil)
}

func RandomHash() ([]byte, error) {
	data := make([]byte, 100)
	if _, err := rand.Read(data); err != nil {
		return nil, err
	}
	return Hash(data), nil
}
