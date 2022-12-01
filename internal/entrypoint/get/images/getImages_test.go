package images

import (
	"io"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	imagesconfiguration "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images"
	imagesgraphtemplate "github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/now"
	plainoutput "github.com/gostevedore/stevedore/internal/infrastructure/output/images/plain"
	treeoutput "github.com/gostevedore/stevedore/internal/infrastructure/output/images/tree"
	defaultreferencename "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/default"
	dockerreferencename "github.com/gostevedore/stevedore/internal/infrastructure/reference/image/docker"
	"github.com/gostevedore/stevedore/internal/infrastructure/render"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// func TestPrepareEntrypointOptions(t *testing.T) {
// 	tests := []struct {
// 		desc              string
// 		entrypoint        *GetImagesEntrypoint
// 		entrypointOptions *Options
// 		handlerOptions    *handler.Options
// 		args              []string
// 		conf              *configuration.Configuration
// 		err               error
// 	}{}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			err := test.entrypoint.Execute(context.TODO(), test.args, test.conf, test.entrypointOptions, test.handlerOptions)
// 			if err != nil {
// 				assert.Equal(t, test.err, err)
// 			}
// 		})
// 	}
// 	assert.True(t, false)
// }

// func TestPrepareHandlerOptions(t *testing.T) {
// 	tests := []struct {
// 		desc              string
// 		entrypoint        *GetImagesEntrypoint
// 		entrypointOptions *Options
// 		handlerOptions    *handler.Options
// 		args              []string
// 		conf              *configuration.Configuration
// 		err               error
// 	}{}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			err := test.entrypoint.Execute(context.TODO(), test.args, test.conf, test.entrypointOptions, test.handlerOptions)
// 			if err != nil {
// 				assert.Equal(t, test.err, err)
// 			}
// 		})
// 	}
// 	assert.True(t, false)
// }

func TestCreateImageRender(t *testing.T) {
	errContext := "(get::images::entrypoint::createImageRender)"

	tests := []struct {
		desc       string
		entrypoint *GetImagesEntrypoint
		now        render.Nower
		res        *render.ImageRender
		err        error
	}{
		{
			desc:       "Testing error creating image render in get images entrypoint when now is not defined",
			entrypoint: NewGetImagesEntrypoint(),
			err:        errors.New(errContext, "To create an image render in get images entrypoint, a nower is required"),
		},
		{
			desc: "Testing create image render in get images entrypoint",
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
	errContext := "(get::images::entrypoint::createImagesStore)"

	baseDir := "/images"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	tests := []struct {
		desc          string
		entrypoint    *GetImagesEntrypoint
		conf          *configuration.Configuration
		render        repository.Renderer
		graph         imagesconfiguration.ImagesGraphTemplatesStorer
		compatibility Compatibilitier
		res           *images.Store
		err           error
	}{
		{
			desc:       "Testing error creating images store in get images entrypoint when fs is not defined",
			entrypoint: NewGetImagesEntrypoint(),
			err:        errors.New(errContext, "To create an images store in get images entrypoint, a filesystem is required"),
		},
		{
			desc: "Testing error creating images store in get images entrypoint when configuration is not defined",
			entrypoint: NewGetImagesEntrypoint(
				WithFileSystem(testFs),
			),
			err: errors.New(errContext, "To create an images store in get images entrypoint, configuration is required"),
		},
		{
			desc: "Testing error creating images store in get images entrypoint when render is not defined",
			entrypoint: NewGetImagesEntrypoint(
				WithFileSystem(testFs),
			),
			conf: &configuration.Configuration{},
			err:  errors.New(errContext, "To create an images store in get images entrypoint, image render is required"),
		},
		{
			desc: "Testing error creating images store in get images entrypoint when graph is not defined",
			entrypoint: NewGetImagesEntrypoint(
				WithFileSystem(testFs),
			),
			conf:   &configuration.Configuration{},
			render: &render.ImageRender{},
			err:    errors.New(errContext, "To create an images store in get images entrypoint, images graph templates storer is required"),
		},
		{
			desc: "Testing error creating images store in get images entrypoint when compatibility is not defined",
			entrypoint: NewGetImagesEntrypoint(
				WithFileSystem(testFs),
			),
			conf:   &configuration.Configuration{},
			render: &render.ImageRender{},
			graph:  &imagesgraphtemplate.ImagesGraphTemplate{},
			err:    errors.New(errContext, "To create an images store in get images entrypoint, compatibility is required"),
		},
		{
			desc: "Testing error creating images store in get images entrypoint when images path is not defined in configuration",
			entrypoint: NewGetImagesEntrypoint(
				WithFileSystem(testFs),
				WithCompatibility(compatibility.NewMockCompatibility()),
			),
			conf:          &configuration.Configuration{},
			render:        &render.ImageRender{},
			graph:         &imagesgraphtemplate.ImagesGraphTemplate{},
			compatibility: &compatibility.Compatibility{},
			err:           errors.New(errContext, "To create an images store in get images entrypoint, images path must be provided in configuration"),
		},
		{
			desc: "Testing create images store",
			entrypoint: NewGetImagesEntrypoint(
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
	errContext := "(get::images::entrypoint::createImagesGraphTemplatesStorer)"

	tests := []struct {
		desc       string
		entrypoint *GetImagesEntrypoint
		factory    *graph.GraphTemplateFactory
		res        *imagesgraphtemplate.ImagesGraphTemplate
		err        error
	}{
		{
			desc:       "Testing error creating images graph templates store in get images entrypoint when factory is not defined",
			entrypoint: NewGetImagesEntrypoint(),
			err:        errors.New(errContext, "To create an images graph templates storer in get images entrypoint, a graph template factory is required"),
		},
		{
			desc:       "Testing create images graph templates storer in get images entrypoint",
			entrypoint: NewGetImagesEntrypoint(),
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
		entrypoint *GetImagesEntrypoint
		res        *graph.GraphTemplateFactory
		err        error
	}{
		{
			desc:       "Testing create graph template factory in get images entrypoint",
			entrypoint: NewGetImagesEntrypoint(),
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

func TestCreateOutput(t *testing.T) {
	tests := []struct {
		desc       string
		entrypoint *GetImagesEntrypoint
		options    *Options
		res        repository.ImagesOutputter
		err        error
	}{
		{
			desc: "Testing create get images entrypoint default (plain text) output",
			entrypoint: NewGetImagesEntrypoint(
				WithWriter(io.Discard),
			),
			options: &Options{},
			res:     plainoutput.NewPlainOutput(),
			err:     &errors.Error{},
		},
		{
			desc:       "Testing error on get images entrypoint when default (plain text) output is not provided by a writer",
			entrypoint: NewGetImagesEntrypoint(),
			options:    &Options{},
			err: errors.New(
				"(get::images::entrypoint::createtOutput)", "",
				errors.New(
					"(get::images::entrypoint::createPlainTextOutput)",
					"Get images entrypoint requires a writer to create the plain text output"),
			),
		},
		{
			desc: "Testing create get images entrypoint tree output",
			entrypoint: NewGetImagesEntrypoint(
				WithWriter(io.Discard),
			),
			options: &Options{
				Tree: true,
			},
			res: treeoutput.NewTreeOutput(),
			err: &errors.Error{},
		},
		{
			desc:       "Testing error on get images entrypoint when tree output is not provided by a writer",
			entrypoint: NewGetImagesEntrypoint(),
			options: &Options{
				Tree: true,
			},
			err: errors.New(
				"(get::images::entrypoint::createtOutput)", "",
				errors.New(
					"(get::images::entrypoint::createTreeOutput)",
					"Get images entrypoint requires a writer to create the tree output"),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			res, err := test.entrypoint.createtOutput(test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.IsType(t, test.res, res)
			}
		})
	}
}

func TestCreateReferenceName(t *testing.T) {
	tests := []struct {
		desc       string
		entrypoint *GetImagesEntrypoint
		options    *Options
		res        repository.ImageReferenceNamer
		err        error
	}{
		{
			desc:       "Testinc create docker reference name on build entrypoint",
			entrypoint: NewGetImagesEntrypoint(),
			options: &Options{
				UseDockerNormalizedName: true,
			},
			res: dockerreferencename.NewDockerNormalizedReferenceName(),
		},
		{
			desc:       "Testinc create default reference name on build entrypoint",
			entrypoint: NewGetImagesEntrypoint(),
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
