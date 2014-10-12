package auth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttp(t *testing.T) {
	wa, err := NewWebAPI("test_data.json")
	if err != nil {
		t.Fatal("Failed to create webapi: ", err)
	}

	ts := httptest.NewServer(wa.Mux)
	defer ts.Close()

	// Test Case 1 - topcoder.com pass - with some extra checks
	// curl -i --data "username=takumi&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost:8080/api/2/domains/topcoder.com/proxyauth
	u := ts.URL + "/api/2/domains/topcoder.com/proxyauth"
	pf := url.Values{}
	pf.Add("username", "takumi")
	pf.Add("password", "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=")
	res, err := http.PostForm(u, pf)
	if err != nil {
		t.Fatalf("Failed to get response from server: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("POST %v - params: %v - status code: %v expected: %v", u, pf, res.StatusCode, http.StatusOK)
	}
	if r := res.Header.Get("Content-Type"); r != "application/json" {
		t.Errorf("POST %v - params: %v - status code: %v", u, pf, r)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("POST %v - params: %v - unable to read body: %v", u, pf, err)
	}
	body := string(bytes[:])
	expected := "{\"access_granted\":true}"
	if body != expected {
		t.Errorf("POST %v - params: %v - body: %v", u, pf, body)
	}

	// Test case 2 - appirio.com pass
	// curl -i --data "username=jun&password={SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=" http://localhost:8080/api/2/domains/appirio.com/proxyauth
	u = ts.URL + "/api/2/domains/appirio.com/proxyauth"
	pf = url.Values{}
	pf.Add("username", "jun")
	pf.Add("password", "{SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=")
	res, err = http.PostForm(u, pf)
	if err != nil {
		t.Fatalf("Failed to get response from server: %v", err)
	}
	defer res.Body.Close()
	bytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("POST %v - params: %v - unable to read body: %v", u, pf, err)
	}
	body = string(bytes[:])
	expected = "{\"access_granted\":true}"
	if res.StatusCode != http.StatusOK || body != expected {
		t.Errorf("POST %v - params: %v - res: %v - body: %v", u, pf, res, body)
	}

	// Test case 3 - appirio.com password fail
	// curl -i --data "username=jun&password={SHA256}p/GZlDycaRbeggsp0wQuvZQ7yV4IzVstwEKKFhGNyGo=" http://localhost:8080/api/2/domains/appirio.com/proxyauth
	u = ts.URL + "/api/2/domains/appirio.com/proxyauth"
	pf = url.Values{}
	pf.Add("username", "jun")
	pf.Add("password", "{SHA256}p/GZlDycaRbeggsp0wQuvZQ7yV4IzVstwEKKFhGNyGo=")
	res, err = http.PostForm(u, pf)
	if err != nil {
		t.Fatalf("Failed to get response from server: %v", err)
	}
	defer res.Body.Close()
	bytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("POST %v - params: %v - unable to read body: %v", u, pf, err)
	}
	body = string(bytes[:])
	expected = "{\"access_granted\":false,\"reason\":\"denied by policy\"}"
	if res.StatusCode != http.StatusOK || body != expected {
		t.Errorf("POST %v - params: %v - res: %v - body: %v", u, pf, res, body)
	}

	// Test case 4 - appirio.com username fail"
	// curl -i --data "username=kyrra&password={SHA256}p/GZlDycaRbeggsp0wQuvZQ7yV4IzVstwEKKFhGNyGo=" http://localhost:8080/api/2/domains/appirio.com/proxyauth
	u = ts.URL + "/api/2/domains/appirio.com/proxyauth"
	pf = url.Values{}
	pf.Add("username", "kyrra")
	pf.Add("password", "{SHA256}p/GZlDycaRbeggsp0wQuvZQ7yV4IzVstwEKKFhGNyGo=")
	res, err = http.PostForm(u, pf)
	if err != nil {
		t.Fatalf("Failed to get response from server: %v", err)
	}
	defer res.Body.Close()
	bytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("POST %v - params: %v - unable to read body: %v", u, pf, err)
	}
	body = string(bytes[:])
	expected = "{\"access_granted\":false,\"reason\":\"denied by policy\"}"
	if res.StatusCode != http.StatusOK || body != expected {
		t.Errorf("POST %v - params: %v - res: %v - body: %v", u, pf, res, body)
	}

	// Test case 5 - domain fail, expect 404"
	// curl -i --data "username=takumi&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost:8080/api/2/domains/example.com/proxyauth
	u = ts.URL + "/api/2/domains/example.com/proxyauth"
	pf = url.Values{}
	pf.Add("username", "takumi")
	pf.Add("password", "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=")
	res, err = http.PostForm(u, pf)
	if err != nil {
		t.Fatalf("Failed to get response from server: %v", err)
	}
	defer res.Body.Close()
	bytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("POST %v - params: %v - unable to read body: %v", u, pf, err)
	}
	body = string(bytes[:])
	expected = ""
	if res.StatusCode != http.StatusNotFound || body != expected {
		t.Errorf("POST %v - params: %v - res: %v - body: %v", u, pf, res, body)
	}
}

func TestHttpExtra(t *testing.T) {
	wa, err := NewWebAPI("test_data.json")
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
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("POST %v - params: %v - unable to read body: %v", u, pf, err)
	}
	body := string(bytes[:])
	expected := "{\"access_granted\":false,\"reason\":\"denied by policy\"}"
	if res.StatusCode != http.StatusOK || body != expected {
		t.Errorf("POST %v - params: %v - res: %v - body: %v", u, pf, res, body)
	}

}
