package auth

import (
// ""
)

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

}
