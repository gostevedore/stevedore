package configuration

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/configuration"
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
			desc:       "Testing run get builders command",
			config:     &configuration.Configuration{},
			entrypoint: entrypoint.NewMockGetConfigurationEntrypoint(),
			args:       []string{},
			prepareMockFunc: func(e Entrypointer, c *configuration.Configuration) {
				e.(*entrypoint.MockGetConfigurationEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{},
					c,
				).Return(nil)
			},
			err: &errors.Error{},
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
