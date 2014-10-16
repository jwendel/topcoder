package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strconv"
)

// CacheRequest represents a single command sent
// to the server
type CacheRequest struct {
	conn   net.Conn
	Cmd    string
	Subcmd string
}

type server struct {
	t      net.Listener
	c      map[string]func(c *CacheRequest)
	cmdReg *regexp.Regexp
}

func NewServer(port, maxItems int) (*server, error) {
	t, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	// a := "a-zA-Z0-9\\!\\#\\$\\%\\&\\'\"\\*\\+\\-\\/\\=\\?\\^\\_\\{\\|\\}\\~\\(\\)\\<\\>\\[\\]\\:\\;\\@\\,\\."
	a := "a-zA-Z0-9"
	// r := regexp.MustCompile("^(" + a + ")+( (" + a + ")*)?$")
	r := regexp.MustCompile("^([" + a + "]+) ?([" + a + "]*)$")
	if err != nil {
		return nil, err
	}

	s := server{t, make(map[string]func(c *CacheRequest)), r}

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
		data = data[0 : len(data)-2]
		cmds := s.cmdReg.FindStringSubmatch(string(data))
		fmt.Println("data: ", data)
		fmt.Println("cmds: ", cmds)
		for _, v := range cmds {
			req.write([]byte(v))
		}

	}
}

func (c *CacheRequest) write(b []byte) {
	data := append(b, []byte("\r\n")...)
	c.conn.Write(data)
}
