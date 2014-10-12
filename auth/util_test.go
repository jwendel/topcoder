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
