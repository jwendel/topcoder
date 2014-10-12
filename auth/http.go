package auth

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type response struct {
	Access_granted bool   `json:"access_granted"`
	Reason         string `json:"reason,omitempty"`
}

var regex *regexp.Regexp

func Start(listenAddr, jsonFilename string) error {
	err := Store.Init(jsonFilename)
	if err != nil {
		return err
	}

	regex, err = regexp.Compile("^/api/2/domains/(.+)/proxyauth$")

	// /api/2/domains/{domain name}/proxyauth

	http.HandleFunc("/api/2/domains/", domainAuth)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(listenAddr, nil)
	return nil
}

func domainAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	matches := regex.FindStringSubmatch(r.URL.Path)
	if matches == nil || len(matches) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	domain := matches[1]
	ok := Store.DomainExists(domain)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	ok = Store.UserPasswordValid(domain, username, password)
	res := response{ok, ""}
	if !ok {
		s := "denied by policy"
		res.Reason = s
	}
	js, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(js)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
