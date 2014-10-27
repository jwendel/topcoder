About
======
The *webapi* package is an implementation of the topcoder Oath2 Golang challenge: http://community.topcoder.com/tc?module=ProjectDetail&pj=30046224

Four .go files come with this package

* **main.go** - Contains a simple main for launching the auth service.  Allows for flags/parameters to be passed to specify auth and token input file and listening address/port.
* **auth/http.go** - Sets up routes and handlers, and does parsing of input and formatting of output.
* **auth/json.go** - datastore for the application.  Handles loading the users.json file, parsing it, and provides lookup functions to see if a domain or user/password are valid.
* **auth/token.go** - Handles much of the logic for the Oauth2 generation and authorizing.
* **auth/util.go** - simple util methods, only contains encryptPassword helper.

The default service starts up ":8080", but the port can be specified to example:

   `./webapi -listen=":80"`

A different data file can be spcified as well, for example:

   `./webapi -listen=":80" -datasource="domains2.json"`

```
Usage of ./webapi:
  -datasource="domains.json": Filename to load JSON user data from
  -listen=":8080": Hostname and address to listen on
  -tokenTimeout=3600: Lifetime of auth tokens in seconds
  -tokensource="": Filename to save and load access_tokens from. Blank to bypass this feature.
```

Path
----
The auth package should be installed to:  **$GOPATH/src/topcoder.com/glc/webapi**


Design Notes
------------
I built Oauth changes into the existing auth package as it is they are a bit intertwined.

Notes:

* access_tokens can be optionally saved to disk using the "-tokensource" parameter.  If not specified, nothing is saved to disk.  The server must be killed with Ctrl-C (SIGINT) to save to disk.
* examples_test.go contains the tests cases outlined in the challenge. They are commented in the file to easily follow.  They also verify the results of the requests.
* The expiration test (case 9) check uses a token loaded from disk.
* Along with the .go files in auth, there are test files for them as well (run with 'go test').
* http.go domainRouter funciton is an internal router based on the URL.  json.go or token.go will handle the bulk of the request based on the URL.
* util.go has helper functions for more easily writing out status to the http client
* All access to the data maps are protected by a mutex (which should allow running this on multiple threads with GOMAXPROCS).

Other
-----
Licensed under BSD-style license.

No external libraries used.
