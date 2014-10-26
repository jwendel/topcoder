// Copyright (c) 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttp(t *testing.T) {
	wa, err := NewWebAPI("test_data.json", 3600)
	if err != nil {
		t.Fatal("Failed to create webapi: ", err)
	}

	ts := httptest.NewServer(wa.Mux)
	defer ts.Close()

	// Test cases 1-4 from original test were removed as they require Authorization data now.
	// See examples_test.go for http tests.

	// Test case 5 - domain fail, expect 404"
	// curl -i --data "username=takumi&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost:8080/api/2/domains/example.com/proxyauth
	u := ts.URL + "/api/2/domains/example.com/proxyauth"
	pf := url.Values{}
	pf.Add("username", "takumi")
	pf.Add("password", "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=")
	res, err := http.PostForm(u, pf)
	if err != nil {
		t.Fatalf("Failed to get response from server: %v", err)
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("POST %v - params: %v - unable to read body: %v", u, pf, err)
	}
	body := string(bytes[:])
	expected := ""
	if res.StatusCode != http.StatusNotFound || body != expected {
		t.Errorf("POST %v - params: %v - res: %v - body: %v", u, pf, res, body)
	}
}

func TestHttpExtra(t *testing.T) {
	wa, err := NewWebAPI("test_data.json", 3600)
	if err != nil {
		t.Fatal("Failed to create webapi: ", err)
	}

	ts := httptest.NewServer(wa.Mux)
	defer ts.Close()

	// GET test for 404
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("Failed to get response from server: %v", err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("GET / - status code: %v expected: %v", res.StatusCode, http.StatusNotFound)
	}

	// Test for invalid POST param
	u := ts.URL + "/api/2/domains/appirio.com/proxyauth"
	pf := url.Values{}
	pf.Add("BADusername", "jun")
	pf.Add("password", "{SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=")
	res, err = http.PostForm(u, pf)
	if err != nil {
		t.Fatalf("Failed to get response from server: %v", err)
	}
	defer res.Body.Close()
}
