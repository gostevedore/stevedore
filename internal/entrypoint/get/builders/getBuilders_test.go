package builders

import (
	"context"
	"io"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	handler "github.com/gostevedore/stevedore/internal/handler/get/builders"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	imagesconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images"
	imagesgraphtemplate "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/now"
	output "github.com/gostevedore/stevedore/internal/infrastructure/output/builders"
	"github.com/gostevedore/stevedore/internal/infrastructure/render"
	buildersstore "github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		desc            string
		entrypoint      *GetBuildersEntrypoint
		args            []string
		conf            *configuration.Configuration
		options         *handler.Options
		prepareMockFunc func()
		err             error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.entrypoint.Execute(context.TODO(), test.args, test.conf, test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			}
		})
	}
}

func TestCreateImageRender(t *testing.T) {
	errContext := "(entrypoint::get::builders::createImageRender)"

	tests := []struct {
		desc       string
		entrypoint *GetBuildersEntrypoint
		now        render.Nower
		res        *render.ImageRender
		err        error
	}{
		{
			desc:       "Testing error creating image render in get builders entrypoint when now is not defined",
			entrypoint: NewGetBuildersEntrypoint(),
			err:        errors.New(errContext, "To create an image render in get builders entrypoint, a nower is required"),
		},
		{
			desc: "Testing create image render in get builders entrypoint",
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
	errContext := "(entrypoint::get::builders::createImagesStore)"

	baseDir := "/images"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	tests := []struct {
		desc          string
		entrypoint    *GetBuildersEntrypoint
		conf          *configuration.Configuration
		render        repository.Renderer
		graph         imagesconfiguration.ImagesGraphTemplatesStorer
		compatibility Compatibilitier
		res           *images.Store
		err           error
	}{
		{
			desc:       "Testing error creating images store in get builders entrypoint when fs is not defined",
			entrypoint: NewGetBuildersEntrypoint(),
			err:        errors.New(errContext, "To create an images store in get builders entrypoint, a filesystem is required"),
		},
		{
			desc: "Testing error creating images store in get builders entrypoint when configuration is not defined",
			entrypoint: NewGetBuildersEntrypoint(
				WithFileSystem(testFs),
			),
			err: errors.New(errContext, "To create an images store in get builders entrypoint, configuration is required"),
		},
		{
			desc: "Testing error creating images store in get builders entrypoint when render is not defined",
			entrypoint: NewGetBuildersEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create an images store in get builders entrypoint, image render is required"),
		},
		{
			desc: "Testing error creating images store in get builders entrypoint when graph is not defined",
			entrypoint: NewGetBuildersEntrypoint(
				WithFileSystem(testFs),
			),
			conf:   &configuration.Configuration{},
			render: &render.ImageRender{},
			err:    errors.New(errContext, "To create an images store in get builders entrypoint, images graph templates storer is required"),
		},
		{
			desc: "Testing error creating images store in get builders entrypoint when compatibility is not defined",
			entrypoint: NewGetBuildersEntrypoint(
				WithFileSystem(testFs),
			),
			conf:   &configuration.Configuration{},
			render: &render.ImageRender{},
			graph:  &imagesgraphtemplate.ImagesGraphTemplate{},
			err:    errors.New(errContext, "To create an images store in get builders entrypoint, compatibility is required"),
		},
		{
			desc: "Testing error creating images store in get builders entrypoint when images path is not defined in configuration",
			entrypoint: NewGetBuildersEntrypoint(
				WithFileSystem(testFs),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf:          &configuration.Configuration{},
			render:        &render.ImageRender{},
			graph:         &imagesgraphtemplate.ImagesGraphTemplate{},
			compatibility: &compatibility.Compatibility{},
			err:           errors.New(errContext, "To create an images store in get builders entrypoint, images path must be provided in configuration"),
		},
		{
			desc: "Testing create images store in get builders entrypoint",
			entrypoint: NewGetBuildersEntrypoint(
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
	errContext := "(entrypoint::get::builders::createImagesGraphTemplatesStorer)"

	tests := []struct {
		desc       string
		entrypoint *GetBuildersEntrypoint
		factory    *graph.GraphTemplateFactory
		res        *imagesgraphtemplate.ImagesGraphTemplate
		err        error
	}{
		{
			desc:       "Testing error creating images graph templates store in get builders entrypoint when factory is not defined",
			entrypoint: NewGetBuildersEntrypoint(),
			err:        errors.New(errContext, "To create an images graph templates storer in get builders entrypoint, a graph template factory is required"),
		},
		{
			desc:       "Testing create images graph templates storer in get builders entrypoint",
			entrypoint: NewGetBuildersEntrypoint(),
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
		entrypoint *GetBuildersEntrypoint
		res        *graph.GraphTemplateFactory
		err        error
	}{
		{
			desc:       "Testing create graph template factory in get builders entrypoint",
			entrypoint: NewGetBuildersEntrypoint(),
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

func TestCreateBuildersStore(t *testing.T) {
	errContext := "(entrypoint::build::createBuildersStore)"

	baseDir := "/builders"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	tests := []struct {
		desc       string
		entrypoint *GetBuildersEntrypoint
		conf       *configuration.Configuration
		res        *buildersstore.Store
		err        error
	}{
		{
			desc:       "Testing error creating builder store in get builders entrypoint when file system is not defined",
			entrypoint: NewGetBuildersEntrypoint(),
			err:        errors.New(errContext, "To create a builders store in build entrypoint, a file system is required"),
		},
		{
			desc: "Testing error creating builder store in get builders entrypoint when configuration is not defined",
			entrypoint: NewGetBuildersEntrypoint(
				WithFileSystem(testFs),
			),
			err: errors.New(errContext, "To create a builders store in build entrypoint, configuration is required"),
		},
		{
			desc: "Testing error creating builder store in get builders entrypoint when builders path is not defined in configuration",
			entrypoint: NewGetBuildersEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create a builders store in build entrypoint, builders path must be provided in configuration"),
		},
		{
			desc: "Testing create builders store in get builders entrypoint",
			entrypoint: NewGetBuildersEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{
				BuildersPath: baseDir,
			},
			res: &buildersstore.Store{},
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

func TestCreateOutput(t *testing.T) {
	tests := []struct {
		desc       string
		entrypoint *GetBuildersEntrypoint
		res        repository.BuildersOutputter
		err        error
	}{
		{
			desc: "Testing create get builders entrypoint default (plain text) output",
			entrypoint: NewGetBuildersEntrypoint(
				WithWriter(io.Discard),
			),
			res: output.NewPlainOutput(),
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			res, err := test.entrypoint.createOutput()
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.IsType(t, test.res, res)
			}
		})
	}
}
