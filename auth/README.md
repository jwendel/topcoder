About
======
The *auth* package is an implementation of the topcoder challenge: http://www.topcoder.com/challenge-details/30046011/?type=develop

Four .go files come with this package

* **example/main.go** - Contains a simple main for launching the auth service.  Allows for flags/parameters to be passed to specify input file and listening address/port.
* **http.go** - Sets up routes and handlers, and does parsing of input and formatting of output.
* **json.go** - datastore for the application.  Handles loading the users.json file, parsing it, and provides lookup functions to see if a domain or user/password are valid.
* **util.go** - simple util methods, only contains encryptPassword helper.

Along with the .go files in auth, there are test files for them as well (run with 'go test').  http_test.go contains the tests outlined in the challenge.  There is an example/test.sh that executes the same test cases with curl but does not validation on the returned data.

The default service starts up ":8080", but the port can be specified to example:

   `./example -listen=":80"`

A different data file can be spcified as well, for example:

   `./example -listen=":80" -datasource="users2.json"`

Path
----
The auth package should be installed to:  $GOPATH/src/bitbucket.org/kyrra/sandbox/auth

Design Notes
------------
For parsing of the domain within the URL passed to the server, I used a simple regex to match it.  If this API was expanded further, it may be better to use a different mux implementation that allows for wildcards in paths.

The datastore (in json.go) has two interesting designs behind it.

* Instead of reading the file on every request we store it in memory.  To handle if the underlying data file changes, a goroutine runs (every 3 seconds) and checks if the modified timestamp on the file has changed.  If it does, it locks the datastore and reloads the source file.
* The data from the loaded json file is stuck into a map of maps (map[DomainName]map[UserName]HashedPassword).  This makes lookups easy and detecting duplicate entries within the input file.  The downside is that maps are slow for small inputs and also use lots of memory.

Other
----- 
Licensed under BSD-style license.

No external libraries used.