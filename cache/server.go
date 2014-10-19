// Copyright 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
)

const (
	MAX_KEY_SIZE  = 250
	MAX_DATA_SIZE = 8192
)

// validChars is used to make sure there are only supports ascii character
// in a given string
var validChars *regexp.Regexp

// server is the TCP server for the cache application.
// It maintains the TCP listener, the map of commands to
// their given function, and keeps the dataCache to pass
// to each new connection.
type server struct {
	l    net.Listener
	cmds map[string]func(c *CacheRequest)
	c    dataCache
}

// NewServer initializes everything needed to handle new
// connections to the cache server.  It attempts the address
// and port to listen on, along with the max items the cache
// can store.
func NewServer(addr string, port, maxItems int) (*server, error) {
	l, err := net.Listen("tcp", addr+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	a := "a-zA-Z0-9!#$%&'\"*+\\-/\\\\=?^_{|}~()<>\\[\\]:;@,. "
	validChars = regexp.MustCompile("^[" + a + "]+$")

	s := server{}
	s.l = l
	s.cmds = make(map[string]func(c *CacheRequest))
	s.c.Cache = make(map[string]string)
	s.c.maxItems = maxItems
	s.c.Stats = &dataStats{}

	return &s, nil
}

// Server will start accepting new connections and
// pass each new connection onto its own goroutine.
func (s *server) Serve() error {

	s.startSigHandler()

	for {
		conn, err := s.l.Accept()
		if err != nil {
			return err
		}
		go s.handle(conn)
	}
}

// Close will shut down the listening socket.  Any open
// connections remain open.
func (s *server) Close() {
	s.l.Close()
}

// AddHandler adds a new command handler for the server to call when
// name is matched to user input
func (s *server) AddHandler(name string, f func(c *CacheRequest)) error {
	_, ok := s.cmds[name]
	if ok {
		return fmt.Errorf("Command '%v' is already registered", name)
	}

	s.cmds[name] = f
	return nil
}

// handle take a connection and reads data from it,
// processing the requests
func (s *server) handle(conn net.Conn) {

	req := CacheRequest{}
	// req.reader = bufio.NewReader(conn)
	req.scanner = bufio.NewScanner(conn)
	req.scanner.Split(scanLines)
	req.Conn = conn
	req.C = s.c

	for {
		data, err := req.Readln()
		if err != nil {
			req.Conn.Close()
			return
		}

		input, err := req.ValidateInput(data)
		if err != nil {
			req.WriteStr(err.Error())
			continue
		}
		if len(input) == 0 {
			continue
		}

		s.processInput(string(data), &req)
	}
}

// scanLines is a copy of bufio.Scanner.ScanLines.  This version removed the
// call to dropCR, as we want the CR there still to validate it does indeed
// end with \r\n.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// processInput takes a string, splits it by space, then calls
// the appropriate cmd function to handle the request
func (s *server) processInput(input string, c *CacheRequest) {
	cmds := strings.Fields(input)
	if len(cmds) == 0 {
		return
	}

	c.Cmd = cmds[0]
	c.Subcmd = cmds[1:]

	f, ok := s.cmds[c.Cmd]
	if !ok {
		c.WriteStr("ERROR unknown command")
		return
	}

	f(c)
}

// startSigHandler create a goroutine to wait for SIGINT calls,
// gets the write lock then shuts down.
func (s *server) startSigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for _ = range c {
			s.c.CacheMutex.Lock()
			fmt.Println("shutting down server")
			s.l.Close()
			os.Exit(0)
		}
	}()
}
