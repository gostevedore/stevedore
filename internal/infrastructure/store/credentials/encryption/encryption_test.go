package encryption

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	var encryptedText, decryptedText, text string
	var err error
	enc := NewEncryption(
		WithKey("encryption-key"),
	)
	text = "here you have a testing message"

	encryptedText, err = enc.Encrypt(text)
	assert.NoError(t, err)

	decryptedText, err = enc.Decrypt(encryptedText)
	assert.NoError(t, err)

	assert.Equal(t, text, decryptedText)
}

func TestEncrypt(t *testing.T) {
	errContext := "(store::credentials::encryption::Encrypt)"

	tests := []struct {
		desc       string
		encryption Encryption
		input      string
		err        error
	}{
		{
			desc:       "Testing error in encryption when key is not provided to encrypt a message",
			encryption: NewEncryption(),
			err:        errors.New(errContext, "Encryption key must be provided to encrypt a message"),
		},
		{
			desc: "Testing credentials encryption",
			encryption: NewEncryption(
				WithKey("key"),
			),
			input: "plaintext",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.encryption.Encrypt(test.input)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				// since encypted text is not always the same, it can't be compared with a fixed value
				assert.NotEmpty(t, res)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	errContext := "(store::credentials::encryption::Decrypt)"

	tests := []struct {
		desc       string
		encryption Encryption
		input      string
		res        string
		err        error
	}{
		{
			desc:       "Testing error in encryption when key is not provided to decrypt a message",
			encryption: NewEncryption(),
			err:        errors.New(errContext, "Encryption key must be provided to decrypt a message"),
		},
		{
			desc: "Testing credentials decryption",
			encryption: NewEncryption(
				WithKey("key"),
			),
			input: "e50990a67be331277dc50b5dcfdd630eef07be89b070d1b5e2ee3091454ab26153811c2ad5",
			res:   "plaintext",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.encryption.Decrypt(test.input)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestHashID(t *testing.T) {

	errContext := "(store::credentials::encryption::hashID)"
	tests := []struct {
		desc string
		id   string
		res  string
		err  error
	}{
		{
			desc: "Testing error when hashing an id with providing the id",
			id:   "",
			err:  errors.New(errContext, "Hash method requires an id"),
		},
		{
			desc: "Testing hashing an id",
			id:   "id",
			res:  "b80bb7740288fda1f201890375a60c8f",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := HashID(test.id)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestGenerateEncryptionKey(t *testing.T) {
	enc := NewEncryption()
	key, err := enc.GenerateEncryptionKey()
	assert.NoError(t, err)
	assert.Equal(t, len(key), 32)
}
