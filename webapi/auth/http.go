// Copyright (c) 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package auth

import (
	"encoding/json"
	"net/http"
	"regexp"
)

var (
	proxyAuthRegex *regexp.Regexp
	tokenRegex     *regexp.Regexp
)

type webapi struct {
	Mux          *http.ServeMux
	store        datastore
	tokenTimeout int
}

type response struct {
	AccessGranted bool   `json:"access_granted"`
	Reason        string `json:"reason,omitempty"`
}

func init() {
	proxyAuthRegex = regexp.MustCompile("^/api/2/domains/(.+)/proxyauth$")
	tokenRegex = regexp.MustCompile("^/api/2/domains/(.+)/oauth/access_token$")
}

// Serve creates a webapi and starts the http server
// Will listen and block unless something goes wrong
func Serve(listenAddr, jsonFilename string, tokenTimeout int) error {
	wa, err := NewWebAPI(jsonFilename, tokenTimeout)
	if err != nil {
		return err
	}

	return http.ListenAndServe(listenAddr, wa.Mux)
}

// NewWebAPI creates a webapi and initialized all fields
// Attach the Mux to a http.Serve to start the listener
func NewWebAPI(jsonFilename string, tokenTimeout int) (*webapi, error) {
	var store datastore
	err := store.Init(jsonFilename)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	wa := webapi{mux, store, tokenTimeout}

	wa.Mux.HandleFunc("/api/2/domains/", wa.domainRouter)
	wa.Mux.HandleFunc("/", notFoundHandler)

	return &wa, nil
}

// domainRouter determines if this is a proxyAuth or access_token request
// and routes to those handlers.  Else it returns a 404 error.
func (wa *webapi) domainRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		notFoundHandler(w, r)
	} else if proxyAuthRegex.MatchString(r.URL.Path) {
		wa.proxyAuthHandler(w, r)
	} else if tokenRegex.MatchString(r.URL.Path) {
		wa.accessTokenHandler(w, r)
	} else {
		notFoundHandler(w, r)
	}
}

// proxyAuthHandler handles domain authentiation based on data in store
func (wa *webapi) proxyAuthHandler(w http.ResponseWriter, r *http.Request) {
	matches := proxyAuthRegex.FindStringSubmatch(r.URL.Path)
	if len(matches) != 2 {
		notFoundHandler(w, r)
		return
	}

	// matches[1] is the domain to lookup
	domain := matches[1]
	ok := wa.store.DomainExists(domain)
	if !ok {
		notFoundHandler(w, r)
		return
	}

	err := wa.store.ValidateAuthHeader(w, r, domain)
	if err != nil {
		badRequestHandler(w, r, err.Error())
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	ok = wa.store.UserPasswordValid(domain, username, password)
	res := response{ok, ""}
	if !ok {
		res.Reason = "denied by policy"
	}

	js, err := json.Marshal(res)
	if err != nil {
		internalErrorHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		internalErrorHandler(w, r)
		return
	}
}

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

type errorJson struct {
	Error string `json:"error"`
}

func badRequestHandler(w http.ResponseWriter, r *http.Request, errStatus string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	if len(errStatus) > 0 {
		e := errorJson{errStatus}
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

func internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
}
