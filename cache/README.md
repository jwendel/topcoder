About Cache
===========
This is an implementation of topcoder challenge: http://www.topcoder.com/challenge-details/30046225/?type=develop

This is a key/value cache server that will take input over TCP following telnet style input/output.  

Notes
-----
* Command input has all whitespace trimmed (beginning/trailing spaces are ignored, and multiple spaces between parameters).
* Server supports multiple connections at once.
* The server will disconnect clients that send 64kb of data without a newline (a property of using bufio.Scanner)
* Mutex is used when accessing the cache.  Almost all locks are write locks (not RLock) as we need to update the dataStats with 4 of the commands.
* examples_test.go has a number of extra tests added to it to verify behavior.
* -addr param is useful for binding only to localhost for unit tests
* server.go and request.go could easily be pulled into their own package if this needed to be a reusable module.  main.go and cmds.go use only public interfaces when working with the server.

Extensibility
-------------
Adding new commands is easy.

1. In main.go registerHandlers function, add a call to register your new function
2. In cmds.go add a function to match the interface of server.AddHandler: func(c *CacheRequest)

The passed in CacheRequest contains everything that a helper should need to process their request.  Be sure to use the mutex in the CacheRequest if you will be reading or writing to the dataCache.


Running
-------
* cd $GOPATH/src/topcoder.com/kyrra/scs/
* go build
* ./scs --help

```
Usage of ./scs:
  -addr="": IP address the server binds to
  -items=65535: Maximum number of items to cache
  -port=11212: Port the server listens on
```

* ./scs

Path
----
The cache package should be installed to:  **$GOPATH/src/topcoder.com/kyrra/scs/**
