// Copyright (c) 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"crypto/sha256"
	"encoding/base64"
)

// EncryptPassword takes a plaintext password and hashes it with SHA256
func EncryptPassword(pw string) string {
	hasher := sha256.New()
	b := []byte(pw)
	hasher.Write(b)
	s := "{SHA256}" + base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return s
}
