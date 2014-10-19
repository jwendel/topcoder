package main

import (
	"bufio"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestExamples runs the examples outlined in the topcoder challenge
func TestExamples(t *testing.T) {

	s, err := NewServer(11212, 65535)
	if err != nil {
		t.Errorf("failed to create server: %v", err)
		return
	}

	err = registerHandlers(s)
	if err != nil {
		t.Errorf("failed to register handlers: %v", err)
		return
	}
	go s.Serve()

	n, err := net.Dial("tcp", "127.0.0.1:11212")
	if err != nil {
		t.Errorf("unable to connect to server: %v", err)
		return
	}

	// All these tests should complete almost immediately
	// Set a timeout incase something goes wrong
	n.SetDeadline(time.Now().Add(3 * time.Second))

	b := bufio.NewReader(n)

	// set sushi
	// delicious
	// STORED
	n.Write([]byte("set sushi\r\ndelicious\r\n"))
	r, err := b.ReadString('\n')
	if err != nil {
		t.Errorf("set sushi read error: %v", err)
	}
	if r != "STORED\r\n" {
		t.Errorf("set sushi failed, got '%v'", r)
	}

	// set topcoder
	// fun
	// STORED
	n.Write([]byte("set topcoder\r\nfun\r\n"))
	r, err = b.ReadString('\n')
	if err != nil {
		t.Errorf("set topcoder read error: %v", err)
	}
	if r != "STORED\r\n" {
		t.Errorf("set topcoder failed, got '%v'", r)
	}

	// get sushi topcoder
	// VALUE sushi
	// delicious
	// VALUE topcoder
	// fun
	// END
	n.Write([]byte("get sushi topcoder\r\n"))
	r, err = b.ReadString('\n')
	if r != "VALUE sushi\r\n" {
		t.Errorf("get sushi/topcoder fail, expected 'VALUE sushi', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "delicious\r\n" {
		t.Errorf("get sushi/topcoder fail, expected 'delicious', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "VALUE topcoder\r\n" {
		t.Errorf("get sushi/topcoder fail, expected 'VALUE topcoder', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "fun\r\n" {
		t.Errorf("get sushi/topcoder fail, expected 'fun', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "END\r\n" {
		t.Errorf("get sushi/topcoder fail, expected 'END', got '%v'", r)
	}

	// get sushi
	// VALUE sushi
	// delicious
	// END
	n.Write([]byte("get sushi\r\n"))
	r, err = b.ReadString('\n')
	if r != "VALUE sushi\r\n" {
		t.Errorf("get sushi fail, expected 'VALUE sushi', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "delicious\r\n" {
		t.Errorf("get sushi fail, expected 'delicious', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "END\r\n" {
		t.Errorf("get sushi fail, expected 'END', got '%v'", r)
	}

	// delete sushi
	// DELETED
	n.Write([]byte("delete sushi\r\n"))
	r, err = b.ReadString('\n')
	if r != "DELETED\r\n" {
		t.Errorf("delete sushi fail, expected 'DELETED', got '%v'", r)
	}

	// get sushi
	// END
	n.Write([]byte("get sushi\r\n"))
	r, err = b.ReadString('\n')
	if r != "END\r\n" {
		t.Errorf("get sushi fail, expected 'END', got '%v'", r)
	}

	// get topcoder
	// VALUE topcoder
	// fun
	// END
	n.Write([]byte("get topcoder\r\n"))
	r, err = b.ReadString('\n')
	if r != "VALUE topcoder\r\n" {
		t.Errorf("get topcoder fail, expected 'VALUE topcoder', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "fun\r\n" {
		t.Errorf("get topcoder fail, expected 'fun', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "END\r\n" {
		t.Errorf("get topcoder fail, expected 'END', got '%v'", r)
	}

	// get topcoder sushi
	// VALUE topcoder
	// fun
	// END
	n.Write([]byte("get topcoder sushi\r\n"))
	r, err = b.ReadString('\n')
	if r != "VALUE topcoder\r\n" {
		t.Errorf("get topcoder/sushi fail, expected 'VALUE topcoder', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "fun\r\n" {
		t.Errorf("get topcoder/sushi fail, expected 'fun', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "END\r\n" {
		t.Errorf("get topcoder/sushi fail, expected 'END', got '%v'", r)
	}

	// delete sushi
	// NOT_FOUND
	n.Write([]byte("delete sushi\r\n"))
	r, err = b.ReadString('\n')
	if r != "NOT_FOUND\r\n" {
		t.Errorf("delete sushi fail, expected 'NOT_FOUND', got '%v'", r)
	}

	// stats
	// cmd_get 7
	// cmd_set 2
	// get_hits 5
	// get_misses 2
	// delete_hits 1
	// delete_misses 1
	// curr_items 1
	// limit_items 65535
	// END
	n.Write([]byte("stats\r\n"))
	r, err = b.ReadString('\n')
	if r != "cmd_get 7\r\n" {
		t.Errorf("stats fail, expected 'cmd_get 7', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "cmd_set 2\r\n" {
		t.Errorf("stats fail, expected 'cmd_set 2', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "get_hits 5\r\n" {
		t.Errorf("stats fail, expected 'get_hits 5', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "get_misses 2\r\n" {
		t.Errorf("stats fail, expected 'get_misses 2', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "delete_hits 1\r\n" {
		t.Errorf("stats fail, expected 'delete_hits 1', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "delete_misses 1\r\n" {
		t.Errorf("stats fail, expected 'delete_misses 1', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "curr_items 1\r\n" {
		t.Errorf("stats fail, expected 'curr_items 1', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "limit_items 65535\r\n" {
		t.Errorf("stats fail, expected 'limit_items 65535', got '%v'", r)
	}
	r, err = b.ReadString('\n')
	if r != "END\r\n" {
		t.Errorf("stats fail, expected 'END', got '%v'", r)
	}

	// quit
	n.Write([]byte("quit\r\n"))
	r, err = b.ReadString('\n')
	if err == nil {
		t.Errorf("quit fail, expected connection to close, got %v", r)
	}

	s.Close()

}

// TestRegex makes sure the regex passes all the characters expected
func TestRegex(t *testing.T) {

	s, err := NewServer(11212, 5)
	if err != nil {
		t.Errorf("failed to create server: %v", err)
		return
	}

	ok := validChars.MatchString("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'\"*+-/=?^_{|}~()<>[]:;@,.")
	if !ok {
		t.Errorf("Passed all valid characters and it failed")
	}

	ok = validChars.MatchString("\r\n")
	if ok {
		t.Errorf("Passed in \\r\\n and it passed.  Expected to fail.")
	}

	s.Close()
}

// TestMultiConnect creates 3 connections to the server
// and runs 'set' calls from all 3 in goroutines.
func TestMultiConnect(t *testing.T) {
	// Server setup
	s, err := NewServer(11212, 65535)
	if err != nil {
		t.Errorf("failed to create server: %v", err)
		return
	}

	err = registerHandlers(s)
	if err != nil {
		t.Errorf("failed to register handlers: %v", err)
		return
	}
	go s.Serve()

	// connection setup
	n1, err := net.Dial("tcp", "127.0.0.1:11212")
	if err != nil {
		t.Errorf("unable to connect to server: %v", err)
		return
	}
	n1.SetDeadline(time.Now().Add(5 * time.Second))

	n2, err := net.Dial("tcp", "127.0.0.1:11212")
	if err != nil {
		t.Errorf("unable to connect to server: %v", err)
		return
	}
	n2.SetDeadline(time.Now().Add(5 * time.Second))

	n3, err := net.Dial("tcp", "127.0.0.1:11212")
	if err != nil {
		t.Errorf("unable to connect to server: %v", err)
		return
	}
	n3.SetDeadline(time.Now().Add(5 * time.Second))

	b1 := bufio.NewReader(n1)
	b2 := bufio.NewReader(n2)
	b3 := bufio.NewReader(n3)

	// Start a goroutine for each connection
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			n1.Write([]byte("set n1-" + strconv.Itoa(i) + "\r\ntopcoder rules 1\r\n"))
			r, err := b1.ReadString('\n')
			if err != nil {
				t.Errorf("set n1-%v: %v", i, err)
			}
			if r != "STORED\r\n" {
				t.Errorf("set n1-%v failed, got '%v'", i, r)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			n2.Write([]byte("set n2-" + strconv.Itoa(i) + "\r\ntopcoder rules 2\r\n"))
			r, err := b2.ReadString('\n')
			if err != nil {
				t.Errorf("set n2-%v: %v", i, err)
			}
			if r != "STORED\r\n" {
				t.Errorf("set n2-%v failed, got '%v'", i, r)
			}
		}
		n2.Write([]byte("quit\r\n"))
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			n3.Write([]byte("set n3-" + strconv.Itoa(i) + "\r\ntopcoder rules 3\r\n"))
			r, err := b3.ReadString('\n')
			if err != nil {
				t.Errorf("set n3-%v: %v", i, err)
			}
			if r != "STORED\r\n" {
				t.Errorf("set n3-%v failed, got '%v'", i, r)
			}
		}
		n3.Write([]byte("quit\r\n"))
	}()

	wg.Wait()

	// verify the stats output after all 3 run
	n1.Write([]byte("stats\r\n"))
	r, err := b1.ReadString('\n')
	if r != "cmd_get 0\r\n" {
		t.Errorf("multi-connect stat error, got '%v'", r)
	}
	r, err = b1.ReadString('\n')
	if r != "cmd_set 3000\r\n" {
		t.Errorf("multi-connect stat error, got '%v'", r)
	}
	r, err = b1.ReadString('\n')
	if r != "get_hits 0\r\n" {
		t.Errorf("multi-connect stat error, got '%v'", r)
	}
	r, err = b1.ReadString('\n')
	if r != "get_misses 0\r\n" {
		t.Errorf("multi-connect stat error, got '%v'", r)
	}
	r, err = b1.ReadString('\n')
	if r != "delete_hits 0\r\n" {
		t.Errorf("multi-connect stat error, got '%v'", r)
	}
	r, err = b1.ReadString('\n')
	if r != "delete_misses 0\r\n" {
		t.Errorf("multi-connect stat error, got '%v'", r)
	}
	r, err = b1.ReadString('\n')
	if r != "curr_items 3000\r\n" {
		t.Errorf("multi-connect stat error, got '%v'", r)
	}
	r, err = b1.ReadString('\n')
	if r != "limit_items 65535\r\n" {
		t.Errorf("multi-connect stat error, got '%v'", r)
	}
	r, err = b1.ReadString('\n')
	if r != "END\r\n" {
		t.Errorf("multi-connect stat error, got '%v'", r)
	}

	s.Close()
}

// TestSizes verifies the max key and data sizes behave properly.
func TestSizes(t *testing.T) {
	// Server setup
	s, err := NewServer(11212, 65535)
	if err != nil {
		t.Errorf("failed to create server: %v", err)
		return
	}

	err = registerHandlers(s)
	if err != nil {
		t.Errorf("failed to register handlers: %v", err)
		return
	}
	go s.Serve()

	// connection setup
	n, err := net.Dial("tcp", "127.0.0.1:11212")
	if err != nil {
		t.Errorf("unable to connect to server: %v", err)
		return
	}
	n.SetDeadline(time.Now().Add(3 * time.Second))
	b := bufio.NewReader(n)

	// 249 character key test - PASS
	var d string
	for i := 0; i < 249; i++ {
		d = d + "a"
	}
	n.Write([]byte("set " + d + "\r\ndata\r\n"))
	r, err := b.ReadString('\n')
	if err != nil {
		t.Errorf("set 250 char key failed: %v", err)
	}
	if r != "STORED\r\n" {
		t.Errorf("set 249 char failed: %v", r)
	}

	// 250 character key test - FAIL
	d = ""
	for i := 0; i < 250; i++ {
		d = d + "a"
	}
	n.Write([]byte("set " + d + "\r\n"))
	r, err = b.ReadString('\n')
	if err != nil {
		t.Errorf("set 250 char key failed: %v", err)
	}
	if r != "ERROR key can only be 250 characters long\r\n" {
		t.Errorf("set 250 char didn't fail: %v", r)
	}

	// 8k-1 character data test - PASS
	d = ""
	for i := 0; i < 8191; i++ {
		d = d + "a"
	}
	n.Write([]byte("set largeData\r\n" + d + "\r\n"))
	r, err = b.ReadString('\n')
	if err != nil {
		t.Errorf("set 8k data failed: %v", err)
	}
	if r != "STORED\r\n" {
		t.Errorf("set 8k-1 data failed: %v", r)
	}

	// 8k character data test - FAIL
	d = ""
	for i := 0; i < 8192; i++ {
		d = d + "a"
	}
	n.Write([]byte("set largeData\r\n" + d + "\r\n"))
	r, err = b.ReadString('\n')
	if err != nil {
		t.Errorf("set 8k data failed: %v", err)
	}
	if r != "ERROR data can only be 8192 characters long\r\n" {
		t.Errorf("set 8k data didn't fail: %v", r)
	}

	s.Close()
}
