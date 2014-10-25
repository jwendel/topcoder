// Copyright (c) 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bitbucket.org/kyrra/sandbox/webapi/auth"
	"flag"
	"fmt"
)

func main() {
	listen := flag.String("listen", ":8080", "Hostname and address to listen on")
	source := flag.String("datasource", "users.json", "Filename to load JSON user data from")
	tokenTimeout := flag.Int("tokenTimeout", 3600, "Lifetime of auth tokens in seconds")
	flag.Parse()

	err := auth.Serve(*listen, *source, *tokenTimeout)
	if err != nil {
		fmt.Println("error starting auth server: ", err)
	}
}
