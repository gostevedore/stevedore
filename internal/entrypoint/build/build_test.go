package build

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	imagesconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images"
	imagesgraphtemplate "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	credentialsfactory "github.com/gostevedore/stevedore/internal/infrastructure/credentials/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/ansible"
	defaultdriver "github.com/gostevedore/stevedore/internal/infrastructure/driver/default"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/dryrun"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/now"
	"github.com/gostevedore/stevedore/internal/infrastructure/plan"
	defaultreferencename "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/default"
	dockerreferencename "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/render"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/job"
	"github.com/gostevedore/stevedore/internal/infrastructure/semver"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
	credentialslocalstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestPrepareEntrypointOptions(t *testing.T) {
	errContext := "(entrypoint::build::prepareEntrypointOptions)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		conf       *configuration.Configuration
		options    *Options
		res        *Options
		err        error
	}{
		{
			desc:       "Testing error preparing build entrypoint options when configuration is not provided",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To prepare build entrypoint options, configuration is required"),
		},
		{
			desc:       "Testing error preparing build entrypoint options when options are not provided",
			entrypoint: &Entrypoint{},
			conf:       &configuration.Configuration{},
			err:        errors.New(errContext, "To prepare build entrypoint options, entrypoint options are required"),
		},
		{
			desc:       "Testing prepare build entrypoint options",
			entrypoint: &Entrypoint{},
			conf: &configuration.Configuration{
				Concurrency: 5,
			},
			options: &Options{
				Concurrency: 10,
				Debug:       true,
			},
			res: &Options{
				Concurrency: 10,
				Debug:       true,
			},
			err: &errors.Error{},
		},
		{
			desc:       "Testing prepare build entrypoint options using configuration concurrency",
			entrypoint: &Entrypoint{},
			conf: &configuration.Configuration{
				Concurrency: 5,
			},
			options: &Options{
				Concurrency: 0,
				Debug:       true,
			},
			res: &Options{
				Concurrency: 5,
				Debug:       true,
			},
			err: &errors.Error{},
		},
		{
			desc:       "Testing prepare build entrypoint options using configuration concurrency and dryrun enabled",
			entrypoint: &Entrypoint{},
			conf: &configuration.Configuration{
				Concurrency: 5,
			},
			options: &Options{
				Concurrency: 0,
				DryRun:      true,
				Debug:       true,
			},
			res: &Options{
				Concurrency: 1,
				Debug:       true,
				DryRun:      true,
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			options, err := test.entrypoint.prepareEntrypointOptions(test.conf, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, options)
			}
		})
	}
}

func TestPrepareImageName(t *testing.T) {

	errContext := "(entrypoint::build::prepareImageName)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		args       []string
		res        string
		err        error
	}{
		{
			desc:       "Testing error preparing image name in build entrypoint when no args is nil",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To execute the build entrypoint, arguments are required"),
		},
		{
			desc:       "Testing error preparing image name in build entrypoint when no args are provided",
			entrypoint: &Entrypoint{},
			args:       []string{},
			err:        errors.New(errContext, "To execute the build entrypoint, arguments are required"),
		},
		{
			desc:       "Testing prepare image name in build entrypoint",
			entrypoint: &Entrypoint{},
			args:       []string{"image"},
			res:        "image",
			err:        &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			name, err := test.entrypoint.prepareImageName(test.args)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, name)
			}
		})
	}
}

func TestPrepareHandlerOptions(t *testing.T) {
	errContext := "(entrypoint::build::prepareHandlerOptions)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		conf       *configuration.Configuration
		options    *handler.Options
		res        *handler.Options
		err        error
	}{
		{
			desc:       "Testing error preparing handler options in build entrypoint when configuration is not provided",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To prepare handler options in build entrypoint, configuration is required"),
		},
		{
			desc:       "Testing error preparing handler options in build entrypoint when options are not provided",
			entrypoint: &Entrypoint{},
			conf:       &configuration.Configuration{},
			err:        errors.New(errContext, "To prepare handler options in build entrypoint, handler options are required"),
		},
		{
			desc:       "Testing prepare handler options in build entrypoint",
			entrypoint: &Entrypoint{},
			conf:       &configuration.Configuration{},
			options: &handler.Options{
				AnsibleConnectionLocal:           true,
				AnsibleIntermediateContainerName: "ansible-intermediate-container",
				AnsibleInventoryPath:             "ansible-inventory-path",
				AnsibleLimit:                     "ansible-limit",
				BuildOnCascade:                   true,
				CascadeDepth:                     3,
				EnableSemanticVersionTags:        true,
				ImageFromName:                    "image-from-name",
				ImageFromRegistryHost:            "image-from-registry-host",
				ImageFromRegistryNamespace:       "image-from-registry-namespace",
				ImageFromVersion:                 "image-from-version",
				ImageName:                        "image-name",
				ImageRegistryHost:                "image-registry-host",
				ImageRegistryNamespace:           "image-registry-namespace",
				Labels:                           []string{"label1", "label2"},
				PersistentLabels:                 []string{"plabel1", "plabel2"},
				PersistentVars:                   []string{"pvar1", "pvar2"},
				PullParentImage:                  true,
				PushImagesAfterBuild:             true,
				RemoveImagesAfterPush:            true,
				SemanticVersionTagsTemplates:     []string{"semantic-version-tags-template1", "semantic-version-tags-template2"},
				Tags:                             []string{"tag1", "tag2"},
				Vars:                             []string{"var1", "var2"},
				Versions:                         []string{"version1", "version2"},
			},
			res: &handler.Options{
				AnsibleConnectionLocal:           true,
				AnsibleIntermediateContainerName: "ansible-intermediate-container",
				AnsibleInventoryPath:             "ansible-inventory-path",
				AnsibleLimit:                     "ansible-limit",
				BuildOnCascade:                   true,
				CascadeDepth:                     3,
				EnableSemanticVersionTags:        true,
				ImageFromName:                    "image-from-name",
				ImageFromRegistryHost:            "image-from-registry-host",
				ImageFromRegistryNamespace:       "image-from-registry-namespace",
				ImageFromVersion:                 "image-from-version",
				ImageName:                        "image-name",
				ImageRegistryHost:                "image-registry-host",
				ImageRegistryNamespace:           "image-registry-namespace",
				Labels:                           []string{"label1", "label2"},
				PersistentLabels:                 []string{"plabel1", "plabel2"},
				PersistentVars:                   []string{"pvar1", "pvar2"},
				PullParentImage:                  true,
				PushImagesAfterBuild:             true,
				RemoveImagesAfterPush:            true,
				SemanticVersionTagsTemplates:     []string{"semantic-version-tags-template1", "semantic-version-tags-template2"},
				Tags:                             []string{"tag1", "tag2"},
				Vars:                             []string{"var1", "var2"},
				Versions:                         []string{"version1", "version2"},
			},
			err: &errors.Error{},
		},
		{
			desc:       "Testing prepare handler options using also configuration options in build entrypoint",
			entrypoint: &Entrypoint{},
			conf: &configuration.Configuration{
				SemanticVersionTagsTemplates: []string{"conf-semantic-version-tags-template1", "conf-semantic-version-tags-template2"},
			},
			options: &handler.Options{
				AnsibleConnectionLocal:           true,
				AnsibleIntermediateContainerName: "ansible-intermediate-container",
				AnsibleInventoryPath:             "ansible-inventory-path",
				AnsibleLimit:                     "ansible-limit",
				BuildOnCascade:                   true,
				CascadeDepth:                     3,
				EnableSemanticVersionTags:        true,
				ImageFromName:                    "image-from-name",
				ImageFromRegistryHost:            "image-from-registry-host",
				ImageFromRegistryNamespace:       "image-from-registry-namespace",
				ImageFromVersion:                 "image-from-version",
				ImageName:                        "image-name",
				ImageRegistryHost:                "image-registry-host",
				ImageRegistryNamespace:           "image-registry-namespace",
				Labels:                           []string{"label1", "label2"},
				PersistentLabels:                 []string{"plabel1", "plabel2"},
				PersistentVars:                   []string{"pvar1", "pvar2"},
				PullParentImage:                  true,
				PushImagesAfterBuild:             true,
				RemoveImagesAfterPush:            true,
				Tags:                             []string{"tag1", "tag2"},
				Vars:                             []string{"var1", "var2"},
				Versions:                         []string{"version1", "version2"},
			},
			res: &handler.Options{
				AnsibleConnectionLocal:           true,
				AnsibleIntermediateContainerName: "ansible-intermediate-container",
				AnsibleInventoryPath:             "ansible-inventory-path",
				AnsibleLimit:                     "ansible-limit",
				BuildOnCascade:                   true,
				CascadeDepth:                     3,
				EnableSemanticVersionTags:        true,
				ImageFromName:                    "image-from-name",
				ImageFromRegistryHost:            "image-from-registry-host",
				ImageFromRegistryNamespace:       "image-from-registry-namespace",
				ImageFromVersion:                 "image-from-version",
				ImageName:                        "image-name",
				ImageRegistryHost:                "image-registry-host",
				ImageRegistryNamespace:           "image-registry-namespace",
				Labels:                           []string{"label1", "label2"},
				PersistentLabels:                 []string{"plabel1", "plabel2"},
				PersistentVars:                   []string{"pvar1", "pvar2"},
				PullParentImage:                  true,
				PushImagesAfterBuild:             true,
				RemoveImagesAfterPush:            true,
				SemanticVersionTagsTemplates:     []string{"conf-semantic-version-tags-template1", "conf-semantic-version-tags-template2"},
				Tags:                             []string{"tag1", "tag2"},
				Vars:                             []string{"var1", "var2"},
				Versions:                         []string{"version1", "version2"},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			options, err := test.entrypoint.prepareHandlerOptions(test.conf, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, options)
			}
		})
	}
}

func TestCreateCredentialsLocalStore(t *testing.T) {

	errContext := "(entrypoint::build::createCredentialsStore)"

	tests := []struct {
		desc          string
		entrypoint    *Entrypoint
		conf          *configuration.CredentialsConfiguration
		compatibility Compatibilitier
		err           error
	}{
		{
			desc:       "Testing error when creating credentials local store in build entrypoint with not defined configuration",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create credentials store in build entrypoint, credentials configuration is required"),
		},
		{
			desc:       "Testing error when creating credentials local store in build entrypoint with undefined format",
			entrypoint: NewEntrypoint(),
			conf: &configuration.CredentialsConfiguration{
				StorageType: credentials.LocalStore,
			},
			err: errors.New(errContext, "To create credentials store in build entrypoint, credentials format must be specified"),
		},
		{
			desc:       "Testing error when creating credentials local store in build entrypoint with undefined compatibilitier",
			entrypoint: NewEntrypoint(),
			conf: &configuration.CredentialsConfiguration{
				StorageType: credentials.LocalStore,
				Format:      credentials.JSONFormat,
			},
			err: errors.New(errContext, "To create credentials store in build entrypoint, compatibility is required"),
		},
		{
			desc: "Testing error when creating credentials local store in build entrypoint with undefined storage path",
			entrypoint: NewEntrypoint(
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.CredentialsConfiguration{
				StorageType: credentials.LocalStore,
				Format:      credentials.JSONFormat,
			},
			err: errors.New(errContext, "To create credentials store in build entrypoint, local storage path is required"),
		},
		{
			desc: "Testing error when creating credentials local store in build entrypoint with unsupported storage type",
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
			desc: "Testing create credentials local store in build entrypoint",
			entrypoint: NewEntrypoint(
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.CredentialsConfiguration{
				StorageType:      "local",
				LocalStoragePath: "local-storage-path",
				Format:           credentials.JSONFormat,
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.entrypoint.createCredentialsLocalStore(test.conf)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.IsType(t, &credentialslocalstore.LocalStore{}, store)
			}
		})
	}
}

func TestCreateCredentialsFactory(t *testing.T) {
	errContext := "(entrypoint::build::createCredentialsFactory)"

	baseDir := "/credentials"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	err := afero.WriteFile(testFs, filepath.Join(baseDir, "30a88abceb172130caa0a565ea982653"), []byte(`
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
		res           repository.CredentialsFactorier
		err           error
	}{
		{
			desc:       "Testing error creating credentials factory in build entrypoint when file system is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create the credentials store in build entrypoint, a file system is required"),
		},
		{
			desc: "Testing error creating credentials factory in build entrypoint when configuration is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			err: errors.New(errContext, "To create the credentials store in build entrypoint, configuration is required"),
		},
		{
			desc: "Testing error creating credentials factory in build entrypoint when credentials configuration is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create the credentials store in build entrypoint, credentials configuration is required"),
		},
		{
			desc: "Testing create credentials factory in build entrypoint",
			entrypoint: NewEntrypoint(
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
			res: &credentialsfactory.CredentialsFactory{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			credentials, err := test.entrypoint.createCredentialsFactory(test.conf)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, credentials)
				assert.IsType(t, test.res, credentials)
			}
		})
	}
}

func TestCreateBuildersStore(t *testing.T) {
	errContext := "(entrypoint::build::createBuildersStore)"

	baseDir := "/builders"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		conf       *configuration.Configuration
		res        *builders.Store
		err        error
	}{
		{
			desc:       "Testing error creating builder store in build entrypoint when file system is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create a builders store in build entrypoint, a file system is required"),
		},
		{
			desc: "Testing error creating builder store in build entrypoint when configuration is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			err: errors.New(errContext, "To create a builders store in build entrypoint, configuration is required"),
		},
		{
			desc: "Testing error creating builder store in build entrypoint when builders path is not defined in configuration",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create a builders store in build entrypoint, builders path must be provided in configuration"),
		},
		{
			desc: "Testing create builders store in build entrypoint",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{
				BuildersPath: baseDir,
			},
			res: &builders.Store{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.entrypoint.createBuildersStore(test.conf)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, store)
				assert.IsType(t, test.res, store)
			}
		})
	}
}

func TestCreateCommandFactory(t *testing.T) {
	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		res        *command.BuildCommandFactory
		err        error
	}{
		{
			desc: "Testing create command factory in build entrypoint",
			res:  &command.BuildCommandFactory{},
			err:  &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			factory, err := test.entrypoint.createCommandFactory()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, factory)
				assert.IsType(t, test.res, factory)
			}
		})
	}
}

func TestCreateJobFactory(t *testing.T) {
	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		res        *job.JobFactory
		err        error
	}{
		{
			desc: "Testing create job factory in build entrypoint",
			res:  &job.JobFactory{},
			err:  &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			factory, err := test.entrypoint.createJobFactory()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, factory)
				assert.IsType(t, test.res, factory)
			}
		})
	}
}

func TestCreateSemVerFactory(t *testing.T) {
	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		res        *semver.SemVerGenerator
		err        error
	}{
		{
			desc: "Testing create semver factory in build entrypoint",
			res:  &semver.SemVerGenerator{},
			err:  &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			factory, err := test.entrypoint.createSemVerFactory()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, factory)
				assert.IsType(t, test.res, factory)
			}
		})
	}
}

func TestCreateImageRender(t *testing.T) {
	errContext := "(entrypoint::build::createImageRender)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		now        render.Nower
		res        *render.ImageRender
		err        error
	}{
		{
			desc:       "Testing error creating image render in build entrypoint when now is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create an image render in build entrypoint, a nower is required"),
		},
		{
			desc: "Testing create image render in build entrypoint",
			now:  now.NewNow(),
			res:  &render.ImageRender{},
			err:  &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			render, err := test.entrypoint.createImageRender(test.now)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, render)
				assert.IsType(t, test.res, render)
			}
		})
	}
}

func TestCreateImagesStore(t *testing.T) {
	errContext := "(entrypoint::build::createImagesStore)"

	baseDir := "/images"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	tests := []struct {
		desc          string
		entrypoint    *Entrypoint
		conf          *configuration.Configuration
		render        repository.Renderer
		graph         imagesconfiguration.ImagesGraphTemplatesStorer
		compatibility Compatibilitier
		res           *images.Store
		err           error
	}{
		{
			desc:       "Testing error creating images store in build entrypoint when fs is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create an images store in build entrypoint, a filesystem is required"),
		},
		{
			desc: "Testing error creating images store in build entrypoint when configuration is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			err: errors.New(errContext, "To create an images store in build entrypoint, configuration is required"),
		},
		{
			desc: "Testing error creating images store in build entrypoint when render is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create an images store in build entrypoint, image render is required"),
		},
		{
			desc: "Testing error creating images store in build entrypoint when graph is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf:   &configuration.Configuration{},
			render: &render.ImageRender{},
			err:    errors.New(errContext, "To create an images store in build entrypoint, images graph templates storer is required"),
		},
		{
			desc: "Testing error creating images store in build entrypoint when compatibility is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf:   &configuration.Configuration{},
			render: &render.ImageRender{},
			graph:  &imagesgraphtemplate.ImagesGraphTemplate{},
			err:    errors.New(errContext, "To create an images store in build entrypoint, compatibility is required"),
		},
		{
			desc: "Testing error creating images store in build entrypoint when images path is not defined in configuration",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf:          &configuration.Configuration{},
			render:        &render.ImageRender{},
			graph:         &imagesgraphtemplate.ImagesGraphTemplate{},
			compatibility: &compatibility.Compatibility{},
			err:           errors.New(errContext, "To create an images store in build entrypoint, images path must be provided in configuration"),
		},
		{
			desc: "Testing create images store in build entrypoint",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf: &configuration.Configuration{
				ImagesPath: baseDir,
			},
			render: &render.ImageRender{},
			graph: imagesgraphtemplate.NewImagesGraphTemplate(
				graph.NewGraphTemplateFactory(false),
			),
			compatibility: &compatibility.Compatibility{},
			res:           &images.Store{},
			err:           &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.entrypoint.createImagesStore(test.conf, test.render, test.graph)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, store)
				assert.IsType(t, test.res, store)
			}

		})
	}
}

func TestCreateImagesGraphTemplatesStorer(t *testing.T) {
	errContext := "(entrypoint::build::createImagesGraphTemplatesStorer)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		factory    *graph.GraphTemplateFactory
		res        *imagesgraphtemplate.ImagesGraphTemplate
		err        error
	}{
		{
			desc:       "Testing error creating images graph templates store in build entrypoint when factory is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create an images graph templates storer in build entrypoint, a graph template factory is required"),
		},
		{
			desc:       "Testing create images graph templates storer in build entrypoint",
			entrypoint: NewEntrypoint(),
			factory:    graph.NewGraphTemplateFactory(false),
			res:        &imagesgraphtemplate.ImagesGraphTemplate{},
			err:        &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			store, err := test.entrypoint.createImagesGraphTemplatesStorer(test.factory)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, store)
				assert.IsType(t, test.res, store)
			}
		})
	}
}

func TestCreateGraphTemplateFactory(t *testing.T) {
	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		res        *graph.GraphTemplateFactory
		err        error
	}{
		{
			desc:       "Testing create graph template factory in build entrypoint",
			entrypoint: NewEntrypoint(),
			res:        &graph.GraphTemplateFactory{},
			err:        &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			factory, err := test.entrypoint.createGraphTemplateFactory()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, factory)
				assert.IsType(t, test.res, factory)
			}
		})
	}
}

func TestCreateBuildDriverFactory(t *testing.T) {

	errContext := "(entrypoint::build::createBuildDriverFactory)"

	tests := []struct {
		desc        string
		entrypoint  *Entrypoint
		credentials repository.CredentialsFactorier
		options     *Options
		err         error
		assertions  func(t *testing.T, driverFactory factory.BuildDriverFactory)
	}{
		{
			desc:        "Testing create build driver factory in build entrypoint with empty credentials",
			entrypoint:  NewEntrypoint(),
			credentials: nil,
			options:     nil,
			err:         errors.New(errContext, "Register drivers requires a credentials store in build entrypoint"),
		},
		{
			desc:        "Testing create build driver factory in build entrypoint with empty options",
			entrypoint:  NewEntrypoint(),
			credentials: credentialsfactory.NewMockCredentialsFactory(),
			options:     nil,
			err:         errors.New(errContext, "Register drivers requires options in build entrypoint"),
		},
		{
			desc:        "Testing create build driver factory in build entrypoint with nil writer",
			entrypoint:  NewEntrypoint(),
			credentials: credentialsfactory.NewMockCredentialsFactory(),
			options:     &Options{},
			err:         errors.New(errContext, "Register drivers requires a writer in build entrypoint"),
		},
		{
			desc:        "Testing create build driver factory in build entrypoint",
			entrypoint:  NewEntrypoint(WithWriter(ioutil.Discard)),
			credentials: credentialsfactory.NewMockCredentialsFactory(),
			options:     &Options{},
			err:         &errors.Error{},
			assertions: func(t *testing.T, f factory.BuildDriverFactory) {
				dDockerFunc, eDocker := f.Get("docker")
				assert.Nil(t, eDocker)
				assert.NotNil(t, dDockerFunc)

				dDocker, eDocker := dDockerFunc()
				assert.Nil(t, eDocker)
				assert.IsType(t, &docker.DockerDriver{}, dDocker)

				dAnsibleFunc, eAnsible := f.Get("ansible-playbook")
				assert.Nil(t, eAnsible)
				assert.NotNil(t, dAnsibleFunc)

				dAnsible, eAnsible := dAnsibleFunc()
				assert.Nil(t, eAnsible)
				assert.IsType(t, &ansible.AnsiblePlaybookDriver{}, dAnsible)

				dDefaultFunc, eDefault := f.Get("default")
				assert.Nil(t, eDefault)
				assert.NotNil(t, dDefaultFunc)

				dDefault, eDefault := dDefaultFunc()
				assert.Nil(t, eDefault)
				assert.IsType(t, &defaultdriver.DefaultDriver{}, dDefault)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			factory, err := test.entrypoint.createBuildDriverFactory(test.credentials, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertions(t, factory)
			}
		})
	}
}

func TestCreateDryRunDriver(t *testing.T) {
	desc := "Testing create dry-run driver"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		res        repository.BuildDriverer
		err        error
	}{
		{
			desc:       "Testing create dry-run driver in build entrypoint",
			entrypoint: NewEntrypoint(),
			res:        &dryrun.DryRunDriver{},
		},
	}

	for _, test := range tests {
		t.Run(desc, func(t *testing.T) {

			driverFunc, err := test.entrypoint.createDryRunDriver()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, driverFunc)
				driver, err := driverFunc()
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}

}

func TestCreateDefaultDriver(t *testing.T) {
	errContext := "(entrypoint::createDefaultDriver)"
	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		options    *Options
		res        repository.BuildDriverer
		err        error
	}{
		{
			desc:       "Testing error creating default driver in build entrypoint when options are not provided",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "Build entrypoint options are required to create default driver"),
		},
		{
			desc:       "Testing create default driver in build entrypoint",
			entrypoint: NewEntrypoint(),
			options:    &Options{},
			res:        &defaultdriver.DefaultDriver{},
		},
		{
			desc:       "Testing create default driver in build entrypoint with dryrun enabled",
			entrypoint: NewEntrypoint(),
			options: &Options{
				DryRun: true,
			},
			res: &dryrun.DryRunDriver{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			driverFunc, err := test.entrypoint.createDefaultDriver(test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, driverFunc)
				driver, err := driverFunc()
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}
}

func TestCreateAnsibleDriver(t *testing.T) {

	errContext := "(entrypoint::build::createAnsibleDriver)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		options    *Options
		res        repository.BuildDriverer
		err        error
	}{
		{
			desc:       "Testing error creating ansible driver in build entrypoint when creating ansible driver with nil options",
			entrypoint: NewEntrypoint(),
			options:    nil,
			err:        errors.New(errContext, "Build entrypoint options are required to create ansible driver"),
		},
		{
			desc:       "Testing create ansible driver in build entrypoint",
			entrypoint: NewEntrypoint(),
			options:    &Options{},
			res:        &ansible.AnsiblePlaybookDriver{},
		},
		{
			desc:       "Testing create ansible driver in build entrypoint with dryrun enabled",
			entrypoint: NewEntrypoint(),
			options: &Options{
				DryRun: true,
			},
			res: &dryrun.DryRunDriver{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			driverFunc, err := test.entrypoint.createAnsibleDriver(test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, driverFunc)
				driver, err := driverFunc()
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}
}

func TestCreateDockerDriver(t *testing.T) {
	errContext := "(entrypoint::build::createDockerDriver)"

	tests := []struct {
		desc        string
		entrypoint  *Entrypoint
		credentials repository.CredentialsFactorier
		options     *Options
		res         repository.BuildDriverer
		err         error
	}{
		{
			desc:        "Testing error creating docker driver in build entrypoint when credentials are empty",
			entrypoint:  NewEntrypoint(),
			credentials: nil,
			err:         errors.New(errContext, "Docker driver requires a credentials store in build entrypoint"),
		},
		{
			desc:        "Testing error creating docker driver in build entrypoint when options are empty",
			entrypoint:  NewEntrypoint(),
			credentials: credentialsfactory.NewMockCredentialsFactory(),
			options:     nil,
			err:         errors.New(errContext, "Build entrypoint options are required to create docker driver"),
		},
		{
			desc:        "Testing create docker driver in build entrypoint",
			entrypoint:  NewEntrypoint(),
			credentials: credentialsfactory.NewMockCredentialsFactory(),
			options:     &Options{},
			res:         &docker.DockerDriver{},
			err:         &errors.Error{},
		},
		{
			desc:        "Testing create docker driver in build entrypoint with dryrun enabled",
			entrypoint:  NewEntrypoint(),
			credentials: credentialsfactory.NewMockCredentialsFactory(),
			options: &Options{
				DryRun: true,
			},
			res: &dryrun.DryRunDriver{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			driverFunc, err := test.entrypoint.createDockerDriver(test.credentials, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NotNil(t, driverFunc)
				driver, err := driverFunc()
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}
}

func TestCreateDispatcher(t *testing.T) {
	desc := "Testing create dispatcher in build entrypoint"

	t.Run(desc, func(t *testing.T) {
		e := NewEntrypoint()
		options := &Options{
			Concurrency: 5,
		}

		dispatch, err := e.createDispatcher(options)

		assert.Nil(t, err)
		assert.NotNil(t, dispatch)
		assert.NotNil(t, dispatch.WorkerPool)
		assert.Equal(t, dispatch.NumWorkers, 5)
	})
}

func TestCreatePlanFactory(t *testing.T) {
	desc := "Testing create build plan factory in build entrypoint"

	t.Run(desc, func(t *testing.T) {
		e := NewEntrypoint()
		options := &Options{}

		imageStore := images.NewStore(nil)
		planFactory, err := e.createPlanFactory(imageStore, options)

		assert.Nil(t, err)
		assert.NotNil(t, planFactory)
		assert.IsType(t, plan.NewPlanFactory(imageStore), planFactory)
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
			desc:       "Testinc create docker reference name on build entrypoint",
			entrypoint: NewEntrypoint(),
			options: &Options{
				UseDockerNormalizedName: true,
			},
			res: dockerreferencename.NewDockerNormalizedReferenceName(),
		},
		{
			desc:       "Testinc create default reference name on build entrypoint",
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
