package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttp(t *testing.T) {
	wa, err := NewWebAPI("test_data.json")
	if err != nil {
		t.Fatal("Failed to create webapi: ", err)
	}

	ts := httptest.NewServer(wa.Mux)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("Failed to get response from server: %v", err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("GET / - status code: %v expected: %v", res.StatusCode, http.StatusNotFound)
	}

}
