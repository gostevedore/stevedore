package credentials

import (
	"io"
)

type Consoler interface {
	io.Writer
	ConsoleWriter
	PasswordReader
}

type ConsoleWriter interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
}

type PasswordReader interface {
	ReadPassword(prompt string) (string, error)
}

// Compatibilitier is the interface for the compatibility checker
type Compatibilitier interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}
