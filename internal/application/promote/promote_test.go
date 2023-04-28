package promote

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	authfactory "github.com/gostevedore/stevedore/internal/infrastructure/auth/factory"
	authmethodbasic "github.com/gostevedore/stevedore/internal/infrastructure/auth/method/basic"
	authmethodkeyfile "github.com/gostevedore/stevedore/internal/infrastructure/auth/method/keyfile"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/dryrun"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/mock"
	reference "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/default"
	"github.com/gostevedore/stevedore/internal/infrastructure/semver"
	"github.com/stretchr/testify/assert"
)

func TestPromote(t *testing.T) {
	errContext := "(application::promote::Promote)"

	tests := []struct {
		desc            string
		service         *Application
		options         *Options
		context         context.Context
		prepareMockFunc func(*Application)
		err             error
	}{
		{
			desc: "Testing promote application source image from local",
			service: NewApplication(
				WithCredentials(authfactory.NewMockAuthFactory()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(factory.NewPromoteFactory()),
				WithReferenceNamer(reference.NewDefaultReferenceName()),
			),
			context: context.TODO(),
			options: &Options{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    false,
				TargetImageName:              image.UndefinedStringValue,
				TargetImageRegistryNamespace: image.UndefinedStringValue,
				TargetImageRegistryHost:      image.UndefinedStringValue,
				TargetImageTags:              nil,
				PromoteSourceImageTag:        false,
				RemoveTargetImageTags:        false,
				RemoteSourceImage:            false,
				SemanticVersionTagsTemplates: nil,
			},
			prepareMockFunc: func(p *Application) {

				options := &image.PromoteOptions{
					TargetImageName:       "registry.test/namespace/image:tag",
					TargetImageTags:       nil,
					RemoveTargetImageTags: false,
					RemoteSourceImage:     false,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "username",
					PullAuthPassword:      "password",
					PushAuthUsername:      "username",
					PushAuthPassword:      "password",
				}

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)
				p.factory.Register(image.DockerPromoterName, mock)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote application source image from remote",
			service: NewApplication(
				WithCredentials(authfactory.NewMockAuthFactory()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(factory.NewPromoteFactory()),
				WithReferenceNamer(reference.NewDefaultReferenceName()),
			),
			context: context.TODO(),
			options: &Options{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    false,
				TargetImageName:              image.UndefinedStringValue,
				TargetImageRegistryNamespace: image.UndefinedStringValue,
				TargetImageRegistryHost:      image.UndefinedStringValue,
				TargetImageTags:              nil,
				PromoteSourceImageTag:        false,
				RemoveTargetImageTags:        false,
				RemoteSourceImage:            true,
				SemanticVersionTagsTemplates: nil,
			},
			prepareMockFunc: func(p *Application) {
				options := &image.PromoteOptions{
					TargetImageName:       "registry.test/namespace/image:tag",
					TargetImageTags:       nil,
					RemoveTargetImageTags: false,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "username",
					PullAuthPassword:      "password",
					PushAuthUsername:      "username",
					PushAuthPassword:      "password",
				}

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := factory.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote application source image using empty target values",
			service: NewApplication(
				WithCredentials(authfactory.NewMockAuthFactory()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(factory.NewPromoteFactory()),
				WithReferenceNamer(reference.NewDefaultReferenceName()),
			),
			context: context.TODO(),
			options: &Options{
				EnableSemanticVersionTags:    false,
				PromoteSourceImageTag:        false,
				RemoteSourceImage:            true,
				RemoveTargetImageTags:        false,
				SemanticVersionTagsTemplates: nil,
				SourceImageName:              "registry.test/namespace/image:tag",
				TargetImageName:              image.UndefinedStringValue,
				TargetImageRegistryHost:      "",
				TargetImageRegistryNamespace: "",
				TargetImageTags:              nil,
			},
			prepareMockFunc: func(p *Application) {
				options := &image.PromoteOptions{
					PullAuthPassword:      "password",
					PullAuthUsername:      "username",
					RemoteSourceImage:     true,
					RemoveTargetImageTags: false,
					SourceImageName:       "registry.test/namespace/image:tag",
					TargetImageName:       "image:tag",
					TargetImageTags:       nil,
				}

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := factory.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing error promoting source image using an empty target image name",
			service: NewApplication(
				WithCredentials(authfactory.NewMockAuthFactory()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(factory.NewPromoteFactory()),
				WithReferenceNamer(reference.NewDefaultReferenceName()),
			),
			context: context.TODO(),
			options: &Options{
				EnableSemanticVersionTags:    false,
				PromoteSourceImageTag:        false,
				RemoteSourceImage:            true,
				RemoveTargetImageTags:        false,
				SemanticVersionTagsTemplates: nil,
				SourceImageName:              "registry.test/namespace/image:tag",
				TargetImageName:              "",
				TargetImageRegistryHost:      "",
				TargetImageRegistryNamespace: "",
				TargetImageTags:              nil,
			},
			prepareMockFunc: func(p *Application) {
				options := &image.PromoteOptions{
					PullAuthPassword:      "password",
					PullAuthUsername:      "username",
					RemoteSourceImage:     true,
					RemoveTargetImageTags: false,
					SourceImageName:       "registry.test/namespace/image:tag",
					TargetImageName:       "image:tag",
					TargetImageTags:       nil,
				}

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := factory.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)
			},
			err: errors.New(errContext, "Error generating target image reference name for 'registry.test/namespace/image:tag'\n Image reference name can not be generated because image name is undefined"),
		},
		{
			desc: "Testing promote application source image with all options",
			service: NewApplication(
				WithCredentials(authfactory.NewMockAuthFactory()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(factory.NewPromoteFactory()),
				WithReferenceNamer(reference.NewDefaultReferenceName()),
			),
			context: context.TODO(),
			options: &Options{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    true,
				TargetImageName:              "targetimage",
				TargetImageRegistryNamespace: "targetnamespace",
				TargetImageRegistryHost:      "targetregistry.test",
				TargetImageTags:              []string{"1.2.3", "tag1", "tag2"},
				PromoteSourceImageTag:        true,
				RemoveTargetImageTags:        true,
				RemoteSourceImage:            true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
			},
			prepareMockFunc: func(p *Application) {
				options := &image.PromoteOptions{
					TargetImageName: "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags: []string{
						"targetregistry.test/targetnamespace/targetimage:tag",
						"targetregistry.test/targetnamespace/targetimage:tag1",
						"targetregistry.test/targetnamespace/targetimage:tag2",
						"targetregistry.test/targetnamespace/targetimage:1",
					},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "username_pull",
					PullAuthPassword:      "password_pull",
					PushAuthUsername:      "username_push",
					PushAuthPassword:      "password_push",
				}

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username_pull",
					Password: "password_pull",
				}, nil)

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "targetregistry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username_push",
					Password: "password_push",
				}, nil)

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := factory.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote application source image with no credentials",
			service: NewApplication(
				WithCredentials(authfactory.NewMockAuthFactory()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(factory.NewPromoteFactory()),
				WithReferenceNamer(reference.NewDefaultReferenceName()),
			),
			context: context.TODO(),
			options: &Options{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    true,
				TargetImageName:              "targetimage",
				TargetImageRegistryNamespace: "targetnamespace",
				TargetImageRegistryHost:      "targetregistry.test",
				TargetImageTags:              []string{"1.2.3", "tag1", "tag2"},
				PromoteSourceImageTag:        true,
				RemoveTargetImageTags:        true,
				RemoteSourceImage:            true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
			},
			prepareMockFunc: func(p *Application) {
				options := &image.PromoteOptions{
					TargetImageName: "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags: []string{
						"targetregistry.test/targetnamespace/targetimage:tag",
						"targetregistry.test/targetnamespace/targetimage:tag1",
						"targetregistry.test/targetnamespace/targetimage:tag2",
						"targetregistry.test/targetnamespace/targetimage:1",
					},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "",
					PullAuthPassword:      "",
					PushAuthUsername:      "",
					PushAuthPassword:      "",
				}

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{}, nil)
				p.credentials.(*authfactory.MockAuthFactory).On("Get", "targetregistry.test").Return(&authmethodbasic.BasicAuthMethod{}, nil)

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := factory.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote application source image with all options and using semver configuration parameters",
			service: NewApplication(
				WithCredentials(authfactory.NewMockAuthFactory()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(factory.NewPromoteFactory()),
				WithReferenceNamer(reference.NewDefaultReferenceName()),
			),
			context: context.TODO(),
			options: &Options{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    true,
				TargetImageName:              "targetimage",
				TargetImageRegistryNamespace: "targetnamespace",
				TargetImageRegistryHost:      "targetregistry.test",
				TargetImageTags:              []string{"1.2.3", "tag1", "tag2"},
				PromoteSourceImageTag:        true,
				RemoveTargetImageTags:        true,
				RemoteSourceImage:            true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
			},
			prepareMockFunc: func(p *Application) {
				options := &image.PromoteOptions{
					TargetImageName: "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags: []string{
						"targetregistry.test/targetnamespace/targetimage:tag",
						"targetregistry.test/targetnamespace/targetimage:tag1",
						"targetregistry.test/targetnamespace/targetimage:tag2",
						"targetregistry.test/targetnamespace/targetimage:1",
					},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "username_pull",
					PullAuthPassword:      "password_pull",
					PushAuthUsername:      "username_push",
					PushAuthPassword:      "password_push",
				}

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username_pull",
					Password: "password_pull",
				}, nil)

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "targetregistry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username_push",
					Password: "password_push",
				}, nil)

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := factory.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote application source image with all options, using semver configuration parameters overridden by service options",
			service: &Application{
				credentials:    authfactory.NewMockAuthFactory(),
				semver:         semver.NewSemVerGenerator(),
				referenceNamer: reference.NewDefaultReferenceName(),
			},
			context: context.TODO(),
			options: &Options{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    true,
				TargetImageName:              "targetimage",
				TargetImageRegistryNamespace: "targetnamespace",
				TargetImageRegistryHost:      "targetregistry.test",
				TargetImageTags:              []string{"1.2.3", "tag1", "tag2"},
				PromoteSourceImageTag:        true,
				RemoveTargetImageTags:        true,
				RemoteSourceImage:            true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
			},
			prepareMockFunc: func(p *Application) {
				options := &image.PromoteOptions{
					TargetImageName: "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags: []string{
						"targetregistry.test/targetnamespace/targetimage:tag",
						"targetregistry.test/targetnamespace/targetimage:tag1",
						"targetregistry.test/targetnamespace/targetimage:tag2",
						"targetregistry.test/targetnamespace/targetimage:1",
					},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "username_pull",
					PullAuthPassword:      "password_pull",
					PushAuthUsername:      "username_push",
					PushAuthPassword:      "password_push",
				}

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username_pull",
					Password: "password_pull",
				}, nil)

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "targetregistry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username_push",
					Password: "password_push",
				}, nil)

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := factory.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing error promote when push credentials are invalid",
			service: &Application{
				credentials:    authfactory.NewMockAuthFactory(),
				semver:         semver.NewSemVerGenerator(),
				referenceNamer: reference.NewDefaultReferenceName(),
			},
			context: context.TODO(),
			options: &Options{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    true,
				TargetImageName:              "targetimage",
				TargetImageRegistryNamespace: "targetnamespace",
				TargetImageRegistryHost:      "targetregistry.test",
				TargetImageTags:              []string{"1.2.3", "tag1", "tag2"},
				PromoteSourceImageTag:        true,
				RemoveTargetImageTags:        true,
				RemoteSourceImage:            true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
			},
			prepareMockFunc: func(p *Application) {
				options := &image.PromoteOptions{
					TargetImageName: "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags: []string{
						"targetregistry.test/targetnamespace/targetimage:tag",
						"targetregistry.test/targetnamespace/targetimage:tag1",
						"targetregistry.test/targetnamespace/targetimage:tag2",
						"targetregistry.test/targetnamespace/targetimage:1",
					},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "username_pull",
					PullAuthPassword:      "password_pull",
					PushAuthUsername:      "username_push",
					PushAuthPassword:      "password_push",
				}

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username_pull",
					Password: "password_pull",
				}, nil)

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "targetregistry.test").Return(&authmethodkeyfile.KeyFileAuthMethod{}, nil)

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := factory.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: errors.New(errContext, "Invalid credentials method for 'targetregistry.test'. Found 'keyfile' when is expected basic auth method"),
		},
		{
			desc: "Testing error promote when pull credentials are invalid",
			service: &Application{
				credentials:    authfactory.NewMockAuthFactory(),
				semver:         semver.NewSemVerGenerator(),
				referenceNamer: reference.NewDefaultReferenceName(),
			},
			context: context.TODO(),
			options: &Options{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    true,
				TargetImageName:              "targetimage",
				TargetImageRegistryNamespace: "targetnamespace",
				TargetImageRegistryHost:      "targetregistry.test",
				TargetImageTags:              []string{"1.2.3", "tag1", "tag2"},
				PromoteSourceImageTag:        true,
				RemoveTargetImageTags:        true,
				RemoteSourceImage:            true,
				SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
			},
			prepareMockFunc: func(p *Application) {
				options := &image.PromoteOptions{
					TargetImageName: "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags: []string{
						"targetregistry.test/targetnamespace/targetimage:tag",
						"targetregistry.test/targetnamespace/targetimage:tag1",
						"targetregistry.test/targetnamespace/targetimage:tag2",
						"targetregistry.test/targetnamespace/targetimage:1",
					},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "username_pull",
					PullAuthPassword:      "password_pull",
					PushAuthUsername:      "username_push",
					PushAuthPassword:      "password_push",
				}

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "targetregistry.test").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username_push",
					Password: "password_push",
				}, nil)

				p.credentials.(*authfactory.MockAuthFactory).On("Get", "registry.test").Return(&authmethodkeyfile.KeyFileAuthMethod{}, nil)

				mock := mock.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := factory.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: errors.New(errContext, "Invalid credentials method for 'registry.test'. Found 'keyfile' when is expected basic auth method"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil {
				test.prepareMockFunc(test.service)
			}

			err := test.service.Promote(test.context, test.options)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				promote, _ := test.service.factory.Get(image.DockerPromoterName)
				promote.(*mock.MockPromote).AssertExpectations(t)
			}
		})
	}
}

func TestGetCredentials(t *testing.T) {
	errContext := "(Service::getCredentials)"

	tests := []struct {
		desc            string
		service         *Application
		registry        string
		prepareMockFunc func(*Application)
		res             repository.AuthMethodReader
		err             error
	}{
		{
			desc:    "Testing error when credentials store is not initialized",
			service: &Application{},
			err:     errors.New(errContext, "Credentials has not been initialized"),
		},
		{
			desc: "Testing get credentials",
			service: &Application{
				credentials: authfactory.NewMockAuthFactory(),
			},
			registry: "myregistry",
			prepareMockFunc: func(p *Application) {
				p.credentials.(*authfactory.MockAuthFactory).On("Get", "myregistry").Return(&authmethodbasic.BasicAuthMethod{
					Username: "username",
					Password: "password",
				}, nil)

			},
			res: &authmethodbasic.BasicAuthMethod{
				Username: "username",
				Password: "password",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil {
				test.prepareMockFunc(test.service)
			}

			res, err := test.service.getCredentials(test.registry)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestGetPromoter(t *testing.T) {

	errContext := "(Handler::getPromoter)"

	tests := []struct {
		desc              string
		service           *Application
		options           *Options
		prepareAssertFunc func(*Application)
		res               repository.Promoter
		err               error
	}{
		{
			desc:    "Testing error when promote factory is nil",
			service: &Application{},
			err:     errors.New(errContext, "Promote factory has not been initialized"),
		},
		{
			desc: "Testing get promoter",
			service: &Application{
				factory: factory.NewPromoteFactory(),
			},
			options: &Options{},
			prepareAssertFunc: func(p *Application) {
				p.factory.Register(image.DockerPromoterName, &docker.DockerPromete{})
			},
			res: &docker.DockerPromete{},
			err: &errors.Error{},
		},
		{
			desc: "Testing get promoter with dry-run",
			service: &Application{
				factory: factory.NewPromoteFactory(),
			},
			options: &Options{
				DryRun: true,
			},
			prepareAssertFunc: func(p *Application) {
				p.factory.Register(image.DockerPromoterName, &docker.DockerPromete{})
				p.factory.Register(image.DryRunPromoterName, &dryrun.DryRunPromote{})
			},
			res: &dryrun.DryRunPromote{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.service)
			}

			res, err := test.service.getPromoter(test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.IsType(t, test.res, res)
			}

		})
	}
}
