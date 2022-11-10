package images

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	entrypoint "github.com/gostevedore/stevedore/internal/entrypoint/get/images"
	handler "github.com/gostevedore/stevedore/internal/handler/get/images"
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
			desc:       "Testing run get images command",
			config:     &configuration.Configuration{},
			entrypoint: entrypoint.NewMockGetImagesEntrypoint(),
			args: []string{
				"--tree",
				"--filter",
				"name=a",
				"--filter",
				"builder=b",
			},
			prepareMockFunc: func(e Entrypointer, c *configuration.Configuration) {
				e.(*entrypoint.MockGetImagesEntrypoint).On(
					"Execute",
					context.TODO(),
					[]string{},
					c,
					&entrypoint.Options{
						Tree: true,
					},
					&handler.Options{
						Filter: []string{
							"name=a",
							"builder=b",
						},
					},
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
			} else {
				test.entrypoint.(*entrypoint.MockGetImagesEntrypoint).AssertExpectations(t)
			}
		})
	}
}
