package envvars

type ConsoleWriter interface {
	Info(msg ...interface{})
	Warn(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
}

type EnvvarsBackender interface {
	Getenv(key string) string
	Environ() []string
}

type Encrypter interface {
	Encrypt(text string) (string, error)
	Decrypt(ciphertext string) (string, error)
}
