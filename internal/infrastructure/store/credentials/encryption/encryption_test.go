package encryption

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	errContext := "(store::credentials::encryption::Encrypt)"

	tests := []struct {
		desc       string
		encryption Encryption
		input      string
		res        string
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
			res:   "3eiA9Ru1oRJkdZGaOyN1bQNERtAsbChquGQebMe6ygZzqNoLUA==",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			_, err := test.encryption.Encrypt(test.input)
			if err != nil {
				assert.Equal(t, test.err, err)
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
			input: "3eiA9Ru1oRJkdZGaOyN1bQNERtAsbChquGQebMe6ygZzqNoLUA==",
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
