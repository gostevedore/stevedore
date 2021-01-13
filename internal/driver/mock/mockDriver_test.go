package mockdriver

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestNewMockDriver(t *testing.T) {
	tests := []struct {
		desc string
		err  error
	}{
		{
			desc: "Dummy test",
			err:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			_, err := NewMockDriver(nil, nil)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			}
		})
	}
}

func TestMockDriverRun(t *testing.T) {
	tests := []struct {
		desc    string
		err     error
		builder *MockDriver
	}{
		{
			desc:    "Dummy test",
			err:     nil,
			builder: &MockDriver{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.builder.Run()
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			}
		})
	}
}

func TestNewMockDriverRunErr(t *testing.T) {
	tests := []struct {
		desc string
		err  error
	}{
		{
			desc: "Dummy test",
			err:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			_, err := NewMockDriverErr(nil, nil)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			}
		})
	}
}

func TestMockDriverErrRun(t *testing.T) {
	tests := []struct {
		desc    string
		err     error
		builder *MockDriverErr
	}{
		{
			desc:    "Dummy test",
			err:     errors.New("(MockDriverRunErr)", "Error"),
			builder: &MockDriverErr{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.builder.Run()
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			}
		})
	}
}
