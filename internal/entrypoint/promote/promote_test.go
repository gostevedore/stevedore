package promote

import (
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/promote"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/dryrun"
	"github.com/gostevedore/stevedore/internal/infrastructure/promote/factory"
	defaultreferencename "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/default"
	dockerreferencename "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/semver"
	credentialsenvvarsstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/envvars"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewEntrypoint(t *testing.T) {
	entrypoint := NewEntrypoint(
		WithWriter(console.NewMockConsole()),
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
			desc: "Testing error preparing handler option in the promote entrypoint when args are nil",
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

	errContext := "(promote::entrypoint::createCredentialsStore)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		conf       *configuration.CredentialsConfiguration
		res        repository.CredentialsStorer
		err        error
	}{
		{
			desc:       "Testing error creating credentials store in the promote entrypoint when configuration is not defined in credentials local store",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create credentials store in the promote entrypoint, credentials configuration is required"),
		},
		{
			desc:       "Testing error creating local credentials store in the promote entrypoint when format is undefined in credentials local store",
			entrypoint: NewEntrypoint(),
			conf: &configuration.CredentialsConfiguration{
				StorageType: credentials.LocalStore,
			},
			err: errors.New(errContext, "To create credentials store in the entrypoint, credentials format must be specified"),
		},
		{
			desc:       "Testing error creating local credentials store in the promote entrypoint when compatibilitier is undefined in the entrypoint",
			entrypoint: NewEntrypoint(),
			conf: &configuration.CredentialsConfiguration{
				StorageType: credentials.LocalStore,
				Format:      credentials.JSONFormat,
			},
			err: errors.New(errContext, "To create credentials store in the promote entrypoint, compatibility is required"),
		},
		{
			desc: "Testing error creating local credentials store in the promote entrypoint when local storage is used and local storage path is undefined",
			entrypoint: NewEntrypoint(
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.CredentialsConfiguration{
				StorageType: credentials.LocalStore,
				Format:      credentials.JSONFormat,
			},
			err: errors.New(errContext, "To create credentials store in the promote entrypoint, local storage path is required"),
		},
		{
			desc: "Testing error creating local credentials store in the promote entrypoint when storage type is not supported in credentials local store",
			entrypoint: NewEntrypoint(
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.CredentialsConfiguration{
				StorageType: "unsupported",
				Format:      credentials.JSONFormat,
			},
			err: errors.New(errContext, "Unsupported credentials storage type 'unsupported'"),
		},
		{
			desc: "Testing create local credentials store in the promote entrypoint",
			entrypoint: NewEntrypoint(
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.CredentialsConfiguration{
				StorageType:      credentials.LocalStore,
				LocalStoragePath: "local-storage-path",
				Format:           credentials.JSONFormat,
			},
			res: &credentialslocalstore.LocalStore{},
			err: &errors.Error{},
		},
		{
			desc:       "Testing create envvars credentials store in the promote entrypoint",
			entrypoint: NewEntrypoint(),
			conf: &configuration.CredentialsConfiguration{
				StorageType: credentials.EnvvarsStore,
				Format:      credentials.JSONFormat,
			},
			res: &credentialsenvvarsstore.EnvvarsStore{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.entrypoint.createCredentialsStore(test.conf)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.IsType(t, test.res, store)
			}
		})
	}
}

func TestCreateAuthFactory(t *testing.T) {
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
		desc          string
		entrypoint    *Entrypoint
		conf          *configuration.Configuration
		compatibility Compatibilitier
		err           error
	}{
		{
			desc:       "Testing error creating credentials factory when file system is not provided",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To create the credentials store in the promote entrypoint, a file system is required"),
		},
		{
			desc: "Testing error creating credentials factory when configuration is not provided",
			entrypoint: NewEntrypoint(
				WithWriter(console.NewMockConsole()),
				WithFileSystem(testFs),
			),
			conf: nil,
			err:  errors.New(errContext, "To create the credentials store in the promote entrypoint, configuration is required"),
		},
		{
			desc: "Testing error creating credentials factory when credentials dir is not provided",
			entrypoint: NewEntrypoint(
				WithWriter(console.NewMockConsole()),
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create the credentials store in the promote entrypoint, credentials configuration is required"),
		},
		{
			desc: "Testing create credentials store in the promote entrypoint",
			entrypoint: NewEntrypoint(
				WithWriter(console.NewMockConsole()),
				WithFileSystem(testFs),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.Configuration{
				Credentials: &configuration.CredentialsConfiguration{
					StorageType:      credentials.LocalStore,
					LocalStoragePath: baseDir,
					Format:           credentials.JSONFormat,
				},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.entrypoint.createAuthFactory(test.conf)
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

	t.Run("Testing create promote factory and not error is returned in the promote entrypoint", func(t *testing.T) {
		assert.Nil(t, err)
		assert.NotNil(t, promoteRepoFactory)
	})

	t.Run("Testing create promote factory and docker promote repository is returned in the promote entrypoint", func(t *testing.T) {
		promoteRepoDocker, err = promoteRepoFactory.Get(image.DockerPromoterName)
		assert.Nil(t, err)
		assert.IsType(t, &docker.DockerPromete{}, promoteRepoDocker)
	})

	t.Run("Testing create promote factory and dry run promote repository is returned in the promote entrypoint", func(t *testing.T) {
		promoteRepoDryRun, err = promoteRepoFactory.Get(image.DryRunPromoterName)
		assert.Nil(t, err)
		assert.IsType(t, &dryrun.DryRunPromote{}, promoteRepoDryRun)
	})
}

func TestCreateSemanticVersionFactory(t *testing.T) {

	t.Run("Testing create semantic version factory in the promote entrypoint", func(t *testing.T) {
		e := NewEntrypoint()

		sv, err := e.createSemanticVersionFactory()

		assert.Nil(t, err)
		assert.NotNil(t, sv)
		assert.IsType(t, &semver.SemVerGenerator{}, sv)
	})
}

func TestCreateReferenceName(t *testing.T) {
	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		options    *Options
		res        repository.ImageReferenceNamer
		err        error
	}{
		{
			desc:       "Testinc create docker reference name on promote entrypoint",
			entrypoint: NewEntrypoint(),
			options: &Options{
				UseDockerNormalizedName: true,
			},
			res: dockerreferencename.NewDockerNormalizedReferenceName(),
		},
		{
			desc:       "Testinc create default reference name on promote entrypoint",
			entrypoint: NewEntrypoint(),
			options:    &Options{},
			res:        defaultreferencename.NewDefaultReferenceName(),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.entrypoint.createReferenceName(test.options)
			if err != nil {
				assert.Equal(t, test.res, res)
			} else {
				assert.IsType(t, test.res, res)
			}
		})
	}
}
