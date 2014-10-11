package main

import (
	"crypto/sha256"
	"encoding/base64"
)

// encryptPassword takes a
func encryptPassword(pw string) string {
	hasher := sha256.New()
	b := []byte(pw)
	hasher.Write(b)
	s := "{SHA256}" + base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return s
}
