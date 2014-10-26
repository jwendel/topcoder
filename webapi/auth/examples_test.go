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
	"strings"
	"testing"
	"time"
)

// tokenData is used to issue an access_token request and represents result data
type tokenData struct {
	url       string
	id        string
	secret    string
	grantType string
	// expected responses
	code        int
	contentType string
	bodyJSON    map[string]interface{}
}

// proxyData is used to issuth a proxyauth request and represents result data
type proxyData struct {
	url        string
	authHeader string
	user       string
	password   string
	// expected responses
	code        int
	contentType string
	bodyJSON    map[string]interface{}
}

// TestOath runs the 9 test cases outlined in the challenge
func TestOath(t *testing.T) {
	wa, err := NewWebAPI("test_data.json", "", 3600)
	if err != nil {
		t.Fatal("Failed to create Webapi: ", err)
	}

	ts := httptest.NewServer(wa.Mux)
	defer ts.Close()

	topTokenURL := ts.URL + "/api/2/domains/topcoder.com/oauth/access_token"
	topProxyURL := ts.URL + "/api/2/domains/topcoder.com/proxyauth"
	appTokenURL := ts.URL + "/api/2/domains/appirio.com/oauth/access_token"
	appProxyURL := ts.URL + "/api/2/domains/appirio.com/proxyauth"

	// Case 1 Success
	// domain : topcoder.com
	// Call oauth/access_token endpoint to obtain an access token
	// Success with status code 200
	// Call proxyauth endpoint with the access token obtained
	// Success with status code 200
	bodyJSON := map[string]interface{}{"token_type": "bearer", "expires_in": float64(3600), "access_token": nil}
	td := tokenData{topTokenURL, "s6BhdRkqt3", "7Fjfp0ZBr1KtDRbnfVdmIw", "client_credentials", 200, "application/json", bodyJSON}
	m := runAccessToken(t, td)
	at := "Bearer " + m["access_token"].(string)

	bodyJSON = map[string]interface{}{"access_granted": true}
	pd := proxyData{topProxyURL, at, "takumi", "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=", 200, "application/json", bodyJSON}
	runProxyAuth(t, pd)

	// Case 2 Success
	// domain : appirio.com
	// Call oauth/access_token endpoint to obtain an access token
	// Success with status code 200
	// Call proxyauth endpoint with the access token obtained
	// Success with status code 200
	bodyJSON = map[string]interface{}{"token_type": "bearer", "expires_in": float64(3600), "access_token": nil}
	td = tokenData{appTokenURL, "MDYyMDI4OD", "NzU1MTQyZWUtYzJhZC00OT", "client_credentials", 200, "application/json", bodyJSON}
	m = runAccessToken(t, td)
	at = "Bearer " + m["access_token"].(string)

	bodyJSON = map[string]interface{}{"access_granted": true}
	pd = proxyData{appProxyURL, at, "jun", "{SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=", 200, "application/json", bodyJSON}
	runProxyAuth(t, pd)

	// Case 3 Failure
	// domain : topcoder.com
	// Call oauth/access_token endpoint to obtain an access token
	// Failure with status code 400
	// error : invalid_request
	bodyJSON = map[string]interface{}{"error": "invalid_request"}
	td = tokenData{topTokenURL, "", "fake", "client_credentials", 400, "application/json", bodyJSON}
	m = runAccessToken(t, td)

	// Case 4 Failure
	// domain : appirio.com
	// Call oauth/access_token endpoint to obtain an access token
	// Failure with status code 400
	// error : invalid_client
	bodyJSON = map[string]interface{}{"error": "invalid_client"}
	td = tokenData{appTokenURL, "FakeClientID", "NzU1MTQyZWUtYzJhZC00OT", "client_credentials", 400, "application/json", bodyJSON}
	m = runAccessToken(t, td)

	// Case 5 Failure
	// domain : appirio.com
	// Call oauth/access_token endpoint to obtain an access token
	// Failure with status code 400
	// error : unsupported_grant_type
	bodyJSON = map[string]interface{}{"error": "unsupported_grant_type"}
	td = tokenData{appTokenURL, "FakeClientID", "NzU1MTQyZWUtYzJhZC00OT", "authorization_code", 400, "application/json", bodyJSON}
	m = runAccessToken(t, td)

	// Case 6 Failure
	// Call oauth/access_token endpoint to obtain an access token
	// Failure with status code 404
	bodyJSON = map[string]interface{}{}
	td = tokenData{ts.URL + "/api/2/domains/google.com/oauth/access_token", "MDYyMDI4OD", "NzU1MTQyZWUtYzJhZC00OT", "client_credentials", 404, "text/plain", bodyJSON}
	m = runAccessToken(t, td)

	// Case 7 Failure
	// Call proxyauth endpoint with no Authorization header.
	// Failure with status code 400
	bodyJSON = map[string]interface{}{"error": "auth_header_missing"}
	pd = proxyData{appProxyURL, "", "jun", "{SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=", 400, "application/json", bodyJSON}
	runProxyAuth(t, pd)

	// Case 8 Failure
	// Call proxyauth endpoint with an invalid access token
	// Failure with status code 400
	bodyJSON = map[string]interface{}{"error": "auth_token_not_found"}
	pd = proxyData{appProxyURL, "Bearer fakeToken", "jun", "{SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=", 400, "application/json", bodyJSON}
	runProxyAuth(t, pd)

	// Case 9 Failure
	// Call proxyauth endpoint with an expired access token
	// Failure with status code 400
	wa2, err := NewWebAPI("test_data.json", "", 1)
	if err != nil {
		t.Fatal("Failed to create Webapi: ", err)
	}
	ts2 := httptest.NewServer(wa2.Mux)
	defer ts2.Close()

	appTokenURL = ts2.URL + "/api/2/domains/appirio.com/oauth/access_token"
	appProxyURL = ts2.URL + "/api/2/domains/appirio.com/proxyauth"

	bodyJSON = map[string]interface{}{"token_type": "bearer", "expires_in": float64(1), "access_token": nil}
	td = tokenData{appTokenURL, "MDYyMDI4OD", "NzU1MTQyZWUtYzJhZC00OT", "client_credentials", 200, "application/json", bodyJSON}
	m = runAccessToken(t, td)
	at = "Bearer " + m["access_token"].(string)

	time.Sleep(2 * time.Second)
	bodyJSON = map[string]interface{}{"error": "access_token_expired"}
	pd = proxyData{appProxyURL, at, "jun", "{SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=", 400, "application/json", bodyJSON}
	runProxyAuth(t, pd)

}

// runAccessToken crafts a access_token request, sends it and validates the result
func runAccessToken(t *testing.T, td tokenData) map[string]interface{} {
	pf := url.Values{}
	pf.Add("client_id", td.id)
	pf.Add("client_secret", td.secret)
	pf.Add("grant_type", td.grantType)

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

	var m map[string]interface{}
	// Lets try to decode the json and compare to the passed in bodyJSON
	if ct == "application/json" {
		m = checkJSONBody(t, bytes, td.bodyJSON)
	}

	return m
}

// runProxyAuth crafts a proxyauth request, sends it and validates the result
func runProxyAuth(t *testing.T, pd proxyData) map[string]interface{} {
	pf := url.Values{}
	pf.Add("username", pd.user)
	pf.Add("password", pd.password)

	dc := http.DefaultClient
	req, err := http.NewRequest("POST", pd.url, strings.NewReader(pf.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v\nproxyData: %v", err, pd)
	}

	req.Header.Set("Authorization", pd.authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := dc.Do(req)
	if err != nil {
		t.Fatalf("Failed to get response from server: %v\nproxyData: %v", err, pd)
	}
	if res.StatusCode != pd.code {
		t.Errorf("StatusCode missmatch.  Got %v  Expected %v\nproxyData: %v", res.StatusCode, pd.code, pd)
	}
	ct := res.Header.Get("Content-Type")
	if ct != pd.contentType {
		t.Errorf("Content-Type missmatch.  Got %v  Expected %v\nproxyData: %v", ct, pd.contentType, pd)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unable to read body: %v\nproxyData: %v", err, pd)
	}

	var m map[string]interface{}
	// Lets try to decode the json and compare to the passed in bodyJSON
	if ct == "application/json" {
		m = checkJSONBody(t, bytes, pd.bodyJSON)
	}

	return m
}

// checkJSONBody will validate data returned in the body of a request.
// It unmarshels JSON data into a map, and compares it to the bodyJSON
// passed in.  It returns the results from the body as a map.
func checkJSONBody(t *testing.T, bytes []byte, bodyJSON map[string]interface{}) map[string]interface{} {
	var m map[string]interface{}
	var f interface{}
	err := json.Unmarshal(bytes, &f)
	if err != nil {
		t.Errorf("Failed to parse json: %v\nbody: %v", err, string(bytes))
	} else {
		m = f.(map[string]interface{})
		if len(m) != len(bodyJSON) {
			t.Errorf("Reponse doesn't have the same number of keys.  Got: %v  Expected %v", m, bodyJSON)
		}

		// Loop over the keys/values expected and comapre to what we got
		for k, tdv := range bodyJSON {
			mv, ok := m[k]
			if !ok {
				t.Errorf("Reponse missing key %v.  Got: %v  Expected %v", k, m, bodyJSON)
				continue
			}
			// t.Logf("testing key %v  value: %v", k, tdv)
			// don't check values if expected data is nil
			if tdv == nil {
				continue
			}
			switch tdvv := tdv.(type) {
			case string:
				mvv, ok := mv.(string)
				if !ok || tdvv != mvv {
					t.Errorf("Types/Values don't match for key %v.  Got: %v  Expected %v", k, m, bodyJSON)
				}
			case float64:
				mvv, ok := mv.(float64)
				if !ok || tdvv != mvv {
					t.Errorf("Types/Values don't match for key %v.  Got: %v  Expected %v", k, m, bodyJSON)
				}
			case bool:
				mvv, ok := mv.(bool)
				if !ok || tdvv != mvv {
					t.Errorf("Types/Values don't match for key %v.  Got: %v  Expected %v", k, m, bodyJSON)
				}
			default:
				t.Errorf("Unhandled type for key %v.  Got: %v  Expected %v", k, m, bodyJSON)
			}
		}
	}

	return m
}
