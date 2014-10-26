package auth

import (
	"encoding/json"
	// "io/ioutil"
	"code.google.com/p/go-uuid/uuid"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type DomainTokens map[string]Tokens

type Tokens map[string]time.Time

type tokenResponse struct {
	Token   string `json:"access_token"`
	Type    string `json:"token_type"`
	Timeout int    `json:"expires_in"`
}

func (wa *webapi) ValidateAuthHeader(w http.ResponseWriter, r *http.Request, domain string) error {
	a := r.Header.Get("Authorization")
	if len(a) == 0 {
		return fmt.Errorf("invalid_request") // TODO
	}

	s := strings.Fields(a)
	if len(s) != 2 || s[0] != "Bearer" {
		return fmt.Errorf("invalid_request") // TODO
	}

	token := s[1]

	t, ok := wa.store.tokenMap[domain][token]
	if !ok {
		return fmt.Errorf("invalid_auth_token") // TODO
	}

	if t.Before(time.Now()) {
		return fmt.Errorf("access_token_expired") // TODO
	}

	return nil
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

	r.ParseForm()
	id := r.Form.Get("client_id")
	secret := r.Form.Get("client_secret")
	grant := r.Form.Get("grant_type")

	err := wa.ValidateClient(domain, id, secret, grant)
	if err != nil {
		badRequestHandler(w, r, err.Error())
		return
	}

	// access_token request is valid, generate token
	token := wa.generateAccessToken(domain)
	t := tokenResponse{token, "bearer", wa.tokenTimeout}
	b, err := json.Marshal(t)
	if err != nil {
		internalErrorHandler(w, r)
		return
	}

	successHandler(w, r, b)
}

func (wa *webapi) ValidateClient(domain, id, secret, grant string) error {
	if len(id) == 0 || len(secret) == 0 || len(grant) == 0 {
		return fmt.Errorf("invalid_request")
	}

	if grant != "client_credentials" {
		return fmt.Errorf("unsupported_grant_type")
	}

	d := wa.store.domainMap[domain]
	s, ok := d.Clients[id]
	if !ok || s != secret {
		return fmt.Errorf("invalid_client")
	}

	return nil
}

func (wa *webapi) generateAccessToken(domain string) string {
	u := uuid.New()
	accessToken := base64.StdEncoding.EncodeToString([]byte(u))

	t := time.Now().Add(time.Duration(wa.tokenTimeout) * time.Second)
	wa.store.tokenMap[domain][accessToken] = t
	return accessToken
}
