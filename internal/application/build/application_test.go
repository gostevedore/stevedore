package build

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/domain/varsmap"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	credentialsfactory "github.com/gostevedore/stevedore/internal/infrastructure/credentials/factory"
	authmethodbasic "github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/basic"
	authmethodkeyfile "github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/keyfile"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/mock"
	"github.com/gostevedore/stevedore/internal/infrastructure/plan"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/dispatch"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/job"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/worker"
	"github.com/gostevedore/stevedore/internal/infrastructure/semver"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
	"github.com/stretchr/testify/assert"
	testmock "github.com/stretchr/testify/mock"
)

func TestBuild(t *testing.T) {
	errContext := "(application::build::Build)"
	_ = errContext
	tests := []struct {
		desc              string
		service           *Application
		buildPlan         Planner
		name              string
		versions          []string
		options           *Options
		prepareAssertFunc func(*Application, Planner)
		assertFunc        func(*Application) bool
		err               error
	}{
		{
			desc:    "Testing error building an image with no options",
			service: &Application{},
			options: nil,
			err:     errors.New(errContext, "To build an image, service options are required"),
		},
		{
			desc:    "Testing error building an image with no execution plan",
			service: &Application{},
			options: &Options{},
			err:     errors.New(errContext, "To build an image, a build plan is required"),
		},
		{
			desc: "Testing build an image",
			service: NewApplication(
				WithBuilders(builders.NewMockStore()),
				WithCommandFactory(command.NewMockBuildCommandFactory()),
				WithDriverFactory(
					&factory.BuildDriverFactory{
						"mock": func() (repository.BuildDriverer, error) {
							return mock.NewMockDriver(), nil
						},
					},
				),
				WithJobFactory(job.NewMockJobFactory()),
				WithDispatch(dispatch.NewMockDispatch()),
				WithSemver(semver.NewSemVerGenerator()),
				WithCredentials(credentialsfactory.NewMockCredentialsFactory()),
			),
			buildPlan: plan.NewMockPlan(),
			name:      "parent",
			versions:  []string{"0.0.0"},
			options: &Options{
				AnsibleConnectionLocal:           true,
				AnsibleIntermediateContainerName: "intermediate_container",
				AnsibleInventoryPath:             "inventory",
				AnsibleLimit:                     "limit",
				EnableSemanticVersionTags:        true,
				ImageFromName:                    image.UndefinedStringValue,
				ImageFromRegistryHost:            image.UndefinedStringValue,
				ImageFromRegistryNamespace:       image.UndefinedStringValue,
				ImageFromVersion:                 image.UndefinedStringValue,
				ImageName:                        image.UndefinedStringValue,
				ImageRegistryHost:                image.UndefinedStringValue,
				ImageRegistryNamespace:           image.UndefinedStringValue,
				PullParentImage:                  true,
				PushImageAfterBuild:              true,
				RemoveImagesAfterPush:            true,
				SemanticVersionTagsTemplates:     []string{"{{.Major}}"},
			},
			err: &errors.Error{},
			assertFunc: func(service *Application) bool {
				return service.credentials.(*credentialsfactory.MockCredentialsFactory).AssertExpectations(t) &&
					service.commandFactory.(*command.MockBuildCommandFactory).AssertExpectations(t) &&
					service.dispatch.(*dispatch.MockDispatch).AssertExpectations(t) &&
					service.jobFactory.(*job.MockJobFactory).AssertExpectations(t)
			},
			prepareAssertFunc: func(service *Application, buildPlan Planner) {

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

				service.credentials.(*credentialsfactory.MockCredentialsFactory).On("Get", "registry").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)
				service.commandFactory.(*command.MockBuildCommandFactory).On("New",
					mock.NewMockDriver(),
					stepParent.Image(),
					&image.BuildDriverOptions{
						AnsibleConnectionLocal:           true,
						AnsibleIntermediateContainerName: "intermediate_container",
						AnsibleInventoryPath:             "inventory",
						AnsibleLimit:                     "limit",
						BuilderOptions:                   &builder.BuilderOptions{},
						BuilderVarMappings:               varsmap.New(),
						PullParentImage:                  true,
						PushAuthPassword:                 "password",
						PushAuthUsername:                 "username",
						PushImageAfterBuild:              true,
						RemoveImageAfterBuild:            true,
					},
				).Return(command.NewMockBuildCommand(), nil)
				service.commandFactory.(*command.MockBuildCommandFactory).On("New",
					mock.NewMockDriver(),
					stepChild.Image(),
					&image.BuildDriverOptions{
						AnsibleConnectionLocal:           true,
						AnsibleIntermediateContainerName: "intermediate_container",
						AnsibleInventoryPath:             "inventory",
						AnsibleLimit:                     "limit",
						BuilderOptions:                   &builder.BuilderOptions{},
						BuilderVarMappings:               varsmap.New(),
						PullParentImage:                  true,
						PushAuthPassword:                 "password",
						PushAuthUsername:                 "username",
						PushImageAfterBuild:              true,
						RemoveImageAfterBuild:            true,
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

	errContext := "(application::build::worker)"

	tests := []struct {
		desc              string
		service           *Application
		image             *image.Image
		options           *Options
		err               error
		prepareAssertFunc func(*Application, *image.Image)
		assertFunc        func(*Application) bool
	}{
		{
			desc:    "Testing error when no options are given to worker",
			service: &Application{},
			options: nil,
			err:     errors.New(errContext, "Build worker requires service options"),
		},
		{
			desc:    "Testing error when no image specification is given to worker",
			service: &Application{},
			options: &Options{},
			err:     errors.New(errContext, "Build worker requires an image specification"),
		},
		{
			desc:    "Testing error when no image specification is given to worker",
			service: &Application{},
			options: &Options{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a dispatcher"),
		},
		{
			desc: "Testing error when no driver factory is given to worker",
			service: &Application{
				dispatch: dispatch.NewDispatch(worker.NewMockWorkerFactory()),
			},
			options: &Options{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a driver factory"),
		},
		{
			desc: "Testing error when no semantic version generator is given to worker",
			service: &Application{
				dispatch:      dispatch.NewDispatch(worker.NewMockWorkerFactory()),
				driverFactory: factory.NewBuildDriverFactory(),
			},
			options: &Options{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a semver generator"),
		},
		{
			desc: "Testing error when no credentials store is given to worker",
			service: &Application{
				dispatch:      dispatch.NewDispatch(worker.NewMockWorkerFactory()),
				driverFactory: factory.NewBuildDriverFactory(),
				semver:        semver.NewSemVerGenerator(),
			},
			options: &Options{},
			image:   &image.Image{},
			err:     errors.New(errContext, "Build worker requires a credentials store"),
		},
		{
			desc: "Testing worker to build an image",
			service: NewApplication(
				WithBuilders(builders.NewMockStore()),
				WithCommandFactory(command.NewMockBuildCommandFactory()),
				WithDriverFactory(
					&factory.BuildDriverFactory{
						"mock": func() (repository.BuildDriverer, error) {
							return mock.NewMockDriver(), nil
						},
					},
				),
				WithJobFactory(job.NewMockJobFactory()),
				WithDispatch(dispatch.NewMockDispatch()),
				WithSemver(semver.NewSemVerGenerator()),
				WithCredentials(credentialsfactory.NewMockCredentialsFactory()),
			),
			options: &Options{
				EnableSemanticVersionTags:    true,
				ImageFromName:                image.UndefinedStringValue,
				ImageFromRegistryHost:        image.UndefinedStringValue,
				ImageFromRegistryNamespace:   image.UndefinedStringValue,
				ImageFromVersion:             image.UndefinedStringValue,
				ImageName:                    image.UndefinedStringValue,
				ImageRegistryHost:            image.UndefinedStringValue,
				ImageRegistryNamespace:       image.UndefinedStringValue,
				Labels:                       map[string]string{"optlabel": "value", "imagelabel_overwritten": "overwritten_value"},
				PersistentLabels:             map[string]string{"optplabel": "value"},
				PersistentVars:               map[string]interface{}{"optpvar": "value", "imagepvar_overwritten": "overwritten_value"},
				PullParentImage:              true,
				PushImageAfterBuild:          true,
				RemoveImagesAfterPush:        true,
				SemanticVersionTagsTemplates: []string{"{{.Major}}"},
				Tags:                         []string{"opttag"},
				Vars:                         map[string]interface{}{"optvar": "value", "imagevar_overwritten": "overwritten_value"},
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
				PersistentVars: map[string]interface{}{"imagepvar": "value", "imagepvar_overwritten": "value"},
				Vars:           map[string]interface{}{"imagevar": "value", "imagevar_overwritten": "value"},
				Labels:         map[string]string{"imagelabel": "value", "imagelabel_overwritten": "value"},
				Tags:           []string{"imagetag"},
				Parent: &image.Image{
					Name:              "parent",
					Version:           "parent_version",
					RegistryHost:      "parent_registry",
					RegistryNamespace: "parent_namespace",
					Builder:           "builder",
					PersistentVars:    map[string]interface{}{},
					Vars:              map[string]interface{}{},
					Labels:            map[string]string{},
					PersistentLabels:  map[string]string{},
				},
			},
			err: &errors.Error{},
			assertFunc: func(service *Application) bool {
				return service.credentials.(*credentialsfactory.MockCredentialsFactory).AssertExpectations(t) &&
					service.commandFactory.(*command.MockBuildCommandFactory).AssertExpectations(t) &&
					service.dispatch.(*dispatch.MockDispatch).AssertExpectations(t) &&
					service.jobFactory.(*job.MockJobFactory).AssertExpectations(t)
			},
			prepareAssertFunc: func(service *Application, i *image.Image) {

				mockJob := job.NewMockJob()
				mockJob.On("Wait").Return(nil)

				service.credentials.(*credentialsfactory.MockCredentialsFactory).On("Get", "registry").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)

				service.credentials.(*credentialsfactory.MockCredentialsFactory).On("Get", "parent_registry").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username_parent",
					Password: "password_parent",
				}, nil)
				service.commandFactory.(*command.MockBuildCommandFactory).On("New",
					//mock.NewMockDriver(),
					testmock.Anything,
					// i,
					&image.Image{
						Name:              "image",
						Version:           "0.0.0",
						RegistryHost:      "registry",
						RegistryNamespace: "namespace",
						Builder: &builder.Builder{
							Name:   "builder",
							Driver: "mock",
						},
						PersistentVars:   map[string]interface{}{"imagepvar": "value", "imagepvar_overwritten": "overwritten_value", "optpvar": "value"},
						PersistentLabels: map[string]string{"optplabel": "value"},
						Vars:             map[string]interface{}{"imagevar": "value", "imagevar_overwritten": "overwritten_value", "optvar": "value"},
						Labels:           map[string]string{"imagelabel": "value", "imagelabel_overwritten": "overwritten_value", "optlabel": "value"},
						Tags:             []string{"imagetag", "0", "opttag"},
						Parent: &image.Image{
							Name:              "parent",
							Version:           "parent_version",
							RegistryHost:      "parent_registry",
							RegistryNamespace: "parent_namespace",
							Builder:           "builder",
							PersistentVars:    map[string]interface{}{},
							Vars:              map[string]interface{}{},
							Labels:            map[string]string{},
							PersistentLabels:  map[string]string{},
						},
					},
					&image.BuildDriverOptions{
						AnsibleConnectionLocal:           false,
						AnsibleIntermediateContainerName: "builder_mock_namespace_image_0.0.0",
						BuilderOptions:                   &builder.BuilderOptions{},
						BuilderVarMappings:               varsmap.New(),
						OutputPrefix:                     "",
						PullAuthPassword:                 "password_parent",
						PullAuthUsername:                 "username_parent",
						PullParentImage:                  true,
						PushAuthPassword:                 "password",
						PushAuthUsername:                 "username",
						PushImageAfterBuild:              true,
						RemoveImageAfterBuild:            true,
					}).Return(command.NewMockBuildCommand(), nil)
				service.jobFactory.(*job.MockJobFactory).On("New", command.NewMockBuildCommand()).Return(mockJob, nil)
				service.dispatch.(*dispatch.MockDispatch).On("Enqueue", mockJob)
			},
		},
		{
			desc: "Testing error build when image credentials are invalid",
			service: NewApplication(
				WithBuilders(builders.NewMockStore()),
				WithCommandFactory(command.NewMockBuildCommandFactory()),
				WithDriverFactory(
					&factory.BuildDriverFactory{
						"mock": func() (repository.BuildDriverer, error) {
							return mock.NewMockDriver(), nil
						},
					},
				),
				WithJobFactory(job.NewMockJobFactory()),
				WithDispatch(dispatch.NewMockDispatch()),
				WithSemver(semver.NewSemVerGenerator()),
				WithCredentials(credentialsfactory.NewMockCredentialsFactory()),
			),
			options: &Options{
				EnableSemanticVersionTags:    true,
				ImageFromName:                image.UndefinedStringValue,
				ImageFromRegistryHost:        image.UndefinedStringValue,
				ImageFromRegistryNamespace:   image.UndefinedStringValue,
				ImageFromVersion:             image.UndefinedStringValue,
				ImageName:                    image.UndefinedStringValue,
				ImageRegistryHost:            image.UndefinedStringValue,
				ImageRegistryNamespace:       image.UndefinedStringValue,
				Labels:                       map[string]string{"optlabel": "value"},
				PersistentVars:               map[string]interface{}{"optpvar": "value"},
				PullParentImage:              true,
				PushImageAfterBuild:          true,
				RemoveImagesAfterPush:        true,
				SemanticVersionTagsTemplates: []string{"{{.Major}}"},
				Tags:                         []string{"opttag"},
				Vars:                         map[string]interface{}{"optvar": "value"},
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
			err: errors.New(errContext, "Invalid credentials method for 'registry'. Found 'keyfile' when is expected basic auth method"),
			prepareAssertFunc: func(service *Application, i *image.Image) {

				mockJob := job.NewMockJob()
				mockJob.On("Wait").Return(nil)

				service.credentials.(*credentialsfactory.MockCredentialsFactory).On("Get", "parent_registry").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)

				service.credentials.(*credentialsfactory.MockCredentialsFactory).On("Get", "registry").Return(&authmethodkeyfile.KeyFileAuthMethod{}, nil)
			},
		},
		{
			desc: "Testing error build when parent credentials are invalid",
			service: NewApplication(
				WithBuilders(builders.NewMockStore()),
				WithCommandFactory(command.NewMockBuildCommandFactory()),
				WithDriverFactory(
					&factory.BuildDriverFactory{
						"mock": func() (repository.BuildDriverer, error) {
							return mock.NewMockDriver(), nil
						},
					},
				),
				WithJobFactory(job.NewMockJobFactory()),
				WithDispatch(dispatch.NewMockDispatch()),
				WithSemver(semver.NewSemVerGenerator()),
				WithCredentials(credentialsfactory.NewMockCredentialsFactory()),
			),
			options: &Options{
				ImageFromName:                image.UndefinedStringValue,
				ImageFromRegistryHost:        image.UndefinedStringValue,
				ImageFromRegistryNamespace:   image.UndefinedStringValue,
				ImageFromVersion:             image.UndefinedStringValue,
				ImageName:                    image.UndefinedStringValue,
				ImageRegistryHost:            image.UndefinedStringValue,
				ImageRegistryNamespace:       image.UndefinedStringValue,
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
			err: errors.New(errContext, "Invalid credentials method for 'parent_registry'. Found 'keyfile' when is expected basic auth method"),
			prepareAssertFunc: func(service *Application, i *image.Image) {

				mockJob := job.NewMockJob()
				mockJob.On("Wait").Return(nil)

				service.credentials.(*credentialsfactory.MockCredentialsFactory).On("Get", "registry").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)

				service.credentials.(*credentialsfactory.MockCredentialsFactory).On("Get", "parent_registry").Return(&authmethodkeyfile.KeyFileAuthMethod{}, nil)
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
	errContext := "(application::build::command)"

	tests := []struct {
		desc              string
		service           *Application
		cmd               job.Commander
		prepareAssertFunc func(*Application, job.Commander)
		err               error
	}{
		{
			desc:    "Testing error when no job factory is defined on service",
			service: &Application{},
			err:     errors.New(errContext, "To create a build job, is required a job factory"),
		},
		{
			desc: "Testing job creation",
			service: &Application{
				jobFactory: job.NewMockJobFactory(),
			},
			cmd: command.NewMockBuildCommand(),
			prepareAssertFunc: func(service *Application, cmd job.Commander) {
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
	errContext := "(application::build::command)"

	driverFactory := factory.NewBuildDriverFactory()
	driverFactory.Register("mock",
		func() (repository.BuildDriverer, error) {
			return mock.NewMockDriver(), nil
		},
	)

	tests := []struct {
		desc              string
		service           *Application
		driver            repository.BuildDriverer
		image             *image.Image
		options           *image.BuildDriverOptions
		res               job.Commander
		prepareAssertFunc func(*Application, *image.Image)
		err               error
	}{
		{
			desc:    "Testing error when no command factory is provided",
			service: &Application{},
			err:     errors.New(errContext, "To create a build command, is required a command factory"),
		},
		{
			desc: "Testing error when no driver is provided",
			service: &Application{
				commandFactory: command.NewMockBuildCommandFactory(),
			},
			err: errors.New(errContext, "To create a build command, is required a driver"),
		},
		{
			desc: "Testing error when no image is provided",
			service: &Application{
				commandFactory: command.NewMockBuildCommandFactory(),
			},
			driver: mock.NewMockDriver(),
			err:    errors.New(errContext, "To create a build command, is required a image"),
		},
		{
			desc: "Testing error when no options are provided",
			service: &Application{
				commandFactory: command.NewMockBuildCommandFactory(),
			},
			driver: mock.NewMockDriver(),
			image:  &image.Image{},
			err:    errors.New(errContext, "To create a build command, is required a service options"),
		},
		{
			desc: "Testing create build command",
			service: NewApplication(
				WithCommandFactory(command.NewMockBuildCommandFactory()),
			),
			driver:  mock.NewMockDriver(),
			image:   &image.Image{},
			options: &image.BuildDriverOptions{},
			prepareAssertFunc: func(s *Application, i *image.Image) {
				s.commandFactory.(*command.MockBuildCommandFactory).On("New", mock.NewMockDriver(), i, &image.BuildDriverOptions{}).Return(command.NewMockBuildCommand(), nil)
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
	errContext := "(application::build::builder)"

	tests := []struct {
		desc              string
		service           *Application
		image             *image.Image
		res               *builder.Builder
		prepareAssertFunc func(*Application)
		assertFunc        func(expected, actual *builder.Builder) bool
		err               error
	}{
		{
			desc:    "Testing error getting a builder to nil image",
			service: &Application{},
			image:   nil,
			err:     errors.New(errContext, "To generate a builder, is required an image definition"),
		},
		{
			desc:    "Testing error getting a builder with no builders store",
			service: &Application{},
			image: &image.Image{
				Builder: "test",
			},
			err: errors.New(errContext, "To generate a builder, is required a builder store defined on build service"),
		},
		{
			desc:    "Testing return a builder defined by an string",
			service: &Application{builders: builders.NewMockStore()},
			image: &image.Image{
				Builder: "test",
			},
			prepareAssertFunc: func(s *Application) {
				s.builders.(*builders.MockStore).On("Find", "test").Return(&builder.Builder{
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
		{
			desc:    "Testing return a builder defined by an interface{}",
			service: &Application{builders: builders.NewMockStore()},
			image: &image.Image{
				Builder: map[interface{}]interface{}{
					"driver": "docker",
					"options": map[interface{}]interface{}{
						"dockerfile": "Dockerfile.test",
						"context": []map[interface{}]interface{}{
							{
								"git": map[interface{}]interface{}{
									"path":       "path",
									"repository": "repository",
									"reference":  "reference",
									"auth": map[interface{}]interface{}{
										"username": "username",
										"password": "password",
									},
								},
							},
						},
					},
				},
			},

			prepareAssertFunc: nil,
			assertFunc: func(expected, actual *builder.Builder) bool {
				return assert.Equal(t, expected, actual)
			},
			res: &builder.Builder{
				Name:   "",
				Driver: "docker",
				Options: &builder.BuilderOptions{
					Context: []*builder.DockerDriverContextOptions{
						{
							Git: &builder.DockerDriverGitContextOptions{
								Path:       "path",
								Repository: "repository",
								Reference:  "reference",
								Auth: &builder.DockerDriverGitContextAuthOptions{
									Username: "username",
									Password: "password",
								},
							},
						},
					},
					Dockerfile: "Dockerfile.test",
				},
				VarMapping: varsmap.New(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service)
			}

			res, err := test.service.getBuilder(test.image)
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

	errContext := "(application::build::getCredentials)"

	tests := []struct {
		desc              string
		service           *Application
		registry          string
		res               repository.AuthMethodReader
		err               error
		prepareAssertFunc func(*Application)
	}{
		{
			desc:    "Testing error when credentials store is nil",
			service: NewApplication(),
			err:     errors.New(errContext, "To get credentials, is required a credentials store"),
		},
		{
			desc: "Testing get credentials",
			service: NewApplication(
				WithCredentials(
					credentialsfactory.NewMockCredentialsFactory(),
				),
			),
			registry: "registry.test",
			res: &authmethodbasic.BasicAuthMethod{
				Username: "username",
				Password: "password",
			},
			prepareAssertFunc: func(service *Application) {
				service.credentials.(*credentialsfactory.MockCredentialsFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing get unexisting credentials",
			service: NewApplication(
				WithCredentials(
					credentialsfactory.NewMockCredentialsFactory(),
				),
			),
			registry: "registry.test",
			res:      nil,
			prepareAssertFunc: func(service *Application) {
				service.credentials.(*credentialsfactory.MockCredentialsFactory).On("Get", "registry.test").Return(nil, errors.New(errContext, "Credentials not found"))
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
	errContext := "(application::build::getDriver)"

	tests := []struct {
		desc    string
		service *Application
		builder *builder.Builder
		options *Options
		res     repository.BuildDriverer
		err     error
	}{
		{
			desc:    "Testing error when driver factory is not defined",
			service: NewApplication(),
			builder: &builder.Builder{
				Driver: "docker",
			},
			options: &Options{},
			res:     &docker.DockerDriver{},
			err:     errors.New(errContext, "To create a build driver, is required a driver factory"),
		},
		{
			desc: "Testing get driver",
			service: NewApplication(
				WithDriverFactory(
					&factory.BuildDriverFactory{
						"docker": func() (repository.BuildDriverer, error) {
							return mock.NewMockDriver(), nil
						},
					},
				),
			),
			builder: &builder.Builder{
				Driver: "docker",
			},
			options: &Options{},
			res:     &mock.MockDriver{},
			err:     &errors.Error{},
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
