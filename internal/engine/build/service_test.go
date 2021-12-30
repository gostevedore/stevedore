package build

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/gostevedore/stevedore/internal/driver"
	mockdriver "github.com/gostevedore/stevedore/internal/driver/mock"
	"github.com/gostevedore/stevedore/internal/engine/build/command"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {

}

func TestWorker(t *testing.T) {

}

func TestGenerateImagesList(t *testing.T) {

}

func TestJob(t *testing.T) {

}

func TestCommand(t *testing.T) {
	errContext := "(build::command)"

	driverFactory := driver.NewBuildDriverFactory()
	driverFactory.Register("mock", mockdriver.NewMockDriver())

	tests := []struct {
		desc              string
		service           *Service
		image             *image.Image
		options           *ServiceOptions
		res               BuildCommander
		prepareAssertFunc func(*Service)
		assertFunc        func(expected, actual BuildCommander) bool
		err               error
	}{
		{
			desc:    "Testing error when no command factory is provided",
			service: &Service{},
			err:     errors.New(errContext, "To create a build command, is required a command factory"),
		},
		{
			desc: "Testing error when no driver factory is provided",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
			},
			err: errors.New(errContext, "To create a build command, is required a driver factory"),
		},
		{
			desc: "Testing error when no image is provided",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
				driverFactory:  driverFactory,
			},
			err: errors.New(errContext, "To create a build command, is required an image"),
		},
		{
			desc: "Testing error when no service options are provided",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
				driverFactory:  driverFactory,
			},
			image: &image.Image{},
			err:   errors.New(errContext, "To create a build command, is required a service options"),
		},
		{
			desc: "Testing error when there is an error preparing the builder",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
				driverFactory:  driverFactory,
			},
			image:   &image.Image{},
			options: &ServiceOptions{},
			err:     errors.New(errContext, "To generate a builder, is required a builder definition"),
		},
		{
			desc: "Testing create build command",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
				driverFactory:  driverFactory,
				builders:       builders.NewMockBuilders(),
			},
			image: &image.Image{
				Builder: "test",
			},
			options: &ServiceOptions{},
			prepareAssertFunc: func(s *Service) {
				s.builders.(*builders.MockBuilders).On("Find", "test").Return(&builder.Builder{
					Name:   "test",
					Driver: "mock", // during test, we use the mock driver
				}, nil)
			},
			assertFunc: func(expected, actual BuildCommander) bool {
				return false
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service)
			}

			res, err := test.service.command(context.TODO(), test.image, test.options)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				if test.assertFunc != nil {
					assert.True(t, test.assertFunc(test.res, res))
				}
			}

		})
	}

}

func TestBuilder(t *testing.T) {
	errContext := "(build::builder)"

	tests := []struct {
		desc              string
		service           *Service
		builderDefinition interface{}
		res               *builder.Builder
		prepareAssertFunc func(*Service)
		assertFunc        func(expected, actual *builder.Builder) bool
		err               error
	}{
		{
			desc:              "Testing error getting a builder with no builder definition",
			service:           &Service{},
			builderDefinition: nil,
			err:               errors.New(errContext, "To generate a builder, is required a builder definition"),
		},
		{
			desc:              "Testing error getting a builder with no builders store",
			service:           &Service{},
			builderDefinition: "test",
			err:               errors.New(errContext, "To generate a builder, is required a builder store defined on build service"),
		},
		{
			desc:              "Testing return a builder defined by an string",
			service:           &Service{builders: builders.NewMockBuilders()},
			builderDefinition: "test",
			prepareAssertFunc: func(s *Service) {
				s.builders.(*builders.MockBuilders).On("Find", "test").Return(&builder.Builder{
					Name: "test",
				}, nil)
			},
			assertFunc: func(expected, actual *builder.Builder) bool {
				//return s.builders.(*builders.MockBuilders).AssertExpectations(t)
				return assert.Equal(t, expected, actual)

			},
			res: &builder.Builder{
				Name: "test",
			},
		},
		// That test to be tested from the caller because is not possible to force builderDerfinetion to be seen as an interface
		// 		{
		// 			desc:    "Testing return a builder defined by an interface{}",
		// 			service: &Service{builders: builders.NewMockBuilders()},
		// 			builderDefinition: `
		// driver: docker
		// options:
		//     dockerfile: Dockerfile.test
		//     context:
		//     - git:
		//         path: path
		//         repository: repository
		//         reference: reference
		//         auth:
		//             username: username
		//             password: password
		// `,
		// 			prepareAssertFunc: nil,
		// 			assertFunc: func(expected, actual *builder.Builder) bool {
		// 				//return s.builders.(*builders.MockBuilders).AssertExpectations(t)
		// 				return assert.Equal(t, expected, actual)
		// 			},
		// 			res: &builder.Builder{
		// 				Name:   "",
		// 				Driver: "docker",
		// 				Options: &builder.BuilderOptions{
		// 					Context: []*builder.DockerDriverContextOptions{
		// 						{
		// 							Git: &builder.DockerDriverGitContextOptions{
		// 								Path:       "path",
		// 								Repository: "repository",
		// 								Reference:  "reference",
		// 								Auth: &builder.DockerDriverGitContextAuthOptions{
		// 									Username: "username",
		// 									Password: "password",
		// 								},
		// 							},
		// 						},
		// 					},
		// 					Dockerfile: "Dockerfile.test",
		// 				},
		// 			},
		// 		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service)
			}

			res, err := test.service.builder(test.builderDefinition)
			_ = res
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				if test.assertFunc != nil {
					assert.True(t, test.assertFunc(test.res, res))
				}
			}
		})
	}
}
