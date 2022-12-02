package images

import (
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	domainimage "github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/images/image"
	imagesgraph "github.com/gostevedore/stevedore/internal/infrastructure/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/render"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCheckCompatibility(t *testing.T) {

	tests := []struct {
		desc              string
		tree              *ImagesConfiguration
		prepareAssertFunc func(*ImagesConfiguration)
	}{
		{
			desc: "Testing check compatibility with deprecated configuration",
			tree: &ImagesConfiguration{
				DEPRECATEDImagesTree: map[string]map[string]*image.Image{
					"image": {
						"version": &image.Image{
							Name:    "image",
							Version: "version",
						},
					},
				},
				graph:         graph.NewMockImagesGraphTemplate(),
				compatibility: compatibility.NewMockCompatibility(),
			},
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "image", "version", &image.Image{
					Name:    "image",
					Version: "version",
				}).Return(nil)
				tree.compatibility.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'images_tree' is deprecated and will be removed on v0.12.0, please use 'images' instead"}).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.tree)
			}

			test.tree.CheckCompatibility()
			test.tree.compatibility.(*compatibility.MockCompatibility).AssertExpectations(t)
		})
	}
}

func TestLoadImagesToStore(t *testing.T) {

	var err error

	baseDir := "/imagestree"
	baseErrorDir := "/imagestree_error"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "images.yaml"), []byte(`
images:
  parent1:
    parent1_version:
      registry: registry.test
      namespace: namespace
      version: "v{{ .Version }}"
      builder: builder
      persistent_labels:
        plabel: plabelvalue
      persistent_vars:
        pvar: pvarvalue
  child:
    version:
      registry: registry.test
      namespace: namespace
      name: child
      version: "{{ .Parent.Version }}"
      builder: builder
      parents:
        parent1:
        - parent1_version
`), 0644)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, filepath.Join(baseErrorDir, "tab_error_file.yaml"), []byte(`
images:
image:
  version:
	registry: registry.test
	namespace: namespace
`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		path              string
		err               error
		images            *ImagesConfiguration
		prepareAssertFunc func(*ImagesConfiguration)
		assertFunc        func(*testing.T, *ImagesConfiguration)
	}{
		{
			desc: "Testing load images to store",
			path: baseDir,
			images: NewImagesConfiguration(
				testFs,
				graph.NewImagesGraphTemplate(
					imagesgraph.NewGraphTemplateFactory(false),
				),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(i *ImagesConfiguration) {

				// Create parent1:parent1_version
				i.render.(*render.MockImageRender).On("Render", "parent1", "parent1_version",
					&domainimage.Image{
						RegistryHost:      "registry.test",
						RegistryNamespace: "namespace",
						Name:              "parent1",
						Version:           "v{{ .Version }}",
						Builder:           "builder",
						PersistentLabels: map[string]string{
							"plabel": "plabelvalue",
						},
						PersistentVars: map[string]interface{}{
							"pvar": "pvarvalue",
						},
					},
				).Return(
					&domainimage.Image{
						RegistryHost:      "registry.test",
						RegistryNamespace: "namespace",
						Name:              "parent1",
						Version:           "vparent1_version",
						Builder:           "builder",
						PersistentLabels: map[string]string{
							"plabel": "plabelvalue",
						},
						PersistentVars: map[string]interface{}{
							"pvar": "pvarvalue",
						},
					}, nil)

				i.store.(*images.MockStore).On("Store", "parent1", "parent1_version",
					&domainimage.Image{
						RegistryHost:      "registry.test",
						RegistryNamespace: "namespace",
						Name:              "parent1",
						Version:           "vparent1_version",
						Builder:           "builder",
						PersistentLabels: map[string]string{
							"plabel": "plabelvalue",
						},
						PersistentVars: map[string]interface{}{
							"pvar": "pvarvalue",
						},
					},
				).Return(nil)

				// Create child:version
				i.store.(*images.MockStore).On("Find", "parent1", "parent1_version").Return([]*domainimage.Image{
					{
						RegistryHost:      "registry.test",
						RegistryNamespace: "namespace",
						Name:              "parent1",
						Version:           "vparent1_version",
						Builder:           "builder",
						PersistentLabels: map[string]string{
							"plabel": "plabelvalue",
						},
						PersistentVars: map[string]interface{}{
							"pvar": "pvarvalue",
						},
					},
				}, nil)

				i.render.(*render.MockImageRender).On("Render", "child", "version",
					&domainimage.Image{
						RegistryHost:      "registry.test",
						RegistryNamespace: "namespace",
						Name:              "child",
						Version:           "{{ .Parent.Version }}",
						Builder:           "builder",
						Children:          []*domainimage.Image{},
						Labels:            map[string]string{},
						PersistentLabels: map[string]string{
							"plabel": "plabelvalue",
						},
						PersistentVars: map[string]interface{}{
							"pvar": "pvarvalue",
						},
						Tags: []string{},
						Vars: map[string]interface{}{},
						Parent: &domainimage.Image{
							RegistryHost:      "registry.test",
							RegistryNamespace: "namespace",
							Name:              "parent1",
							Version:           "vparent1_version",
							Builder:           "builder",
							PersistentLabels: map[string]string{
								"plabel": "plabelvalue",
							},
							PersistentVars: map[string]interface{}{
								"pvar": "pvarvalue",
							},
						},
					},
				).Return(
					&domainimage.Image{
						RegistryHost:      "registry.test",
						RegistryNamespace: "namespace",
						Name:              "child",
						Version:           "vparent_version",
						Builder:           "builder",
						Children:          []*domainimage.Image{},
						Labels:            map[string]string{},
						PersistentLabels: map[string]string{
							"plabel": "plabelvalue",
						},
						PersistentVars: map[string]interface{}{
							"pvar": "pvarvalue",
						},
						Tags: []string{},
						Vars: map[string]interface{}{},
						Parent: &domainimage.Image{
							RegistryHost:      "registry.test",
							RegistryNamespace: "namespace",
							Name:              "parent1",
							Version:           "v{{ .Version }}",
							Builder:           "builder",
							PersistentLabels: map[string]string{
								"plabel": "plabelvalue",
							},
							PersistentVars: map[string]interface{}{
								"pvar": "pvarvalue",
							},
						},
					}, nil)

				i.store.(*images.MockStore).On("Store", "child", "version",
					&domainimage.Image{
						RegistryHost:      "registry.test",
						RegistryNamespace: "namespace",
						Name:              "child",
						Version:           "vparent_version",
						Builder:           "builder",
						Children:          []*domainimage.Image{},
						Labels:            map[string]string{},
						PersistentLabels: map[string]string{
							"plabel": "plabelvalue",
						},
						PersistentVars: map[string]interface{}{
							"pvar": "pvarvalue",
						},
						Tags: []string{},
						Vars: map[string]interface{}{},
						Parent: &domainimage.Image{
							RegistryHost:      "registry.test",
							RegistryNamespace: "namespace",
							Name:              "parent1",
							Version:           "v{{ .Version }}",
							Builder:           "builder",
							PersistentLabels: map[string]string{
								"plabel": "plabelvalue",
							},
							PersistentVars: map[string]interface{}{
								"pvar": "pvarvalue",
							},
						},
					},
				).Return(nil)
			},
			assertFunc: func(t *testing.T, i *ImagesConfiguration) {
				i.store.(*images.MockStore).AssertExpectations(t)
				i.render.(*render.MockImageRender).AssertExpectations(t)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.images)
			}

			err := test.images.LoadImagesToStore(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.images)
			}
		})
	}
}

func TestLoadImagesConfigurationFromFile(t *testing.T) {
	var err error

	errContext := "(tree::LoadImagesConfigurationFromFile)"
	_ = errContext

	baseDir := "/imagestree"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "single_image.yaml"), []byte(`
images:
  image:
    version:
      registry: registry
      namespace: namespace
      name: image
      version: version
      builder: builder
      children:
        child1:
          - child1.1
      parents:
        parent1:
          - parent1.1
      tags:
        - tag1
      vars:
        var1: value1
      persistent_vars:
        pvar1: pvalue1
      labels:
        label1: value1
`), 0644)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "deprecated_definition.yaml"), []byte(`
images_tree:
  deprecated_image:
    deprecated_version:
      registry: registry
      namespace: namespace
      name: image
      version: version
      builder: builder
      children:
        child1:
          - child1.1
      parents:
        parent1:
          - parent1.1
      tags:
        - tag1
      vars:
        var1: value1
      persistent_vars:
        pvar1: pvalue1
      labels:
        label1: value1
`), 0644)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "tab_error_file.yaml"), []byte(`
images:
image:
  version:
	registry: registry
	namespace: namespace
`), 0644)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "multiple_images.yaml"), []byte(`
images:
  parent2:
    parent2_version:
      registry: registry.test
      namespace: namespace
      name: parent2
      version: parent2_version
      builder: builder
      children:
        other_child:
        - other_child_version
  child:
    version:
      registry: registry.test
      namespace: namespace
      name: child
      version: version
      builder: builder
      parents:
        parent1:
        - parent1_version
        parent2:
        - parent2_version
  other_child:
    other_child_version:
      registry: registry.test
      namespace: namespace
      name: other_child
      version: other_child_version
      builder: builder
  parent1:
    parent1_version:
      registry: registry.test
      namespace: namespace
      name: parent1
      version: parent1_version
      builder: builder
`), 0644)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "multiple_parents.yaml"), []byte(`
images:
  parent1:
    parent1_version:
      registry: registry.test
      namespace: namespace
      name: parent1
      version: parent1_version
      builder: builder

  parent2:
    parent2_version:
      registry: registry.test
      namespace: namespace
      name: parent2
      version: parent2_version
      builder: builder

  child:
    child_version:
      registry: registry.test
      namespace: namespace
      name: child
      version: "{{ .Parent.Version }}"
      builder: builder
      parents:
        parent1:
          - parent1_version
        parent2:
          - parent2_version
`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		path              string
		tree              *ImagesConfiguration
		prepareAssertFunc func(*ImagesConfiguration)
		err               error
	}{
		{
			desc: "Testing error on load images tree from file",
			path: filepath.Join(baseDir, "tab_error_file.yaml"),
			tree: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			err: errors.New(errContext, "Error loading images tree from file '/imagestree/tab_error_file.yaml'\nfound:\n\nimages:\nimage:\n  version:\n\tregistry: registry\n\tnamespace: namespace\n\n\tyaml: line 5: found character that cannot start any token"),
		},
		{
			desc: "Testing load images tree from file",
			path: filepath.Join(baseDir, "single_image.yaml"),
			tree: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "image", "version", &image.Image{
					Name:              "image",
					Version:           "version",
					RegistryHost:      "registry",
					RegistryNamespace: "namespace",
					Builder:           "builder",
					Children: map[string][]string{
						"child1": {"child1.1"},
					},
					Parents: map[string][]string{
						"parent1": {"parent1.1"},
					},
					Tags: []string{"tag1"},
					Vars: map[string]interface{}{
						"var1": "value1",
					},
					PersistentVars: map[string]interface{}{
						"pvar1": "pvalue1",
					},
					Labels: map[string]string{
						"label1": "value1",
					},
				}).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing load images tree from file with multiple images an relationships",
			path: filepath.Join(baseDir, "multiple_images.yaml"),
			tree: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "parent1", "parent1_version", &image.Image{
					Name:              "parent1",
					Version:           "parent1_version",
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Builder:           "builder",
				}).Return(nil)
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "parent2", "parent2_version", &image.Image{
					Name:              "parent2",
					Version:           "parent2_version",
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Builder:           "builder",
					Children: map[string][]string{
						"other_child": {"other_child_version"},
					},
				}).Return(nil)
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "child", "version", &image.Image{
					Name:              "child",
					Version:           "version",
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Builder:           "builder",
					Parents: map[string][]string{
						"parent1": {"parent1_version"},
						"parent2": {"parent2_version"},
					},
				}).Return(nil)
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "other_child", "other_child_version", &image.Image{
					Name:              "other_child",
					Version:           "other_child_version",
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Builder:           "builder",
				}).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing load images with one child and multiple parents",
			path: filepath.Join(baseDir, "multiple_parents.yaml"),
			tree: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "parent1", "parent1_version", &image.Image{
					Name:              "parent1",
					Version:           "parent1_version",
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Builder:           "builder",
				}).Return(nil)
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "parent2", "parent2_version", &image.Image{
					Name:              "parent2",
					Version:           "parent2_version",
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Builder:           "builder",
				}).Return(nil)
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "child", "child_version", &image.Image{
					Name:              "child",
					Version:           "{{ .Parent.Version }}",
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Builder:           "builder",
					Parents: map[string][]string{
						"parent1": {"parent1_version"},
						"parent2": {"parent2_version"},
					},
				}).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing load images tree from file with deprecated definition",
			path: filepath.Join(baseDir, "deprecated_definition.yaml"),
			tree: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.compatibility.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"'images_tree' is deprecated and will be removed on v0.12.0, please use 'images' instead"}).Return(nil)
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "deprecated_image", "deprecated_version", &image.Image{
					Name:              "image",
					Version:           "version",
					RegistryHost:      "registry",
					RegistryNamespace: "namespace",
					Builder:           "builder",
					Children: map[string][]string{
						"child1": {"child1.1"},
					},
					Parents: map[string][]string{
						"parent1": {"parent1.1"},
					},
					Tags: []string{"tag1"},
					Vars: map[string]interface{}{
						"var1": "value1",
					},
					PersistentVars: map[string]interface{}{
						"pvar1": "pvalue1",
					},
					Labels: map[string]string{
						"label1": "value1",
					},
				}).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing error when adding image to images graph store",
			path: filepath.Join(baseDir, "single_image.yaml"),
			tree: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "image", "version", &image.Image{
					Name:              "image",
					Version:           "version",
					RegistryHost:      "registry",
					RegistryNamespace: "namespace",
					Builder:           "builder",
					Children: map[string][]string{
						"child1": {"child1.1"},
					},
					Parents: map[string][]string{
						"parent1": {"parent1.1"},
					},
					Tags: []string{"tag1"},
					Vars: map[string]interface{}{
						"var1": "value1",
					},
					PersistentVars: map[string]interface{}{
						"pvar1": "pvalue1",
					},
					Labels: map[string]string{
						"label1": "value1",
					},
				}).Return(
					errors.New(errContext, "Error adding image to images graph store"),
				)
			},
			err: errors.New(errContext, "\n\tError adding image to images graph store"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.tree)
			}

			err := test.tree.LoadImagesConfigurationFromFile(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.tree.graph.(*graph.MockImagesGraphTemplate).AssertExpectations(t)
			}
		})
	}
}

func TestLoadImagesConfigurationFromDir(t *testing.T) {
	var err error
	errContext := "(tree::LoadImagesConfigurationFromDir)"

	baseDir := "/imagestree"
	baseErrorDir := "/imagestree_error"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "file1.yaml"), []byte(`
images:
  image:
    version:
      registry: registry
      namespace: namespace
      name: image
      version: version
      builder: builder
      children:
        child1:
          - child1.1
      parents:
        parent1:
          - parent1.1
      tags:
        - tag1
      vars:
        var1: value1
      persistent_vars:
        pvar1: pvalue1
      labels:
        label1: value1
`), 0644)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "file2.yaml"), []byte(`
images:
  image2:
    version:
      registry: registry
      namespace: namespace
      name: image2
      version: version
      builder: builder
`), 0644)
	if err != nil {
		t.Log(err)
	}
	err = afero.WriteFile(testFs, filepath.Join(baseDir, "empty_image_tree.yaml"), []byte(`
images:
`), 0644)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, filepath.Join(baseErrorDir, "tab_error_file.yaml"), []byte(`
images:
image:
  version:
	registry: registry
	namespace: namespace
`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		path              string
		tree              *ImagesConfiguration
		prepareAssertFunc func(tree *ImagesConfiguration)
		err               error
	}{
		{
			desc: "Testing load images tree from directory",
			path: baseDir,
			tree: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "image", "version", &image.Image{
					Name:              "image",
					Version:           "version",
					RegistryHost:      "registry",
					RegistryNamespace: "namespace",
					Builder:           "builder",
					Children: map[string][]string{
						"child1": {"child1.1"},
					},
					Parents: map[string][]string{
						"parent1": {"parent1.1"},
					},
					Tags: []string{"tag1"},
					Vars: map[string]interface{}{
						"var1": "value1",
					},
					PersistentVars: map[string]interface{}{
						"pvar1": "pvalue1",
					},
					Labels: map[string]string{
						"label1": "value1",
					},
				}).Return(nil)

				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "image2", "version", &image.Image{
					Name:              "image2",
					Version:           "version",
					RegistryHost:      "registry",
					RegistryNamespace: "namespace",
					Builder:           "builder",
				}).Return(nil)

			},
			err: &errors.Error{},
		},
		{
			desc: "Testing error when adding and existing image on images tree",
			path: baseDir,
			tree: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "image", "version", &image.Image{
					Name:              "image",
					Version:           "version",
					RegistryHost:      "registry",
					RegistryNamespace: "namespace",
					Builder:           "builder",
					Children: map[string][]string{
						"child1": {"child1.1"},
					},
					Parents: map[string][]string{
						"parent1": {"parent1.1"},
					},
					Tags: []string{"tag1"},
					Vars: map[string]interface{}{
						"var1": "value1",
					},
					PersistentVars: map[string]interface{}{
						"pvar1": "pvalue1",
					},
					Labels: map[string]string{
						"label1": "value1",
					},
				}).Return(nil)

				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "image2", "version", &image.Image{
					Name:              "image2",
					Version:           "version",
					RegistryHost:      "registry",
					RegistryNamespace: "namespace",
					Builder:           "builder",
				}).Return(
					errors.New(errContext, "Error adding image2"),
				)

			},
			err: errors.New(errContext, "\n\tError adding image2\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.tree)
			}

			err := test.tree.LoadImagesConfigurationFromDir(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.tree.graph.(*graph.MockImagesGraphTemplate).AssertExpectations(t)
			}
		})
	}
}

func TestLoadImagesConfiguration(t *testing.T) {
	var err error
	errContext := "(tree::LoadImagesConfiguration)"
	_ = errContext

	baseDir := "/imagestree"
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(baseDir, 0755)

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "file1.yaml"), []byte(`
images:
  image:
    version:
      registry: registry
      namespace: namespace
      name: image
      version: version
      builder: builder
`), 0644)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		path              string
		images            *ImagesConfiguration
		prepareAssertFunc func(tree *ImagesConfiguration)
		err               error
	}{
		{
			desc: "Testing load images tree from file",
			path: filepath.Join(baseDir, "file1.yaml"),
			images: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "image", "version", &image.Image{
					Name:              "image",
					Version:           "version",
					RegistryHost:      "registry",
					RegistryNamespace: "namespace",
					Builder:           "builder",
				}).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing load images tree from dir",
			path: baseDir,
			images: NewImagesConfiguration(
				testFs,
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "image", "version", &image.Image{
					Name:              "image",
					Version:           "version",
					RegistryHost:      "registry",
					RegistryNamespace: "namespace",
					Builder:           "builder",
				}).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.images)
			}

			err := test.images.LoadImagesConfiguration(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.images.graph.(*graph.MockImagesGraphTemplate).AssertExpectations(t)
			}
		})
	}
}

func TestIsAValidName(t *testing.T) {
	tests := []struct {
		desc string
		name string
		res  bool
	}{
		{
			desc: "Testing valid name",
			name: "valid",
			res:  true,
		},
		{
			desc: "Testing invalid name",
			name: "in:valid",
			res:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			assert.Equal(t, test.res, isAValidName(test.name))
		})
	}
}

func TestIsAValidVersion(t *testing.T) {
	tests := []struct {
		desc    string
		version string
		res     bool
	}{
		{
			desc:    "Testing valid version",
			version: "valid",
			res:     true,
		},
		{
			desc:    "Testing invalid version",
			version: "in:valid",
			res:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			assert.Equal(t, test.res, isAValidVersion(test.version))
		})
	}
}

func TestRenderImage(t *testing.T) {

	errContext := "(images::renderImage)"

	tests := []struct {
		desc              string
		name              string
		version           string
		image             *domainimage.Image
		err               error
		images            *ImagesConfiguration
		res               *domainimage.Image
		prepareAssertFunc func(*ImagesConfiguration)
		assertFunc        func(*testing.T, *ImagesConfiguration)
	}{
		{
			desc:    "Testing render image on image configuration when version is a wildcard symbol",
			name:    "image",
			version: domainimage.ImageWildcardVersionSymbol,
			image: &domainimage.Image{
				Name:    "image",
				Version: "{{ .Version }}",
			},
			res: &domainimage.Image{
				Name:    "image",
				Version: "{{ .Version }}",
			},
		},
		{
			desc: "Testing error rendering an image on image configuration when name is not provided",
			err:  errors.New(errContext, "Image name must be provided to render an image"),
		},
		{
			desc: "Testing error rendering an image on image configuration when renderer is not provided",
			name: "image",
			err:  errors.New(errContext, "Image version must be provided to render an image"),
		},
		{
			desc:    "Testing error rendering an image on image configuration when renderer is not provided",
			name:    "image",
			version: "version",
			images: NewImagesConfiguration(
				afero.NewMemMapFs(),
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				nil,
				compatibility.NewMockCompatibility(),
			),
			err: errors.New(errContext, "Image 'image:version' could not be rendered because renderer must by provided"),
		},
		{
			desc:    "Testing render an image on image configuration",
			name:    "image",
			version: "version",
			images: NewImagesConfiguration(
				afero.NewMemMapFs(),
				graph.NewMockImagesGraphTemplate(),
				images.NewMockStore(),
				render.NewMockImageRender(),
				compatibility.NewMockCompatibility(),
			),
			err: &errors.Error{},
			image: &domainimage.Image{
				Name:    "image",
				Version: "{{ .Version }}",
			},
			prepareAssertFunc: func(i *ImagesConfiguration) {
				i.render.(*render.MockImageRender).On("Render", "image", "version",
					&domainimage.Image{
						Name:    "image",
						Version: "{{ .Version }}",
					},
				).Return(
					&domainimage.Image{
						Name:    "image",
						Version: "version",
					},
					nil,
				)
			},
			assertFunc: func(t *testing.T, i *ImagesConfiguration) {
				i.render.(*render.MockImageRender).AssertExpectations(t)
			},
			res: &domainimage.Image{
				Name:    "image",
				Version: "version",
				// RegistryHost:      "docker.io",
				// RegistryNamespace: "library",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.images)
			}

			res, err := test.images.renderImage(test.name, test.version, test.image)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)

				if test.assertFunc != nil {
					test.assertFunc(t, test.images)
				}
			}
		})
	}
}

func TestPropagatePersistentAttributes(t *testing.T) {
	errContext := "(images::propagatePersistentAttributes)"

	tests := []struct {
		desc   string
		images *ImagesConfiguration
		image  *domainimage.Image
		res    *domainimage.Image
		err    error
	}{
		{
			desc:   "Testing error progatating persistent attributes when image is not provided",
			images: &ImagesConfiguration{},
			err:    errors.New(errContext, "To propagate persistent attributs, an image must be provided"),
		},
		{
			desc:   "Testing propagate persistent attributes to an image",
			images: &ImagesConfiguration{},
			image: &domainimage.Image{
				PersistentVars: map[string]interface{}{
					"pvar1": "pvalue1",
					"pvar3": "pvalue3",
				},
				PersistentLabels: map[string]string{
					"label1": "value1",
					"label3": "value3",
				},
				Parent: &domainimage.Image{
					PersistentVars: map[string]interface{}{
						"pvar1": "parent_pvalue1",
						"pvar2": "parent_pvalue2",
					},
					PersistentLabels: map[string]string{
						"label1": "parent_value1",
						"label2": "parent_value2",
					},
				},
			},
			res: &domainimage.Image{
				PersistentVars: map[string]interface{}{
					"pvar1": "parent_pvalue1",
					"pvar2": "parent_pvalue2",
					"pvar3": "pvalue3",
				},
				PersistentLabels: map[string]string{
					"label1": "parent_value1",
					"label2": "parent_value2",
					"label3": "value3",
				},
				Parent: &domainimage.Image{
					PersistentVars: map[string]interface{}{
						"pvar1": "parent_pvalue1",
						"pvar2": "parent_pvalue2",
					},
					PersistentLabels: map[string]string{
						"label1": "parent_value1",
						"label2": "parent_value2",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			err := test.images.propagatePersistentAttributes(test.image)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, test.image)
			}
		})
	}

}
