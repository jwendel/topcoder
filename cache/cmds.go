package main

import (
	"fmt"
)

func cmdSet(c *CacheRequest) {

	if len(c.Subcmd) != 1 || len(c.Subcmd[0]) == 0 {
		c.writeStr("ERROR key required with set comamand")
		return
	}

	d, err := c.readln()
	if err != nil {
		c.writeStr("ERROR invalid data for set")
		return
	}

	// Check that it was \r\n
	if len(d) <= 2 || d[len(d)-2] != '\r' {
		c.writeStr("ERROR invalid data in set")
		return
	}

	// Trim \r\n
	d = d[:len(d)-2]
	input := string(d)

	if !validChars.MatchString(input) {
		c.writeStr("ERROR data contains invalid characters")
		return
	}

	c.s.cacheMutex.Lock()
	defer c.s.cacheMutex.Unlock()

	c.s.cache[c.Subcmd[0]] = input
	c.writeStr("STORED")
}

func cmdGet(c *CacheRequest) {
	if len(c.Subcmd) == 0 {
		c.writeStr("ERROR key required with get command")
		return
	}

	c.s.cacheMutex.RLock()
	defer c.s.cacheMutex.RUnlock()

	for _, v := range c.Subcmd {
		d, ok := c.s.cache[v]
		if !ok {
			continue
		}

		c.writeStr(fmt.Sprintf("VALUE %v", v))
		c.writeStr(d)

	}
	c.writeStr("END")

}

func cmdDelete(c *CacheRequest) {
	c.s.cacheMutex.Lock()
	defer c.s.cacheMutex.Unlock()

}

func cmdStats(c *CacheRequest) {

}

func cmdQuit(c *CacheRequest) {

}
