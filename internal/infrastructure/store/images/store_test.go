package images

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/render"

	// "github.com/gostevedore/stevedore/internal/infrastructure/render"
	"github.com/stretchr/testify/assert"
)

func TestAddImageNameDefinitionVersionList(t *testing.T) {
	errContext := "(store::images::Store::addImageNameDefinitionVersionList)"
	tests := []struct {
		desc    string
		store   *Store
		name    string
		version string
		res     map[string]map[string]struct{}
		err     error
	}{
		{
			desc:  "Testing error when adding a definition version into list without providing image name",
			store: &Store{},
			err:   errors.New(errContext, "Image name must be provided to add to add a definition version into list"),
		},
		{
			desc:  "Testing error when adding a definition version into list without providing image version",
			store: &Store{},
			name:  "name",
			err:   errors.New(errContext, "Image version must be provided to add to add a definition version into list"),
		},
		{
			desc:    "Testing add image name definition version into store",
			store:   &Store{},
			name:    "image1",
			version: "v1",
			res: map[string]map[string]struct{}{
				"image1": {"v1": struct{}{}},
			},
		},
		{
			desc: "Testing append image name definition version into store",
			store: &Store{
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}},
				},
			},
			name:    "image1",
			version: "v2",
			res: map[string]map[string]struct{}{
				"image1": {"v1": struct{}{}, "v2": struct{}{}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.store.addImageNameDefinitionVersionList(test.name, test.version)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.store.imageNameDefinitionVersionList)
			}
		})
	}

}

func TestAddImageNameVersionRenderedVersionsList(t *testing.T) {

	errContext := "(store::images::Store::addImageNameVersionRenderedVersionsList)"

	tests := []struct {
		desc    string
		store   *Store
		name    string
		version string
		image   *image.Image
		res     map[string]map[string]map[string]struct{}
		err     error
	}{
		{
			desc:  "Testing error when adding a rendered version into list without providing image name",
			store: &Store{},
			err:   errors.New(errContext, "Image name must be provided to add to add a rendered version into list"),
		},
		{
			desc:  "Testing error when adding a rendered version into list without providing image version",
			store: &Store{},
			name:  "name",
			err:   errors.New(errContext, "Image version must be provided to add to add a rendered version into list"),
		},
		{
			desc:    "Testing error when adding a rendered version into list without providing image",
			store:   &Store{},
			name:    "name",
			version: "version",
			err:     errors.New(errContext, "Image must be provided to add to add a rendered version into list"),
		},
		{
			desc:    "Testing add image name rendered version into store",
			store:   &Store{},
			name:    "image1",
			version: "v1",
			image: &image.Image{
				Version: "v1-parent",
			},
			res: map[string]map[string]map[string]struct{}{
				"image1": {
					"v1": {"v1-parent": struct{}{}},
				},
			},
		},
		{
			desc: "Testing append image name rendered version into store",
			store: &Store{
				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-parent": struct{}{}},
					},
				},
			},
			name:    "image1",
			version: "v1",
			image: &image.Image{
				Version: "v1-another-parent",
			},
			res: map[string]map[string]map[string]struct{}{
				"image1": {
					"v1": {"v1-parent": struct{}{}, "v1-another-parent": struct{}{}},
				},
			},
		},
		{
			desc: "Testing add a second image name rendered version into store",
			store: &Store{
				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-parent": struct{}{}},
					},
				},
			},
			name:    "image1",
			version: "v2",
			image: &image.Image{
				Version: "v2-parent",
			},
			res: map[string]map[string]map[string]struct{}{
				"image1": {
					"v1": {"v1-parent": struct{}{}},
					"v2": {"v2-parent": struct{}{}},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.store.addImageNameVersionRenderedVersionsList(test.name, test.version, test.image)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.store.imageNameVersionRenderedVersionsList)
			}
		})
	}
}

func TestAddImageToIndex(t *testing.T) {

	errContext := "(store::images::Store::addImageToIndex)"

	tests := []struct {
		desc  string
		store *Store
		name  string
		image *image.Image
		res   map[string]map[string]*image.Image
		err   error
	}{
		{
			desc:  "Testing error when adding image into index without providing image name",
			store: &Store{},
			err:   errors.New(errContext, "Image name must be provided to add to add image to index"),
		},
		{
			desc:  "Testing error when adding image into index without providing image",
			store: &Store{},
			name:  "name",
			err:   errors.New(errContext, "Image must be provided to add to add to add image to index"),
		},

		{
			desc:  "Testing add image to index",
			store: &Store{},
			name:  "image1",
			image: &image.Image{
				Version: "v1-parent",
				Tags:    []string{"latest"},
			},
			res: map[string]map[string]*image.Image{
				"image1": {
					"v1-parent": &image.Image{
						Version: "v1-parent",
						Tags:    []string{"latest"},
					},
					"latest": &image.Image{
						Version: "v1-parent",
						Tags:    []string{"latest"},
					},
				},
			},
		},
		{
			desc: "Testing add image with another version into index",
			store: &Store{
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-parent": &image.Image{
							Version: "v1-parent",
							Tags:    []string{"latest"},
						},
						"latest": &image.Image{
							Version: "v1-parent",
							Tags:    []string{"latest"},
						},
					},
				},
			},
			name: "image1",
			image: &image.Image{
				Version: "v2-parent",
				Tags:    []string{"beta"},
			},
			res: map[string]map[string]*image.Image{
				"image1": {
					"v1-parent": &image.Image{
						Version: "v1-parent",
						Tags:    []string{"latest"},
					},
					"latest": &image.Image{
						Version: "v1-parent",
						Tags:    []string{"latest"},
					},
					"v2-parent": &image.Image{
						Version: "v2-parent",
						Tags:    []string{"beta"},
					},
					"beta": &image.Image{
						Version: "v2-parent",
						Tags:    []string{"beta"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.store.addImageToIndex(test.name, test.image)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.store.imagesIndex)
			}
		})
	}
}

func TestStore(t *testing.T) {

	errContext := "(store::Store)"

	tests := []struct {
		desc              string
		store             *Store
		name              string
		version           string
		image             *image.Image
		err               error
		prepareAssertFunc func(*Store, *image.Image)
		assertFunc        func(*testing.T, *Store)
	}{
		{
			desc:  "Testing error when render is not defined",
			store: NewStore(nil),
			err:   errors.New(errContext, "To add an image to the store an image render is required"),
		},
		{
			desc:  "Testing error when name is not defined",
			store: NewStore(render.NewMockImageRender()),
			err:   errors.New(errContext, "To add an image to the store a name is required"),
		},
		{
			desc:  "Testing error when version is not defined",
			store: NewStore(render.NewMockImageRender()),
			name:  "image_name",
			err:   errors.New(errContext, "To add an image to the store a version is required"),
		},
		{
			desc:    "Testing error when image is not defined",
			store:   NewStore(render.NewMockImageRender()),
			name:    "image_name",
			version: "image_version",
			err:     errors.New(errContext, "To add an image to the store an image is required"),
		},

		// Commented because store does not render images anymore it is done by imagesConfiguration
		//
		// {
		// 	desc:    "Testing add a image to an empty store",
		// 	store:   NewStore(render.NewMockImageRender()),
		// 	name:    "image_name",
		// 	version: "image_version",
		// 	image: &image.Image{
		// 		Name:              "{{.Name}}-{{.Parent.Name}}",
		// 		Version:           "{{.Version}}-{{.Parent.Version}}",
		// 		RegistryNamespace: "{{.Parent.RegistryNamespace}}",
		// 		Parent: &image.Image{
		// 			Name:              "parent_name",
		// 			Version:           "parent_version",
		// 			RegistryNamespace: "parent_registry_namespace",
		// 		},
		// 		Tags: []string{"tag1", "tag2"},
		// 	},
		// 	prepareAssertFunc: func(s *Store, i *image.Image) {
		// 		s.render.(*render.MockImageRender).On("Render", "image_name", "image_version", i).Return(
		// 			&image.Image{
		// 				Name:              "image_name-parent_name",
		// 				Version:           "image_version-parent_version",
		// 				RegistryNamespace: "parent_registry_namespace",
		// 				Parent: &image.Image{
		// 					Name:              "parent_name",
		// 					Version:           "parent_version",
		// 					RegistryNamespace: "parent_registry_namespace",
		// 				},
		// 				Tags: []string{"tag1", "tag2"},
		// 			},
		// 			nil,
		// 		)
		// 	},
		// 	assertFunc: func(t *testing.T, s *Store) {
		// 		assert.Equal(t, 1, len(s.store), "Unexpected number of images in the store")
		// 		assert.Equal(t,
		// 			map[string]map[string]struct{}{
		// 				"image_name": {"image_version": struct{}{}},
		// 			},
		// 			s.imageNameDefinitionVersionList)
		// 		assert.Equal(t,
		// 			map[string]map[string]map[string]struct{}{
		// 				"image_name": {
		// 					"image_version": {"image_version-parent_version": struct{}{}},
		// 				},
		// 			},
		// 			s.imageNameVersionRenderedVersionsList)
		// 		assert.Equal(t,
		// 			map[string]map[string]map[string]struct{}{
		// 				"image_name": {
		// 					"image_version": {"image_version-parent_version": struct{}{}},
		// 				},
		// 			},
		// 			s.imageNameVersionRenderedVersionsList)

		// 	},
		// 	err: &errors.Error{},
		// },

		{
			desc:    "Testing add a image to an empty store",
			store:   NewStore(render.NewMockImageRender()),
			name:    "image_name",
			version: "image_version",
			image: &image.Image{
				Name:              "image_name",
				Version:           "image_version",
				RegistryNamespace: "image_namespace",
				Parent: &image.Image{
					Name:              "parent_name",
					Version:           "parent_version",
					RegistryNamespace: "parent_registry_namespace",
				},
				Tags: []string{"tag1", "tag2"},
			},
			prepareAssertFunc: func(s *Store, i *image.Image) {
				s.render.(*render.MockImageRender).On("Render", "image_name", "image_version", i).Return(
					&image.Image{
						Name:              "image_name",
						Version:           "image_version",
						RegistryNamespace: "image_namespace",
						Parent: &image.Image{
							Name:              "parent_name",
							Version:           "parent_version",
							RegistryNamespace: "parent_registry_namespace",
						},
						Tags: []string{"tag1", "tag2"},
					},
					nil,
				)
			},
			assertFunc: func(t *testing.T, s *Store) {
				assert.Equal(t, 1, len(s.store), "Unexpected number of images in the store")
				assert.Equal(t,
					map[string]map[string]struct{}{
						"image_name": {"image_version": struct{}{}},
					},
					s.imageNameDefinitionVersionList)
				assert.Equal(t,
					map[string]map[string]map[string]struct{}{
						"image_name": {
							"image_version": {"image_version": struct{}{}},
						},
					},
					s.imageNameVersionRenderedVersionsList)
			},
			err: &errors.Error{},
		},

		{
			desc:    "Testing add a wildcard image to an empty store",
			store:   NewStore(render.NewMockImageRender()),
			name:    "image_name",
			version: "*",
			image: &image.Image{
				Name:              "{{.Name}}-{{.Parent.Name}}",
				Version:           "{{.Version}}-{{.Parent.Version}}",
				RegistryNamespace: "{{.Parent.RegistryNamespace}}",
				Parent: &image.Image{
					Name:              "parent_name",
					Version:           "parent_version",
					RegistryNamespace: "parent_registry_namespace",
				},
				Tags: []string{"tag1", "tag2"},
			},
			prepareAssertFunc: func(s *Store, i *image.Image) {},
			assertFunc: func(t *testing.T, s *Store) {
				image := &image.Image{
					Name:              "{{.Name}}-{{.Parent.Name}}",
					Version:           "{{.Version}}-{{.Parent.Version}}",
					RegistryNamespace: "{{.Parent.RegistryNamespace}}",
					Parent: &image.Image{
						Name:              "parent_name",
						Version:           "parent_version",
						RegistryNamespace: "parent_registry_namespace",
					},
					Tags: []string{"tag1", "tag2"},
				}
				storedImage, exist := s.imageWildcardIndex["image_name"]
				assert.True(t, exist, "Image is not on the wildcard index")
				assert.Equal(t, image, storedImage, "Unexpected image in the wildcard index")
				assert.Equal(t, 0, len(s.store))
				assert.Equal(t, 0, len(s.imageNameDefinitionVersionList))
				assert.Equal(t, 0, len(s.imageNameVersionRenderedVersionsList))
				assert.Equal(t, 0, len(s.imagesIndex))
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.store, test.image)
			}

			err := test.store.Store(test.name, test.version, test.image)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.store)
			}
		})
	}

}

func TestStoreWildcardImage(t *testing.T) {
	errContext := "(store::images::Store::storeWildcardImage)"

	tests := []struct {
		desc              string
		store             *Store
		name              string
		image             *image.Image
		err               error
		prepareAssertFunc func(s *Store)
		assertFunc        func(t *testing.T, s *Store)
	}{
		{
			desc:       "Testing error storing a wildcard image without providing a name",
			store:      NewStore(render.NewMockImageRender()),
			err:        errors.New(errContext, "Image name must be provided to store a wildcard image"),
			assertFunc: func(t *testing.T, s *Store) {},
		},
		{
			desc:       "Testing error storing a wildcard image without providing an image",
			store:      NewStore(render.NewMockImageRender()),
			name:       "image_name",
			err:        errors.New(errContext, "Image must be provided to store 'image_name' wildcard image"),
			assertFunc: func(t *testing.T, s *Store) {},
		},
		{
			desc:  "Testing store a wildcard image",
			store: NewStore(render.NewMockImageRender()),
			name:  "image_name",
			image: &image.Image{
				Name:              "{{.Name}}-{{.Parent.Name}}",
				Version:           "{{.Version}}-{{.Parent.Version}}",
				RegistryNamespace: "{{.Parent.RegistryNamespace}}",
				Parent: &image.Image{
					Name:              "parent_name",
					Version:           "parent_version",
					RegistryNamespace: "parent_registry_namespace",
				},
				Tags: []string{"tag1", "tag2"},
			},
			err: &errors.Error{},
			assertFunc: func(t *testing.T, s *Store) {
				expected := &image.Image{
					Name:              "{{.Name}}-{{.Parent.Name}}",
					Version:           "{{.Version}}-{{.Parent.Version}}",
					RegistryNamespace: "{{.Parent.RegistryNamespace}}",
					Parent: &image.Image{
						Name:              "parent_name",
						Version:           "parent_version",
						RegistryNamespace: "parent_registry_namespace",
					},
					Tags: []string{"tag1", "tag2"},
				}

				image, exist := s.imageWildcardIndex["image_name"]
				assert.True(t, exist, "Image name is not on the index")
				assert.Equal(t, image, expected, "Unexpected image in the index")
			},
		},
		{
			desc:  "Testing store an existing image",
			store: NewStore(render.NewMockImageRender()),
			name:  "image_name",
			image: &image.Image{
				Name:              "{{.Name}}-{{.Parent.Name}}",
				Version:           "{{.Version}}-{{.Parent.Version}}",
				RegistryNamespace: "{{.Parent.RegistryNamespace}}",
				Parent: &image.Image{
					Name:              "parent_name",
					Version:           "parent_version",
					RegistryNamespace: "parent_registry_namespace",
				},
				Tags: []string{"tag1", "tag2"},
			},
			err: errors.New(errContext, "Image 'image_name' already exists on wildcard images index"),
			prepareAssertFunc: func(s *Store) {
				s.storeWildcardImage("image_name", &image.Image{
					Name:              "{{.Name}}-{{.Parent.Name}}",
					Version:           "{{.Version}}-{{.Parent.Version}}",
					RegistryNamespace: "{{.Parent.RegistryNamespace}}",
					Parent: &image.Image{
						Name:              "parent_name",
						Version:           "parent_version",
						RegistryNamespace: "parent_registry_namespace",
					},
					Tags: []string{"tag1", "tag2"},
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.store)
			}

			err := test.store.storeWildcardImage(test.name, test.image)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.store)
			}
		})
	}
}

func TestList(t *testing.T) {

	errContext := "(store::images::Store::List)"

	tests := []struct {
		desc              string
		store             *Store
		prepareAssertFunc func(*Store)
		assertFunc        func(*testing.T, *Store, []*image.Image)
		err               error
	}{
		{
			desc:  "Testing error listing images when images store is not initialized",
			store: &Store{},
			err:   errors.New(errContext, "To list images, store must be initialized"),
		},
		{
			desc:  "Testing list images",
			store: NewStore(render.NewMockImageRender()),
			prepareAssertFunc: func(s *Store) {

				s.store = []*image.Image{
					{
						Name:    "image_1",
						Version: "version_1",
					},
					{
						Name:    "image_2",
						Version: "version_2",
					},
					{
						Name:    "image_3",
						Version: "version_3",
					},
					{
						Name:    "image_1",
						Version: "version_12",
					},
				}
			},
			assertFunc: func(t *testing.T, s *Store, list []*image.Image) {
				expected := []*image.Image{
					{
						Name:    "image_1",
						Version: "version_1",
					},
					{
						Name:    "image_1",
						Version: "version_12",
					},
					{
						Name:    "image_2",
						Version: "version_2",
					},
					{
						Name:    "image_3",
						Version: "version_3",
					},
				}

				assert.Equal(t, expected, list)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.store)
			}

			list, err := test.store.List()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.store, list)
			}
		})
	}
}

func TestFindByName(t *testing.T) {

	errContext := "(store::images::Store::FindByName)"

	tests := []struct {
		desc  string
		store *Store
		name  string
		res   []*image.Image
		err   error
	}{
		{
			desc:  "Testing error finding images by name when list strucutres are not initialized - imageNameDefinitionVersionList",
			store: &Store{},
			err:   errors.New(errContext, "To find images by name into images store, list structures must be initialized"),
		},
		{
			desc: "Testing error finding images by name when list strucutres are not initialized - imageNameVersionRenderedVersionsList",
			store: &Store{
				imageNameDefinitionVersionList: make(map[string]map[string]struct{}),
			},
			err: errors.New(errContext, "To find images by name into images store, list structures must be initialized"),
		},
		{
			desc: "Testing error finding images by name when list strucutres are not initialized - imagesIndex",
			store: &Store{
				imageNameDefinitionVersionList:       make(map[string]map[string]struct{}),
				imageNameVersionRenderedVersionsList: make(map[string]map[string]map[string]struct{}),
			},
			err: errors.New(errContext, "To find images by name into images store, list structures must be initialized"),
		},
		{
			desc: "Testing find images by name into images store",
			store: &Store{
				render: render.NewMockImageRender(),
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}, "v2": struct{}{}},
					"image2": {"v1": struct{}{}},
				},

				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
						"v2": {"v2-a": struct{}{}, "v2-b": struct{}{}},
					},
					"image2": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
					},
				},
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
						"v2-a": &image.Image{
							Version: "v2-a",
						},
						"v2-b": &image.Image{
							Version: "v2-b",
						},
					},
					"image2": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
					},
				},
			},
			name: "image1",
			res: []*image.Image{
				{
					Version: "v1-a",
				},
				{
					Version: "v1-b",
				},
				{
					Version: "v2-a",
				},
				{
					Version: "v2-b",
				},
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			list, err := test.store.FindByName(test.name)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, list)
			}
		})
	}
}

func TestFind(t *testing.T) {

	errContext := "(store::images::Store::Find)"
	_ = errContext

	tests := []struct {
		desc    string
		store   *Store
		name    string
		version string
		res     []*image.Image
		err     error
	}{
		{
			desc:  "Testing error finding images by when list strucutres are not initialized - imageNameVersionRenderedVersionsList",
			store: &Store{},
			err:   errors.New(errContext, "To find images into images store, list structures must be initialized"),
		},
		{
			desc: "Testing error finding images by when list strucutres are not initialized - imagesIndex",
			store: &Store{
				imageNameVersionRenderedVersionsList: make(map[string]map[string]map[string]struct{}),
			},
			err: errors.New(errContext, "To find images into images store, list structures must be initialized"),
		},
		{
			desc: "Testing error finding images by when list strucutres are not initialized - imageWildcardIndex",
			store: &Store{
				imageNameVersionRenderedVersionsList: make(map[string]map[string]map[string]struct{}),
				imagesIndex:                          make(map[string]map[string]*image.Image),
			},
			err: errors.New(errContext, "To find images into images store, list structures must be initialized"),
		},
		{
			desc: "Testing find images into images store from rendered version list",
			store: &Store{
				render:             render.NewMockImageRender(),
				imageWildcardIndex: map[string]*image.Image{},
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}, "v2": struct{}{}},
					"image2": {"v1": struct{}{}},
				},

				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
						"v2": {"v2-a": struct{}{}, "v2-b": struct{}{}},
					},
					"image2": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
					},
				},
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
						"v2-a": &image.Image{
							Version: "v2-a",
						},
						"v2-b": &image.Image{
							Version: "v2-b",
						},
					},
					"image2": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
					},
				},
			},

			name:    "image1",
			version: "v1",
			res: []*image.Image{
				{
					Version: "v1-a",
				},
				{
					Version: "v1-b",
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing find images into images store from images index",
			store: &Store{
				render:             render.NewMockImageRender(),
				imageWildcardIndex: map[string]*image.Image{},
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}, "v2": struct{}{}},
					"image2": {"v1": struct{}{}},
				},

				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
						"v2": {"v2-a": struct{}{}, "v2-b": struct{}{}},
					},
					"image2": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
					},
				},
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
						"v2-a": &image.Image{
							Version: "v2-a",
						},
						"v2-b": &image.Image{
							Version: "v2-b",
						},
					},
					"image2": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
					},
				},
			},
			name:    "image1",
			version: "v1-a",
			res: []*image.Image{
				{
					Version: "v1-a",
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing find wildcard image into images store",
			store: &Store{
				render: render.NewMockImageRender(),
				imageWildcardIndex: map[string]*image.Image{
					"image1": {
						Version: "{{ .Version }}",
					},
				},
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}, "v2": struct{}{}},
					"image2": {"v1": struct{}{}},
				},
				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
						"v2": {"v2-a": struct{}{}, "v2-b": struct{}{}},
					},
					"image2": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
					},
				},
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
						"v2-a": &image.Image{
							Version: "v2-a",
						},
						"v2-b": &image.Image{
							Version: "v2-b",
						},
					},
					"image2": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
					},
				},
			},
			name:    "image1",
			version: "*",
			res: []*image.Image{
				{
					Version: "{{ .Version }}",
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing find unexisting image into images store",
			store: &Store{
				render:             render.NewMockImageRender(),
				imageWildcardIndex: map[string]*image.Image{},
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}, "v2": struct{}{}},
					"image2": {"v1": struct{}{}},
				},
				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
						"v2": {"v2-a": struct{}{}, "v2-b": struct{}{}},
					},
					"image2": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
					},
				},
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
						"v2-a": &image.Image{
							Version: "v2-a",
						},
						"v2-b": &image.Image{
							Version: "v2-b",
						},
					},
					"image2": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
					},
				},
			},
			name:    "image1",
			version: "unexisting",
			res:     []*image.Image{},
			err:     &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			list, err := test.store.Find(test.name, test.version)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, list)
			}
		})
	}
}

func TestFindGuaranteed(t *testing.T) {

	errContext := "(store::images::Store::FindGuaranteed)"

	tests := []struct {
		desc  string
		store *Store
		// findName          string
		// findVersion       string
		imageName         string
		imageVersion      string
		prepareAssertFunc func(*Store)
		assertFunc        func(*testing.T, *Store)
		res               []*image.Image
		err               error
	}{
		{
			desc: "Testing find images in a guaranteed finding mode",
			store: &Store{
				render: render.NewMockImageRender(),
				imageWildcardIndex: map[string]*image.Image{
					"image1": {
						Version: "{{ .Version }}",
					},
				},
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}, "v2": struct{}{}},
					"image2": {"v1": struct{}{}},
				},
				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
						"v2": {"v2-a": struct{}{}, "v2-b": struct{}{}},
					},
					"image2": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
					},
				},
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
						"v2-a": &image.Image{
							Version: "v2-a",
						},
						"v2-b": &image.Image{
							Version: "v2-b",
						},
					},
					"image2": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
					},
				},
			},
			prepareAssertFunc: func(*Store) {},
			// findName:     "image1",
			// findVersion:  "v1",
			imageName:    "image1",
			imageVersion: "v1",
			res: []*image.Image{
				{
					Version: "v1-a",
				},
				{
					Version: "v1-b",
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing find a wildcard image to be rendered",
			store: &Store{
				render: render.NewMockImageRender(),
				imageWildcardIndex: map[string]*image.Image{
					"image_wildcard": {
						Name:    "image_wildcard",
						Version: "{{ .Version }}-{{ .Parent.Version }}",
						Parent: &image.Image{
							Name:    "image_parent",
							Version: "version_parent",
						},
					},
				},
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}, "v2": struct{}{}},
					"image2": {"v1": struct{}{}},
				},
				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
						"v2": {"v2-a": struct{}{}, "v2-b": struct{}{}},
					},
					"image2": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
					},
				},
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
						"v2-a": &image.Image{
							Version: "v2-a",
						},
						"v2-b": &image.Image{
							Version: "v2-b",
						},
					},
					"image2": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
					},
				},
			},
			// findName:     "image_wildcard",
			// findVersion:  "wildcard",
			imageName:    "image_wildcard",
			imageVersion: "wildcard",
			prepareAssertFunc: func(s *Store) {
				s.render.(*render.MockImageRender).On("Render", "image_wildcard", "wildcard", &image.Image{
					Children:         []*image.Image{},
					Labels:           map[string]string{},
					PersistentLabels: map[string]string{},
					PersistentVars:   map[string]interface{}{},
					Name:             "image_wildcard",
					Parent: &image.Image{
						Name:    "image_parent",
						Version: "version_parent",
					},
					Tags:    []string{},
					Version: "{{ .Version }}-{{ .Parent.Version }}",
					Vars:    map[string]interface{}{},
				}).Return(
					&image.Image{
						Children:         []*image.Image{},
						Labels:           map[string]string{},
						PersistentLabels: map[string]string{},
						PersistentVars:   map[string]interface{}{},
						Name:             "image_wildcard",
						Parent: &image.Image{
							Name:    "image_parent",
							Version: "version_parent",
						},
						Version: "wildcard-version_parent",
						Tags:    []string{},
						Vars:    map[string]interface{}{},
					},
					nil,
				)
			},
			res: []*image.Image{
				{
					Children:         []*image.Image{},
					Labels:           map[string]string{},
					PersistentLabels: map[string]string{},
					PersistentVars:   map[string]interface{}{},
					Name:             "image_wildcard",
					Parent: &image.Image{
						Name:    "image_parent",
						Version: "version_parent",
					},
					Version: "wildcard-version_parent",
					Tags:    []string{},
					Vars:    map[string]interface{}{},
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing find the wildcard image",
			store: &Store{
				render: render.NewMockImageRender(),
				imageWildcardIndex: map[string]*image.Image{
					"image_wildcard": {
						Name:    "image_wildcard",
						Version: "{{ .Version }}-{{ .Parent.Version }}",
						Parent: &image.Image{
							Name:    "image_parent",
							Version: "version_parent",
						},
					},
				},
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}, "v2": struct{}{}},
					"image2": {"v1": struct{}{}},
				},
				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
						"v2": {"v2-a": struct{}{}, "v2-b": struct{}{}},
					},
					"image2": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
					},
				},
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
						"v2-a": &image.Image{
							Version: "v2-a",
						},
						"v2-b": &image.Image{
							Version: "v2-b",
						},
					},
					"image2": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
					},
				},
			},
			imageName:    "image_wildcard",
			imageVersion: "*",
			res: []*image.Image{
				{
					Name:    "image_wildcard",
					Version: "{{ .Version }}-{{ .Parent.Version }}",
					Parent: &image.Image{
						Name:    "image_parent",
						Version: "version_parent",
					},
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing find unexisting image",
			store: &Store{
				render: render.NewMockImageRender(),
				imageWildcardIndex: map[string]*image.Image{
					"image_wildcard": {
						Name:    "image_wildcard",
						Version: "{{ .Version }}-{{ .Parent.Version }}",
						Parent: &image.Image{
							Name:    "image_parent",
							Version: "version_parent",
						},
					},
				},
				imageNameDefinitionVersionList: map[string]map[string]struct{}{
					"image1": {"v1": struct{}{}, "v2": struct{}{}},
					"image2": {"v1": struct{}{}},
				},
				imageNameVersionRenderedVersionsList: map[string]map[string]map[string]struct{}{
					"image1": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
						"v2": {"v2-a": struct{}{}, "v2-b": struct{}{}},
					},
					"image2": {
						"v1": {"v1-a": struct{}{}, "v1-b": struct{}{}},
					},
				},
				imagesIndex: map[string]map[string]*image.Image{
					"image1": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
						"v2-a": &image.Image{
							Version: "v2-a",
						},
						"v2-b": &image.Image{
							Version: "v2-b",
						},
					},
					"image2": {
						"v1-a": &image.Image{
							Version: "v1-a",
						},
						"v1-b": &image.Image{
							Version: "v1-b",
						},
					},
				},
			},
			// findName:     "image",
			// findVersion:  "unexisting",
			imageName:    "image",
			imageVersion: "unexisting",
			err:          errors.New(errContext, "Image 'image:unexisting' does not exist on the store"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.store)
			}

			list, err := test.store.FindGuaranteed(test.imageName, test.imageVersion)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, list)
			}
		})
	}
}

func TestFindWildcardImage(t *testing.T) {
	errContext := "(store::images::Store::FindWildcardImage)"

	tests := []struct {
		desc  string
		store *Store
		name  string
		res   *image.Image
		err   error
	}{
		{
			desc: "Testing find wild card image error when wildcard index has not been initialized on images store",
			store: &Store{
				store: []*image.Image{},
			},
			err: errors.New(errContext, "To find a wildcard image, Wildcard index must be initialized"),
		},
		{
			desc: "Testing find the wildcard image into images store",
			store: &Store{
				store: []*image.Image{},
				imageWildcardIndex: map[string]*image.Image{
					"image": {
						Name:    "image_wildcard",
						Version: "{{ .Version }}",
					},
				},
			},
			name: "image",
			res: &image.Image{
				Name:    "image_wildcard",
				Version: "{{ .Version }}",
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			i, err := test.store.FindWildcardImage(test.name)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, i)
			}
		})
	}
}

func TestGenerateImageFromWildcard(t *testing.T) {
	errContext := "(store::images::Store::GenerateImageFromWildcard)"

	tests := []struct {
		desc              string
		store             *Store
		image             *image.Image
		name              string
		version           string
		prepareAssertFunc func(*Store)
		assertFunc        func(*testing.T, *Store, *image.Image)
		err               error
	}{
		{
			desc:  "Testing error when generating an image when wildcard image is nil",
			store: &Store{},
			err:   errors.New(errContext, "Provided wildcard image is nil"),
		},
		{
			desc:    "Testing when generating an image from wildcard image",
			store:   NewStore(render.NewMockImageRender()),
			name:    "image_wildcard",
			version: "wildcard",
			image: &image.Image{
				Name:    "image_wildcard",
				Version: "{{ .Version }}-{{ .Parent.Version }}",
				Parent: &image.Image{
					Name:    "image_parent",
					Version: "{{ .Parent.Version }}",
					Parent: &image.Image{
						Name:    "image_grandparent",
						Version: "version_grandparent",
					},
				},
			},
			prepareAssertFunc: func(s *Store) {
				s.imageWildcardIndex = map[string]*image.Image{
					"image_wildcard": {
						Name:    "image_wildcard",
						Version: "{{ .Version }}-{{ .Parent.Version }}",
						Parent: &image.Image{
							Name:    "image_parent",
							Version: "version_parent",
						},
					},
					"image_parent": {
						Name:    "image_parent",
						Version: "{{ .Parent.Version }}",
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
					},
				}
				s.render.(*render.MockImageRender).On("Render", "image_parent", "wildcard", &image.Image{
					Children:         []*image.Image{},
					Labels:           map[string]string{},
					Name:             "image_parent",
					PersistentLabels: map[string]string{},
					PersistentVars:   map[string]interface{}{},
					Parent: &image.Image{
						Name:    "image_grandparent",
						Version: "version_grandparent",
					},
					Tags:    []string{},
					Vars:    map[string]interface{}{},
					Version: "{{ .Parent.Version }}",
				}).Return(
					&image.Image{
						Children:         []*image.Image{},
						Name:             "image_parent",
						Labels:           map[string]string{},
						PersistentLabels: map[string]string{},
						PersistentVars:   map[string]interface{}{},
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
						Tags:    []string{},
						Vars:    map[string]interface{}{},
						Version: "version_grandparent",
					},
					nil,
				)

				s.render.(*render.MockImageRender).On("Render", "image_wildcard", "wildcard", &image.Image{
					Children:         []*image.Image{},
					Labels:           map[string]string{},
					PersistentLabels: map[string]string{},
					PersistentVars:   map[string]interface{}{},
					Name:             "image_wildcard",
					Parent: &image.Image{
						Children:         []*image.Image{},
						Labels:           map[string]string{},
						Name:             "image_parent",
						PersistentLabels: map[string]string{},
						PersistentVars:   map[string]interface{}{},
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
						Tags:    []string{},
						Vars:    map[string]interface{}{},
						Version: "version_grandparent",
					},
					Tags:    []string{},
					Vars:    map[string]interface{}{},
					Version: "{{ .Version }}-{{ .Parent.Version }}",
				}).Return(
					&image.Image{
						Children:         []*image.Image{},
						Labels:           map[string]string{},
						Name:             "image_wildcard",
						PersistentLabels: map[string]string{},
						PersistentVars:   map[string]interface{}{},
						Parent: &image.Image{
							Children:         []*image.Image{},
							Labels:           map[string]string{},
							Name:             "image_parent",
							PersistentLabels: map[string]string{},
							PersistentVars:   map[string]interface{}{},
							Parent: &image.Image{
								Name:    "image_grandparent",
								Version: "version_grandparent",
							},
							Tags:    []string{},
							Vars:    map[string]interface{}{},
							Version: "version_grandparent",
						},
						Tags:    []string{},
						Vars:    map[string]interface{}{},
						Version: "wildcard-version_grandparent",
					},
					nil,
				)
			},
			assertFunc: func(t *testing.T, s *Store, i *image.Image) {
				expected := &image.Image{
					Children:         []*image.Image{},
					Labels:           map[string]string{},
					Name:             "image_wildcard",
					PersistentLabels: map[string]string{},
					PersistentVars:   map[string]interface{}{},
					Parent: &image.Image{
						Children:         []*image.Image{},
						Labels:           map[string]string{},
						Name:             "image_parent",
						PersistentLabels: map[string]string{},
						PersistentVars:   map[string]interface{}{},
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
						Tags:    []string{},
						Vars:    map[string]interface{}{},
						Version: "version_grandparent",
					},
					Tags:    []string{},
					Vars:    map[string]interface{}{},
					Version: "wildcard-version_grandparent",
				}

				assert.Equal(t, expected, i)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.store)
			}

			i, err := test.store.GenerateImageFromWildcard(test.image, test.name, test.version)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.store, i)
			}
		})
	}
}
