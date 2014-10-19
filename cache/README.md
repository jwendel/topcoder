About Cache
===========
This is an implementation of topcoder challenge: http://www.topcoder.com/challenge-details/30046225/?type=develop

This is a key/value cache server that will take input over TCP following telnet style input/output.  

Notes
-----
* Command input has all whitespace trimmed (beginning/trailing spaces are ignored, and multiple spaces between parameters).
* Server supports multiple connections at once.
* Mutex is used when accessing the cache.  Almost all locks are write locks (not RLock) as we need to update the dataStats with 4 of the commands.

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
  -items=65535: Maximum number of items to cache
  -port=11212: Port the server listens on
```

* ./scs

Path
----
The cache package should be installed to:  **$GOPATH/src/topcoder.com/kyrra/scs/**
