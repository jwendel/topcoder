// Copyright (c) 2014 James Wendel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type datastore struct {
	mutex         sync.RWMutex
	authFilename  string
	authFileinfo  os.FileInfo
	tokenFilename string
	// map[DomainName]
	domainMap DomainAuths
	tokenMap  DomainTokens
}

// DomainAuths is map[domainName]Domain.  Maps a domain name to data
// about that domain
type DomainAuths map[string]Domain

// Domain data for proxyauth and oauth lookup
type Domain struct {
	// map[username]HashedPassword
	Users map[string]string
	// map[client_id]client_secret
	Clients map[string]string
}

// domainJSON structure to read domains from JSON input file
type domainJSON struct {
	Name    string       `json:"domain"`
	Users   []userJSON   `json:"users"`
	Clients []clientJSON `json:"clients"`
}

// user structure to read users from JSON input file
type userJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// client TODO
type clientJSON struct {
	ID     string `json:"client_id"`
	Secret string `json:"client_secret"`
}

// Init loads the passed in json file, unmarshels the data,
// and starts a fileWatcher to look for changes to the file
func (ds *datastore) Init(authFile, tokenFile string) error {
	ds.authFilename = authFile
	ds.tokenFilename = tokenFile

	b, err := ds.loadAuthFile()
	if err != nil {
		return err
	}
	err = ds.unmarshal(b)
	if err != nil {
		return err
	}
	ds.authFileinfo, err = os.Stat(ds.authFilename)
	if err != nil {
		return err
	}
	go ds.fileWatcher()
	go ds.startSigHandler()
	return nil
}

// DomainExists checks if the given domain exists in the data store.
func (ds *datastore) DomainExists(domain string) bool {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	_, ok := ds.domainMap[domain]
	return ok
}

// UserPasswordValid returns true when the password is valid for a given domain/user
// else it just returns false.  Password is expected to be in encrypted form.
func (ds *datastore) UserPasswordValid(domain, username, password string) bool {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	d, ok := ds.domainMap[domain]
	if !ok {
		return false
	}
	pass, ok := d.Users[username]
	if !ok {
		return false
	}

	if pass == password {
		return true
	}
	return false
}

// loadAuthFile loads the full file from disk
func (ds *datastore) loadAuthFile() ([]byte, error) {
	// Load the data source from disk
	b, err := ioutil.ReadFile(ds.authFilename)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// unmarshal converts bytes to a JSON structure then populates the
// datastore.dataMap with the results.
func (ds *datastore) unmarshal(bytes []byte) error {
	// Updating the user database, write lock needed
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	var domains []domainJSON
	err := json.Unmarshal(bytes, &domains)
	if err != nil {
		return err
	}

	domainMap := make(map[string]Domain)
	tokenMap := make(map[string]Tokens)

	// Loop over all domains and users, inserting them into the domainMap
	// The user password will be encrypted with this step
	for _, d := range domains {
		_, ok := domainMap[d.Name]
		if !ok {
			var domain Domain
			domain.Users = make(map[string]string)
			for _, u := range d.Users {
				_, ok := domain.Users[u.Username]
				if !ok {
					domain.Users[u.Username] = EncryptPassword(u.Password)
				} else {
					return fmt.Errorf("duplicate username '%v' for domain '%v'", u.Username, d.Name)
				}
			}

			domain.Clients = make(map[string]string)
			for _, u := range d.Clients {
				_, ok := domain.Clients[u.ID]
				if !ok {
					domain.Clients[u.ID] = u.Secret
				} else {
					return fmt.Errorf("duplicate client_id '%v' for domain '%v'", u.ID, d.Name)
				}
			}

			domainMap[d.Name] = domain
			tokenMap[d.Name] = make(map[string]time.Time)

		} else {
			return fmt.Errorf("duplicate domain '%v' in input file", d.Name)
		}
	}

	ds.domainMap = domainMap
	ds.tokenMap = tokenMap

	return nil
}

// fileWatcher checks once every 3 seconds if the source json file has changed
// based on it's timestamp.  If it chagnes it will reload the user data.
func (ds *datastore) fileWatcher() {
	for {
		time.Sleep(3 * time.Second)
		fi, err := os.Stat(ds.authFilename)
		if err != nil {
			fmt.Printf("Failed watching file '%v' for updates\n", ds.authFilename)
			return
		}

		if !fi.ModTime().Equal(ds.authFileinfo.ModTime()) {
			// file modified time changed, reload data
			b, err := ds.loadAuthFile()
			if err != nil {
				fmt.Printf("Error loading file '%v': %v", ds.authFilename, err)
				return
			}
			err = ds.unmarshal(b)
			if err != nil {
				fmt.Printf("Error unmarshling '%v': %v", ds.authFilename, err)
				return
			}
			ds.authFileinfo = fi
		}
	}
}
