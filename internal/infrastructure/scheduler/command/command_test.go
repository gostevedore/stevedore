package command

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/mock"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {

	tests := []struct {
		desc              string
		command           *BuildCommand
		prepareAssertFunc func(command *BuildCommand)
		assertFunc        func(*testing.T, *BuildCommand)
		err               error
	}{
		{
			desc:    "Testing error when image is nil",
			command: &BuildCommand{},
			err:     errors.New("(command::Execute)", "An image is required to execute a command"),
		},
		{
			desc: "Testing error when options are nil",
			command: &BuildCommand{
				image: &image.Image{},
			},
			err: errors.New("(command::Execute)", "Options are required to execute a command"),
		},
		{
			desc: "Testing execute command",
			command: &BuildCommand{
				image:   &image.Image{},
				options: &image.BuildDriverOptions{},
				driver:  mock.NewMockDriver(),
			},
			prepareAssertFunc: func(command *BuildCommand) {
				command.driver.(*mock.MockDriver).On("Build", context.TODO(), command.image, command.options).Return(nil)
			},
			assertFunc: func(t *testing.T, command *BuildCommand) {
				assert.True(t, command.driver.(*mock.MockDriver).AssertExpectations(t))
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.command)
			}

			err := test.command.Execute(context.TODO())
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.command)
			}
		})
	}

}
