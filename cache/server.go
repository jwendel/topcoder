package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var validChars *regexp.Regexp

// CacheRequest represents a single command sent
// to the server
type CacheRequest struct {
	conn   net.Conn
	Cmd    string
	Subcmd []string
	s      *server
	reader *bufio.Reader
}

type server struct {
	l          net.Listener
	cmds       map[string]func(c *CacheRequest)
	cache      map[string]string
	cacheMutex sync.RWMutex
}

func NewServer(port, maxItems int) (*server, error) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	// abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'\"*+-/=?^_{|}~()<>[]:;@,.
	a := "a-zA-Z0-9!#$%&'\"*+\\-/\\\\=?^_{|}~()<>\\[\\]:;@,. "
	validChars = regexp.MustCompile("^[" + a + "]+$")

	s := server{}
	s.l = l
	s.cmds = make(map[string]func(c *CacheRequest))
	s.cache = make(map[string]string)

	return &s, nil
}

func (s *server) Serve() error {

	for {
		conn, err := s.l.Accept()
		if err != nil {
			return err
		}
		go s.handle(conn)
	}
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

func (s *server) handle(conn net.Conn) {

	// Do we want to enable read timeout?
	// conn.SetReadDeadline(t)

	req := CacheRequest{}
	req.reader = bufio.NewReader(conn)
	req.conn = conn
	req.s = s

	for {
		data, err := req.readln()
		if err != nil {
			fmt.Println("error reading data: ", err)
			return
		}

		// Check that it was \r\n
		if len(data) < 2 || data[len(data)-2] != '\r' {
			req.writeStr("ERROR invalid input")
			continue
		}

		// Trim \r\n
		data = data[:len(data)-2]
		// Validate string data
		if len(data) == 0 {
			continue
		}

		s.processInput(string(data), req)
	}
}

func (s *server) processInput(input string, c CacheRequest) {
	if !validChars.MatchString(input) {
		c.writeStr("ERROR invalid input characters")
		return
	}

	cmds := strings.Split(input, " ")
	if len(cmds) == 0 { // TODO: remove?
		return
	}

	c.Cmd = cmds[0]
	c.Subcmd = cmds[1:]

	f, ok := s.cmds[c.Cmd]
	if !ok {
		c.writeStr("ERROR unknown command")
		return
	}

	f(&c)
}

func (c *CacheRequest) readln() ([]byte, error) {
	data, err := c.reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("error reading data: ", err)
		c.conn.Close()
		return nil, err
	}
	return data, nil
}

func (c *CacheRequest) writeStr(s string) {
	data := append([]byte(s), []byte("\r\n")...)
	c.conn.Write(data)
}

func (c *CacheRequest) write(b []byte) {
	data := append(b, []byte("\r\n")...)
	c.conn.Write(data)
}
