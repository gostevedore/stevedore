package audit

// Logger interface to log errors
type Logger interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
	Fatal(msg ...interface{})
	Panic(msg ...interface{})
}
