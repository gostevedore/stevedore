package promote

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	repodocker "github.com/gostevedore/stevedore/internal/promote/docker"
	repodryrun "github.com/gostevedore/stevedore/internal/promote/dryrun"
	"github.com/gostevedore/stevedore/internal/promote/factory"
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

func TestPrepareHandlerOptions(t *testing.T) {
	errContext := "(Entrypoint::prepareHandlerOptions)"

	tests := []struct {
		desc           string
		entrypoint     *Entrypoint
		args           []string
		configuration  *configuration.Configuration
		handlerOptions *handler.Options
		res            *handler.Options
		err            error
	}{
		{
			desc: "Testing error when args are nil",
			args: nil,
			err:  errors.New(errContext, "To execute the promote entrypoint, promote image argument is required"),
		},
		{
			desc: "Testing error when promote image is not provided",
			args: []string{},
			err:  errors.New(errContext, "To execute the promote entrypoint, promote image argument is required"),
		},
		{
			desc:       "Testing error when handler options are not provided",
			entrypoint: &Entrypoint{},
			args:       []string{"image"},
			err:        errors.New(errContext, "To execute the promote entrypoint, handler options are required"),
		},
		{
			desc:           "Testing error when configuration is not provided",
			args:           []string{"image"},
			entrypoint:     &Entrypoint{},
			handlerOptions: &handler.Options{},
			configuration:  nil,
			err:            errors.New(errContext, "To execute the promote entrypoint, configuration is required"),
		},
		{
			desc:       "Testing prepare handler options",
			entrypoint: &Entrypoint{},
			args:       []string{"image"},
			err:        &errors.Error{},
			configuration: &configuration.Configuration{
				SemanticVersionTagsTemplates: []string{"template"},
			},
			handlerOptions: &handler.Options{
				DryRun:                       true,
				EnableSemanticVersionTags:    true,
				TargetImageName:              "target_image_name",
				TargetImageRegistryNamespace: "target_image_regsitry_namespace",
				TargetImageRegistryHost:      "target_image_registry_host",
				TargetImageTags:              []string{"target_image_tag"},
				RemoveTargetImageTags:        true,
				PromoteSourceImageTag:        true,
				RemoteSourceImage:            true,
			},
			res: &handler.Options{
				DryRun:                       true,
				EnableSemanticVersionTags:    true,
				SourceImageName:              "image",
				TargetImageName:              "target_image_name",
				TargetImageRegistryNamespace: "target_image_regsitry_namespace",
				TargetImageRegistryHost:      "target_image_registry_host",
				TargetImageTags:              []string{"target_image_tag"},
				RemoveTargetImageTags:        true,
				SemanticVersionTagsTemplates: []string{"template"},
				PromoteSourceImageTag:        true,
				RemoteSourceImage:            true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			options, err := test.entrypoint.prepareHandlerOptions(test.args, test.configuration, test.handlerOptions)
			if err != nil {
				assert.Equal(t, err.Error(), test.err.Error())
			} else {
				assert.Equal(t, test.res, options)
			}
		})
	}
}

func TestCreateCredentialsStore(t *testing.T) {
	var err error

	errContext := "(Entrypoint::createPromoteRepoFactory)"

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
		conf       *configuration.Configuration
		err        error
	}{
		{
			desc:       "Testing error when fs is not provided",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To create the credentials store, a file system is required"),
		},
		{
			desc: "Testing error when conf is not provided",
			entrypoint: NewEntrypoint(
				WithWriter(ioutil.Discard),
				WithFileSystem(testFs),
			),
			conf: nil,
			err:  errors.New(errContext, "To execute the promote entrypoint, configuration is required"),
		},
		{
			desc: "Testing error when credentials dir is not provided",
			entrypoint: NewEntrypoint(
				WithWriter(ioutil.Discard),
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "Docker credentials path must be provided in the configuration"),
		},
		{
			desc: "Testing create credentials store",
			entrypoint: NewEntrypoint(
				WithWriter(ioutil.Discard),
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{
				DockerCredentialsDir: baseDir,
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.entrypoint.createCredentialsStore(test.conf)
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
	var promoteRepoFactory factory.PromoteFactory
	var promoteRepoDocker, promoteRepoDryRun repository.Promoter

	e := NewEntrypoint()
	promoteRepoFactory, err = e.createPromoteFactory()

	t.Run("Testing not error is returned", func(t *testing.T) {
		assert.Nil(t, err)
		assert.NotNil(t, promoteRepoFactory)
	})

	t.Run("Testing docker promote repository is returned", func(t *testing.T) {
		promoteRepoDocker, err = promoteRepoFactory.Get(image.DockerPromoterName)
		assert.Nil(t, err)
		assert.IsType(t, &repodocker.DockerPromete{}, promoteRepoDocker)
	})

	t.Run("Testing dry run promote repository is returned", func(t *testing.T) {
		promoteRepoDryRun, err = promoteRepoFactory.Get(image.DryRunPromoterName)
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
