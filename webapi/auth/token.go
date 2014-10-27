package auth

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

// DomainTokens is a map[domainName]Tokens.  It maps domain names to
// access_tokens.
type DomainTokens map[string]Tokens

// Tokens is a map[access_token]expireTime. It makes a generated access_token
// to when that token expires
type Tokens map[string]time.Time

// tokenResponse represents the JSON data sent to a client on successful
// access_token generation.
type tokenResponse struct {
	Token   string `json:"access_token"`
	Type    string `json:"token_type"`
	Timeout int    `json:"expires_in"`
}

// ValidateAuthHeader will look at an HTTP request for a given domain and see if the
// Authorization header is present and valid.  Returns an error if any check fails.
func (ds *datastore) ValidateAuthHeader(w http.ResponseWriter, r *http.Request, domain string) error {
	a := r.Header.Get("Authorization")
	if len(a) == 0 {
		return fmt.Errorf("auth_header_missing")
	}

	s := strings.Fields(a)
	if len(s) != 2 || s[0] != "Bearer" {
		return fmt.Errorf("auth_header_invalid")
	}
	token := s[1]

	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	t, ok := ds.tokenMap[domain][token]
	if !ok {
		return fmt.Errorf("auth_token_not_found")
	}

	if t.Before(time.Now()) {
		return fmt.Errorf("access_token_expired")
	}

	return nil
}

// SaveTokens will take all access_tokens stored in the server and write them
// to disk if the tokenFilename is specified.
func (ds *datastore) SaveTokens() error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	if len(ds.tokenFilename) == 0 {
		return fmt.Errorf("Not saving access data")
	}

	fmt.Println("Saving access_tokens to", ds.tokenFilename)

	b, err := json.Marshal(ds.tokenMap)
	if err != nil {
		return fmt.Errorf("failed to marshel access_token data: %v", err)
	}
	ioutil.WriteFile(ds.tokenFilename, b, 0644)
	return nil
}

// loadTokens will attempt to load tokens from the given file.
func (ds *datastore) loadTokens() error {
	b, err := ds.loadFile(ds.tokenFilename)
	if err != nil {
		return err
	}

	var d DomainTokens
	err = json.Unmarshal(b, &d)
	if err != nil {
		return err
	}

	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	ds.tokenMap = d

	return nil
}

// ValidateClient takes a domain and passed int client_id, client_secrent, and greant_type
// then verifies the formats are correct and looks up the client information.  An error
// is returned if any check fails.
func (ds *datastore) ValidateClient(domain, id, secret, grant string) error {
	if len(id) == 0 || len(secret) == 0 || len(grant) == 0 {
		return fmt.Errorf("invalid_request")
	}

	if grant != "client_credentials" {
		return fmt.Errorf("unsupported_grant_type")
	}

	d := ds.domainMap[domain]
	s, ok := d.Clients[id]
	if !ok || s != secret {
		return fmt.Errorf("invalid_client")
	}

	return nil
}

// generateAccessToken will generate a new UUID and add it to the tokenMap.
func (ds *datastore) generateAccessToken(domain string) string {
	u := uuid.New()
	accessToken := base64.StdEncoding.EncodeToString([]byte(u))

	t := time.Now().Add(time.Duration(ds.tokenTimeout) * time.Second)

	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.tokenMap[domain][accessToken] = t
	return accessToken
}

// startSigHandler create a goroutine to wait for SIGINT calls,
// gets the write lock then shuts down.
func (ds *datastore) startSigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for _ = range c {
			fmt.Println("shutting down server")
			err := ds.SaveTokens()
			if err != nil {
				fmt.Println(err)
			}
			ds.mutex.Lock()
			os.Exit(0)
		}
	}()
}
