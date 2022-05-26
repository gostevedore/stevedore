package version

// Consoler interface to show messages through console
type Consoler interface {
	Write(data []byte) (int, error)
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
}
