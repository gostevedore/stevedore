package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"

	errors "github.com/apenella/go-common-utils/error"
)

// OptionsFunc defines the signature for an option function to set encryption
type OptionsFunc func(opts *Encryption)

type Encryption struct {
	key string
}

func NewEncryption(opts ...OptionsFunc) Encryption {
	e := &Encryption{}
	e.Options(opts...)
	return *e
}

// Options provides the options to encryption
func (e *Encryption) Options(opts ...OptionsFunc) {
	for _, opt := range opts {
		opt(e)
	}
}

// WithKey sets the encryption key
func WithKey(key string) OptionsFunc {
	return func(e *Encryption) {
		e.key = key
	}
}

// Encrypt return the input text encripted
func (e Encryption) Encrypt(text string) (string, error) {

	var err error
	var key []byte
	var block cipher.Block
	var gcm cipher.AEAD

	errContext := "(store::credentials::encryption::Encrypt)"
	if e.key == "" {
		return "", errors.New(errContext, "Encryption key must be provided to encrypt a message")
	}

	key, err = hashKey(e.key)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	block, err = aes.NewCipher(key)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New(errContext, "", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)

	return hex.EncodeToString(ciphertext), nil
}

// Decrypt return the input text decrypted
func (e Encryption) Decrypt(ciphertext string) (string, error) {

	var err error
	var key, enc []byte
	var block cipher.Block
	var gcm cipher.AEAD

	errContext := "(store::credentials::encryption::Decrypt)"
	if e.key == "" {
		return "", errors.New(errContext, "Encryption key must be provided to decrypt a message")
	}

	key, err = hashKey(e.key)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	block, err = aes.NewCipher(key)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	enc, err = hex.DecodeString(ciphertext)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, bytedCiphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, bytedCiphertext, nil)
	if err != nil {
		return "", errors.New(errContext, "", err)
	}

	return string(plaintext), nil

}

func hashKey(key string) ([]byte, error) {

	errContext := "(store::credentials::encryption::hashKey)"
	if key == "" {
		return nil, errors.New(errContext, "To get the hash for the key, the key must be provided")
	}

	hashFunc := sha256.New()
	hashFunc.Write([]byte(key))

	return hashFunc.Sum(nil), nil
}

// HashID generates a hash for the id
func HashID(id string) (string, error) {

	errContext := "(store::credentials::encryption::hashID)"

	if id == "" {
		return "", errors.New(errContext, "Hash method requires an id")
	}

	hashFunc := md5.New()
	hashFunc.Write([]byte(id))
	registryHashed := hex.EncodeToString(hashFunc.Sum(nil))

	return registryHashed, nil
}
