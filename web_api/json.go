package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var store datastore

type datastore struct {
	mutex sync.RWMutex
}

type domain struct {
	Domain string `json:"domain"`
	Users  []user `json:"users"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s datastore) load(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nill {
		return err
	}

	json.Unmarshal(b, datastore)

	return nil
}
