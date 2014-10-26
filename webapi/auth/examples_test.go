// Copyright (c) 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type tokenData struct {
	url        string
	id         string
	secret     string
	grant_type string
	// expected responses
	code        int
	contentType string
	// body        string
}

func TestOath(t *testing.T) {
	wa, err := NewWebAPI("test_data.json", 3600)
	if err != nil {
		t.Fatal("Failed to create webapi: ", err)
	}

	ts := httptest.NewServer(wa.Mux)
	defer ts.Close()

	topcoderUrl := ts.URL + "/api/2/domains/topcoder.com/oauth/access_token"

	td := tokenData{topcoderUrl, "s6BhdRkqt3", "7Fjfp0ZBr1KtDRbnfVdmIw", "client_credentials", 200, "application/json"}
	b, m := runAccessToken(t, td)
	t.Log("body:", b, "\nmap:", m)
}

func runAccessToken(t *testing.T, td tokenData) (string, map[string]interface{}) {
	pf := url.Values{}
	pf.Add("client_id", td.id)
	pf.Add("client_secret", td.secret)
	pf.Add("grant_type", td.grant_type)

	res, err := http.PostForm(td.url, pf)
	if err != nil {
		t.Fatalf("Failed to get response from server: %v\ntokenData: %v", err, td)
	}
	if res.StatusCode != td.code {
		t.Errorf("StatusCode missmatch.  Got %v  Expected %v\ntokenData: %v", res.StatusCode, td.code, td)
	}
	ct := res.Header.Get("Content-Type")
	if ct != td.contentType {
		t.Errorf("Content-Type missmatch.  Got %v  Expected %v\ntokenData: %v", ct, td.contentType, td)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unable to read body: %v\ntokenData: %v", err, td)
	}

	body := string(bytes[:])
	var m map[string]interface{}

	if ct == "application/json" {
		var f interface{}
		err := json.Unmarshal(bytes, &f)
		if err != nil {
			t.Errorf("Failed to parse json: %v\ntokenData: %v\nbody: %v", err, td, body)
		} else {
			m = f.(map[string]interface{})
		}
	}
	// if body != td.body {
	// 	t.Errorf("Body data does not match, got %v  Expected %v\ntokenData: %v", body, td.body, td)
	// }

	return body, m
}
