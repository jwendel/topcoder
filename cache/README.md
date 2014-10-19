About Cache
===========
This is an implementation of topcoder challenge: http://www.topcoder.com/challenge-details/30046225/?type=develop

This is a key/value cache server that will take input over TCP following telnet style input/output.  

Notes
-----
* Command input has all whitespace trimmed (beginning/trailing spaces are ignored, and multiple spaces between parameters).
* Server supports multiple connections at once.

Running
-------
* cd $GOPATH/src/bitbucket.org/kyrra/sandbox/cache
* go build
* ./cache --help
   `Usage of ./cache:
  -items=65535: Maximum number of items to cache
  -port=11212: Port the server listens on`
* ./cache

Path
----
The cache package should be installed to:  **$GOPATH/src/bitbucket.org/kyrra/sandbox/cache**
