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
	"github.com/gostevedore/stevedore/internal/schedule/job"
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
		// {
		// 	desc: "Testing error when no dispatcher is given to worker",
		// 	service: &Service{
		// 		dispatch: nil,
		// 	},
		// 	options: &ServiceOptions{},
		// 	image:   &image.Image{},
		// 	err:     errors.New(errContext, "Build worker requires a dispatcher"),
		// },
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
			service: NewService(
				nil,
				builders.NewMockBuilders(),
				command.NewMockBuildCommandFactory(),
				&driver.BuildDriverFactory{
					"mock": mockdriver.NewMockDriver(),
				},
				job.NewMockJobFactory(),
				dispatch.NewMockDispatch(),
				semver.NewSemVerGenerator(),
				credentials.NewCredentialsStoreMock(),
			),
			// service: &Service{
			// 	builders:       builders.NewMockBuilders(),
			// 	commandFactory: command.NewMockBuildCommandFactory(),
			// 	dispatch:       dispatch.NewDispatch(1, worker.NewMockWorkerFactory()),
			// 	driverFactory: driver.BuildDriverFactory{
			// 		"mock": mockdriver.NewMockDriver(),
			// 	},
			// 	jobFactory:  job.NewMockJobFactory(),
			// 	semver:      semver.NewSemVerGenerator(),
			// 	credentials: credentials.NewCredentialsStoreMock(),
			// },
			options: &ServiceOptions{
				EnableSemanticVersionTags:    true,
				PushImageAfterBuild:          true,
				PullParentImage:              true,
				SemanticVersionTagsTemplates: []string{"{{.Major}}", "{{.Major}}.{{.Minor}}"},
				RemoveAfterBuild:             true,
			},
			image: &image.Image{
				Name:              "image",
				Version:           "0.0.0",
				RegistryHost:      "registry",
				RegistryNamespace: "namespace",
				Builder:           "builder",
			},
			err: &errors.Error{},
			assertFunc: func(service *Service) bool {
				return service.credentials.(*credentials.CredentialsStoreMock).AssertExpectations(t) &&
					service.builders.(*builders.MockBuilders).AssertExpectations(t) &&
					service.commandFactory.(*command.MockBuildCommandFactory).AssertExpectations(t) &&
					service.dispatch.(*dispatch.MockDispatch).AssertExpectations(t) &&
					service.jobFactory.(*job.MockJobFactory).AssertExpectations(t)
			},
			prepareAssertFunc: func(service *Service) {

				mockJob := job.NewMockJob()
				mockJob.On("Wait").Return(nil)

				service.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "registry").Return(&credentials.RegistryUserPassAuth{
					Username: "user",
					Password: "pass",
				}, nil)

				service.builders.(*builders.MockBuilders).On("Find", "builder").Return(&builder.Builder{
					Driver: "mock",
				}, nil)
				service.commandFactory.(*command.MockBuildCommandFactory).On("New",
					mockdriver.NewMockDriver(),
					&driver.BuildDriverOptions{
						BuilderName:                "builder_mock_namespace_image_0.0.0",
						BuilderOptions:             nil,
						BuilderVarMappings:         nil,
						ConnectionLocal:            false,
						ImageFromName:              "",
						ImageFromRegistryNamespace: "",
						ImageFromRegistryHost:      "",
						ImageFromVersion:           "",
						ImageName:                  "image",
						ImageVersion:               "0.0.0",
						Labels:                     nil,
						OutputPrefix:               "",
						PersistentVars:             nil,
						RegistryNamespace:          "namespace",
						RegistryHost:               "registry",
						PullAuthUsername:           "",
						PullAuthPassword:           "",
						PullParentImage:            true,
						PushAuthUsername:           "user",
						PushAuthPassword:           "pass",
						PushImageAfterBuild:        true,
						RemoveImageAfterBuild:      true,
						Tags:                       []string{"0", "0.0"},
						Vars:                       nil,
					}).Return(command.NewMockBuildCommand(), nil)
				service.jobFactory.(*job.MockJobFactory).On("New", command.NewMockBuildCommand()).Return(mockJob, nil)
				service.dispatch.(*dispatch.MockDispatch).On("Enqueue", mockJob)
			},
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

			//time.Sleep(time.Millisecond * 100)
		})
	}

}

func TestJob(t *testing.T) {
	errContext := "(build::command)"

	tests := []struct {
		desc              string
		service           *Service
		cmd               job.Commander
		prepareAssertFunc func(*Service, job.Commander)
		err               error
	}{
		{
			desc:    "Testing error when no job factory is defined on service",
			service: &Service{},
			err:     errors.New(errContext, "To create a build job, is required a job factory"),
		},
		{
			desc: "Testing job creation",
			service: &Service{
				jobFactory: job.NewMockJobFactory(),
			},
			cmd: command.NewMockBuildCommand(),
			prepareAssertFunc: func(service *Service, cmd job.Commander) {
				service.jobFactory.(*job.MockJobFactory).On("New", cmd).Return(job.NewMockJob(), nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service, test.cmd)
			}

			_, err := test.service.job(context.TODO(), test.cmd)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.service.jobFactory.(*job.MockJobFactory).AssertExpectations(t)
			}

		})
	}
}

func TestCommand(t *testing.T) {
	errContext := "(build::command)"

	driverFactory := driver.NewBuildDriverFactory()
	driverFactory.Register("mock", mockdriver.NewMockDriver())

	tests := []struct {
		desc              string
		service           *Service
		driver            driver.BuildDriverer
		options           *driver.BuildDriverOptions
		res               job.Commander
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
			service: NewService(
				nil,
				nil,
				command.NewMockBuildCommandFactory(),
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			options: &driver.BuildDriverOptions{},
			driver:  mockdriver.NewMockDriver(),
			prepareAssertFunc: func(s *Service) {
				s.commandFactory.(*command.MockBuildCommandFactory).On("New", mockdriver.NewMockDriver(), &driver.BuildDriverOptions{}).Return(command.NewMockBuildCommand(), nil)
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
