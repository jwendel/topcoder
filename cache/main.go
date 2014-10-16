package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hi Rachel!")
	s, err := NewServer()
	if err != nil {
		fmt.Println("It borked: ", err)
	}
	s.Serve()
}
