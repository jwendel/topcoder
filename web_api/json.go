package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var store datastore

type datastore struct {
	mutex    sync.RWMutex
	filename string
	fileinfo os.FileInfo
	// map[DomainName]map[UserName]EncryptedPassword
	domainMap map[string]map[string]string
}

type domain struct {
	Domain string `json:"domain"`
	Users  []user `json:"users"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *datastore) Init(filename string) error {
	s.filename = filename
	err := s.load()
	if err != nil {
		return err
	}
	s.fileinfo, err = os.Stat(s.filename)
	if err != nil {
		return err
	}
	go s.fileWatcher()
	return nil
}

func (s *datastore) load() error {
	// Load the data source from disk
	defer 
	b, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return err
	}

	var domains []domain
	json.Unmarshal(b, &domains)

	domainMap := make(map[string]map[string]string)

	// Loop over all domains and users, inserting them into the domainMap
	// The user password will be encrypted with this step
	for _, d := range domains {
		_, ok := s.domainMap[d.Domain]
		if !ok {
			userMap := make(map[string]string)
			for _, u := range d.Users {
				_, ok := userMap[u.Username]
				if !ok {
					userMap[u.Username] = encryptPassword(u.Password)
				} else {
					return fmt.Errorf("duplicate username '%v' for domain '%v'", u.Username, d.Domain)
				}
			}
			s.domainMap[d.Domain] = userMap
		} else {
			return fmt.Errorf("duplicate domains '%v' in input file", d.Domain)
		}
	}

	s.domainMap = domainMap

	return nil
}

func (s *datastore) fileWatcher() {

	for {
		time.Sleep(1 * time.Second)
		fi, err := os.Stat(s.filename)
		if err != nil {
			fmt.Printf("Failed watching file '%v' for updates\n", s.filename)
		}

		if !fi.ModTime().Equal(s.fileinfo.ModTime()) {
			// file modified time changed, reload data
			s.load()
			s.fileinfo = fi
		}
	}
}
