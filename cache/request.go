// Copyright 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package main

package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

// CacheRequest represents a single command sent
// to the server
type CacheRequest struct {
	C      dataCache
	Cmd    string
	Subcmd []string
	conn   net.Conn
	reader *bufio.Reader
}

type dataCache struct {
	Cache      map[string]string
	CacheMutex sync.RWMutex
	Stats      *dataStats
	maxItems   int
}

type dataStats struct {
	get       int
	set       int
	getHits   int
	getMisses int
	delHits   int
	delMisses int
}

func (c *CacheRequest) ValidateInput(data []byte) (string, error) {
	// Check that it was \r\n
	if len(data) < 2 || data[len(data)-2] != '\r' {
		return "", fmt.Errorf("ERROR invalid input")
	}

	// Trim \r\n
	data = data[:len(data)-2]
	// Validate string data
	if len(data) == 0 {
		return "", nil
	}

	input := string(data)
	if !validChars.MatchString(input) {
		return "", fmt.Errorf("ERROR invalid input characters")
	}

	return input, nil
}

func (c *CacheRequest) Readln() ([]byte, error) {
	data, err := c.reader.ReadBytes('\n')
	if err != nil {
		c.conn.Close()
		return nil, err
	}
	return data, nil
}

func (c *CacheRequest) WriteStr(s string) {
	data := append([]byte(s), []byte("\r\n")...)
	c.conn.Write(data)
}

func (c *CacheRequest) Write(b []byte) {
	data := append(b, []byte("\r\n")...)
	c.conn.Write(data)
}
