package credentials

import (
	"fmt"
	"os"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type CredentialsStore struct {
	Store map[string]*UserPasswordAuth
	mutex sync.RWMutex
	wg    sync.WaitGroup
	fs    afero.Fs
}

func NewCredentialsStore(fs afero.Fs) *CredentialsStore {
	return &CredentialsStore{
		Store: make(map[string]*UserPasswordAuth),

		fs: fs,
	}
}

func (s *CredentialsStore) LoadCredentials(path string) error {
	var err error
	var isDir bool

	errContext := "(credentials::LoadCredentials)"

	// if s == nil {
	// 	return errors.New(errContext, "Unable to load credentials because store is not initialized")
	// }

	if s.Store == nil {
		s.Store = make(map[string]*UserPasswordAuth)
	}

	isDir, err = afero.IsDir(s.fs, path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	if isDir {
		return s.LoadCredentialsFromDir(path)
	} else {
		return s.LoadCredentialsFromFile(path)
	}

	// _, err = os.Stat(dir)
	// if err == nil {

	// 	files, err := ioutil.ReadDir(dir)
	// 	if err != nil {
	// 		return errors.New(errContext, fmt.Sprintf("Error reading directory '%s'", dir), err)
	// 	}

	// 	for _, file := range files {
	// 		userpass := &UserPasswordAuth{}
	// 		if file.Mode().IsRegular() {
	// 			filename := file.Name()
	// 			err := data.LoadJSONFile(strings.Join([]string{dir, filename}, string(os.PathSeparator)), userpass)
	// 			if err == nil {
	// 				AddCredential(filename, userpass)
	// 			}
	// 		}
	// 	}
	// }

	// return nil
}

func (s *CredentialsStore) LoadCredentialsFromDir(path string) error {
	var err error
	errFuncs := []func() error{}
	errContext := "(credentials::LoadCredentialsFromDir)"

	credFiles, err := afero.Glob(s.fs, path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	loadCredentialsFromFile := func(path string) func() error {
		var err error

		c := make(chan struct{}, 1)
		go func() {
			defer close(c)
			err = s.LoadCredentialsFromFile(path)
			s.wg.Done()
		}()

		return func() error {
			<-c
			return err
		}
	}

	for _, file := range credFiles {
		s.wg.Add(1)
		f := loadCredentialsFromFile(file)
		errFuncs = append(errFuncs, f)
	}

	s.wg.Wait()

	errMsg := ""
	for _, f := range errFuncs {
		err = f()
		if err != nil {
			errMsg = fmt.Sprintf("%s%s\n", errMsg, err.Error())
		}
	}
	if errMsg != "" {
		return errors.New(errContext, errMsg)
	}

	return nil
}

func (s *CredentialsStore) LoadCredentialsFromFile(path string) error {

	var err error
	var fileData []byte
	var fileInfo os.FileInfo

	errContext := "(credentials::LoadCredentialsFromFile)"

	fileData, err = afero.ReadFile(s.fs, path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	userpass := &UserPasswordAuth{}
	err = yaml.Unmarshal(fileData, userpass)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error loading credentials from file '%s'", path), err)
	}

	fileInfo, err = s.fs.Stat(path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	err = s.AddCredentials(fileInfo.Name(), userpass)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}

// AddCredential
func (s *CredentialsStore) AddCredentials(id string, auth *UserPasswordAuth) error {

	errContext := "(credentials::AddCredential)"

	if s == nil {
		return errors.New(errContext, "Unable to add new credential because store is not initialized")
	}

	if s.Store == nil {
		s.Store = make(map[string]*UserPasswordAuth)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.Store[id]
	if exists {
		return errors.New(errContext, fmt.Sprintf("Auth method with id '%s' already exist", id))
	}

	s.Store[id] = auth

	return nil
}

// GetCredentials
func (s *CredentialsStore) Get(registry string) (*UserPasswordAuth, error) {

	errContext := "(credentials::GetCredential)"

	if s == nil {
		return nil, errors.New(errContext, "Unable to get credential because credentials store is not initialized")
	}

	if s.Store == nil {
		return nil, errors.New(errContext, "Unable to get credential because store is not initialized")
	}

	hashedRegisty := hashRegistryName(registry)

	credential, exists := s.Store[hashedRegisty]
	if !exists {
		return nil, errors.New(errContext, fmt.Sprintf("No credential found for '%s'", registry))
	}

	return credential, nil
}
