// Copyright 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

// cmdSet takes a single key then will read one more line
// from the connection and add the data to the cache
func cmdSet(c *CacheRequest) {
	if len(c.Subcmd) != 1 {
		c.WriteStr("ERROR set command requires a single key to be specified")
		return
	}

	if len(c.Subcmd[0]) >= MAX_KEY_SIZE {
		c.WriteStr(fmt.Sprintf("ERROR key can only be %v characters long", MAX_KEY_SIZE))
		return
	}

	d, err := c.Readln()
	if err != nil {
		c.WriteStr("ERROR invalid data for set")
		return
	}

	input, err := c.ValidateInput(d)
	if err != nil {
		c.WriteStr(err.Error())
		return
	}

	if len(input) >= MAX_DATA_SIZE {
		c.WriteStr(fmt.Sprintf("ERROR data can only be %v characters long", MAX_DATA_SIZE))
		return
	}

	c.C.CacheMutex.Lock()
	defer c.C.CacheMutex.Unlock()

	_, ok := c.C.Cache[c.Subcmd[0]]
	if !ok && len(c.C.Cache) == c.C.maxItems {
		c.WriteStr("ERROR cache is full")
		return
	}

	c.C.Stats.set++
	c.C.Cache[c.Subcmd[0]] = input
	c.WriteStr("STORED")
}

// cmdGet takes 1 or more keys and will return the data
// for each key found in the cache.
func cmdGet(c *CacheRequest) {
	if len(c.Subcmd) == 0 {
		c.WriteStr("ERROR key required with get command")
		return
	}

	c.C.CacheMutex.Lock()
	defer c.C.CacheMutex.Unlock()

	for _, v := range c.Subcmd {
		c.C.Stats.get++
		d, ok := c.C.Cache[v]
		if !ok {
			c.C.Stats.getMisses++
			continue
		}

		c.C.Stats.getHits++
		c.WriteStr(fmt.Sprintf("VALUE %v", v))
		c.WriteStr(d)

	}
	c.WriteStr("END")

}

// cmdDelete takes a single key and will attempt to remove
// the key from the cache.
func cmdDelete(c *CacheRequest) {
	if len(c.Subcmd) != 1 {
		c.WriteStr("ERROR delete command requires a single key to be specified")
		return
	}

	key := c.Subcmd[0]

	c.C.CacheMutex.Lock()
	defer c.C.CacheMutex.Unlock()

	_, ok := c.C.Cache[key]
	if !ok {
		c.C.Stats.delMisses++
		c.WriteStr("NOT_FOUND")
		return
	}

	c.C.Stats.delHits++
	delete(c.C.Cache, key)
	c.WriteStr("DELETED")
}

// cmdStats prints the current usage statistics for the cache.
func cmdStats(c *CacheRequest) {
	if len(c.Subcmd) != 0 {
		c.WriteStr("ERROR stats does not take any parameters")
		return
	}

	c.C.CacheMutex.RLock()
	defer c.C.CacheMutex.RUnlock()

	c.WriteStr(fmt.Sprintf("cmd_get %v", c.C.Stats.get))
	c.WriteStr(fmt.Sprintf("cmd_set %v", c.C.Stats.set))
	c.WriteStr(fmt.Sprintf("get_hits %v", c.C.Stats.getHits))
	c.WriteStr(fmt.Sprintf("get_misses %v", c.C.Stats.getMisses))
	c.WriteStr(fmt.Sprintf("delete_hits %v", c.C.Stats.delHits))
	c.WriteStr(fmt.Sprintf("delete_misses %v", c.C.Stats.delMisses))
	c.WriteStr(fmt.Sprintf("curr_items %v", len(c.C.Cache)))
	c.WriteStr(fmt.Sprintf("limit_items %v", c.C.maxItems))
	c.WriteStr("END")
}

// cmdQuit closes the connection with the client.
func cmdQuit(c *CacheRequest) {
	if len(c.Subcmd) != 0 {
		c.WriteStr("ERROR quit does not take any parameters")
		return
	}

	c.conn.Close()
}
