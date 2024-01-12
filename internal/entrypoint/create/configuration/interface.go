package configuration

import "io"

type EncryptionKeyGenerator interface {
	GenerateEncryptionKey() (string, error)
}

type Consoler interface {
	io.Writer
	ConsoleWriter
}

type ConsoleWriter interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
}
