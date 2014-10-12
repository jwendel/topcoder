package auth

import (
	"crypto/sha256"
	"encoding/base64"
)

// encryptPassword takes a plaintext password and hashes it with SHA256
func encryptPassword(pw string) string {
	hasher := sha256.New()
	b := []byte(pw)
	hasher.Write(b)
	s := "{SHA256}" + base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return s
}
