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
	console   ConsoleWriter
	backend   EnvvarsBackender
	encyption Encrypter
	formater  repository.Formater
	//loader *viper.Viper
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
		s.encyption = e
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

	if s.encyption == nil {
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

	strBadge, err = s.encyption.Encrypt(strBadge)
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
	var hashedID, key, encryptedBadge, strBadge string
	var badge *credentials.Badge

	if s.backend == nil {
		return nil, errors.New(errContext, "Envvars credentials store requires a backend to get credentials badge")
	}

	if s.formater == nil {
		return nil, errors.New(errContext, "Envvars credentials store requires a formater to get credentials badge")
	}

	if s.encyption == nil {
		return nil, errors.New(errContext, "Envvars credentials store requires encryption to get credentials badge")
	}

	if id == "" {
		return nil, errors.New(errContext, "To get credentials badge, is required an ID")
	}

	hashedID, err = encryption.HashID(id)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error hashing the id '%s'", id), err)
	}
	key = generateEnvvarKey(envvarsCredentialsPrefix, hashedID)
	encryptedBadge = s.backend.Getenv(key)

	strBadge, err = s.encyption.Decrypt(encryptedBadge)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error decrypting the '%s''s badge", id), err)
	}

	badge, err = s.formater.Unmarshal([]byte(strBadge))
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error unmarshaling the '%s''s badge", id), err)
	}

	return badge, nil
}

// All returns all badges
func (s *EnvvarsStore) All() []*credentials.Badge {

	badges := []*credentials.Badge{}
	IDs := map[string]struct{}{}
	allVars := s.backend.Environ()
	prefix := generateEnvvarKey(envvarsCredentialsPrefix)

	for _, envvar := range allVars {
		if strings.HasPrefix(envvar, prefix) {
			id := strings.Split(strings.ToLower(envvar), envvarsCredentialsAttributePrefix)[0][len(envvarsCredentialsPrefix)+1 : len(strings.Split(strings.ToLower(envvar), envvarsCredentialsAttributePrefix)[0])-1]
			IDs[id] = struct{}{}
		}
	}

	for id, _ := range IDs {
		badge, _ := s.Get(id)
		badges = append(badges, badge)
	}

	return badges
}

// func convertFieldValueToString(field interface{}) (string, error) {
// 	errContext := "(store::credentials::envvars::convertFieldValueToString)"

// 	switch field.(type) {
// 	case string:
// 		return fmt.Sprintf("%s", field), nil
// 	case []string:
// 		return strings.Join(field.([]string), ","), nil
// 	case bool:
// 		val := "0"
// 		if field.(bool) {
// 			val = "1"
// 		}
// 		return fmt.Sprintf("%s", val), nil
// 	case int:
// 		return fmt.Sprintf("%d", field), nil
// 	default:
// 		return "", errors.New(errContext, "Field could not be converted to string")
// 	}
// }

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
