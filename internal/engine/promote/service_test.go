package promote

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/credentials"
	promoterepository "github.com/gostevedore/stevedore/internal/promote"
	mockpromote "github.com/gostevedore/stevedore/internal/promote/mock"
	"github.com/gostevedore/stevedore/internal/semver"
	"github.com/stretchr/testify/assert"
)

func TestPromote(t *testing.T) {

	promoteMockID := "mock"

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
			service: &Service{
				credentials:   credentials.NewCredentialsStoreMock(),
				semver:        semver.NewSemVerGenerator(),
				configuration: &configuration.Configuration{},
			},
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

				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "registry.test").Return(&credentials.RegistryUserPassAuth{
					Username: "name",
					Password: "pass",
				}, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(promoteMockID, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote source image from remote",
			service: &Service{
				credentials:   credentials.NewCredentialsStoreMock(),
				semver:        semver.NewSemVerGenerator(),
				configuration: &configuration.Configuration{},
			},
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
				factory.Register(promoteMockID, mock)
				p.factory = factory

				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "registry.test").Return(&credentials.RegistryUserPassAuth{
					Username: "name",
					Password: "pass",
				}, nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote source image with all options",
			service: &Service{
				credentials:   credentials.NewCredentialsStoreMock(),
				semver:        semver.NewSemVerGenerator(),
				configuration: &configuration.Configuration{},
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

				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "registry.test").Return(&credentials.RegistryUserPassAuth{
					Username: "pullname",
					Password: "pullpass",
				}, nil)

				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "targetregistry.test").Return(&credentials.RegistryUserPassAuth{
					Username: "pushname",
					Password: "pushpass",
				}, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(promoteMockID, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote source image with no credentials",
			service: &Service{
				credentials:   credentials.NewCredentialsStoreMock(),
				semver:        semver.NewSemVerGenerator(),
				configuration: &configuration.Configuration{},
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
					PullAuthUsername:      "",
					PullAuthPassword:      "",
					PushAuthUsername:      "",
					PushAuthPassword:      "",
				}

				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "registry.test").Return(nil, nil)
				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "targetregistry.test").Return(nil, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(promoteMockID, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing promote source image with all options and using semver configuration parameters",
			service: &Service{
				credentials: credentials.NewCredentialsStoreMock(),
				semver:      semver.NewSemVerGenerator(),
				configuration: &configuration.Configuration{
					EnableSemanticVersionTags:    true,
					SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
				},
			},
			context: context.TODO(),
			options: &ServiceOptions{
				SourceImageName:              "registry.test/namespace/image:tag",
				EnableSemanticVersionTags:    false,
				TargetImageName:              "targetimage",
				TargetImageRegistryNamespace: "targetnamespace",
				TargetImageRegistryHost:      "targetregistry.test",
				TargetImageTags:              []string{"1.2.3", "tag1", "tag2"},
				PromoteSourceImageTag:        true,
				RemoveTargetImageTags:        true,
				RemoteSourceImage:            true,
				SemanticVersionTagsTemplates: nil,
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

				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "registry.test").Return(&credentials.RegistryUserPassAuth{
					Username: "pullname",
					Password: "pullpass",
				}, nil)

				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "targetregistry.test").Return(&credentials.RegistryUserPassAuth{
					Username: "pushname",
					Password: "pushpass",
				}, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(promoteMockID, mock)
				p.factory = factory
			},
			err: &errors.Error{},
		},

		{
			desc: "Testing promote source image with all options, using semver configuration parameters overridden by service options",
			service: &Service{
				credentials: credentials.NewCredentialsStoreMock(),
				semver:      semver.NewSemVerGenerator(),
				configuration: &configuration.Configuration{
					EnableSemanticVersionTags:    true,
					SemanticVersionTagsTemplates: []string{"{{ .Major }}"},
				},
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

				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "registry.test").Return(&credentials.RegistryUserPassAuth{
					Username: "pullname",
					Password: "pullpass",
				}, nil)

				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "targetregistry.test").Return(&credentials.RegistryUserPassAuth{
					Username: "pushname",
					Password: "pushpass",
				}, nil)

				mock := mockpromote.NewMockPromote()
				mock.On("Promote", context.TODO(), options).Return(nil)

				factory := promoterepository.NewPromoteFactory()
				factory.Register(promoteMockID, mock)
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

			err := test.service.Promote(test.context, test.options, promoteMockID)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				promote, _ := test.service.factory.GetPromoter(promoteMockID)
				promote.(*mockpromote.MockPromote).AssertExpectations(t)
			}

		})
	}
}

func TestGetCredentials(t *testing.T) {
	tests := []struct {
		desc            string
		service         *Service
		registry        string
		prepareMockFunc func(*Service)
		res             *credentials.RegistryUserPassAuth
	}{
		{
			desc: "Testing get credentials",
			service: &Service{
				credentials: credentials.NewCredentialsStoreMock(),
			},
			registry: "myregistry",
			prepareMockFunc: func(p *Service) {
				p.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "myregistry").Return(&credentials.RegistryUserPassAuth{
					Username: "name",
					Password: "pass",
				}, nil)

			},
			res: &credentials.RegistryUserPassAuth{
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

			res := test.service.getCredentials(test.registry)
			assert.Equal(t, test.res, res)
		})
	}
}
