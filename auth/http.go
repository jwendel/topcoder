package auth

import (
	"fmt"
	"html"
	"net/http"
)

func Start(listenAddr, jsonFilename string) error {
	err := Store.Init(jsonFilename)
	if err != nil {
		return err
	}

	http.HandleFunc("/", asdf)
	http.ListenAndServe(listenAddr, nil)
	return nil
}

func asdf(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
