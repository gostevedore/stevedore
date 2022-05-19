package panic

// Logger interface to log errors
type Logger interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
	Fatal(msg ...interface{})
	Panic(msg ...interface{})
}

// Consoler interface to show messages through console
type Consoler interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
}
