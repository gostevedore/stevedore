package build

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/driver"
	mockdriver "github.com/gostevedore/stevedore/internal/driver/mock"
	"github.com/gostevedore/stevedore/internal/engine/build/command"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/schedule/dispatch"
	"github.com/gostevedore/stevedore/internal/schedule/worker"
	"github.com/gostevedore/stevedore/internal/semver"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {

}

func TestWorker(t *testing.T) {

	errContext := "(build::worker)"

	tests := []struct {
		desc              string
		service           *Service
		image             *image.Image
		options           *ServiceOptions
		err               error
		prepareAssertFunc func(*Service)
		assertFunc        func(*Service) bool
	}{
		{
			desc:    "Testing error when no options are given to worker",
			service: &Service{},
			options: nil,
			err:     errors.New(errContext, "Build worker requires service options"),
		},
		{
			desc:    "Testing error when no image specification is given to worker",
			service: &Service{},
			options: &ServiceOptions{},
			err:     errors.New(errContext, "Build worker requires an image specification"),
		},
		{
			desc: "Testing error when no dispatcher is given to worker",
			service: &Service{
				dispatch: nil,
			},
			options: &ServiceOptions{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a dispatcher"),
		},
		{
			desc: "Testing error when no driver factory is given to worker",
			service: &Service{
				dispatch: dispatch.NewDispatch(1, worker.NewMockWorkerFactory()),
			},
			options: &ServiceOptions{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a driver factory"),
		},
		{
			desc: "Testing error when no semantic version generator is given to worker",
			service: &Service{
				dispatch:      dispatch.NewDispatch(1, worker.NewMockWorkerFactory()),
				driverFactory: driver.NewBuildDriverFactory(),
			},
			options: &ServiceOptions{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a semver generator"),
		},
		{
			desc: "Testing error when no credentials store is given to worker",
			service: &Service{
				dispatch:      dispatch.NewDispatch(1, worker.NewMockWorkerFactory()),
				driverFactory: driver.NewBuildDriverFactory(),
				semver:        semver.NewSemVerGenerator(),
			},
			options: &ServiceOptions{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a credentials store"),
		},

		{
			desc: "Testing build an image",
			service: &Service{
				dispatch:      dispatch.NewDispatch(1, worker.NewMockWorkerFactory()),
				driverFactory: driver.NewBuildDriverFactory(),
				semver:        semver.NewSemVerGenerator(),
				credentials:   credentials.NewCredentialsStoreMock(),
			},
			options:           &ServiceOptions{},
			image:             &image.Image{},
			err:               &errors.Error{},
			prepareAssertFunc: func(service *Service) {},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service)
			}

			err := test.service.worker(context.TODO(), test.image, test.options)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.True(t, test.assertFunc(test.service))
			}
		})
	}

}

// func TestGenerateImagesList(t *testing.T) {
// 	errContext := "(build::generateImagesList)"

// 	tests := []struct {
// 		desc              string
// 		service           *Service
// 		name              string
// 		versions          []string
// 		res               []*image.Image
// 		prepareAssertFunc func(*Service)
// 		assertFunc        func(*Service) bool
// 		err               error
// 	}{
// 		{
// 			desc: "Testing error when no image name is provided",
// 			name: "",
// 			err:  errors.New(errContext, "Image name is required to build an image"),
// 		},
// 		{
// 			desc: "Testing generate images list",
// 			service: &Service{
// 				images: imagestore.NewMockImageStore(),
// 			},
// 			name:     "image",
// 			versions: []string{"version1", "version2"},
// 			res: []*image.Image{
// 				{
// 					Name:    "image",
// 					Version: "version1",
// 				},
// 				{
// 					Name:    "image",
// 					Version: "version2",
// 				},
// 			},
// 			err: &errors.Error{},
// 			prepareAssertFunc: func(s *Service) {
// 				s.images.(*imagestore.MockImageStore).On("Find", "image", "version1").Return(&image.Image{
// 					Name:    "image",
// 					Version: "version1",
// 				}, nil)
// 				s.images.(*imagestore.MockImageStore).On("Find", "image", "version2").Return(&image.Image{
// 					Name:    "image",
// 					Version: "version2",
// 				}, nil)
// 			},
// 			assertFunc: func(s *Service) bool {
// 				return s.images.(*imagestore.MockImageStore).AssertExpectations(t)
// 			},
// 		},

// 		{
// 			desc: "Testing generate images list when no version is provided",
// 			service: &Service{
// 				images: imagestore.NewMockImageStore(),
// 			},
// 			name:     "image",
// 			versions: []string{},
// 			res: []*image.Image{
// 				{
// 					Name:    "image",
// 					Version: "version1",
// 				},
// 				{
// 					Name:    "image",
// 					Version: "version2",
// 				},
// 			},
// 			err: &errors.Error{},
// 			prepareAssertFunc: func(s *Service) {
// 				s.images.(*imagestore.MockImageStore).On("All", "image").Return([]*image.Image{
// 					{
// 						Name:    "image",
// 						Version: "version1",
// 					},
// 					{
// 						Name:    "image",
// 						Version: "version2",
// 					},
// 				}, nil)
// 			},
// 			assertFunc: func(s *Service) bool {
// 				return s.images.(*imagestore.MockImageStore).AssertExpectations(t)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			if test.prepareAssertFunc != nil {
// 				test.prepareAssertFunc(test.service)
// 			}

// 			res, err := test.service.generateImagesList(test.name, test.versions)

// 			if err != nil {
// 				assert.Equal(t, test.err.Error(), err.Error())
// 			} else {
// 				assert.True(t, test.assertFunc(test.service))
// 				assert.Equal(t, test.res, res)
// 			}
// 		})
// 	}

// }

func TestCommand(t *testing.T) {
	errContext := "(build::command)"

	driverFactory := driver.NewBuildDriverFactory()
	driverFactory.Register("mock", mockdriver.NewMockDriver())

	tests := []struct {
		desc              string
		service           *Service
		driver            driver.BuildDriverer
		options           *driver.BuildDriverOptions
		res               BuildCommander
		prepareAssertFunc func(*Service)
		err               error
	}{
		{
			desc:    "Testing error when no command factory is provided",
			service: &Service{},
			err:     errors.New(errContext, "To create a build command, is required a command factory"),
		},
		{
			desc: "Testing error when no driver is provided",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
			},
			err: errors.New(errContext, "To create a build command, is required a driver"),
		},
		{
			desc: "Testing error when no options are provided",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
			},
			driver: mockdriver.NewMockDriver(),
			err:    errors.New(errContext, "To create a build command, is required a service options"),
		},

		{
			desc: "Testing create build command",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
				driverFactory:  driverFactory,
				builders:       builders.NewMockBuilders(),
			},
			options: &driver.BuildDriverOptions{},
			driver:  mockdriver.NewMockDriver(),
			prepareAssertFunc: func(s *Service) {
				s.builders.(*builders.MockBuilders).On("Find", "test").Return(&builder.Builder{
					Name:   "test",
					Driver: "mock", // during test, we use the mock driver
				}, nil)
			},
			res: &command.MockBuildCommand{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service)
			}

			res, err := test.service.command(test.driver, test.options)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
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
