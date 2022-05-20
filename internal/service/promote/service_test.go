package promote

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	credentialsstore "github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/promote"
	promoterepository "github.com/gostevedore/stevedore/internal/promote"
	dockerpromote "github.com/gostevedore/stevedore/internal/promote/docker"
	dryrunpromote "github.com/gostevedore/stevedore/internal/promote/dryrun"
	mockpromote "github.com/gostevedore/stevedore/internal/promote/mock"
	"github.com/gostevedore/stevedore/internal/semver"
	"github.com/stretchr/testify/assert"
)

func TestPromote(t *testing.T) {

	tests := []struct {
		desc            string
		service         *Service
		options         *ServiceOptions
		context         context.Context
		prepareMockFunc func(*Service)
		err             error
	}{
		{
			desc: "Testing promote source image from local",
			service: NewService(
				WithCredentials(credentialsstore.NewCredentialsStoreMock()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(promoterepository.NewPromoteFactory()),
			),
			context: context.TODO(),
			options: &ServiceOptions{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    false,
				TargetImageName:              "",
				TargetImageRegistryNamespace: "",
				TargetImageRegistryHost:      "",
				TargetImageTags:              nil,
				PromoteSourceImageTag:        false,
				RemoveTargetImageTags:        false,
				RemoteSourceImage:            false,
				SemanticVersionTagsTemplates: nil,
			},
			prepareMockFunc: func(p *Service) {

				options := &promoterepository.PromoteOptions{
					TargetImageName:       "registry.test/namespace/image:tag",
					TargetImageTags:       nil,
					RemoveTargetImageTags: false,
					RemoteSourceImage:     false,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "name",
					PullAuthPassword:      "pass",
					PushAuthUsername:      "name",
					PushAuthPassword:      "pass",
				}

				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "registry.test").Return(&credentials.UserPasswordAuth{
					Username: "name",
					Password: "pass",
				}, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)
				p.factory.Register(image.DockerPromoterName, mock)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote source image from remote",
			service: NewService(
				WithCredentials(credentialsstore.NewCredentialsStoreMock()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(promoterepository.NewPromoteFactory()),
			),
			context: context.TODO(),
			options: &ServiceOptions{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    false,
				TargetImageName:              "",
				TargetImageRegistryNamespace: "",
				TargetImageRegistryHost:      "",
				TargetImageTags:              nil,
				PromoteSourceImageTag:        false,
				RemoveTargetImageTags:        false,
				RemoteSourceImage:            true,
				SemanticVersionTagsTemplates: nil,
			},
			prepareMockFunc: func(p *Service) {
				options := &promoterepository.PromoteOptions{
					TargetImageName:       "registry.test/namespace/image:tag",
					TargetImageTags:       nil,
					RemoveTargetImageTags: false,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "name",
					PullAuthPassword:      "pass",
					PushAuthUsername:      "name",
					PushAuthPassword:      "pass",
				}

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory

				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "registry.test").Return(&credentials.UserPasswordAuth{
					Username: "name",
					Password: "pass",
				}, nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote source image with all options",
			service: NewService(
				WithCredentials(credentialsstore.NewCredentialsStoreMock()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(promoterepository.NewPromoteFactory()),
			),
			context: context.TODO(),
			options: &ServiceOptions{
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
			prepareMockFunc: func(p *Service) {
				options := &promoterepository.PromoteOptions{
					TargetImageName:       "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags:       []string{"tag", "tag1", "tag2", "1"},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "pullname",
					PullAuthPassword:      "pullpass",
					PushAuthUsername:      "pushname",
					PushAuthPassword:      "pushpass",
				}

				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "registry.test").Return(&credentials.UserPasswordAuth{
					Username: "pullname",
					Password: "pullpass",
				}, nil)

				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "targetregistry.test").Return(&credentials.UserPasswordAuth{
					Username: "pushname",
					Password: "pushpass",
				}, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote source image with no credentials",
			service: NewService(
				WithCredentials(credentialsstore.NewCredentialsStoreMock()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(promoterepository.NewPromoteFactory()),
			),
			context: context.TODO(),
			options: &ServiceOptions{
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
			prepareMockFunc: func(p *Service) {
				options := &promoterepository.PromoteOptions{
					TargetImageName:       "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags:       []string{"tag", "tag1", "tag2", "1"},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "",
					PullAuthPassword:      "",
					PushAuthUsername:      "",
					PushAuthPassword:      "",
				}

				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "registry.test").Return(&credentials.UserPasswordAuth{}, nil)
				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "targetregistry.test").Return(&credentials.UserPasswordAuth{}, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote source image with all options and using semver configuration parameters",
			service: NewService(
				WithCredentials(credentialsstore.NewCredentialsStoreMock()),
				WithSemver(semver.NewSemVerGenerator()),
				WithPromoteFactory(promoterepository.NewPromoteFactory()),
			),
			context: context.TODO(),
			options: &ServiceOptions{
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
			prepareMockFunc: func(p *Service) {
				options := &promoterepository.PromoteOptions{
					TargetImageName:       "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags:       []string{"tag", "tag1", "tag2", "1"},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "pullname",
					PullAuthPassword:      "pullpass",
					PushAuthUsername:      "pushname",
					PushAuthPassword:      "pushpass",
				}

				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "registry.test").Return(&credentials.UserPasswordAuth{
					Username: "pullname",
					Password: "pullpass",
				}, nil)

				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "targetregistry.test").Return(&credentials.UserPasswordAuth{
					Username: "pushname",
					Password: "pushpass",
				}, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote source image with all options, using semver configuration parameters overridden by service options",
			service: &Service{
				credentials: credentialsstore.NewCredentialsStoreMock(),
				semver:      semver.NewSemVerGenerator(),
			},
			context: context.TODO(),
			options: &ServiceOptions{
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
			prepareMockFunc: func(p *Service) {
				options := &promoterepository.PromoteOptions{
					TargetImageName:       "targetregistry.test/targetnamespace/targetimage:1.2.3",
					TargetImageTags:       []string{"tag", "tag1", "tag2", "1"},
					RemoveTargetImageTags: true,
					RemoteSourceImage:     true,
					SourceImageName:       "registry.test/namespace/image:tag",
					PullAuthUsername:      "pullname",
					PullAuthPassword:      "pullpass",
					PushAuthUsername:      "pushname",
					PushAuthPassword:      "pushpass",
				}

				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "registry.test").Return(&credentials.UserPasswordAuth{
					Username: "pullname",
					Password: "pullpass",
				}, nil)

				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "targetregistry.test").Return(&credentials.UserPasswordAuth{
					Username: "pushname",
					Password: "pushpass",
				}, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(image.DockerPromoterName, mock)
				p.factory = factory
			},
			err: &errors.Error{},
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
				promote.(*mockpromote.MockPromote).AssertExpectations(t)
			}
		})
	}
}

func TestGetCredentials(t *testing.T) {
	errContext := "(Service::getCredentials)"

	tests := []struct {
		desc            string
		service         *Service
		registry        string
		prepareMockFunc func(*Service)
		res             *credentials.UserPasswordAuth
		err             error
	}{
		{
			desc:    "Testing error when credentials store is not initialized",
			service: &Service{},
			err:     errors.New(errContext, "Credentials has not been initialized"),
		},
		{
			desc: "Testing get credentials",
			service: &Service{
				credentials: credentialsstore.NewCredentialsStoreMock(),
			},
			registry: "myregistry",
			prepareMockFunc: func(p *Service) {
				p.credentials.(*credentialsstore.CredentialsStoreMock).On("Get", "myregistry").Return(&credentials.UserPasswordAuth{
					Username: "name",
					Password: "pass",
				}, nil)

			},
			res: &credentials.UserPasswordAuth{
				Username: "name",
				Password: "pass",
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
		service           *Service
		options           *ServiceOptions
		prepareAssertFunc func(*Service)
		res               promote.Promoter
		err               error
	}{
		{
			desc:    "Testing error when promote factory is nil",
			service: &Service{},
			err:     errors.New(errContext, "Promote factory has not been initialized"),
		},
		{
			desc: "Testing get promoter",
			service: &Service{
				factory: promoterepository.NewPromoteFactory(),
			},
			options: &ServiceOptions{},
			prepareAssertFunc: func(p *Service) {
				p.factory.Register(image.DockerPromoterName, &dockerpromote.DockerPromete{})
			},
			res: &dockerpromote.DockerPromete{},
			err: &errors.Error{},
		},
		{
			desc: "Testing get promoter with dry-run",
			service: &Service{
				factory: promoterepository.NewPromoteFactory(),
			},
			options: &ServiceOptions{
				DryRun: true,
			},
			prepareAssertFunc: func(p *Service) {
				p.factory.Register(image.DockerPromoterName, &dockerpromote.DockerPromete{})
				p.factory.Register(image.DryRunPromoterName, &dryrunpromote.DryRunPromote{})
			},
			res: &dryrunpromote.DryRunPromote{},
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
