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
	bodyJson map[string]interface{}
}

func TestOath(t *testing.T) {
	wa, err := NewWebAPI("test_data.json", 3600)
	if err != nil {
		t.Fatal("Failed to create webapi: ", err)
	}

	ts := httptest.NewServer(wa.Mux)
	defer ts.Close()

	topcoderUrl := ts.URL + "/api/2/domains/topcoder.com/oauth/access_token"

	bodyJson := map[string]interface{}{"token_type": "bearer", "expires_in": float64(3600), "access_token": nil}
	td := tokenData{topcoderUrl, "s6BhdRkqt3", "7Fjfp0ZBr1KtDRbnfVdmIw", "client_credentials", 200, "application/json", bodyJson}
	b, m := runAccessToken(t, td)
	at := m["access_token"].(string)
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

	// Lets try to decode the json and compare to the passed in bodyJson
	if ct == "application/json" {
		var f interface{}
		err := json.Unmarshal(bytes, &f)
		if err != nil {
			t.Errorf("Failed to parse json: %v\ntokenData: %v\nbody: %v", err, td, body)
		} else {
			m = f.(map[string]interface{})
			if len(m) != len(td.bodyJson) {
				t.Errorf("Reponse doesn't have the same number of keys.  Got: %v  Expected %v\ntokenData: %v", m, td.bodyJson, td)
			}

			// Loop over the keys/values expected and comapre to what we got
			for k, tdv := range td.bodyJson {
				mv, ok := m[k]
				if !ok {
					t.Errorf("Reponse missing key %v.  Got: %v  Expected %v\ntokenData: %v", k, m, td.bodyJson, td)
					continue
				}
				t.Logf("testing key %v  value: %v", k, tdv)
				// don't check values if expected data is nil
				if tdv == nil {
					continue
				}
				switch tdvv := tdv.(type) {
				case string:
					mvv, ok := mv.(string)
					if !ok || tdvv != mvv {
						t.Errorf("Types/Values don't match for key %v.  Got: %v  Expected %v\ntokenData: %v", k, m, td.bodyJson, td)
					}
				case float64:
					mvv, ok := mv.(float64)
					if !ok || tdvv != mvv {
						t.Errorf("Types/Values don't match for key %v.  Got: %v  Expected %v\ntokenData: %v", k, m, td.bodyJson, td)
					}
				default:
					t.Errorf("Unhandled type for key %v.  Got: %v  Expected %v\ntokenData: %v", k, m, td.bodyJson, td)
				}
			}
		}
	}

	return body, m
}

func runProxyAuth(t testing.T) {
	// curl --header "Authorization: Bearer MmU3ZWI2YzgtMDY3YS00NjM5LTg1MjEtYzcyYzc1NjU3ODEw"
	// --data "username=takumi&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU="
	// http://localhost:8080/api/2/domains/topcoder.com/proxyauth ; echo
}
