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

// Store stores a badge
func (s *EnvvarsStore) Store(id string, badge *credentials.Badge) error {

	errContext := "(store::credentials::envvars::Store)"

	var err error
	var hashedID string
	var strBadge string

	if s.console == nil {
		return errors.New(errContext, "Envvars credentials store requires a console writer to store a badge")
	}

	if s.formater == nil {
		return errors.New(errContext, "Envvars credentials store requires a formater to store a badge")
	}

	if s.encryption == nil {
		return errors.New(errContext, "Envvars credentials store requires encryption to store a badge")
	}

	if id == "" {
		return errors.New(errContext, "To store credentials badege, is required an ID")
	}

	if badge.ID == "" {
		badge.ID = id
	}

	strBadge, err = s.formater.Marshal(badge)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error formating '%s''s badege", id), err)
	}

	strBadge, err = s.encryption.Encrypt(strBadge)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	hashedID, err = encryption.HashID(id)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error hashing the id '%s'", id), err)
	}

	key := generateEnvvarKey(envvarsCredentialsPrefix, hashedID)
	s.console.Warn("You must create the following environment variable to use the recently created credentials:")
	s.console.Warn(fmt.Sprintf(" %s=%s", key, strBadge))

	return nil
}

// Get returns a auth for the badge id
func (s *EnvvarsStore) Get(id string) (*credentials.Badge, error) {
	errContext := "(store::credentials::envvars::Get)"

	var err error
	var hashedID string
	var badge *credentials.Badge

	if id == "" {
		return nil, errors.New(errContext, "To get credentials badge, is required an ID")
	}

	hashedID, err = encryption.HashID(id)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error hashing the id '%s'", id), err)
	}

	badge, err = s.get(hashedID)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error getting credentials badge '%s'", id), err)
	}

	return badge, nil
}

func (s *EnvvarsStore) get(hashedID string) (*credentials.Badge, error) {
	var err error
	var key, encryptedBadge, strBadge string
	var badge *credentials.Badge

	errContext := "(store::credentials::envvars::get)"
	if s.backend == nil {
		return nil, errors.New(errContext, "Envvars credentials store requires a backend to get credentials badge")
	}

	if s.formater == nil {
		return nil, errors.New(errContext, "Envvars credentials store requires a formater to get credentials badge")
	}

	if s.encryption == nil {
		return nil, errors.New(errContext, "Envvars credentials store requires encryption to get credentials badge")
	}

	key = generateEnvvarKey(envvarsCredentialsPrefix, hashedID)
	encryptedBadge = s.backend.Getenv(key)

	if encryptedBadge == "" {
		return nil, nil
	}

	strBadge, err = s.encryption.Decrypt(encryptedBadge)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error decrypting the '%s''s badge", hashedID), err)
	}

	badge, err = s.formater.Unmarshal([]byte(strBadge))
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error unmarshaling the '%s''s badge", hashedID), err)
	}

	return badge, nil
}

// All returns all badges
func (s *EnvvarsStore) All() ([]*credentials.Badge, error) {
	errContext := "(store::credentials::envvars::All)"
	badges := []*credentials.Badge{}
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
		badge, err := s.get(id)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		badges = append(badges, badge)
	}

	return badges, nil
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
