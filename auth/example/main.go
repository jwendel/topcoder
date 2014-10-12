package main

import (
	"bitbucket.org/kyrra/sandbox/auth"
	"fmt"
)

func main() {
	err := auth.Start(":8080", "users.json")
	if err != nil {
		fmt.Println("Error starting auth server: ", err)
	}
}
