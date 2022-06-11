package local

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// LocalStore is a local store for credentials
type LocalStore struct {
	// auth  repository.AuthProviderer
	store map[string]*credentials.Badge
	fs    afero.Fs
	mutex sync.RWMutex
	wg    sync.WaitGroup
}

// NewLocalStore creates a new local store for credentials
func NewLocalStore(fs afero.Fs) *LocalStore {
	return &LocalStore{
		store: make(map[string]*credentials.Badge),
		fs:    fs,
	}
}

// Store stores a badge
func (s *LocalStore) Store(id string, badge *credentials.Badge) error {

	errContext := "(local::Store)"
	hashedID := hashID(id)

	_, exists := s.store[hashedID]
	if exists {
		return errors.New(errContext, fmt.Sprintf("Badge with id '%s' already exists", id))
	}

	s.store[hashedID] = badge

	return nil
}

// Get returns a auth for the badge id
func (s *LocalStore) Get(id string) (*credentials.Badge, error) {

	errContext := "(local::GetAuth)"

	hashedID := hashID(id)

	badge, exists := s.store[hashedID]
	if !exists {
		return nil, errors.New(errContext, fmt.Sprintf("Badge with id '%s' not found", id))
	}

	return badge, nil
}

// // Get returns a badge
// func (s *LocalStore) GetBadge(id string) (*credentials.Badge, error) {
// 	errContext := "(local::GetAuth)"

// 	hashedID := hashID(id)

// 	badge, exists := s.store[hashedID]
// 	if !exists {
// 		return nil, errors.New(errContext, fmt.Sprintf("Badge with id '%s' not found", id))
// 	}

// 	return badge, nil
// }

// hashID
func hashID(id string) string {
	hasher := md5.New()
	hasher.Write([]byte(id))
	registryHashed := hex.EncodeToString(hasher.Sum(nil))

	return registryHashed
}

// LoadCredentials to the store
func (s *LocalStore) LoadCredentials(path string) error {
	var err error
	var isDir bool

	errContext := "(local::LoadCredentials)"

	if s.store == nil {
		s.store = make(map[string]*credentials.Badge)
	}

	isDir, err = afero.IsDir(s.fs, path)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	if isDir {
		return s.LoadCredentialsFromDir(path)
	} else {
		return s.LoadCredentialsFromFile(path)
	}
}

func (s *LocalStore) LoadCredentialsFromDir(path string) error {
	var err error
	errFuncs := []func() error{}
	errContext := "(local::LoadCredentialsFromDir)"

	credFiles, err := afero.Glob(s.fs, filepath.Join(path, "*"))
	if err != nil {
		return errors.New(errContext, "", err)
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

func (s *LocalStore) LoadCredentialsFromFile(path string) error {

	var err error
	var fileData []byte
	var fileInfo os.FileInfo

	errContext := "(local::LoadCredentialsFromFile)"

	fileData, err = afero.ReadFile(s.fs, path)
	if err != nil {
		return errors.New(errContext, "", err)
	}
	badge := &credentials.Badge{}
	err = yaml.Unmarshal(fileData, badge)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error loading credentials from file '%s'", path), err)
	}

	fileInfo, err = s.fs.Stat(path)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = s.AddCredentials(fileInfo.Name(), badge)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

// AddCredential
func (s *LocalStore) AddCredentials(id string, auth *credentials.Badge) error {

	errContext := "(local::AddCredential)"

	if s == nil {
		return errors.New(errContext, "Unable to add new credential because store is not initialized")
	}

	if s.store == nil {
		s.store = make(map[string]*credentials.Badge)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.store[id]
	if exists {
		return errors.New(errContext, fmt.Sprintf("Auth method with id '%s' already exist", id))
	}

	s.store[id] = auth

	return nil
}
