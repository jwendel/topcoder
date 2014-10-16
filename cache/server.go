package main

import (
	"bufio"
	"fmt"
	"net"
)

type CacheRequest struct {
	conn net.Conn
}

type server struct {
	t net.Listener
	c map[string]func(c *CacheRequest)
}

func NewServer() (*server, error) {
	t, err := net.Listen("tcp", ":9000")
	if err != nil {
		return nil, err
	}
	s := server{t, make(map[string]func(c *CacheRequest))}

	return &s, nil
}

func (s *server) Serve() error {

	for {
		conn, err := s.t.Accept()
		if err != nil {
			return err
		}
		go s.handle(conn)
	}
}

func (s *server) handle(conn net.Conn) {

	// Do we want to enable read timeout?
	// conn.SetReadDeadline(t)

	reader := bufio.NewReader(conn)
	req := CacheRequest{}
	req.conn = conn

	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("error reading data: ", err)
			conn.Close()
			return
		}

		// Check that it was \r\n
		if len(data) < 2 || data[len(data)-2] != '\r' {
			conn.Write([]byte("ERROR invalid input\r\n"))
			continue
		}

		req.write(data)
	}
}

func (c *CacheRequest) write(b []byte) {
	data := append(b, []byte("\r\n")...)
	c.conn.Write(data)
}
