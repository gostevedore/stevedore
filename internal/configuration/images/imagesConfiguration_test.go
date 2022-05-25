package images

import (
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/configuration/images/image"
	domainimage "github.com/gostevedore/stevedore/internal/core/domain/image"
	imagesgraph "github.com/gostevedore/stevedore/internal/images/graph"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	err = afero.WriteFile(testFs, filepath.Join(baseDir, "file1.yaml"), []byte(`
images:
  parent1:
    parent1_version:
      registry: registry.test
      namespace: namespace
      builder: builder
      persistent_labels:
        plabel: plabelvalue
  parent2:
    parent2_version:
      registry: registry.test
      namespace: namespace
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
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(i *ImagesConfiguration) {

				i.store.(*images.MockStore).On("Find", "parent1", "parent1_version").Return(&domainimage.Image{
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Name:              "parent1",
					Version:           "parent1_version",
					Builder:           "builder",
				}, nil)
				i.store.(*images.MockStore).On("Find", "parent2", "parent2_version").Return(&domainimage.Image{
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Name:              "parent2",
					Version:           "parent2_version",
					Builder:           "builder",
				}, nil)

				parent1 := &domainimage.Image{
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Name:              "parent1",
					Version:           "parent1_version",
					Builder:           "builder",
				}
				parent2 := &domainimage.Image{
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Name:              "parent2",
					Version:           "parent2_version",
					Builder:           "builder",
				}
				childParent1 := &domainimage.Image{
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Name:              "child",
					Version:           "version",
					Builder:           "builder",
					Labels:            map[string]string{},
					Tags:              []string{},
					Vars:              map[string]interface{}{},
				}
				childParent2 := &domainimage.Image{
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Name:              "child",
					Version:           "version",
					Builder:           "builder",
					Labels:            map[string]string{},
					Tags:              []string{},
					Vars:              map[string]interface{}{},
				}
				addOtherChild := &domainimage.Image{
					RegistryHost:      "registry.test",
					RegistryNamespace: "namespace",
					Name:              "other_child",
					Version:           "other_child_version",
					Builder:           "builder",
				}

				addOtherChild.Options(
					domainimage.WithParent(parent2),
				)
				// parent2.AddChild(addOtherChild)
				childParent1.Options(
					domainimage.WithParent(parent1),
				)
				// parent1.AddChild(childParent1)
				childParent2.Options(
					domainimage.WithParent(parent2),
				)
				// parent2.AddChild(childParent2)

				i.store.(*images.MockStore).On("Store", "parent1", "parent1_version", mock.AnythingOfType("*image.Image")).Return(nil)
				i.store.(*images.MockStore).On("Store", "parent2", "parent2_version", mock.AnythingOfType("*image.Image")).Return(nil)
				i.store.(*images.MockStore).On("Store", "child", "version", mock.AnythingOfType("*image.Image")).Return(nil)
				i.store.(*images.MockStore).On("Store", "child", "version", mock.AnythingOfType("*image.Image")).Return(nil)
				i.store.(*images.MockStore).On("Store", "other_child", "other_child_version", mock.AnythingOfType("*image.Image")).Return(nil)
			},
			assertFunc: func(t *testing.T, i *ImagesConfiguration) {
				i.store.(*images.MockStore).AssertExpectations(t)
				i.store.(*images.MockStore).AssertNumberOfCalls(t, "Store", 5)
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
    children:
      child:
        - child_version
  parent2:
    parent2_version:
      registry: registry.test
      namespace: namespace
      name: parent2
      version: parent2_version
      builder: builder
      children:
        child:
        - child_version
  child:
    child_version:
      registry: registry.test
      namespace: namespace
      name: child
      version: {{ .Parent.Version }}
      builder: builder
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
				compatibility.NewMockCompatibility(),
			),
			prepareAssertFunc: func(tree *ImagesConfiguration) {
				// tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "parent1", "parent1_version", &image.Image{
				// 	Name:              "parent1",
				// 	Version:           "parent1_version",
				// 	RegistryHost:      "registry.test",
				// 	RegistryNamespace: "namespace",
				// 	Builder:           "builder",
				// }).Return(nil)
				// tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "parent2", "parent2_version", &image.Image{
				// 	Name:              "parent2",
				// 	Version:           "parent2_version",
				// 	RegistryHost:      "registry.test",
				// 	RegistryNamespace: "namespace",
				// 	Builder:           "builder",
				// 	Children: map[string][]string{
				// 		"other_child": {"other_child_version"},
				// 	},
				// }).Return(nil)
				// tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "child", "version", &image.Image{
				// 	Name:              "child",
				// 	Version:           "version",
				// 	RegistryHost:      "registry.test",
				// 	RegistryNamespace: "namespace",
				// 	Builder:           "builder",
				// 	Parents: map[string][]string{
				// 		"parent1": {"parent1_version"},
				// 		"parent2": {"parent2_version"},
				// 	},
				// }).Return(nil)
				// tree.graph.(*graph.MockImagesGraphTemplate).On("AddImage", "other_child", "other_child_version", &image.Image{
				// 	Name:              "other_child",
				// 	Version:           "other_child_version",
				// 	RegistryHost:      "registry.test",
				// 	RegistryNamespace: "namespace",
				// 	Builder:           "builder",
				// }).Return(nil)
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

// LoadImagesConfiguration test
// func TestLoadImagesConfiguration(t *testing.T) {

// 	testBaseDir := "test"

// 	tests := []struct {
// 		desc       string
// 		file       string
// 		err        error
// 		imagesTree *ImagesConfiguration
// 	}{
// 		{
// 			desc:       "Testing an unexistent file",
// 			file:       "nofile",
// 			err:        errors.New("(tree::LoadImagesConfiguration)", "Error loading images tree configuration", errors.New("", "(LoadYAMLFile) Error loading file nofile. open nofile: no such file or directory")),
// 			imagesTree: &ImagesConfiguration{},
// 		},
// 		{
// 			desc: "Testing a simple tree",
// 			file: filepath.Join(testBaseDir, "stevedore_multiple_images.yml"),
// 			err:  nil,
// 			imagesTree: &ImagesConfiguration{
// 				Images: map[string]map[string]*image.Image{
// 					"php-fpm": {
// 						"7.1": &image.Image{
// 							Builder: "mock-builder",
// 							Tags: []string{
// 								"7.1",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "php-fpm",
// 								"source_image_tag": "16.04",
// 							},
// 							Children: map[string][]string{
// 								"php-fpm-dev": {
// 									"7.1",
// 								},
// 							},
// 						},
// 						"7.2": &image.Image{
// 							Builder: "mock-builder",
// 							Tags: []string{
// 								"7.2",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "php-fpm",
// 								"source_image_tag": "16.04",
// 							},
// 							Children: map[string][]string{
// 								"php-fpm-dev": {
// 									"7.2",
// 								},
// 							},
// 						},
// 					},
// 					"php-fpm-dev": {
// 						"7.1": &image.Image{
// 							Builder: "mock-builder",
// 							Tags: []string{
// 								"7.1",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "php-fpm-dev",
// 								"source_image_tag": "16.04",
// 							},
// 						},
// 						"7.2": &image.Image{
// 							Builder: "mock-builder",
// 							Tags: []string{
// 								"7.2",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "php-fpm-dev",
// 								"source_image_tag": "16.04",
// 							},
// 						},
// 					},
// 					"ubuntu": {
// 						"16.04": &image.Image{
// 							Builder: "mock-builder",
// 							Tags: []string{
// 								"16.04",
// 								"xenial",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "ubuntu",
// 								"source_image_tag": "16.04",
// 							},
// 							Children: map[string][]string{
// 								"php-builder": {
// 									"7.1",
// 								},
// 								"php-fpm": {
// 									"7.1",
// 									"7.2",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			desc:       "Testing a simple tree",
// 			file:       filepath.Join(testBaseDir, "stevedore_nil.yml"),
// 			err:        errors.New("(tree::LoadImagesConfiguration)", "Image tree is not defined properly on "+filepath.Join(testBaseDir, "stevedore_nil.yml")),
// 			imagesTree: nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Log(test.desc)

// 		imagesTree, err := LoadImagesConfiguration(test.file)

// 		if err != nil && assert.Error(t, err) {
// 			assert.Equal(t, test.err.Error(), err.Error())
// 		} else {
// 			assert.Equal(t, test.imagesTree, imagesTree, "Unexpected value")
// 		}
// 	}
// }

// func TestGenerateGraph(t *testing.T) {
// 	tests := []struct {
// 		desc       string
// 		file       string
// 		err        error
// 		imagesTree *ImagesConfiguration
// 	}{
// 		{
// 			desc: "Testing a simple tree",
// 			file: "../test/images/simpleImagesConfiguration.yml",
// 			err:  nil,
// 			imagesTree: &ImagesConfiguration{
// 				Images: map[string]map[string]*image.Image{
// 					"php-fpm": {
// 						"7.1": &image.Image{
// 							Builder: "infrastructure",
// 							Tags: []string{
// 								"7.1",
// 							},
// 							PersistentVars: map[string]interface{}{
// 								"php_version": "7.1",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "php-fpm",
// 								"source_image_tag": "16.04",
// 							},
// 							Children: map[string][]string{
// 								"php-fpm-dev": {
// 									"7.1",
// 								},
// 							},
// 						},
// 						"7.2": &image.Image{
// 							Builder: "infrastructure",
// 							Tags: []string{
// 								"7.2",
// 							},
// 							PersistentVars: map[string]interface{}{
// 								"php_version": "7.1",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "php-fpm",
// 								"source_image_tag": "16.04",
// 							},
// 							Children: map[string][]string{
// 								"php-fpm-dev": {
// 									"7.2",
// 								},
// 							},
// 						},
// 					},
// 					"php-fpm-dev": {
// 						"7.1": &image.Image{
// 							Builder: "infrastructure",
// 							Tags: []string{
// 								"7.1",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "php-fpm-dev",
// 								"source_image_tag": "16.04",
// 							},
// 						},
// 						"7.2": &image.Image{
// 							Builder: "infrastructure",
// 							Tags: []string{
// 								"7.2",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "php-fpm-dev",
// 								"source_image_tag": "16.04",
// 							},
// 						},
// 					},
// 					"ubuntu": {
// 						"16.04": &image.Image{
// 							Builder: "infrastructure",
// 							Tags: []string{
// 								"16.04",
// 								"xenial",
// 							},
// 							Vars: map[string]interface{}{
// 								"container_name":   "ubuntu",
// 								"source_image_tag": "16.04",
// 							},
// 							Children: map[string][]string{
// 								"php-builder": {
// 									"7.1",
// 								},
// 								"php-fpm": {
// 									"7.1",
// 									"7.2",
// 								},
// 								"php-cli": {
// 									"7.1",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Log(test.desc)

// 		_, _, err := test.imagesTree.GenerateGraph()
// 		if err != nil {
// 			t.Fatal(err.Error())
// 		}
// 	}

// }

// func TestGenerateTemplateGraph(t *testing.T) {

// 	imagesTree := &ImagesConfiguration{
// 		Images: map[string]map[string]*image.Image{
// 			"php-fpm": {
// 				"7.1": &image.Image{
// 					Builder: "infrastructure",
// 					Tags: []string{
// 						"7.1",
// 					},
// 					Vars: map[string]interface{}{
// 						"container_name":   "php-fpm",
// 						"source_image_tag": "16.04",
// 					},
// 					Children: map[string][]string{
// 						"php-fpm-dev": {
// 							"7.1",
// 						},
// 					},
// 				},
// 				"7.2": &image.Image{
// 					Builder: "infrastructure",
// 					Tags: []string{
// 						"7.2",
// 					},
// 					Vars: map[string]interface{}{
// 						"container_name":   "php-fpm",
// 						"source_image_tag": "16.04",
// 					},
// 					Children: map[string][]string{
// 						"php-fpm-dev": {
// 							"7.2",
// 						},
// 					},
// 				},
// 			},
// 			"php-fpm-dev": {
// 				"7.1": &image.Image{
// 					Builder: "infrastructure",
// 					Tags: []string{
// 						"7.1",
// 					},
// 					Vars: map[string]interface{}{
// 						"container_name":   "php-fpm-dev",
// 						"source_image_tag": "16.04",
// 					},
// 				},
// 				"7.2": &image.Image{
// 					Builder: "infrastructure",
// 					Tags: []string{
// 						"7.2",
// 					},
// 					Vars: map[string]interface{}{
// 						"container_name":   "php-fpm-dev",
// 						"source_image_tag": "16.04",
// 					},
// 				},
// 			},
// 			"ubuntu": {
// 				"16.04": &image.Image{
// 					Builder: "infrastructure",
// 					Tags: []string{
// 						"16.04",
// 						"xenial",
// 					},
// 					Vars: map[string]interface{}{
// 						"container_name":   "ubuntu",
// 						"source_image_tag": "16.04",
// 					},
// 					Children: map[string][]string{
// 						"php-builder": {
// 							"7.1",
// 						},
// 						"php-fpm": {
// 							"7.1",
// 							"7.2",
// 						},
// 						"php-cli": {
// 							"7.1",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	tests := []struct {
// 		desc         string
// 		nodeImage    *Image
// 		imageGraph   *gdsexttree.Graph
// 		imageName    string
// 		imageVersion string
// 		parent       *gdsexttree.Node
// 		res          *gdstree.Graph
// 		err          error
// 	}{
// 		{
// 			desc:         "Generate graph with an empty node image",
// 			nodeImage:    nil,
// 			imageGraph:   &gdsexttree.Graph{},
// 			imageName:    "",
// 			imageVersion: "",
// 			parent:       nil,
// 			res:          nil,
// 			err:          errors.New("(tree::generateGraphRec)", "Node Image is null"),
// 		},
// 		{
// 			desc: "Adding Image to an existing graph",
// 			nodeImage: &image.Image{
// 				Name:    "nginx",
// 				Version: "1.10",
// 			},
// 			imageName:    "nginx",
// 			imageVersion: "1.10",
// 			parent: &gdsexttree.Node{
// 				Name: "ubuntu:16.04",
// 				Item: &image.Image{},
// 			},
// 			imageGraph: &gdsexttree.Graph{
// 				Root: []*gdsexttree.Node{
// 					{
// 						Name: "ubuntu:16.04",
// 					},
// 				},
// 				NodesIndex: map[string]*gdsexttree.Node{
// 					"ubuntu:16.04": {
// 						Name: "ubuntu:16.04",
// 					},
// 					"php-fpm:7.1": {
// 						Name: "php-fpm:7.1",
// 						Parents: []*gdsexttree.Node{
// 							{
// 								Name: "ubuntu:16.04",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			res: &gdstree.Graph{
// 				Root: []*gdstree.Node{
// 					{
// 						Name: "ubuntu:16.04",
// 					},
// 				},
// 				NodesIndex: map[string]*gdstree.Node{
// 					"ubuntu:16.04": {
// 						Name: "ubuntu:16.04",
// 					},
// 					"php-fpm:7.1": {
// 						Name: "php-fpm:7.1",
// 						Parent: &gdstree.Node{
// 							Name: "ubuntu:16.04",
// 						},
// 					},
// 					"nginx:1.10": {
// 						Name: "nginx:1.10",
// 						Parent: &gdstree.Node{
// 							Name: "ubuntu:16.04",
// 						},
// 					},
// 				},
// 			},
// 			err: nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Log(test.desc)

// 		err := imagesTree.generateTemplateGraph(test.imageName, test.imageVersion, test.nodeImage, test.imageGraph, test.parent)

// 		if err != nil && assert.Error(t, err) {
// 			assert.Equal(t, test.err.Error(), err.Error())
// 		} else {
// 			assert.Equal(t, len(test.res.Root), len(test.imageGraph.Root), "Unexpected lenght of root elements")
// 			assert.Equal(t, len(test.res.NodesIndex), len(test.imageGraph.NodesIndex), "Unexpected lenght of nodes index elements")
// 		}
// 	}
// }

// func TestGenerateNodeName(t *testing.T) {
// 	t.Log("Testing node name generation")

// 	name := "name"
// 	version := "version"
// 	res := name + ImageNodeNameSeparator + version

// 	i := &image.Image{
// 		Name:    name,
// 		Version: version,
// 	}
// 	nodename := GenerateNodeName(i)

// 	assert.Equal(t, nodename, res, "Nodename is not valid")
// }

// func TestRenderizeGraph(t *testing.T) {

// 	phpFpmDev71 := &gdsexttree.Node{
// 		Name: "php-fpm-dev:7.1",
// 		Item: &image.Image{
// 			Name:    "php-fpm-dev",
// 			Version: "{{ .Parent.Version }}",
// 		},
// 	}
// 	phpCliDev71 := &gdsexttree.Node{
// 		Name: "php-cli-dev:7.1",
// 		Item: &image.Image{
// 			Name:    "php-cli-dev",
// 			Version: "{{ .Parent.Version }}",
// 		},
// 	}
// 	phpFpm71 := &gdsexttree.Node{
// 		Name: "php-fpm:7.1",
// 		Item: &image.Image{
// 			Name:    "php-fpm",
// 			Version: "7.1-ubuntu{{ .Parent.Version }}",
// 		},
// 		Children: []*gdsexttree.Node{
// 			phpFpmDev71,
// 		},
// 	}
// 	phpFpm72 := &gdsexttree.Node{
// 		Name: "php-fpm:7.2",
// 		Item: &image.Image{
// 			Name:    "php-fpm",
// 			Version: "7.2-ubuntu{{ .Parent.Version }}",
// 		},
// 	}
// 	phpCli71 := &gdsexttree.Node{
// 		Name: "php-cli:7.1",
// 		Item: &image.Image{
// 			Name:    "php-cli",
// 			Version: "7.1-ubuntu{{ .Parent.Version }}",
// 		},
// 		Children: []*gdsexttree.Node{
// 			phpCliDev71,
// 		},
// 	}
// 	phpCli72 := &gdsexttree.Node{
// 		Name: "php-cli:7.2",
// 		Item: &image.Image{
// 			Name:    "php-cli",
// 			Version: "7.2-ubuntu{{ .Parent.Version }}",
// 		},
// 	}
// 	phpBuilder71 := &gdsexttree.Node{
// 		Name: "php-builder:7.1",
// 		Item: &image.Image{
// 			Name:    "php-builder",
// 			Version: "7.1-ubuntu{{ .Parent.Version }}",
// 		},
// 	}
// 	ubuntu16 := &gdsexttree.Node{
// 		Name: "ubuntu:16.04",
// 		Item: &image.Image{
// 			Name:    "ubuntu",
// 			Version: "16.04",
// 		},
// 		Children: []*gdsexttree.Node{
// 			phpFpm71,
// 			phpFpm72,
// 			phpCli71,
// 			phpCli72,
// 			phpBuilder71,
// 		},
// 	}
// 	ubuntu18 := &gdsexttree.Node{
// 		Name: "ubuntu:18.04",
// 		Item: &image.Image{
// 			Name:    "ubuntu",
// 			Version: "18.04",
// 		},
// 		Children: []*gdsexttree.Node{
// 			phpFpm71,
// 			phpFpm72,
// 			phpCli71,
// 			phpCli72,
// 			phpBuilder71,
// 		},
// 	}

// 	phpBuilder71.AddParent(ubuntu16)
// 	phpBuilder71.AddParent(ubuntu18)
// 	phpCli71.AddParent(ubuntu16)
// 	phpCli71.AddParent(ubuntu18)
// 	phpCli72.AddParent(ubuntu16)
// 	phpCli72.AddParent(ubuntu18)
// 	phpFpm71.AddParent(ubuntu16)
// 	phpFpm71.AddParent(ubuntu18)
// 	phpFpm72.AddParent(ubuntu16)
// 	phpFpm72.AddParent(ubuntu18)
// 	phpFpmDev71.AddParent(phpFpm71)
// 	phpCliDev71.AddParent(phpCli71)

// 	imagesGraph := &gdsexttree.Graph{
// 		Root: []*gdsexttree.Node{
// 			ubuntu16,
// 			ubuntu18,
// 		},
// 		NodesIndex: map[string]*gdsexttree.Node{
// 			"ubuntu:16.04":    ubuntu16,
// 			"ubuntu:18.04":    ubuntu18,
// 			"php-fpm:7.1":     phpFpm71,
// 			"php-fpm:7.2":     phpFpm72,
// 			"php-cli:7.1":     phpCli71,
// 			"php-cli:7.2":     phpCli72,
// 			"php-builder:7.1": phpBuilder71,
// 			"php-fpm-dev:7.1": phpFpmDev71,
// 			"php-cli-dev:7.1": phpCliDev71,
// 		},
// 	}

// 	ResPhpFpmDev7116 := &gdstree.Node{
// 		Name: "php-fpm-dev:7.1-ubuntu16.04@php-fpm:7.1-ubuntu16.04",
// 		Item: &image.Image{
// 			Name:    "php-fpm-dev",
// 			Version: "7.1-ubuntu16.04",
// 		},
// 	}
// 	ResPhpFpmDev7118 := &gdstree.Node{
// 		Name: "php-fpm-dev:7.1-ubuntu18.04@php-fpm:7.1-ubuntu18.04",
// 		Item: &image.Image{
// 			Name:    "php-fpm-dev",
// 			Version: "7.1-ubuntu18.04",
// 		},
// 	}
// 	ResPhpCliDev7116 := &gdstree.Node{
// 		Name: "php-cli-dev:7.1-ubuntu16.04@php-fpm:7.1-ubuntu16.04",
// 		Item: &image.Image{
// 			Name:    "php-cli-dev",
// 			Version: "7.1-ubuntu16.04",
// 		},
// 	}
// 	ResPhpCliDev7118 := &gdstree.Node{
// 		Name: "php-cli-dev:7.1-ubuntu18.04@php-fpm:7.1-ubuntu18.04",
// 		Item: &image.Image{
// 			Name:    "php-cli-dev",
// 			Version: "7.1-ubuntu18.04",
// 		},
// 	}
// 	ResPhpFpm7116 := &gdstree.Node{
// 		Name: "php-fpm:7.1-ubuntu16.04@ubuntu:16.04",
// 		Item: &image.Image{
// 			Name:    "php-fpm",
// 			Version: "7.1-ubuntu16.04",
// 		},
// 	}
// 	ResPhpFpm7118 := &gdstree.Node{
// 		Name: "php-fpm:7.1-ubuntu18.04@ubuntu:18.04",
// 		Item: &image.Image{
// 			Name:    "php-fpm",
// 			Version: "7.1-ubuntu18.04",
// 		},
// 	}
// 	ResPhpFpm7216 := &gdstree.Node{
// 		Name: "php-fpm:7.2-ubuntu16.04@ubuntu:16.04",
// 		Item: &image.Image{
// 			Name:    "php-fpm",
// 			Version: "7.2-ubuntu16.04",
// 		},
// 	}
// 	ResPhpFpm7218 := &gdstree.Node{
// 		Name: "php-fpm:7.2-ubuntu18.04@ubuntu:18.04",
// 		Item: &image.Image{
// 			Name:    "php-fpm",
// 			Version: "7.2-ubuntu18.04",
// 		},
// 	}
// 	ResPhpCli7116 := &gdstree.Node{
// 		Name: "php-cli:7.1-ubuntu16.04@ubuntu:16.04",
// 		Item: &image.Image{
// 			Name:    "php-cli",
// 			Version: "7.1-ubuntu16.04",
// 		},
// 	}
// 	ResPhpCli7118 := &gdstree.Node{
// 		Name: "php-cli:7.1-ubuntu18.04@ubuntu:18.04",
// 		Item: &image.Image{
// 			Name:    "php-cli",
// 			Version: "7.1-ubuntu18.04",
// 		},
// 	}
// 	ResPhpCli7216 := &gdstree.Node{
// 		Name: "php-cli:7.2-ubuntu16.04@ubuntu:16.04",
// 		Item: &image.Image{
// 			Name:    "php-cli",
// 			Version: "7.2-ubuntu16.04",
// 		},
// 	}
// 	ResPhpCli7218 := &gdstree.Node{
// 		Name: "php-cli:7.2-ubuntu18.04@ubuntu:18.04",
// 		Item: &image.Image{
// 			Name:    "php-cli",
// 			Version: "7.2-ubuntu18.04",
// 		},
// 	}
// 	ResPhpBuilder7116 := &gdstree.Node{
// 		Name: "php-builder:7.1-ubuntu16.04@ubuntu:16.04",
// 		Item: &image.Image{
// 			Name:    "php-builder",
// 			Version: "7.1-ubuntu16.04",
// 		},
// 	}
// 	ResPhpBuilder7118 := &gdstree.Node{
// 		Name: "php-builder:7.1-ubuntu18.04@ubuntu:16.04",
// 		Item: &image.Image{
// 			Name:    "php-builder",
// 			Version: "7.1-ubuntu18.04",
// 		},
// 	}
// 	ResUbuntu16 := &gdstree.Node{
// 		Name: "ubuntu:16.04",
// 		Item: &image.Image{
// 			Name:    "ubuntu",
// 			Version: "16.04",
// 		},
// 	}
// 	ResUbuntu18 := &gdstree.Node{
// 		Name: "ubuntu:18.04",
// 		Item: &image.Image{
// 			Name:    "ubuntu",
// 			Version: "18.04",
// 		},
// 	}

// 	graphRes := &gdstree.Graph{}
// 	// bases
// 	graphRes.AddNode(ResUbuntu16)
// 	graphRes.AddNode(ResUbuntu18)
// 	// ubuntu16 base
// 	graphRes.AddNode(ResPhpBuilder7116)
// 	graphRes.AddRelationship(ResUbuntu16, ResPhpBuilder7116)
// 	graphRes.AddNode(ResPhpFpm7116)
// 	graphRes.AddRelationship(ResUbuntu16, ResPhpFpm7116)
// 	graphRes.AddNode(ResPhpFpm7216)
// 	graphRes.AddRelationship(ResUbuntu16, ResPhpFpm7216)
// 	graphRes.AddNode(ResPhpCli7116)
// 	graphRes.AddRelationship(ResUbuntu16, ResPhpCli7116)
// 	graphRes.AddNode(ResPhpCli7216)
// 	graphRes.AddRelationship(ResUbuntu16, ResPhpCli7216)
// 	graphRes.AddNode(ResPhpFpmDev7116)
// 	graphRes.AddRelationship(ResPhpFpm7116, ResPhpFpmDev7116)
// 	graphRes.AddNode(ResPhpCliDev7116)
// 	graphRes.AddRelationship(ResPhpCli7116, ResPhpCliDev7116)
// 	// ubuntu18 base
// 	graphRes.AddNode(ResPhpBuilder7118)
// 	graphRes.AddRelationship(ResUbuntu16, ResPhpBuilder7118)
// 	graphRes.AddNode(ResPhpFpm7118)
// 	graphRes.AddRelationship(ResUbuntu18, ResPhpFpm7118)
// 	graphRes.AddNode(ResPhpFpm7218)
// 	graphRes.AddRelationship(ResUbuntu18, ResPhpFpm7218)
// 	graphRes.AddNode(ResPhpCli7118)
// 	graphRes.AddRelationship(ResUbuntu18, ResPhpCli7118)
// 	graphRes.AddNode(ResPhpCli7218)
// 	graphRes.AddRelationship(ResUbuntu18, ResPhpCli7218)
// 	graphRes.AddNode(ResPhpFpmDev7118)
// 	graphRes.AddRelationship(ResPhpFpm7118, ResPhpFpmDev7118)
// 	graphRes.AddNode(ResPhpCliDev7118)
// 	graphRes.AddRelationship(ResPhpCli7118, ResPhpCliDev7118)

// 	indexhRes := &ImageIndex{
// 		NameIndex: map[string][]string{
// 			"ubuntu":      {"ubuntu:16.04", "ubuntu:18.04"},
// 			"php-fpm":     {""},
// 			"php-fpm-dev": {""},
// 			"php-cli":     {""},
// 			"php-cli-dev": {""},
// 			"php-builder": {""},
// 		},
// 		NameVersionIndex: map[string][]*gdstree.Node{
// 			"ubuntu:16.04":    nil,
// 			"ubuntu:18.04":    nil,
// 			"php-fpm:7.1":     nil,
// 			"php-fpm-dev:7.1": nil,
// 			"php-cli:7.1":     nil,
// 			"php-cli-dev:7.1": nil,
// 			"php-builder:7.1": nil,
// 			"php-fpm:7.2":     nil,
// 			"php-cli:7.2":     nil,
// 		},
// 		NameVersionAlternativeIndex: map[string][]*gdstree.Node{
// 			"php-fpm:7.1-ubuntu16.04":     nil,
// 			"php-fpm-dev:7.1-ubuntu16.04": nil,
// 			"php-cli:7.1-ubuntu16.04":     nil,
// 			"php-cli-dev:7.1-ubuntu16.04": nil,
// 			"php-builder:7.1-ubuntu16.04": nil,
// 			"php-fpm:7.2-ubuntu16.04":     nil,
// 			"php-cli:7.2-ubuntu16.04":     nil,
// 			"php-fpm:7.1-ubuntu18.04":     nil,
// 			"php-fpm-dev:7.1-ubuntu18.04": nil,
// 			"php-cli:7.1-ubuntu18.04":     nil,
// 			"php-cli-dev:7.1-ubuntu18.04": nil,
// 			"php-builder:7.1-ubuntu18.04": nil,
// 			"php-fpm:7.2-ubuntu18.04":     nil,
// 			"php-cli:7.2-ubuntu18.04":     nil,
// 		},
// 	}

// 	g, i, _ := RenderizeGraph(imagesGraph)
// 	t.Log("Testing graph elements")
// 	assert.Equal(t, len(graphRes.Root), len(g.Root), "Unexpected lenght of root elements")
// 	assert.Equal(t, len(graphRes.NodesIndex), len(g.NodesIndex), "Unexpected lenght of nodes index elements")
// 	t.Log("Testing images index elements")
// 	assert.Equal(t, len(indexhRes.NameIndex), len(i.NameIndex), "Unexpected lenght of index names")
// 	assert.Equal(t, len(indexhRes.NameVersionIndex), len(i.NameVersionIndex), "Unexpected lenght of name-version elements")
// 	assert.Equal(t, len(indexhRes.NameVersionAlternativeIndex), len(i.NameVersionAlternativeIndex), "Unexpected lenght of alternatives elements")
// }

// func TestRenderizeGraphRec(t *testing.T) {
// 	//TODO
// }

// func TestGetNodeImage(t *testing.T) {

// 	imageNode := &image.Image{}

// 	tests := []struct {
// 		desc string
// 		node *gdstree.Node
// 		res  *Image
// 		err  error
// 	}{
// 		{
// 			desc: "Testing get node image from a nil node",
// 			node: nil,
// 			res:  nil,
// 			err:  errors.New("(tree::GetNodeImage)", "Node is nil"),
// 		},
// 		{
// 			desc: "Testing get node image from a node with nil image",
// 			node: &gdstree.Node{
// 				Item: nil,
// 			},
// 			res: nil,
// 			err: errors.New("(tree::GetNodeImage)", "Node item is nil"),
// 		},
// 		{
// 			desc: "Testing get node image from a node",
// 			node: &gdstree.Node{
// 				Item: imageNode,
// 			},
// 			res: imageNode,
// 			err: errors.New("(tree::GetNodeImage)", "Node item is nil"),
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Log(test.desc)

// 		res, err := GetNodeImage(test.node)
// 		if err != nil && assert.Error(t, err) {
// 			assert.Equal(t, test.err.Error(), err.Error())
// 		} else {
// 			assert.Equal(t, test.res, res, "Unexpected Image")
// 		}
// 	}
// }

// func TestGenerateWilcardVersionNode(t *testing.T) {

// 	phpFpm71 := &image.Image{
// 		Name:    "php-fpm",
// 		Builder: "infrastructure",
// 		Tags: []string{
// 			"7.1",
// 		},
// 		PersistentVars: map[string]interface{}{
// 			"php_version": "7.1",
// 		},
// 		Vars: map[string]interface{}{
// 			"container_name": "php-fpm",
// 		},
// 		Children: map[string][]string{
// 			"php-fpm-dev": {
// 				"7.1",
// 			},
// 		},
// 	}

// 	phpFpmWilcard := &image.Image{
// 		Name:    "php-fpm",
// 		Builder: "infrastructure",
// 		Tags: []string{
// 			"{{ .Version }}",
// 		},
// 		PersistentVars: map[string]interface{}{
// 			"php_version": "{{ .Version }}",
// 		},
// 		Vars: map[string]interface{}{
// 			"container_name": "php-fpm",
// 		},
// 		Children: map[string][]string{
// 			"php-fpm-dev": {
// 				"{{ .Version }}",
// 			},
// 		},
// 	}

// 	phpFpmDev71 := &image.Image{
// 		Name:    "php-fpm-dev",
// 		Builder: "infrastructure",
// 		Tags: []string{
// 			"7.1",
// 		},
// 		Vars: map[string]interface{}{
// 			"container_name": "php-fpm-dev",
// 		},
// 	}

// 	phpFpmDevWilcard := &image.Image{
// 		Name:    "php-fpm-dev",
// 		Builder: "infrastructure",
// 		Tags: []string{
// 			"{{ .Version }}",
// 		},
// 		Vars: map[string]interface{}{
// 			"container_name": "php-fpm-dev",
// 		},
// 	}

// 	ubuntu16 := &image.Image{
// 		Name:    "ubuntu",
// 		Builder: "infrastructure",
// 		Tags: []string{
// 			"16.04",
// 			"xenial",
// 		},
// 		PersistentVars: map[string]interface{}{
// 			"ubuntu_version": "16.04",
// 		},
// 		Vars: map[string]interface{}{
// 			"container_name":   "ubuntu",
// 			"source_image_tag": "16.04",
// 		},
// 		Children: map[string][]string{
// 			"php-fpm": {
// 				"7.1",
// 				"*",
// 			},
// 		},
// 	}

// 	tree := &ImagesConfiguration{
// 		Images: map[string]map[string]*image.Image{
// 			"php-fpm": {
// 				"7.1": phpFpm71,
// 				"*":   phpFpmWilcard,
// 			},
// 			"php-fpm-dev": {
// 				"7.1": phpFpmDev71,
// 				"*":   phpFpmDevWilcard,
// 			},
// 			"ubuntu": {
// 				"16.04": ubuntu16,
// 			},
// 		},
// 	}

// 	tests := []struct {
// 		desc     string
// 		tree     *ImagesConfiguration
// 		version  string
// 		nodeBase *gdstree.Node
// 		res      *gdstree.Node
// 		err      error
// 	}{
// 		{
// 			desc:     "Testing generate wildcard version gave a nil images tree",
// 			tree:     nil,
// 			version:  "",
// 			nodeBase: nil,
// 			res:      nil,
// 			err:      errors.New("(tree::GenerateNodeWithWilcardVersion)", "Images tree is nil"),
// 		},
// 		{
// 			desc:     "Testing generate wildcard version gave a nil node",
// 			tree:     tree,
// 			version:  "",
// 			nodeBase: nil,
// 			res:      nil,
// 			err:      errors.New("(tree::GenerateNodeWithWilcardVersion)", "Node is nil"),
// 		},
// 		{
// 			desc:    "Testing generate wildcard version node",
// 			version: "version",
// 			nodeBase: &gdstree.Node{
// 				Name:   "php-fpm:*",
// 				Item:   phpFpmWilcard,
// 				Parent: nil,
// 			},
// 			tree: tree,
// 			res: &gdstree.Node{
// 				Name: "php-fpm:version",
// 				Item: &image.Image{
// 					Name:      "php-fpm",
// 					Namespace: "",
// 					Version:   "version",
// 					Builder:   "infrastructure",
// 					Tags: []string{
// 						"version",
// 					},
// 					PersistentVars: map[string]interface{}{
// 						"php_version": "version",
// 					},
// 					Vars: map[string]interface{}{
// 						"container_name": "php-fpm",
// 					},
// 					Children: map[string][]string{
// 						"php-fpm-dev": {
// 							"version",
// 						},
// 					},
// 					Childs: map[string][]string{},
// 				},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Log(test.desc)

// 		res, err := test.tree.GenerateWilcardVersionNode(test.nodeBase, test.version)
// 		if err != nil && assert.Error(t, err) {
// 			assert.Equal(t, test.err.Error(), err.Error())
// 		} else {
// 			assert.Equal(t, test.res, res, "Unexpected Node")
// 		}
// 	}
// }
