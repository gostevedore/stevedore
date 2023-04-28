package envvars

import (
	"fmt"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	//	"github.com/fatih/structs"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/encryption"
	//"github.com/spf13/viper"
)

const (
	envvarsCredentialsPrefix          = "stevedore_envvars_credentials"
	envvarsCredentialsAttributePrefix = "attr"
)

// OptionsFunc defines the signature for an option function to set envvars credentials store
type OptionsFunc func(opts *EnvvarsStore)

// EnvvarsStore is a store for credentials
type EnvvarsStore struct {
	console    ConsoleWriter
	backend    EnvvarsBackender
	encryption Encrypter
	formater   repository.Formater
}

// NewEnvvarsStore creates a new mocked store for credentials
func NewEnvvarsStore(opts ...OptionsFunc) *EnvvarsStore {
	store := &EnvvarsStore{}
	store.Options(opts...)

	return store
}

// WithBackend sets the writer to envvars credentials store
func WithBackend(backend EnvvarsBackender) OptionsFunc {
	return func(s *EnvvarsStore) {
		s.backend = backend
	}
}

// WithConsole sets the writer to envvars credentials store
func WithConsole(console ConsoleWriter) OptionsFunc {
	return func(s *EnvvarsStore) {
		s.console = console
	}
}

// WithFormater sets the formater to envvars credentials store
func WithFormater(formater repository.Formater) OptionsFunc {
	return func(s *EnvvarsStore) {
		s.formater = formater
	}
}

func WithEncryption(e Encrypter) OptionsFunc {
	return func(s *EnvvarsStore) {
		s.encryption = e
	}
}

// Options provides the options to envvars credentials store
func (s *EnvvarsStore) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(s)
	}
}

// Store stores a credential
func (s *EnvvarsStore) Store(id string, credential *credentials.Credential) error {

	errContext := "(store::credentials::envvars::Store)"

	var err error
	var hashedID string
	var strCredential string

	if s.console == nil {
		return errors.New(errContext, "Envvars credentials store requires a console writer to store a credential")
	}

	if s.formater == nil {
		return errors.New(errContext, "Envvars credentials store requires a formater to store a credential")
	}

	if s.encryption == nil {
		return errors.New(errContext, "Envvars credentials store requires encryption to store a credential")
	}

	if id == "" {
		return errors.New(errContext, "To store credentials badege, is required an ID")
	}

	if credential.ID == "" {
		credential.ID = id
	}

	strCredential, err = s.formater.Marshal(credential)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error formating '%s''s badege", id), err)
	}

	strCredential, err = s.encryption.Encrypt(strCredential)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	hashedID, err = encryption.HashID(id)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error hashing the id '%s'", id), err)
	}

	key := generateEnvvarKey(envvarsCredentialsPrefix, hashedID)
	s.console.Warn("You must create the following environment variable to use the recently created credentials:")
	s.console.Warn(fmt.Sprintf(" %s=%s", key, strCredential))

	return nil
}

// Get returns a auth for the credential id
func (s *EnvvarsStore) Get(id string) (*credentials.Credential, error) {
	errContext := "(store::credentials::envvars::Get)"

	var err error
	var hashedID string
	var credential *credentials.Credential

	if id == "" {
		return nil, errors.New(errContext, "To get credentials credential, is required an ID")
	}

	hashedID, err = encryption.HashID(id)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error hashing the id '%s'", id), err)
	}

	credential, err = s.get(hashedID)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error getting credentials credential '%s'", id), err)
	}

	return credential, nil
}

func (s *EnvvarsStore) get(hashedID string) (*credentials.Credential, error) {
	var err error
	var key, encryptedCredential, strCredential string
	var credential *credentials.Credential

	errContext := "(store::credentials::envvars::get)"
	if s.backend == nil {
		return nil, errors.New(errContext, "Envvars credentials store requires a backend to get credentials credential")
	}

	if s.formater == nil {
		return nil, errors.New(errContext, "Envvars credentials store requires a formater to get credentials credential")
	}

	if s.encryption == nil {
		return nil, errors.New(errContext, "Envvars credentials store requires encryption to get credentials credential")
	}

	key = generateEnvvarKey(envvarsCredentialsPrefix, hashedID)
	encryptedCredential = s.backend.Getenv(key)

	if encryptedCredential == "" {
		return nil, nil
	}

	strCredential, err = s.encryption.Decrypt(encryptedCredential)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error decrypting the '%s''s credential", hashedID), err)
	}

	credential, err = s.formater.Unmarshal([]byte(strCredential))
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error unmarshaling the '%s''s credential", hashedID), err)
	}

	return credential, nil
}

// All returns all credentials
func (s *EnvvarsStore) All() ([]*credentials.Credential, error) {
	errContext := "(store::credentials::envvars::All)"
	credentials := []*credentials.Credential{}
	IDs := map[string]struct{}{}
	allVars := s.backend.Environ()
	prefix := generateEnvvarKey(envvarsCredentialsPrefix)

	for _, envvar := range allVars {
		if strings.HasPrefix(envvar, prefix) {
			envvarTokens := strings.Split(envvar, "=")
			if len(envvarTokens) != 2 {
				continue
			}

			id := envvarTokens[0][len(prefix)+1:]
			IDs[id] = struct{}{}
		}
	}

	for id := range IDs {
		credential, err := s.get(id)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		credentials = append(credentials, credential)
	}

	return credentials, nil
}

func generateEnvvarKey(items ...string) string {
	var key string

	envvarKeySeparator := "_"
	replacements := []string{".", ":", "-"}

	key = strings.Join(append([]string{}, items...), envvarKeySeparator)

	for _, toReplace := range replacements {
		key = strings.ReplaceAll(key, toReplace, envvarKeySeparator)
	}

	key = strings.ToUpper(key)

	return key
}
