package build

import (
	"context"
	"io/ioutil"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/gostevedore/stevedore/internal/builders/store"
	"github.com/gostevedore/stevedore/internal/builders/varsmap"
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/driver"
	dockerdriver "github.com/gostevedore/stevedore/internal/driver/docker"
	dryrundriver "github.com/gostevedore/stevedore/internal/driver/dryrun"
	mockdriver "github.com/gostevedore/stevedore/internal/driver/mock"
	"github.com/gostevedore/stevedore/internal/images/image"
	"github.com/gostevedore/stevedore/internal/schedule/dispatch"
	"github.com/gostevedore/stevedore/internal/schedule/job"
	"github.com/gostevedore/stevedore/internal/schedule/worker"
	"github.com/gostevedore/stevedore/internal/semver"
	"github.com/gostevedore/stevedore/internal/service/build/command"
	"github.com/gostevedore/stevedore/internal/service/build/plan"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	errContext := "(build::Build)"
	_ = errContext
	tests := []struct {
		desc              string
		service           *Service
		buildPlan         Planner
		name              string
		versions          []string
		options           *ServiceOptions
		prepareAssertFunc func(*Service, Planner)
		assertFunc        func(*Service) bool
		err               error
	}{
		{
			desc:    "Testing error building an image with no options",
			service: &Service{},
			options: nil,
			err:     errors.New(errContext, "To build an image, service options are required"),
		},
		{
			desc:    "Testing error building an image with no execution plan",
			service: &Service{},
			options: &ServiceOptions{},
			err:     errors.New(errContext, "To build an image, a build plan is required"),
		},
		{
			desc: "Testing build an image",
			service: NewService(
				WithBuilders(store.NewMockBuildersStore()),
				WithCommandFactory(command.NewMockBuildCommandFactory()),
				WithDriverFactory(
					&driver.BuildDriverFactory{
						"mock": mockdriver.NewMockDriver(),
					},
				),
				WithJobFactory(job.NewMockJobFactory()),
				WithDispatch(dispatch.NewMockDispatch()),
				WithSemver(semver.NewSemVerGenerator()),
				WithCredentials(credentials.NewCredentialsStoreMock()),
			),
			buildPlan: plan.NewMockPlan(),
			name:      "parent",
			versions:  []string{"0.0.0"},
			options: &ServiceOptions{
				AnsibleConnectionLocal:           true,
				AnsibleIntermediateContainerName: "intermediate_container",
				AnsibleInventoryPath:             "inventory",
				AnsibleLimit:                     "limit",
				EnableSemanticVersionTags:        true,
				PushImageAfterBuild:              true,
				PullParentImage:                  true,
				SemanticVersionTagsTemplates:     []string{"{{.Major}}"},
				RemoveImagesAfterPush:            true,
			},
			err: &errors.Error{},
			assertFunc: func(service *Service) bool {
				return service.credentials.(*credentials.CredentialsStoreMock).AssertExpectations(t) &&
					service.commandFactory.(*command.MockBuildCommandFactory).AssertExpectations(t) &&
					service.dispatch.(*dispatch.MockDispatch).AssertExpectations(t) &&
					service.jobFactory.(*job.MockJobFactory).AssertExpectations(t)
			},
			prepareAssertFunc: func(service *Service, buildPlan Planner) {

				mockJob := job.NewMockJob()
				mockJob.On("Wait").Return(nil)

				childSyncChan := make(chan struct{})
				stepChild := plan.NewStep(
					&image.Image{
						Name:              "child",
						Version:           "0.0.0",
						RegistryHost:      "registry",
						RegistryNamespace: "namespace",
						Builder: &builder.Builder{
							Name:   "builder",
							Driver: "mock",
						},
						Tags: []string{"0"},
					}, "child_image", childSyncChan)
				stepParent := plan.NewStep(
					&image.Image{
						Name:              "parent", // ERROR: this is the parent image
						Version:           "0.0.0",
						RegistryHost:      "registry",
						RegistryNamespace: "namespace",
						Builder: &builder.Builder{
							Name:   "builder",
							Driver: "mock",
						},
						Tags: []string{"0"},
					}, "parent_image", nil)
				stepParent.Subscribe(childSyncChan)

				buildPlan.(*plan.MockPlan).On("Plan", "parent", []string{"0.0.0"}).Return([]*plan.Step{
					stepParent,
					stepChild,
				}, nil)

				service.credentials.(*credentials.CredentialsStoreMock).On("Get", "registry").Return(&credentials.UserPasswordAuth{
					Username: "user",
					Password: "pass",
				}, nil)
				service.commandFactory.(*command.MockBuildCommandFactory).On("New",
					mockdriver.NewMockDriver(),
					stepParent.Image(),
					&driver.BuildDriverOptions{
						AnsibleConnectionLocal:           true,
						AnsibleIntermediateContainerName: "intermediate_container",
						AnsibleInventoryPath:             "inventory",
						AnsibleLimit:                     "limit",
						PullParentImage:                  true,
						PushAuthUsername:                 "user",
						PushAuthPassword:                 "pass",
						PushImageAfterBuild:              true,
						RemoveImageAfterBuild:            true,
						BuilderVarMappings:               varsmap.New(),
						BuilderOptions:                   &builder.BuilderOptions{},
					},
				).Return(command.NewMockBuildCommand(), nil)
				service.commandFactory.(*command.MockBuildCommandFactory).On("New",
					mockdriver.NewMockDriver(),
					stepChild.Image(),
					&driver.BuildDriverOptions{
						AnsibleConnectionLocal:           true,
						AnsibleIntermediateContainerName: "intermediate_container",
						AnsibleInventoryPath:             "inventory",
						AnsibleLimit:                     "limit",
						PullParentImage:                  true,
						PushAuthUsername:                 "user",
						PushAuthPassword:                 "pass",
						PushImageAfterBuild:              true,
						RemoveImageAfterBuild:            true,
						BuilderVarMappings:               varsmap.New(),
						BuilderOptions:                   &builder.BuilderOptions{},
					},
				).Return(command.NewMockBuildCommand(), nil)

				service.jobFactory.(*job.MockJobFactory).On("New", command.NewMockBuildCommand()).Return(mockJob, nil)
				service.dispatch.(*dispatch.MockDispatch).On("Enqueue", mockJob)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service, test.buildPlan)
			}

			err := test.service.Build(context.TODO(), test.buildPlan, test.name, test.versions, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(test.service)
			}
		})
	}
}

func TestBuildWorker(t *testing.T) {

	errContext := "(build::worker)"

	tests := []struct {
		desc              string
		service           *Service
		image             *image.Image
		options           *ServiceOptions
		err               error
		prepareAssertFunc func(*Service, *image.Image)
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
			desc:    "Testing error when no image specification is given to worker",
			service: &Service{},
			options: &ServiceOptions{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a dispatcher"),
		},
		{
			desc: "Testing error when no driver factory is given to worker",
			service: &Service{
				dispatch: dispatch.NewDispatch(worker.NewMockWorkerFactory()),
			},
			options: &ServiceOptions{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a driver factory"),
		},
		{
			desc: "Testing error when no semantic version generator is given to worker",
			service: &Service{
				dispatch:      dispatch.NewDispatch(worker.NewMockWorkerFactory()),
				driverFactory: driver.NewBuildDriverFactory(),
			},
			options: &ServiceOptions{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a semver generator"),
		},
		{
			desc: "Testing error when no credentials store is given to worker",
			service: &Service{
				dispatch:      dispatch.NewDispatch(worker.NewMockWorkerFactory()),
				driverFactory: driver.NewBuildDriverFactory(),
				semver:        semver.NewSemVerGenerator(),
			},
			options: &ServiceOptions{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a credentials store"),
		},
		{
			desc: "Testing worker to build an image",
			service: NewService(
				WithBuilders(store.NewMockBuildersStore()),
				WithCommandFactory(command.NewMockBuildCommandFactory()),
				WithDriverFactory(
					&driver.BuildDriverFactory{
						"mock": mockdriver.NewMockDriver(),
					},
				),
				WithJobFactory(job.NewMockJobFactory()),
				WithDispatch(dispatch.NewMockDispatch()),
				WithSemver(semver.NewSemVerGenerator()),
				WithCredentials(credentials.NewCredentialsStoreMock()),
			),
			options: &ServiceOptions{
				EnableSemanticVersionTags:    true,
				PushImageAfterBuild:          true,
				PullParentImage:              true,
				SemanticVersionTagsTemplates: []string{"{{.Major}}"},
				RemoveImagesAfterPush:        true,
				PersistentVars:               map[string]interface{}{"optpvar": "value"},
				Vars:                         map[string]interface{}{"optvar": "value"},
				Labels:                       map[string]string{"optlabel": "value"},
				Tags:                         []string{"opttag"},
			},
			image: &image.Image{
				Name:              "image",
				Version:           "0.0.0",
				RegistryHost:      "registry",
				RegistryNamespace: "namespace",
				Builder: &builder.Builder{
					Name:   "builder",
					Driver: "mock",
				},
				PersistentVars: map[string]interface{}{"imagepvar": "value"},
				Vars:           map[string]interface{}{"imagevar": "value"},
				Labels:         map[string]string{"imagelabel": "value"},
				Tags:           []string{"imagetag"},
				Parent: &image.Image{
					Name:              "parent",
					Version:           "parent_version",
					RegistryHost:      "parent_registry",
					RegistryNamespace: "parent_namespace",
					Builder:           "builder",
					PersistentVars:    map[string]interface{}{"parentpvar": "value"},
					Vars:              map[string]interface{}{"parentvar": "value"},
					Labels:            map[string]string{"parentlabel": "value"},
				},
			},
			err: &errors.Error{},
			assertFunc: func(service *Service) bool {
				return service.credentials.(*credentials.CredentialsStoreMock).AssertExpectations(t) &&
					service.commandFactory.(*command.MockBuildCommandFactory).AssertExpectations(t) &&
					service.dispatch.(*dispatch.MockDispatch).AssertExpectations(t) &&
					service.jobFactory.(*job.MockJobFactory).AssertExpectations(t)
			},
			prepareAssertFunc: func(service *Service, image *image.Image) {

				mockJob := job.NewMockJob()
				mockJob.On("Wait").Return(nil)

				service.credentials.(*credentials.CredentialsStoreMock).On("Get", "registry").Return(&credentials.UserPasswordAuth{
					Username: "user",
					Password: "pass",
				}, nil)

				service.credentials.(*credentials.CredentialsStoreMock).On("Get", "parent_registry").Return(&credentials.UserPasswordAuth{
					Username: "parent_user",
					Password: "parent_pass",
				}, nil)
				service.commandFactory.(*command.MockBuildCommandFactory).On("New",
					mockdriver.NewMockDriver(),
					image,
					&driver.BuildDriverOptions{
						AnsibleConnectionLocal:           false,
						AnsibleIntermediateContainerName: "builder_mock_namespace_image_0.0.0",
						OutputPrefix:                     "",
						PullAuthUsername:                 "parent_user",
						PullAuthPassword:                 "parent_pass",
						PullParentImage:                  true,
						PushAuthUsername:                 "user",
						PushAuthPassword:                 "pass",
						PushImageAfterBuild:              true,
						RemoveImageAfterBuild:            true,
						BuilderVarMappings:               varsmap.New(),
						BuilderOptions:                   &builder.BuilderOptions{},
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
				test.prepareAssertFunc(test.service, test.image)
			}

			err := test.service.build(context.TODO(), test.image, test.options)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.True(t, test.assertFunc(test.service))
			}
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
		image             *image.Image
		options           *driver.BuildDriverOptions
		res               job.Commander
		prepareAssertFunc func(*Service, *image.Image)
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
			desc: "Testing error when no image is provided",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
			},
			driver: mockdriver.NewMockDriver(),
			err:    errors.New(errContext, "To create a build command, is required a image"),
		},
		{
			desc: "Testing error when no options are provided",
			service: &Service{
				commandFactory: command.NewMockBuildCommandFactory(),
			},
			driver: mockdriver.NewMockDriver(),
			image:  &image.Image{},
			err:    errors.New(errContext, "To create a build command, is required a service options"),
		},
		{
			desc: "Testing create build command",
			service: NewService(
				WithCommandFactory(command.NewMockBuildCommandFactory()),
			),
			driver:  mockdriver.NewMockDriver(),
			image:   &image.Image{},
			options: &driver.BuildDriverOptions{},
			prepareAssertFunc: func(s *Service, image *image.Image) {
				s.commandFactory.(*command.MockBuildCommandFactory).On("New", mockdriver.NewMockDriver(), image, &driver.BuildDriverOptions{}).Return(command.NewMockBuildCommand(), nil)
			},
			res: &command.MockBuildCommand{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service, test.image)
			}

			res, err := test.service.command(test.driver, test.image, test.options)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
			}

		})
	}
}

func TestGetBuilder(t *testing.T) {
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
			service:           &Service{builders: store.NewMockBuildersStore()},
			builderDefinition: "test",
			prepareAssertFunc: func(s *Service) {
				s.builders.(*store.MockBuildersStore).On("Find", "test").Return(&builder.Builder{
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

			res, err := test.service.getBuilder(test.builderDefinition)
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

func TestGetCredentials(t *testing.T) {

	errContext := "(build::getCredentials)"

	tests := []struct {
		desc              string
		service           *Service
		registry          string
		res               *credentials.UserPasswordAuth
		err               error
		prepareAssertFunc func(*Service)
	}{
		{
			desc:    "Testing error when credentials store is nil",
			service: NewService(),
			err:     errors.New(errContext, "To get credentials, is required a credentials store"),
		},
		{
			desc: "Testing get credentials",
			service: NewService(
				WithCredentials(
					credentials.NewCredentialsStoreMock(),
				),
			),
			registry: "registry.test",
			res: &credentials.UserPasswordAuth{
				Username: "user",
				Password: "pass",
			},
			prepareAssertFunc: func(service *Service) {
				service.credentials.(*credentials.CredentialsStoreMock).On("Get", "registry.test").Return(&credentials.UserPasswordAuth{
					Username: "user",
					Password: "pass",
				}, nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing get unexisting credentials",
			service: NewService(
				WithCredentials(
					credentials.NewCredentialsStoreMock(),
				),
			),
			registry: "registry.test",
			res:      nil,
			prepareAssertFunc: func(service *Service) {
				service.credentials.(*credentials.CredentialsStoreMock).On("Get", "registry.test").Return(nil, errors.New(errContext, "Credentials not found"))
			},
			err: errors.New(errContext, "Credentials not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service)
			}

			cred, err := test.service.getCredentials(test.registry)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, cred)
			}
		})
	}
}

func TestGetDriver(t *testing.T) {
	errContext := "(build::getDriver)"

	tests := []struct {
		desc    string
		service *Service
		builder *builder.Builder
		options *ServiceOptions
		res     driver.BuildDriverer
		err     error
	}{
		{
			desc:    "Testing error when driver factory is not defined",
			service: NewService(),
			builder: &builder.Builder{
				Driver: "docker",
			},
			options: &ServiceOptions{},
			res:     &dockerdriver.DockerDriver{},
			err:     errors.New(errContext, "To create a build driver, is required a driver factory"),
		},
		{
			desc: "Testing get driver",
			service: NewService(
				WithDriverFactory(
					&driver.BuildDriverFactory{
						"docker": mockdriver.NewMockDriver(),
					},
				),
			),
			builder: &builder.Builder{
				Driver: "docker",
			},
			options: &ServiceOptions{},
			res:     &mockdriver.MockDriver{},
			err:     &errors.Error{},
		},
		{
			desc: "Testing get driver with dry-run",
			service: NewService(
				WithDriverFactory(
					&driver.BuildDriverFactory{
						"docker":  mockdriver.NewMockDriver(),
						"dry-run": dryrundriver.NewDryRunDriver(ioutil.Discard),
					},
				),
			),
			builder: &builder.Builder{
				Driver: "docker",
			},
			options: &ServiceOptions{
				DryRun: true,
			},
			res: &dryrundriver.DryRunDriver{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			driver, err := test.service.getDriver(test.builder, test.options)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}
}