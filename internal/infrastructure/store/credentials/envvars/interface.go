package envvars

import "github.com/gostevedore/stevedore/internal/core/domain/credentials"

type ConsoleWriter interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
}

type EnvvarsBackender interface {
	Setenv(key, value string)
	//Getenv(key string) string
	//LookupEnv(key string) (string, bool)
	Environ() []string
	AchieveBadge(id string) (*credentials.Badge, error)
}
