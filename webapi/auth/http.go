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

// Webapi represents the auth http server and the data
// needed to serve requests.
type Webapi struct {
	Mux   *http.ServeMux
	store datastore
}

type response struct {
	AccessGranted bool   `json:"access_granted"`
	Reason        string `json:"reason,omitempty"`
}

type errorJSON struct {
	Error string `json:"error"`
}

func init() {
	proxyAuthRegex = regexp.MustCompile("^/api/2/domains/(.+)/proxyauth$")
	tokenRegex = regexp.MustCompile("^/api/2/domains/(.+)/oauth/access_token$")
}

// Serve creates a Webapi and starts the http server
// Will listen and block unless something goes wrong
func Serve(listenAddr, jsonFilename, tokenFilename string, tokenTimeout int) error {
	wa, err := NewWebAPI(jsonFilename, tokenFilename, tokenTimeout)
	if err != nil {
		return err
	}

	return http.ListenAndServe(listenAddr, wa.Mux)
}

// NewWebAPI creates a Webapi and initialized all fields
// Attach the Mux to a http.Serve to start the listener
func NewWebAPI(jsonFilename, tokenFilename string, tokenTimeout int) (*Webapi, error) {
	var store datastore
	store.tokenTimeout = tokenTimeout
	err := store.Init(jsonFilename, tokenFilename)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	wa := Webapi{mux, store}

	wa.Mux.HandleFunc("/api/2/domains/", wa.domainRouter)
	wa.Mux.HandleFunc("/", notFoundHandler)

	return &wa, nil
}

// domainRouter determines if this is a proxyAuth or access_token request
// and routes to those handlers.  Else it returns a 404 error.
func (wa *Webapi) domainRouter(w http.ResponseWriter, r *http.Request) {
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
func (wa *Webapi) proxyAuthHandler(w http.ResponseWriter, r *http.Request) {
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

// accessTokenHandler is the route for generating access tokens.  It will load the client id
// and secret from the request, validate it for the domain and return an UUID.
// An error is written if anything goes wrong.
func (wa *Webapi) accessTokenHandler(w http.ResponseWriter, r *http.Request) {
	matches := tokenRegex.FindStringSubmatch(r.URL.Path)
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

	r.ParseForm()
	id := r.Form.Get("client_id")
	secret := r.Form.Get("client_secret")
	grant := r.Form.Get("grant_type")

	err := wa.store.ValidateClient(domain, id, secret, grant)
	if err != nil {
		badRequestHandler(w, r, err.Error())
		return
	}

	// access_token request is valid, generate token
	token := wa.store.generateAccessToken(domain)
	t := tokenResponse{token, "bearer", wa.store.tokenTimeout}
	b, err := json.Marshal(t)
	if err != nil {
		internalErrorHandler(w, r)
		return
	}

	successHandler(w, r, b)
}
