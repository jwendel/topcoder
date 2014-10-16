package main

import (
	"flag"
	"fmt"
)

func main() {
	p := flag.Int("port", 11212, "Port the server listens on")
	i := flag.Int("items", 65535, "Maximum number of items to cache")
	flag.Parse()

	s, err := NewServer(*p, *i)
	if err != nil {
		fmt.Println("It borked: ", err)
	}
	fmt.Println("Ready to accept cache requests")
	s.Serve()
}
