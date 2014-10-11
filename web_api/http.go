package main

import (
	"fmt"
	"html"
	"net/http"
)

func main() {
	err := store.Load("users.json")
	if err != nil {
		fmt.Println("Error loading user data: ", err)
		return
	}

	fmt.Println(store.domainMap)

	http.HandleFunc("/", asdf)
	// http.ListenAndServe(":8080", nil)
}

func asdf(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
