package crypto

import "crypto/sha256"

// Hash with SHA256.
func Hash(data []byte) []byte {
	hash := sha256.New()

	hash.Write(data)

	return hash.Sum(nil)
}