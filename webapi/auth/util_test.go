// Copyright (c) 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"testing"
)

func TestEncryptPassword(t *testing.T) {

	p := "ilovego"
	expect := "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU="
	s := EncryptPassword(p)
	if s != expect {
		t.Errorf("EncryptPassword('%v') = %v, expected %v", p, s, expect)
	}
}
