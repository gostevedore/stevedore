package compatibility

// Consoler is the interface for the console output
type Consoler interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
}
