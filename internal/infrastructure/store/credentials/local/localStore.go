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

func (s *LocalStore) SafeStore(id string, credential *credentials.Credential) error {
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

	err = s.Store(id, credential)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

// Store save a credential in the local store
func (s *LocalStore) Store(id string, credential *credentials.Credential) error {

	var err error
	var formatedCredential string
	var credentialFile afero.File

	errContext := "(store::credentials::local::Store)"

	if s.path == "" {
		return errors.New(errContext, "To store a credential into local store, local store path must be provided")
	}

	if id == "" {
		return errors.New(errContext, "To store a credential into local store, id must be provided")
	}

	if credential == nil {
		return errors.New(errContext, fmt.Sprintf("To store a credential for '%s' into local store, credentials credential must be provided", id))
	}

	if credential.ID == "" {
		credential.ID = id
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

	err = s.compatibility.CheckCompatibility(credential)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	formatedCredential, err = s.formater.Marshal(credential)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error formatting '%s' credential before to be persisted on '%s'", id, s.path), err)
	}

	if s.encryption != nil {
		formatedCredential, err = s.encryption.Encrypt(formatedCredential)
		if err != nil {
			return errors.New(errContext, "", err)
		}
	}

	_, err = credentialFile.WriteString(formatedCredential)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

// Get returns a auth for the credential id
func (s *LocalStore) Get(id string) (*credentials.Credential, error) {
	var err error
	var credential *credentials.Credential

	errContext := "(store::credentials::local::Get)"

	if id == "" {
		return nil, errors.New(errContext, "To get a credential from the store, id must be provided")
	}

	hashedID, err := encryption.HashID(id)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	credential, err = s.get(hashedID)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return credential, nil
}

// get return a credential from the store using the hashed id
func (s *LocalStore) get(id string) (*credentials.Credential, error) {
	var err error
	var fileData []byte
	var strFileData string
	var credential *credentials.Credential

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

	credential, err = s.formater.Unmarshal(fileData)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error unmarshaling credentials from file '%s'", filepath.Join(s.path, id)), err)
	}

	if credential.ID == "" {
		credential.ID = id
	}

	err = s.compatibility.CheckCompatibility(credential)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return credential, nil
}

// All returns all credentials from the store
func (s *LocalStore) All() ([]*credentials.Credential, error) {

	var credential *credentials.Credential
	credentials := []*credentials.Credential{}

	afero.Walk(s.fs, s.path, func(path string, info os.FileInfo, err error) error {

		errContext := "(store::credentials::local::All::walk)"

		_, err = s.fs.Stat(path)
		if err != nil {
			return errors.New(errContext, fmt.Sprintf("Error reading credentials file '%s'", path), err)
		}

		if !info.IsDir() {
			credential, _ = s.get(info.Name())
			credentials = append(credentials, credential)
		}

		return nil
	})

	return credentials, nil
}
