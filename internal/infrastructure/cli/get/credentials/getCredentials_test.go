package credentials

import (
	"context"
	"testing"

	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		desc            string
		config          *configuration.Configuration
		entrypoint      Entrypointer
		prepareMockFunc func(Entrypointer, *configuration.Configuration)
		args            []string
		err             error
	}{
		{
			desc:       "Testing run promote command",
			config:     &configuration.Configuration{},
			entrypoint: entrypoint.NewMockEntrypoint(),
			args:       []string{},
			prepareMockFunc: func(ep Entrypointer, config *configuration.Configuration) {
				ep.(*entrypoint.MockEntrypoint).On("Execute", context.TODO(), []string{}, config).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil {
				test.prepareMockFunc(test.entrypoint, test.config)
			}

			cmd := NewCommand(context.TODO(), test.config, test.entrypoint)
			cmd.Command.ParseFlags(test.args)
			err := cmd.Command.RunE(cmd.Command, test.args)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			}

		})
	}
}
