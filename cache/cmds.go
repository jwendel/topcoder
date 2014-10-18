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

	c.C.CacheMutex.Lock()
	defer c.C.CacheMutex.Unlock()

	c.C.Cache[c.Subcmd[0]] = input
	c.writeStr("STORED")
}

func cmdGet(c *CacheRequest) {
	if len(c.Subcmd) == 0 {
		c.writeStr("ERROR key required with get command")
		return
	}

	c.C.CacheMutex.RLock()
	defer c.C.CacheMutex.RUnlock()

	for _, v := range c.Subcmd {
		d, ok := c.C.Cache[v]
		if !ok {
			continue
		}

		c.writeStr(fmt.Sprintf("VALUE %v", v))
		c.writeStr(d)

	}
	c.writeStr("END")

}

func cmdDelete(c *CacheRequest) {
	if len(c.Subcmd) != 1 {
		c.writeStr("ERROR key required with delete command")
		return
	}

	key := c.Subcmd[0]

	c.C.CacheMutex.Lock()
	defer c.C.CacheMutex.Unlock()

	_, ok := c.C.Cache[key]
	if !ok {
		c.writeStr("NOT_FOUND")
		return
	}

	delete(c.C.Cache, key)
	c.writeStr("DELETED")
}

func cmdStats(c *CacheRequest) {
	c.C.CacheMutex.RLock()
	defer c.C.CacheMutex.RUnlock()

	c.writeStr(fmt.Sprintf("cmd_get %v", c.C.Stats.get))
	c.writeStr(fmt.Sprintf("cmd_set %v", c.C.Stats.set))
	c.writeStr(fmt.Sprintf("get_hits %v", c.C.Stats.getHits))
	c.writeStr(fmt.Sprintf("get_misses %v", c.C.Stats.getMiss))
	c.writeStr(fmt.Sprintf("delete_hits %v", c.C.Stats.delHits))
	c.writeStr(fmt.Sprintf("delete_misses %v", c.C.Stats.delMiss))
	c.writeStr(fmt.Sprintf("curr_items %v", len(c.C.Cache)))
	c.writeStr(fmt.Sprintf("limit_items %v", c.C.maxItems))
}

func cmdQuit(c *CacheRequest) {
	c.conn.Close()
}
