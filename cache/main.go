// Copyright 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Implementation of a simple cache server for TopCoder.  See url for details:
http://www.topcoder.com/challenge-details/30046225/?type=develop

This program takes a subset of ascii data input can can be uased to trigger
command handlers to further work on the input.
*/
package main

import (
	"flag"
	"fmt"
)

// main Entrypoint to the application.  Defines command line flags,
// creates cache server, registers handlers, and starts serving.
func main() {
	p := flag.Int("port", 11212, "Port the server listens on")
	i := flag.Int("items", 65535, "Maximum number of items to cache")
	flag.Parse()

	s, err := NewServer(*p, *i)
	if err != nil {
		fmt.Println("failed to create server: ", err)
		return
	}

	err = registerHandlers(s)
	if err != nil {
		fmt.Println("failed to register handlers: ", err)
		return
	}

	fmt.Println("ready to accept cache requests")
	s.Serve()
}

// registerHandlers will associate command handlers to their
// given command available via the server.
func registerHandlers(s *server) error {
	err := s.AddHandler("set", cmdSet)
	if err != nil {
		return err
	}
	err = s.AddHandler("get", cmdGet)
	if err != nil {
		return err
	}
	err = s.AddHandler("delete", cmdDelete)
	if err != nil {
		return err
	}
	err = s.AddHandler("stats", cmdStats)
	if err != nil {
		return err
	}
	err = s.AddHandler("quit", cmdQuit)
	if err != nil {
		return err
	}
	return nil
}
