package promote

// Compatibilitier is the interface for the compatibility checker
type Compatibilitier interface {
	AddDeprecated(deprecated ...string)
	AddRemoved(removed ...string)
	AddChanged(changed ...string)
}

type ConsoleWriter interface {
	Debug(msg ...interface{})
	Error(msg ...interface{})
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Write(data []byte) (int, error)
}
