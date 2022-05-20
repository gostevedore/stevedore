package godockerbuilder

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/build"
	godockerbuilderbuildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
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
				cmd:            &MockDockerBuildCmd{},
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
				cmd:            &MockDockerBuildCmd{},
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
				cmd:            &MockDockerBuildCmd{},
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
				test.prepareAssertFunc(test.driver.cmd)
			}

			err := test.driver.AddBuildContext(test.options...)

			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				if test.assertFunc != nil {
					assert.True(t, test.assertFunc(test.driver.cmd))
				} else {
					t.Error(test.desc, "missing assertFunc")
				}
			}

		})
	}
}

func TestWithDockerfile(t *testing.T) {
	t.Log("Testing WithDockerfile")

	driver := &GoDockerBuildDriver{
		cmd:            &MockDockerBuildCmd{},
		contextFactory: nil,
	}
	driver.cmd.(*MockDockerBuildCmd).On("WithDockerfile", "my-dockerfile").Return(&MockDockerBuildCmd{})
	driver.WithDockerfile("my-dockerfile")

	assert.Equal(t, driver.cmd.(*build.DockerBuildCmd).ImageBuildOptions.Dockerfile, "my-dockerfile")
}
func TestWithImageName(t *testing.T) {
	t.Log("Testing WithImageName")

	driver := &GoDockerBuildDriver{
		cmd:            &MockDockerBuildCmd{},
		contextFactory: nil,
	}
	driver.cmd.(*MockDockerBuildCmd).On("WithImageName", "image-name").Return(&MockDockerBuildCmd{})
	driver.WithImageName("image-name")

	assert.Equal(t, driver.cmd.(*build.DockerBuildCmd).ImageName, "image-name")
}
func TestWithPullParentImage(t *testing.T) {
	t.Log("Testing WithPullParentImage")

	driver := &GoDockerBuildDriver{
		cmd:            &MockDockerBuildCmd{},
		contextFactory: nil,
	}
	driver.cmd.(*MockDockerBuildCmd).On("WithPullParentImage").Return(&MockDockerBuildCmd{})
	driver.WithPullParentImage()

	assert.True(t, driver.cmd.(*build.DockerBuildCmd).PullParentImage)
}
func TestWithPushAfterBuild(t *testing.T) {
	t.Log("Testing WithPushAfterBuild")

	driver := &GoDockerBuildDriver{
		cmd:            &MockDockerBuildCmd{},
		contextFactory: nil,
	}
	driver.cmd.(*MockDockerBuildCmd).On("WithPushAfterBuild").Return(&MockDockerBuildCmd{})
	driver.WithPushAfterBuild()

	assert.True(t, driver.cmd.(*build.DockerBuildCmd).PushAfterBuild)
}
func TestWithUseNormalizedNamed(t *testing.T) {
	t.Log("Testing WithUseNormalizedNamed")

	driver := &GoDockerBuildDriver{
		cmd:            &MockDockerBuildCmd{},
		contextFactory: nil,
	}
	driver.cmd.(*MockDockerBuildCmd).On("WithUseNormalizedNamed").Return(&MockDockerBuildCmd{})
	driver.WithUseNormalizedNamed()

	assert.True(t, driver.cmd.(*build.DockerBuildCmd).UseNormalizedNamed)
}
func TestWithRemoveAfterPush(t *testing.T) {
	t.Log("Testing WithRemoveAfterPush")

	driver := &GoDockerBuildDriver{
		cmd:            &MockDockerBuildCmd{},
		contextFactory: nil,
	}
	driver.cmd.(*MockDockerBuildCmd).On("WithRemoveAfterPush").Return(&MockDockerBuildCmd{})
	driver.WithRemoveAfterPush()

	assert.True(t, driver.cmd.(*build.DockerBuildCmd).RemoveAfterPush)
}
