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
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/spf13/afero"
)

// LocalStore is a local store for credentials
type LocalStore struct {
	// auth  repository.AuthProviderer
	store    map[string]*credentials.Badge
	fs       afero.Fs
	mutex    sync.RWMutex
	wg       sync.WaitGroup
	formater repository.Formater
}

// NewLocalStore creates a new local store for credentials
func NewLocalStore(fs afero.Fs, f repository.Formater) *LocalStore {
	return &LocalStore{
		store:    make(map[string]*credentials.Badge),
		fs:       fs,
		formater: f,
	}
}

// Store stores a badge on memory store
func (s *LocalStore) Store(id string, badge *credentials.Badge) error {

	errContext := "(store::credentials::local::Store)"
	hashedID := hashID(id)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.store[hashedID]
	if exists {
		return errors.New(errContext, fmt.Sprintf("Badge with id '%s' already exists", id))
	}

	s.store[hashedID] = badge

	return nil
}

// Persist save a badge in the local store
func (s *LocalStore) Persist(path, id string, badge *credentials.Badge) error {

	var err error
	var formatedBadge string
	var credentialFile afero.File

	errContext := "(store::credentials::local::Persist)"
	hashedID := hashID(id)

	err = s.Store(id, badge)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Badge with id '%s', could not be persisted to '%s'", id, path), err)
	}

	err = s.fs.MkdirAll(path, 0755)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error creating directory '%s'", path), err)
	}

	credentialFile, err = s.fs.OpenFile(filepath.Join(path, hashedID), os.O_RDWR|os.O_CREATE, 0600)
	defer credentialFile.Close()

	formatedBadge, err = s.formater.Marshal(badge)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error formatting '%s' badge before to be persisted on '%s'", id, path), err)
	}

	_, err = credentialFile.WriteString(formatedBadge)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

// Get returns a auth for the badge id
func (s *LocalStore) Get(id string) (*credentials.Badge, error) {

	errContext := "(store::credentials::local::GetAuth)"

	hashedID := hashID(id)

	badge, exists := s.store[hashedID]
	if !exists {
		return nil, errors.New(errContext, fmt.Sprintf("Badge with id '%s' not found", id))
	}

	return badge, nil
}

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

	errContext := "(store::credentials::local::LoadCredentials)"

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
	errContext := "(store::credentials::local::LoadCredentialsFromDir)"

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
	var badge *credentials.Badge

	errContext := "(store::credentials::local::LoadCredentialsFromFile)"

	fileData, err = afero.ReadFile(s.fs, path)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error reading credentials file '%s'", path), err)
	}

	badge, err = s.formater.Unmarshal(fileData)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error unmarshaling credentials from file '%s'", path), err)
	}

	fileInfo, err = s.fs.Stat(path)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error getting the stat of credentials file '%s'", path), err)
	}

	err = s.AddCredentials(fileInfo.Name(), badge)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

// AddCredential
func (s *LocalStore) AddCredentials(id string, auth *credentials.Badge) error {

	errContext := "(store::credentials::local::AddCredential)"

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
