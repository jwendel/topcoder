// Copyright (c) 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// EncryptPassword takes a plaintext password and hashes it with SHA256
func EncryptPassword(pw string) string {
	hasher := sha256.New()
	b := []byte(pw)
	hasher.Write(b)
	s := "{SHA256}" + base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return s
}

// successHandler returns status 200 with the body populated with response.
func successHandler(w http.ResponseWriter, r *http.Request, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(response)
	if err != nil {
		internalErrorHandler(w, r)
		return
	}
}

// notFoundHandler returns 404 error
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
}

// badRequestHandler returns error 400 with json body with 1 field of "error", with
// a value of errStatus
func badRequestHandler(w http.ResponseWriter, r *http.Request, errStatus string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	if len(errStatus) > 0 {
		e := errorJSON{errStatus}
		js, err := json.Marshal(e)
		if err != nil {
			internalErrorHandler(w, r)
			return
		}

		_, err = w.Write(js)
		if err != nil {
			internalErrorHandler(w, r)
			return
		}
	}
}

// internalErrorHandler returns an error 500
func internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
}
