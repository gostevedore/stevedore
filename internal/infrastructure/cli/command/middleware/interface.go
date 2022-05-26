package middleware

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
