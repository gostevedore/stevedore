package envvars

import (
	"fmt"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/fatih/structs"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
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
	console ConsoleWriter
	backend EnvvarsBackender
	//loader *viper.Viper
}

// NewEnvvarsStore creates a new mocked store for credentials
func NewEnvvarsStore(opts ...OptionsFunc) *EnvvarsStore {
	store := &EnvvarsStore{}
	store.Options(opts...)

	return store
}

// WithConsole sets the writer to envvars credentials store
func WithConsole(console ConsoleWriter) OptionsFunc {
	return func(s *EnvvarsStore) {
		s.console = console
	}
}

// WithBackend sets the writer to envvars credentials store
func WithBackend(backend EnvvarsBackender) OptionsFunc {
	return func(s *EnvvarsStore) {
		s.backend = backend
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
	output := []string{}

	if s.console == nil {
		return errors.New(errContext, "Envvars credentials store requires a console writer")
	}

	if s.backend == nil {
		return errors.New(errContext, "Envvars credentials store requires a backend to store envvars")
	}

	badgeFields := structs.Fields(badge)
	for _, field := range badgeFields {
		if !field.IsZero() {
			attribute := field.Tag("mapstructure")
			if attribute == "" {
				continue
			}

			value, err := convertFieldValueToString(field.Value())
			if err != nil {
				return errors.New(errContext, fmt.Sprintf("Error converting to '%s''s value", attribute), err)
			}

			if value != "" {
				key := generateEnvvarKey(envvarsCredentialsPrefix, id, envvarsCredentialsAttributePrefix, attribute)
				output = append(output, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}

	if len(output) > 0 {
		s.console.Warn("You must create the following environment variables to use the recently created credentials:")
		for _, key := range output {
			s.console.Warn(fmt.Sprintf(" %s", key))
		}
	}

	return nil
}

// Get returns a auth for the badge id
func (s *EnvvarsStore) Get(id string) (*credentials.Badge, error) {
	errContext := "(store::credentials::envvars::Get)"

	badge, err := s.backend.AchieveBadge(generateEnvvarKey(envvarsCredentialsPrefix, id, envvarsCredentialsAttributePrefix))
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}
	if badge.ID == "" {
		badge.ID = id
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

func convertFieldValueToString(field interface{}) (string, error) {
	errContext := "(store::credentials::envvars::convertFieldValueToString)"

	switch field.(type) {
	case string:
		return fmt.Sprintf("%s", field), nil
	case []string:
		return strings.Join(field.([]string), ","), nil
	case bool:
		val := "0"
		if field.(bool) {
			val = "1"
		}
		return fmt.Sprintf("%s", val), nil
	case int:
		return fmt.Sprintf("%d", field), nil
	default:
		return "", errors.New(errContext, "Field could not be converted to string")
	}
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
