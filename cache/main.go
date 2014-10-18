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
		fmt.Println("Failed to create server: ", err)
		return
	}

	err = registerHandlers(s)
	if err != nil {
		fmt.Println("Failed to start server: ", err)
		return
	}

	fmt.Println("Ready to accept cache requests")
	s.Serve()
}

func registerHandlers(s *server) error {
	err := s.AddHandler("set", cmdSet)
	if err != nil {
		return err
	}
	err = s.AddHandler("get", cmdGet)
	if err != nil {
		return err
	}
	err = s.AddHandler("delete", cmdDelete)
	if err != nil {
		return err
	}
	err = s.AddHandler("stats", cmdStats)
	if err != nil {
		return err
	}
	err = s.AddHandler("quit", cmdQuit)
	if err != nil {
		return err
	}
	return nil
}
