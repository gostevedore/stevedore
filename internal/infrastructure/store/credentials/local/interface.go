package local

import "github.com/gostevedore/stevedore/internal/core/domain/credentials"

type CredentialsCompatibilier interface {
	CheckCompatibility(credential *credentials.Credential) error
}

type Encrypter interface {
	Encrypt(text string) (string, error)
	Decrypt(ciphertext string) (string, error)
}
