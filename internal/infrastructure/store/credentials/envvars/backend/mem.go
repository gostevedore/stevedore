package backend

import "os"

// MemEnvvarsBackend is a backend for envvars that uses the OS envvars
type MemEnvvarsBackend struct {
	mem map[string]string
}

func _NewMemEnvvarsBackend() *MemEnvvarsBackend {
	return &MemEnvvarsBackend{
		mem: make(map[string]string),
	}
}

// Setenv sets the value of the environment variable named by the key. It returns an error, if any
func (b *MemEnvvarsBackend) Setenv(key, value string) {
	if b.mem == nil {
		b.mem = make(map[string]string)
	}

	os.Setenv(key, value)
}

// Getenv retrieves the value of the environment variable named by the key. It returns the value, which will be empty if the variable is not present. To distinguish between an empty value and an unset value, use LookupEnv
func (b *MemEnvvarsBackend) Getenv(key string) string {
	if b.mem == nil {
		b.mem = make(map[string]string)
		return ""
	}

	return os.Getenv(key)
}

// LookupEnv retrieves the value of the environment variable named by the key. If the variable is set the value (which may be empty) is returned and the boolean is true. Otherwise the returned value will be empty and the boolean will be false
func (b *MemEnvvarsBackend) LookupEnv(key string) (string, bool) {
	if b.mem == nil {
		b.mem = make(map[string]string)
		return "", false
	}

	return os.LookupEnv(key)
}

// Environ returns a copy of strings representing the environment, in the form "key=value"
func (b *MemEnvvarsBackend) Environ() []string {
	if b.mem == nil {
		b.mem = make(map[string]string)
		return nil
	}

	return os.Environ()
}
