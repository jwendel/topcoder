package auth

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type webapi struct {
	Mux         *http.ServeMux
	domainRegex *regexp.Regexp
	store       datastore
}

type response struct {
	Access_granted bool   `json:"access_granted"`
	Reason         string `json:"reason,omitempty"`
}

// Serve creates a webapi and starts the http server
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
	domainRegex, err := regexp.Compile("^/api/2/domains/(.+)/proxyauth$")
	if err != nil {
		return nil, err
	}

	var store datastore
	err = store.Init(jsonFilename)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	wa := webapi{mux, domainRegex, store}

	wa.Mux.HandleFunc("/api/2/domains/", wa.domainAuth)
	wa.Mux.HandleFunc("/", wa.defaultHandler)

	return &wa, nil
}

// domainAuth handles domain authentiation based on data in store
func (wa *webapi) domainAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	matches := wa.domainRegex.FindStringSubmatch(r.URL.Path)
	if matches == nil || len(matches) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// matches[1] is the domain to lookup
	domain := matches[1]
	ok := wa.store.DomainExists(domain)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
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

// defaultHandler returns 404 for the all paths not explicity specified
func (wa *webapi) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
