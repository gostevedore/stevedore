package build

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	imagesconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images"
	imagesgraphtemplate "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/ansible"
	defaultdriver "github.com/gostevedore/stevedore/internal/infrastructure/driver/default"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/dryrun"
	"github.com/gostevedore/stevedore/internal/infrastructure/driver/factory"
	"github.com/gostevedore/stevedore/internal/infrastructure/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/now"
	"github.com/gostevedore/stevedore/internal/infrastructure/plan"
	"github.com/gostevedore/stevedore/internal/infrastructure/render"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/scheduler/job"
	"github.com/gostevedore/stevedore/internal/infrastructure/semver"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
	credentialsstorelocal "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/local"
	credentialsstoremock "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/mock"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {

}

func TestPrepareEntrypointOptions(t *testing.T) {
	errContext := "(Entrypoint::prepareEntrypointOptions)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		conf       *configuration.Configuration
		options    *Options
		res        *Options
		err        error
	}{
		{
			desc:       "Testing error when configuration is not provided",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To prepare entrypoint options, configuration is required"),
		},
		{
			desc:       "Testing error when options are not provided",
			entrypoint: &Entrypoint{},
			conf:       &configuration.Configuration{},
			err:        errors.New(errContext, "To prepare entrypoint options, entrypoint options are required"),
		},
		{
			desc:       "Testing prepare entrypoint options",
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
			desc:       "Testing prepare entrypoint options using configuration concurrency",
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

	errContext := "(Entrypoint::prepareImageName)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		args       []string
		res        string
		err        error
	}{
		{
			desc:       "Testing error when no args is nil",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To execute the build entrypoint, arguments are required"),
		},
		{
			desc:       "Testing error when no args are provided",
			entrypoint: &Entrypoint{},
			args:       []string{},
			err:        errors.New(errContext, "To execute the build entrypoint, arguments are required"),
		},
		{
			desc:       "Testing prepare image name",
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
	errContext := "(Entrypoint::prepareHandlerOptions)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		conf       *configuration.Configuration
		options    *handler.Options
		res        *handler.Options
		err        error
	}{
		{
			desc:       "Testing error when configuration is not provided",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To prepare handler options, configuration is required"),
		},
		{
			desc:       "Testing error when options are not provided",
			entrypoint: &Entrypoint{},
			conf:       &configuration.Configuration{},
			err:        errors.New(errContext, "To prepare handler options, handler options are required"),
		},
		{
			desc:       "Testing prepare handler options",
			entrypoint: &Entrypoint{},
			conf:       &configuration.Configuration{},
			options: &handler.Options{
				AnsibleConnectionLocal:           true,
				AnsibleIntermediateContainerName: "ansible-intermediate-container",
				AnsibleInventoryPath:             "ansible-inventory-path",
				AnsibleLimit:                     "ansible-limit",
				BuildOnCascade:                   true,
				CascadeDepth:                     3,
				DryRun:                           true,
				EnableSemanticVersionTags:        true,
				ImageFromName:                    "image-from-name",
				ImageFromRegistryHost:            "image-from-registry-host",
				ImageFromRegistryNamespace:       "image-from-registry-namespace",
				ImageFromVersion:                 "image-from-version",
				ImageName:                        "image-name",
				ImageRegistryHost:                "image-registry-host",
				ImageRegistryNamespace:           "image-registry-namespace",
				Labels:                           []string{"label1", "label2"},
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
				DryRun:                           true,
				EnableSemanticVersionTags:        true,
				ImageFromName:                    "image-from-name",
				ImageFromRegistryHost:            "image-from-registry-host",
				ImageFromRegistryNamespace:       "image-from-registry-namespace",
				ImageFromVersion:                 "image-from-version",
				ImageName:                        "image-name",
				ImageRegistryHost:                "image-registry-host",
				ImageRegistryNamespace:           "image-registry-namespace",
				Labels:                           []string{"label1", "label2"},
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
			desc:       "Testing prepare handler options using also configuration options",
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
				DryRun:                           true,
				EnableSemanticVersionTags:        true,
				ImageFromName:                    "image-from-name",
				ImageFromRegistryHost:            "image-from-registry-host",
				ImageFromRegistryNamespace:       "image-from-registry-namespace",
				ImageFromVersion:                 "image-from-version",
				ImageName:                        "image-name",
				ImageRegistryHost:                "image-registry-host",
				ImageRegistryNamespace:           "image-registry-namespace",
				Labels:                           []string{"label1", "label2"},
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
				DryRun:                           true,
				EnableSemanticVersionTags:        true,
				ImageFromName:                    "image-from-name",
				ImageFromRegistryHost:            "image-from-registry-host",
				ImageFromRegistryNamespace:       "image-from-registry-namespace",
				ImageFromVersion:                 "image-from-version",
				ImageName:                        "image-name",
				ImageRegistryHost:                "image-registry-host",
				ImageRegistryNamespace:           "image-registry-namespace",
				Labels:                           []string{"label1", "label2"},
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

func TestCreateCredentialsStore(t *testing.T) {
	errContext := "(Entrypoint::createCredentialsStore)"

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
		desc       string
		entrypoint *Entrypoint
		conf       *configuration.Configuration
		res        repository.CredentialsStorer
		err        error
	}{
		{
			desc:       "Testing error when file system is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create the credentials store, a file system is required"),
		},
		{
			desc: "Testing error when configuration is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			err: errors.New(errContext, "To create the credentials store, configuration is required"),
		},
		{
			desc: "Testing error when credentials path is not defined in configuration",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "\n\tDocker credentials directory is not defined on configuration"),
		},
		{
			desc: "Testing create credentials store",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{
				DockerCredentialsDir: baseDir,
			},
			res: &credentialsstorelocal.LocalStore{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			credentials, err := test.entrypoint.createCredentialsStore(test.conf)
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
	errContext := "(Entrypoint::createBuildersStore)"

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
			desc:       "Testing error when file system is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create a builders store, a file system is required"),
		},
		{
			desc: "Testing error when configuration is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			err: errors.New(errContext, "To create a builders store, configuration is required"),
		},
		{
			desc: "Testing error when builders path is not defined in configuration",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create a builders store, builders path must be provided in configuration"),
		},
		{
			desc: "Testing create builders store",
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
			desc: "Testing create command factory",
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
			desc: "Testing create job factory",
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
			desc: "Testing create semver factory",
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
	errContext := "(Entrypoint::createImageRender)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		now        render.Nower
		res        *render.ImageRender
		err        error
	}{
		{
			desc:       "Testing error when now is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create an image render, a nower is required"),
		},
		{
			desc: "Testing create image render",
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
	errContext := "(Entrypoint::createImagesStore)"

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
			desc:       "Testing error when fs is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create an images store, a filesystem is required"),
		},
		{
			desc: "Testing error when configuration is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			err: errors.New(errContext, "To create an images store, configuration is required"),
		},
		{
			desc: "Testing error when render is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create an images store, image render is required"),
		},
		{
			desc: "Testing error when graph is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf:   &configuration.Configuration{},
			render: &render.ImageRender{},
			err:    errors.New(errContext, "To create an images store, images graph templates storer is required"),
		},
		{
			desc: "Testing error when compatibility is not defined",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf:   &configuration.Configuration{},
			render: &render.ImageRender{},
			graph:  &imagesgraphtemplate.ImagesGraphTemplate{},
			err:    errors.New(errContext, "To create an images store, compatibility is required"),
		},
		{
			desc: "Testing error when images path is not defined in configuration",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
			),
			conf:          &configuration.Configuration{},
			render:        &render.ImageRender{},
			graph:         &imagesgraphtemplate.ImagesGraphTemplate{},
			compatibility: &compatibility.Compatibility{},
			err:           errors.New(errContext, "To create an images store, images path must be provided in configuration"),
		},
		{
			desc: "Testing create images store",
			entrypoint: NewEntrypoint(
				WithFileSystem(testFs),
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

			store, err := test.entrypoint.createImagesStore(test.conf, test.render, test.graph, test.compatibility)
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
	errContext := "(Entrypoint::createImagesGraphTemplatesStorer)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		factory    *graph.GraphTemplateFactory
		res        *imagesgraphtemplate.ImagesGraphTemplate
		err        error
	}{
		{
			desc:       "Testing error when factory is not defined",
			entrypoint: NewEntrypoint(),
			err:        errors.New(errContext, "To create an images graph templates storer, a graph template factory is required"),
		},
		{
			desc:       "Testing create images graph templates storer",
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
			desc:       "Testing create graph template factory",
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

	errContext := "(entrypoint::createBuildDriverFactory)"

	tests := []struct {
		desc        string
		entrypoint  *Entrypoint
		credentials repository.CredentialsStorer
		options     *Options
		err         error
		assertions  func(t *testing.T, driverFactory factory.BuildDriverFactory)
	}{
		{
			desc:        "Testing create build driver factory with empty credentials",
			entrypoint:  NewEntrypoint(),
			credentials: nil,
			options:     nil,
			err:         errors.New(errContext, "Register drivers requires a credentials store"),
		},
		{
			desc:        "Testing create build driver factory with empty options",
			entrypoint:  NewEntrypoint(),
			credentials: credentialsstoremock.NewMockStore(),
			options:     nil,
			err:         errors.New(errContext, "Register drivers requires options"),
		},
		{
			desc:        "Testing create build driver factory with nil writer",
			entrypoint:  NewEntrypoint(),
			credentials: credentialsstoremock.NewMockStore(),
			options:     &Options{},
			err:         errors.New(errContext, "Register drivers requires a writer"),
		},
		{
			desc:        "Testing create build driver factory",
			entrypoint:  NewEntrypoint(WithWriter(ioutil.Discard)),
			credentials: credentialsstoremock.NewMockStore(),
			options:     &Options{},
			err:         &errors.Error{},
			assertions: func(t *testing.T, f factory.BuildDriverFactory) {
				dDocker, eDocker := f.Get("docker")
				assert.Nil(t, eDocker)
				assert.NotNil(t, dDocker)
				assert.IsType(t, &docker.DockerDriver{}, dDocker)

				dAnsible, eAnsible := f.Get("ansible-playbook")
				assert.Nil(t, eAnsible)
				assert.NotNil(t, dAnsible)
				assert.IsType(t, &ansible.AnsiblePlaybookDriver{}, dAnsible)

				dDefault, eDefault := f.Get("default")
				assert.Nil(t, eDefault)
				assert.NotNil(t, dDefault)
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
			desc:       "Testing create dry-run driver",
			entrypoint: NewEntrypoint(),
			res:        &dryrun.DryRunDriver{},
		},
	}

	for _, test := range tests {
		t.Run(desc, func(t *testing.T) {

			driver, err := test.entrypoint.createDryRunDriver()

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}

}

func TestCreateDefaultDriver(t *testing.T) {

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		res        repository.BuildDriverer
		err        error
	}{
		{
			desc:       "Testing create default driver",
			entrypoint: NewEntrypoint(),
			res:        &defaultdriver.DefaultDriver{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			driver, err := test.entrypoint.createDefaultDriver()

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}
}

func TestCreateAnsibleDriver(t *testing.T) {

	errContext := "(entrypoint::createAnsibleDriver)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		options    *Options
		res        repository.BuildDriverer
		err        error
	}{
		{
			desc:       "Testing error when creating ansible driver with nil options",
			entrypoint: NewEntrypoint(),
			options:    nil,
			err:        errors.New(errContext, "Entrypoint options are required to create ansible driver"),
		},
		{
			desc:       "Testing create ansible driver",
			entrypoint: NewEntrypoint(),
			options:    &Options{},
			res:        &ansible.AnsiblePlaybookDriver{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			driver, err := test.entrypoint.createAnsibleDriver(test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}
}

func TestCreateDockerDriver(t *testing.T) {
	errContext := "(entrypoint::createDockerDriver)"

	tests := []struct {
		desc        string
		entrypoint  *Entrypoint
		credentials repository.CredentialsStorer
		options     *Options
		res         repository.BuildDriverer
		err         error
	}{
		{
			desc:        "Testing error when creating docker driver with empty credentials",
			entrypoint:  NewEntrypoint(),
			credentials: nil,
			err:         errors.New(errContext, "Docker driver requires a credentials store"),
		},
		{
			desc:        "Testing error when creating docker driver with empty options",
			entrypoint:  NewEntrypoint(),
			credentials: credentialsstoremock.NewMockStore(),
			options:     nil,
			err:         errors.New(errContext, "Entrypoint options are required to create docker driver"),
		},
		{
			desc:        "Testing create docker driver",
			entrypoint:  NewEntrypoint(),
			credentials: credentialsstoremock.NewMockStore(),
			options:     &Options{},
			res:         &docker.DockerDriver{},
			err:         &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			driver, err := test.entrypoint.createDockerDriver(test.credentials, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}

		})
	}
}

func TestCreateDispatcher(t *testing.T) {
	desc := "Testing create dispatcher"

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
	desc := "Testing create build plan factory"

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
