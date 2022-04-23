package entrypoint

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	"github.com/gostevedore/stevedore/internal/promote"
	repodocker "github.com/gostevedore/stevedore/internal/promote/docker"
	repodryrun "github.com/gostevedore/stevedore/internal/promote/dryrun"
	"github.com/gostevedore/stevedore/internal/semver"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewEntrypoint(t *testing.T) {
	entrypoint := NewEntrypoint(
		WithWriter(ioutil.Discard),
	)

	assert.NotNil(t, entrypoint.writer)
}
func TestExecute(t *testing.T) {
	errContext := "(Entrypoint::Execute)"
	var err error

	baseDir := "/credentials"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "30a88abceb172130caa0a565ea982653"), []byte(`
{
	"docker_login_username": "username",
	"docker_login_password": "password"
}
`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		entrypoint        *Entrypoint
		args              []string
		configuration     *configuration.Configuration
		entrypointOptions *Options
		handlerOptions    *handler.Options
		err               error
		assertions        func(*testing.T, *Entrypoint, []string, *Options, *handler.Options)
	}{
		{
			desc:       "Testing error when configuration is not provided",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To execute the promote entrypoint, configuration is required"),
		},
		{
			desc:          "Testing error when arguments are not provided",
			entrypoint:    &Entrypoint{},
			configuration: &configuration.Configuration{},
			err:           errors.New(errContext, "To execute the promote entrypoint, arguments are required"),
		},
		{
			desc:          "Testing error when handler options are not provided",
			entrypoint:    &Entrypoint{},
			configuration: &configuration.Configuration{},
			args:          []string{"image"},
			err:           errors.New(errContext, "To execute the promote entrypoint, handler options are required"),
		},
		{
			desc: "Testing execute entrypoint",
			entrypoint: NewEntrypoint(
				WithWriter(ioutil.Discard),
				WithFileSystem(testFs),
			),
			configuration: &configuration.Configuration{
				DockerCredentialsDir:      baseDir,
				EnableSemanticVersionTags: true,
				SemanticVersionTagsTemplates: []string{
					"template",
				},
			},
			args: []string{"image"},
			handlerOptions: &handler.Options{
				DryRun: true,
			},
			//			err:            errors.New(errContext, "Image 'image' could not be promoted\n\tError tagging image 'image' to 'docker.io/library/image:latest'\n\tError response from daemon: No such image: image:latest"),
			err: &errors.Error{},
			assertions: func(t *testing.T, e *Entrypoint, args []string, entrypointOptions *Options, handlerOptions *handler.Options) {
				assert.True(t, handlerOptions.EnableSemanticVersionTags, "Enable semantic version tags should be true")
				assert.Equal(t, []string{"template"}, handlerOptions.SemanticVersionTagsTemplates, "Semantic version tags templates is not as expected")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.entrypoint.Execute(context.TODO(), test.args, test.configuration, test.handlerOptions)
			if err != nil {
				assert.Equal(t, err.Error(), test.err.Error())
			} else {
				test.assertions(t, test.entrypoint, test.args, test.entrypointOptions, test.handlerOptions)
			}
		})
	}
}

func TestCreateCredentialsStore(t *testing.T) {
	var err error

	baseDir := "/credentials"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "30a88abceb172130caa0a565ea982653"), []byte(`
{
	"docker_login_username": "username",
	"docker_login_password": "password"
}
`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		fs         afero.Fs
		path       string
		err        error
	}{
		{
			desc: "Testing create credentials store",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			fs:   testFs,
			path: baseDir,
			err:  &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.entrypoint.createCredentialsStore(test.path)
			if err != nil {
				assert.Equal(t, err.Error(), test.err.Error())
			} else {
				assert.NotNil(t, store)
			}
		})
	}

}

func TestCreatePromoteFactory(t *testing.T) {

	var err error
	var promoteRepoFactory promote.PromoteFactory
	var promoteRepoDocker, promoteRepoDryRun promote.Promoter

	e := NewEntrypoint()
	promoteRepoFactory, err = e.createPromoteFactory()

	t.Run("Testing not error is returned", func(t *testing.T) {
		assert.Nil(t, err)
		assert.NotNil(t, promoteRepoFactory)
	})

	t.Run("Testing docker promote repository is returned", func(t *testing.T) {
		promoteRepoDocker, err = promoteRepoFactory.Get(promote.DockerPromoterName)
		assert.Nil(t, err)
		assert.IsType(t, &repodocker.DockerPromete{}, promoteRepoDocker)
	})

	t.Run("Testing dry run promote repository is returned", func(t *testing.T) {
		promoteRepoDryRun, err = promoteRepoFactory.Get(promote.DryRunPromoterName)
		assert.Nil(t, err)
		assert.IsType(t, &repodryrun.DryRunPromote{}, promoteRepoDryRun)
	})
}

func TestCreateSemanticVersionFactory(t *testing.T) {

	t.Run("Testing create semantic version factory", func(t *testing.T) {
		e := NewEntrypoint()

		sv, err := e.createSemanticVersionFactory()

		assert.Nil(t, err)
		assert.NotNil(t, sv)
		assert.IsType(t, &semver.SemVerGenerator{}, sv)
	})
}
