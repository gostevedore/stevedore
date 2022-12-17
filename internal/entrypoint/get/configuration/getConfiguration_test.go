package configuration

import (
	"context"
	"testing"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		desc            string
		entrypoint      *GetConfigurationEntrypoint
		args            []string
		conf            *configuration.Configuration
		prepareMockFunc func()
		err             error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.entrypoint.Execute(context.TODO(), test.args, test.conf)
			if err != nil {
				assert.Equal(t, test.err, err)
			}
		})
	}
}
