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
	Mux   *http.ServeMux
	store datastore
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
func Serve(listenAddr, jsonFilename string) error {
	wa, err := NewWebAPI(jsonFilename)
	if err != nil {
		return err
	}

	return http.ListenAndServe(listenAddr, wa.Mux)
}

// NewWebAPI creates a webapi and initialized all fields
// Attach the Mux to a http.Serve to start the listener
func NewWebAPI(jsonFilename string) (*webapi, error) {
	var store datastore
	err := store.Init(jsonFilename)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	wa := webapi{mux, store}

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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// accessTokenHandler TODO
func (wa *webapi) accessTokenHandler(w http.ResponseWriter, r *http.Request) {
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
}

// notFoundHandler returns 404 error
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
}
