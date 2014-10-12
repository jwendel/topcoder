package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var Store datastore

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

func (s *datastore) DomainExists(domain string) bool {
	_, ok := s.domainMap[domain]
	return ok
}

func (s *datastore) UserPasswordValid(domain, username, password string) bool {
	d, ok := s.domainMap[domain]
	if !ok {
		return false
	}
	pass, ok := d[username]
	if !ok {
		return false
	}

	if pass == password {
		return true
	}
	return false
}

func (s *datastore) load() error {
	// Updating the user database, write lock needed
	defer s.mutex.Unlock()
	s.mutex.Lock()

	// Load the data source from disk
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
		_, ok := domainMap[d.Domain]
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
			domainMap[d.Domain] = userMap
		} else {
			return fmt.Errorf("duplicate domains '%v' in input file", d.Domain)
		}
	}

	s.domainMap = domainMap

	return nil
}

// fileWatcher checks once a second if the source json file has changed
// based on it's timestamp.  If it chagnes it will reload the user data.
func (s *datastore) fileWatcher() {
	for {
		time.Sleep(1 * time.Second)
		fi, err := os.Stat(s.filename)
		if err != nil {
			fmt.Printf("Failed watching file '%v' for updates\n", s.filename)
			return
		}

		if !fi.ModTime().Equal(s.fileinfo.ModTime()) {
			// file modified time changed, reload data
			s.load()
			s.fileinfo = fi
		}
	}
}
