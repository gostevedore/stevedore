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
	fs            afero.Fs
	path          string
	mutex         sync.RWMutex
	wg            sync.WaitGroup
	compatibility CredentialsCompatibilier
	formater      repository.Formater
}

// NewLocalStore creates a new local store for credentials
func NewLocalStore(fs afero.Fs, path string, f repository.Formater, compatibility CredentialsCompatibilier) *LocalStore {
	return &LocalStore{
		path:          path,
		fs:            fs,
		formater:      f,
		compatibility: compatibility,
	}
}

// Store save a badge in the local store
func (s *LocalStore) Store(id string, badge *credentials.Badge) error {

	var err error
	var formatedBadge string
	var credentialFile afero.File

	errContext := "(store::credentials::local::Store)"

	if s.path == "" {
		return errors.New(errContext, "To store a badge into local store, local store path must be provided")
	}

	if id == "" {
		return errors.New(errContext, "To store a badge into local store, id must be provided")
	}

	if badge == nil {
		return errors.New(errContext, fmt.Sprintf("To store a badge for '%s' into local store, credentials badge must be provided", id))
	}

	hashedID, err := hashID(id)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	err = s.fs.MkdirAll(s.path, 0755)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error creating directory '%s'", s.path), err)
	}

	credentialFile, err = s.fs.OpenFile(filepath.Join(s.path, hashedID), os.O_RDWR|os.O_CREATE, 0600)
	defer credentialFile.Close()

	err = s.compatibility.CheckCompatibility(badge)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	formatedBadge, err = s.formater.Marshal(badge)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error formatting '%s' badge before to be persisted on '%s'", id, s.path), err)
	}

	_, err = credentialFile.WriteString(formatedBadge)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

// Get returns a auth for the badge id
func (s *LocalStore) Get(id string) (*credentials.Badge, error) {
	var err error
	var badge *credentials.Badge

	errContext := "(store::credentials::local::Get)"

	if id == "" {
		return nil, errors.New(errContext, "To get a badge from the store, id must be provided")
	}

	hashedID, err := hashID(id)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	badge, err = s.get(hashedID)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return badge, nil
}

// get return a badge from the store using the hashed id
func (s *LocalStore) get(id string) (*credentials.Badge, error) {
	var err error
	var fileData []byte
	var badge *credentials.Badge

	errContext := "(store::credentials::local::get)"

	fileData, err = afero.ReadFile(s.fs, filepath.Join(s.path, id))
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error reading credentials file '%s'", filepath.Join(s.path, id)), err)
	}

	badge, err = s.formater.Unmarshal(fileData)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error unmarshaling credentials from file '%s'", filepath.Join(s.path, id)), err)
	}

	if badge.ID == "" {
		badge.ID = id
	}

	err = s.compatibility.CheckCompatibility(badge)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return badge, nil
}

// All returns all badges from the store
func (s *LocalStore) All() []*credentials.Badge {
	var badge *credentials.Badge
	badges := []*credentials.Badge{}

	afero.Walk(s.fs, s.path, func(path string, info os.FileInfo, err error) error {

		errContext := "(store::credentials::local::get::walk)"

		_, err = s.fs.Stat(path)
		if err != nil {
			return errors.New(errContext, fmt.Sprintf("Error reading credentials file '%s'", path), err)
		}

		if !info.IsDir() {
			badge, _ = s.get(info.Name())
			badges = append(badges, badge)
		}

		return nil
	})

	return badges
}

// hashID generates a hash for the id
func hashID(id string) (string, error) {

	errContext := "(store::credentials::local::hashID)"

	if id == "" {
		return "", errors.New(errContext, "Hash method requires an id")
	}

	hasher := md5.New()
	hasher.Write([]byte(id))
	registryHashed := hex.EncodeToString(hasher.Sum(nil))

	return registryHashed, nil
}
