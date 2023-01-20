package local

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/encryption"
	"github.com/spf13/afero"
)

// OptionsFunc defines the signature for an option function to set local credentials store
type OptionsFunc func(opts *LocalStore)

// LocalStore is a local store for credentials
type LocalStore struct {
	fs            afero.Fs
	path          string
	mutex         sync.RWMutex
	wg            sync.WaitGroup
	compatibility CredentialsCompatibilier
	formater      repository.Formater
	encryption    Encrypter
}

// NewLocalStore creates a new local store for credentials
func NewLocalStore(opts ...OptionsFunc) *LocalStore {
	store := &LocalStore{}
	store.Options(opts...)
	return store
}

func WithCompatibility(compatibility CredentialsCompatibilier) OptionsFunc {
	return func(s *LocalStore) {
		s.compatibility = compatibility
	}
}

func WithFilesystem(fs afero.Fs) OptionsFunc {
	return func(s *LocalStore) {
		s.fs = fs
	}
}

// WithFormater sets the formater to envvars credentials store
func WithFormater(formater repository.Formater) OptionsFunc {
	return func(s *LocalStore) {
		s.formater = formater
	}
}

func WithPath(path string) OptionsFunc {
	return func(s *LocalStore) {
		s.path = path
	}
}

func WithEncryption(e Encrypter) OptionsFunc {
	return func(s *LocalStore) {
		s.encryption = e
	}
}

// Options provides the options to envvars credentials store
func (s *LocalStore) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *LocalStore) SafeStore(id string, badge *credentials.Badge) error {
	errContext := "(store::credentials::local::SafeStore)"
	var err error
	var credentialsStat os.FileInfo

	hashedID, err := encryption.HashID(id)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	credentialsStat, _ = s.fs.Stat(filepath.Join(s.path, hashedID))

	if credentialsStat != nil && credentialsStat.Name() != "" {
		return errors.New(errContext, fmt.Sprintf("Credentials '%s' already exist", id))
	}

	err = s.Store(id, badge)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
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

	if badge.ID == "" {
		badge.ID = id
	}

	hashedID, err := encryption.HashID(id)
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

	if s.encryption != nil {
		formatedBadge, err = s.encryption.Encrypt(formatedBadge)
		if err != nil {
			return errors.New(errContext, "", err)
		}
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

	hashedID, err := encryption.HashID(id)
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
	var strFileData string
	var badge *credentials.Badge

	errContext := "(store::credentials::local::get)"

	fileData, err = afero.ReadFile(s.fs, filepath.Join(s.path, id))
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error reading credentials file '%s'", filepath.Join(s.path, id)), err)
	}

	if s.encryption != nil {
		strFileData, err = s.encryption.Decrypt(string(fileData))
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		fileData = []byte(strFileData)
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
func (s *LocalStore) All() ([]*credentials.Badge, error) {

	var badge *credentials.Badge
	badges := []*credentials.Badge{}

	afero.Walk(s.fs, s.path, func(path string, info os.FileInfo, err error) error {

		errContext := "(store::credentials::local::All::walk)"

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

	return badges, nil
}

// // hashID generates a hash for the id
// func hashID(id string) (string, error) {

// 	errContext := "(store::credentials::local::hashID)"

// 	if id == "" {
// 		return "", errors.New(errContext, "Hash method requires an id")
// 	}

// 	hasher := md5.New()
// 	hasher.Write([]byte(id))
// 	registryHashed := hex.EncodeToString(hasher.Sum(nil))

// 	return registryHashed, nil
// }
