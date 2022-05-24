package docker

import (
	"context"
	"io/ioutil"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/docker/godockerbuilder"
	"github.com/stretchr/testify/assert"
)

func TestNewDockerPromote(t *testing.T) {
	p := NewDockerPromote(godockerbuilder.NewPromoteMock(), nil)

	assert.NotNil(t, p.cmd, "Failed because copier is nil")
	assert.NotNil(t, p.writer, "Failed because writer is nil")
}

func TestPromote(t *testing.T) {

	contextError := "(docker::Promote)"

	dummyWriter := ioutil.Discard

	tests := []struct {
		desc              string
		options           *image.PromoteOptions
		prom              *DockerPromete
		prepareAssertFunc func(*DockerPromete, *image.PromoteOptions)
		assertFunc        func(*DockerPromete) bool
		err               error
	}{
		{
			desc: "Testing promote with nil copy command",
			prom: &DockerPromete{
				cmd: nil,
			},
			err: errors.New(contextError, "Command to copy docker images must be initialized before promote an image to docker registry"),
		},
		{
			desc: "Testing promote with nil writer on copy command",
			prom: &DockerPromete{
				cmd: godockerbuilder.NewPromoteMock(),
			},
			err: errors.New(contextError, "Writer must be initialized before promote an image to docker registry"),
		},
		{
			desc: "Testing promote with nil options on copy command",
			prom: &DockerPromete{
				cmd:    godockerbuilder.NewPromoteMock(),
				writer: dummyWriter,
			},
			options: nil,
			err:     errors.New(contextError, "Image could not be promoted because options must be defined"),
		},
		{
			desc: "Testing promote with undefined image source name on promote options",
			prom: &DockerPromete{
				cmd:    godockerbuilder.NewPromoteMock(),
				writer: dummyWriter,
			},
			options: &image.PromoteOptions{},
			err:     errors.New(contextError, "Image could not be promoted because source image name must be defined on promote options"),
		},
		{
			desc: "Testing promote with undefined target image name on promote options",
			prom: &DockerPromete{
				cmd:    godockerbuilder.NewPromoteMock(),
				writer: dummyWriter,
			},
			options: &image.PromoteOptions{
				SourceImageName: "image",
			},
			err: errors.New(contextError, "Image could not be promoted because target image name must be defined on promote options"),
		},
		{
			desc: "Testing promote remote run failure",
			prom: &DockerPromete{
				cmd:    godockerbuilder.NewPromoteMock(),
				writer: dummyWriter,
			},
			options: &image.PromoteOptions{
				SourceImageName: "image",
				TargetImageName: "promoteImage",
				TargetImageTags: []string{"tag1", "tag2"},
			},
			err: errors.New(contextError, "Image 'image' could not be promoted", errors.New(contextError, "error from mock")),
			prepareAssertFunc: func(m *DockerPromete, o *image.PromoteOptions) {
				m.cmd.(*godockerbuilder.PromoteMock).On("WithSourceImage", o.SourceImageName)
				m.cmd.(*godockerbuilder.PromoteMock).On("WithTargetImage", o.TargetImageName)
				m.cmd.(*godockerbuilder.PromoteMock).On("WithResponse", m.writer, o.TargetImageName)
				m.cmd.(*godockerbuilder.PromoteMock).On("WithTags", o.TargetImageTags)
				m.cmd.(*godockerbuilder.PromoteMock).On("WithUseNormalizedNamed")
				m.cmd.(*godockerbuilder.PromoteMock).On("Run", context.TODO()).Return(errors.New(contextError, "error from mock"))
			},
			assertFunc: nil,
		},
		{
			desc: "Testing promote remote image",
			prom: &DockerPromete{
				cmd:    godockerbuilder.NewPromoteMock(),
				writer: dummyWriter,
			},
			options: &image.PromoteOptions{
				SourceImageName:       "sourceRegistry/namespace/image",
				TargetImageName:       "targetRegistry/namespace/image",
				TargetImageTags:       []string{"tag1", "tag2"},
				RemoveTargetImageTags: true,
				RemoteSourceImage:     true,
				PullAuthUsername:      "pullname",
				PullAuthPassword:      "pullpass",
				PushAuthUsername:      "pushname",
				PushAuthPassword:      "pushpass",
			},
			err: &errors.Error{},
			prepareAssertFunc: func(m *DockerPromete, o *image.PromoteOptions) {

				// m.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "sourceRegistry").Return(&credentials.RegistryUserPassAuth{
				// 	Username: "name",
				// 	Password: "pass",
				// }, nil)

				// m.credentials.(*credentials.CredentialsStoreMock).On("GetCredentials", "targetRegistry").Return(&credentials.RegistryUserPassAuth{
				// 	Username: "name",
				// 	Password: "pass",
				// }, nil)

				m.cmd.(*godockerbuilder.PromoteMock).On("WithSourceImage", o.SourceImageName)
				m.cmd.(*godockerbuilder.PromoteMock).On("WithTargetImage", o.TargetImageName)
				m.cmd.(*godockerbuilder.PromoteMock).On("WithResponse", m.writer, o.TargetImageName)
				m.cmd.(*godockerbuilder.PromoteMock).On("WithTags", o.TargetImageTags)
				m.cmd.(*godockerbuilder.PromoteMock).On("WithUseNormalizedNamed")
				m.cmd.(*godockerbuilder.PromoteMock).On("WithRemoteSource")
				m.cmd.(*godockerbuilder.PromoteMock).On("WithRemoveAfterPush")
				m.cmd.(*godockerbuilder.PromoteMock).On("WithTags", o.TargetImageTags)
				m.cmd.(*godockerbuilder.PromoteMock).On("AddPullAuth", "pullname", "pullpass").Return(nil)
				m.cmd.(*godockerbuilder.PromoteMock).On("AddPushAuth", "pushname", "pushpass").Return(nil)
				m.cmd.(*godockerbuilder.PromoteMock).On("Run", context.TODO()).Return(nil)
			},
			assertFunc: func(m *DockerPromete) bool {
				return m.cmd.(*godockerbuilder.PromoteMock).AssertNumberOfCalls(t, "WithResponse", 1) &&
					m.cmd.(*godockerbuilder.PromoteMock).AssertNumberOfCalls(t, "WithSourceImage", 1) &&
					m.cmd.(*godockerbuilder.PromoteMock).AssertNumberOfCalls(t, "WithTargetImage", 1) &&
					m.cmd.(*godockerbuilder.PromoteMock).AssertNumberOfCalls(t, "WithResponse", 1) &&
					m.cmd.(*godockerbuilder.PromoteMock).AssertNumberOfCalls(t, "WithTags", 1) &&
					m.cmd.(*godockerbuilder.PromoteMock).AssertNumberOfCalls(t, "WithUseNormalizedNamed", 1) &&
					m.cmd.(*godockerbuilder.PromoteMock).AssertNumberOfCalls(t, "AddPushAuth", 1) &&
					m.cmd.(*godockerbuilder.PromoteMock).AssertNumberOfCalls(t, "AddPullAuth", 1) &&
					m.cmd.(*godockerbuilder.PromoteMock).AssertNumberOfCalls(t, "Run", 1)
				// m.credentials.(*credentials.CredentialsStoreMock).AssertNumberOfCalls(t, "GetCredentials", 2)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.prom, test.options)
			}

			err := test.prom.Promote(context.TODO(), test.options)

			// if err != nil && assert.Error(t, err) {
			// 	assert.Equal(t, test.err.Error(), err.Error())
			// }
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				if test.assertFunc != nil {
					assert.True(t, test.assertFunc(test.prom))
				} else {
					t.Error(test.desc, "missing assertFunc")
				}
			}

		})
	}

}
