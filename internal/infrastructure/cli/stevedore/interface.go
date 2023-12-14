package stevedore

import "io"

// CompatibilityStorer is the interface for the compatibility checker
type CompatibilityStorer interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}

// CompatibilityReporter is the interface to report compatibilities
type CompatibilityReporter interface {
	Report()
}

// Logger interface to log errors
type Logger interface {
	ReloadWithWriter(w io.Writer)
	Sync() error
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
	Fatal(msg ...interface{})
	Panic(msg ...interface{})
}

// Consoler interface to show messages through console
type Consoler interface {
	ConsoleWriter
	ConsoleReader
}

type ConsoleWriter interface {
	Write(data []byte) (int, error)
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
}

type ConsoleReader interface {
	Read() string
	ReadPassword(prompt string) (string, error)
}
