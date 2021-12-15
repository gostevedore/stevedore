package godockerbuilder

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	godockerbuilderbuildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	dockerbuildcontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context"
	gitcontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/git"
	pathcontext "github.com/gostevedore/stevedore/internal/driver/docker/godockerbuilder/context/path"
	"github.com/stretchr/testify/assert"
)

func TestAddBuildContext(t *testing.T) {
	errContext := "(godockerbuilder::AddBuildContext)"
	tests := []struct {
		desc              string
		driver            *GoDockerBuildDriver
		options           []*builder.DockerDriverContextOptions
		prepareAssertFunc func(DockerBuilder)
		assertFunc        func(DockerBuilder) bool
		err               error
	}{
		{
			desc: "Testing error when no options are passed to the method",
			driver: &GoDockerBuildDriver{
				docker:         &MockDockerBuildCmd{},
				contextFactory: nil,
			},
			options:           nil,
			prepareAssertFunc: nil,
			assertFunc:        nil,
			err:               errors.New(errContext, "Docker build context options are missing"),
		},
		{
			desc: "Testing error when options are nil",
			driver: &GoDockerBuildDriver{
				docker:         &MockDockerBuildCmd{},
				contextFactory: nil,
			},
			options:           []*builder.DockerDriverContextOptions{},
			prepareAssertFunc: nil,
			assertFunc:        nil,
			err:               errors.New(errContext, "No Docker build context is defined"),
		},
		{
			desc: "Testing add Docker build context",
			driver: &GoDockerBuildDriver{
				docker:         &MockDockerBuildCmd{},
				contextFactory: &dockerbuildcontext.DockerBuildContextFactory{},
			},
			options: []*builder.DockerDriverContextOptions{
				{
					Path: "my-path",
				},
				{
					Git: &builder.DockerDriverGitContextOptions{
						Repository: "my-repository",
						Reference:  "main",
					},
				},
			},
			prepareAssertFunc: func(b DockerBuilder) {
				pathContext := pathcontext.NewPathBuildContext()
				pathContext.WithPath("my-path")

				gitContext := gitcontext.NewGitBuildContext()
				gitContext.WithRepository("my-repository")
				gitContext.WithReference("main")

				contextList := []godockerbuilderbuildcontext.DockerBuildContexter{
					pathContext,
					gitContext,
				}

				b.(*MockDockerBuildCmd).On("AddBuildContext", contextList).Return(nil)
			},
			assertFunc: func(b DockerBuilder) bool {
				return b.(*MockDockerBuildCmd).AssertExpectations(t)
			},
			err: errors.New(errContext, "Docker build context is missing"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.driver.docker)
			}

			err := test.driver.AddBuildContext(test.options...)

			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				if test.assertFunc != nil {
					assert.True(t, test.assertFunc(test.driver.docker))
				} else {
					t.Error(test.desc, "missing assertFunc")
				}
			}

		})
	}
}
