package main

import (
	"bitbucket.org/kyrra/sandbox/auth"
	"flag"
	"fmt"
)

func main() {
	listen := flag.String("listen", ":8080", "Hostname and address to listen on")
	source := flag.String("datasource", "users.json", "Filename to load JSON user data from")
	flag.Parse()

	err := auth.Serve(*listen, *source)
	if err != nil {
		fmt.Println("error starting auth server: ", err)
	}
}
