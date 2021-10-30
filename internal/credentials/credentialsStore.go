package credentials

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	data "github.com/apenella/go-common-utils/data"
	errors "github.com/apenella/go-common-utils/error"
)

type CredentialsStore struct {
	store map[string]*RegistryUserPassAuth
	//	backend afero.Fs
}

func NewCredentialsStore() *CredentialsStore {
	return &CredentialsStore{
		store: make(map[string]*RegistryUserPassAuth),
	}
}

func (s *CredentialsStore) LoadCredentials(dir string) error {
	var err error

	errContext := "(credentials::LoadCredentials)"

	if s == nil {
		return errors.New(errContext, "Unable to load credentials because store is not initialized")
	}

	if s.store == nil {
		s.store = make(map[string]*RegistryUserPassAuth)
	}

	_, err = os.Stat(dir)
	if err == nil {

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return errors.New(errContext, fmt.Sprintf("Error reading directory '%s'", dir), err)
		}

		for _, file := range files {
			userpass := &RegistryUserPassAuth{}
			if file.Mode().IsRegular() {
				filename := file.Name()
				err := data.LoadJSONFile(strings.Join([]string{dir, filename}, string(os.PathSeparator)), userpass)
				if err == nil {
					AddCredential(filename, userpass)
				}
			}
		}
	}

	return nil
}

// AddCredential
func (s *CredentialsStore) AddCredentials(id string, auth *RegistryUserPassAuth) error {

	errContext := "(credentials::AddCredential)"

	if s == nil {
		return errors.New(errContext, "Unable to add new credential because store is not initialized")
	}

	if s.store == nil {
		s.store = make(map[string]*RegistryUserPassAuth)
	}

	_, exists := s.store[id]
	if exists {
		return errors.New(errContext, fmt.Sprintf("Auth method with id '%s' already exist", id))
	}

	s.store[id] = auth

	return nil
}

// GetCredentials
func (s *CredentialsStore) GetCredentials(registry string) (*RegistryUserPassAuth, error) {

	errContext := "(credentials::GetCredential)"

	if s == nil {
		return nil, errors.New(errContext, "Unable to get credential because credentials store is not initialized")
	}

	if s.store == nil {
		return nil, errors.New(errContext, "Unable to get credential because store is not initialized")
	}

	hashedRegisty := hashRegistryName(registry)

	credential, exists := s.store[hashedRegisty]
	if !exists {
		return nil, errors.New(errContext, fmt.Sprintf("No credential found for '%s'", registry))
	}

	return credential, nil
}
