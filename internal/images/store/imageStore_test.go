package store

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
	"github.com/gostevedore/stevedore/internal/images/image/render"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {

	errContext := "(store::Store)"

	tests := []struct {
		desc              string
		store             *ImageStore
		name              string
		version           string
		image             *image.Image
		err               error
		prepareAssertFunc func(*ImageStore, *image.Image)
		assertFunc        func(*testing.T, *ImageStore)
	}{
		{
			desc:  "Testing error when render is not defined",
			store: NewImageStore(nil),
			err:   errors.New(errContext, "To add an image to the store an image render is required"),
		},
		{
			desc:  "Testing error when name is not defined",
			store: NewImageStore(render.NewMockImageRender()),
			err:   errors.New(errContext, "To add an image to the store a name is required"),
		},
		{
			desc:  "Testing error when version is not defined",
			store: NewImageStore(render.NewMockImageRender()),
			name:  "image_name",
			err:   errors.New(errContext, "To add an image to the store a version is required"),
		},
		{
			desc:    "Testing error when image is not defined",
			store:   NewImageStore(render.NewMockImageRender()),
			name:    "image_name",
			version: "image_version",
			err:     errors.New(errContext, "To add an image to the store an image is required"),
		},
		{
			desc:    "Testing add a image to an empty store",
			store:   NewImageStore(render.NewMockImageRender()),
			name:    "image_name",
			version: "image_version",
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
			prepareAssertFunc: func(s *ImageStore, i *image.Image) {
				s.render.(*render.MockImageRender).On("Render", "image_name", "image_version", i).Return(
					&image.Image{
						Name:              "image_name-parent_name",
						Version:           "image_version-parent_version",
						RegistryNamespace: "parent_registry_namespace",
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
			assertFunc: func(t *testing.T, s *ImageStore) {
				assert.Equal(t, 1, len(s.store), "Unexpected number of images in the store")
				assert.Equal(t, 1, len(s.imageNameVersionIndex), "Unexpected number of images in the index")
				assert.Equal(t, 0, len(s.imageWildcardIndex), "Unexpected number of images in the wildcard index")
				assert.Equal(t, 4, len(s.imageNameVersionIndex["image_name"]), "Unexpected number of 'image_name' items")
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing add a wildcard image to an empty store",
			store:   NewImageStore(render.NewMockImageRender()),
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
			prepareAssertFunc: func(s *ImageStore, i *image.Image) {},
			assertFunc: func(t *testing.T, s *ImageStore) {

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

func TestStoreImage(t *testing.T) {

	errContext := "(store::storeImage)"

	tests := []struct {
		desc              string
		store             *ImageStore
		name              string
		version           string
		image             *image.Image
		err               error
		prepareAssertFunc func(s *ImageStore)
		assertFunc        func(t *testing.T, s *ImageStore)
	}{
		{
			desc:       "Testing store an image that is nil",
			store:      NewImageStore(render.NewMockImageRender()),
			name:       "image_name",
			version:    "image_version",
			image:      nil,
			err:        errors.New(errContext, "Provided image for 'image_name:image_version' is nil"),
			assertFunc: func(t *testing.T, s *ImageStore) {},
		},
		{
			desc:    "Testing store an image",
			store:   NewImageStore(render.NewMockImageRender()),
			name:    "image_name",
			version: "image_version",
			image: &image.Image{
				Name:              "image_name-parent_name",
				Version:           "image_version-parent_version",
				RegistryNamespace: "parent_registry_namespace",
				Parent: &image.Image{
					Name:              "parent_name",
					Version:           "parent_version",
					RegistryNamespace: "parent_registry_namespace",
				},
				Tags: []string{"tag1", "tag2"},
			},
			err: &errors.Error{},
			assertFunc: func(t *testing.T, s *ImageStore) {
				expected := &image.Image{
					Name:              "image_name-parent_name",
					Version:           "image_version-parent_version",
					RegistryNamespace: "parent_registry_namespace",
					Parent: &image.Image{
						Name:              "parent_name",
						Version:           "parent_version",
						RegistryNamespace: "parent_registry_namespace",
					},
					Tags: []string{"tag1", "tag2"},
				}

				subImageNameVersionIndex, exist := s.imageNameVersionIndex["image_name"]
				assert.True(t, exist, "Image name is not on the index")
				image, exist := subImageNameVersionIndex["image_version"]
				assert.True(t, exist, "Image version is not on the index")
				assert.Equal(t, image, expected, "Unexpected image in the index")
			},
		},
		{
			desc:    "Testing store an existing image",
			store:   NewImageStore(render.NewMockImageRender()),
			name:    "image_name",
			version: "image_version",
			image: &image.Image{
				Name:              "image_name-parent_name",
				Version:           "image_version-parent_version",
				RegistryNamespace: "parent_registry_namespace",
				Parent: &image.Image{
					Name:              "parent_name",
					Version:           "parent_version",
					RegistryNamespace: "parent_registry_namespace",
				},
				Tags: []string{"tag1", "tag2"},
			},
			err: errors.New(errContext, "Image 'image_name:image_version' already exists"),
			prepareAssertFunc: func(s *ImageStore) {
				s.storeImage("image_name", "image_version", &image.Image{
					Name:              "image_name-parent_name",
					Version:           "image_version-parent_version",
					RegistryNamespace: "parent_registry_namespace",
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

			err := test.store.storeImage(test.name, test.version, test.image)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.store)
			}
		})
	}
}

func TestStoreWildcardImage(t *testing.T) {
	errContext := "(store::storeWildcardImage)"

	tests := []struct {
		desc              string
		store             *ImageStore
		name              string
		image             *image.Image
		err               error
		prepareAssertFunc func(s *ImageStore)
		assertFunc        func(t *testing.T, s *ImageStore)
	}{
		{
			desc:       "Testing store an image that is nil",
			store:      NewImageStore(render.NewMockImageRender()),
			name:       "image_name",
			image:      nil,
			err:        errors.New(errContext, "Provided wildcard image for 'image_name' is nil"),
			assertFunc: func(t *testing.T, s *ImageStore) {},
		},

		{
			desc:  "Testing store an image",
			store: NewImageStore(render.NewMockImageRender()),
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
			assertFunc: func(t *testing.T, s *ImageStore) {
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
			store: NewImageStore(render.NewMockImageRender()),
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
			prepareAssertFunc: func(s *ImageStore) {
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

	errContext := "(store::List)"

	tests := []struct {
		desc              string
		store             *ImageStore
		prepareAssertFunc func(*ImageStore)
		assertFunc        func(*testing.T, *ImageStore, []*image.Image)
		err               error
	}{
		{
			desc:  "Testing error listing images when store is not initialized",
			store: &ImageStore{},
			err:   errors.New(errContext, "Store has not been initialized"),
		},
		{
			desc:  "Testing list images",
			store: NewImageStore(render.NewMockImageRender()),
			prepareAssertFunc: func(s *ImageStore) {

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
			assertFunc: func(t *testing.T, s *ImageStore, list []*image.Image) {
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
	errContext := "(store::FindByName)"

	tests := []struct {
		desc              string
		store             *ImageStore
		name              string
		prepareAssertFunc func(*ImageStore)
		assertFunc        func(*testing.T, *ImageStore, []*image.Image)
		err               error
	}{
		{
			desc:  "Testing error finding images by name when store is not initialized",
			store: &ImageStore{},
			err:   errors.New(errContext, "Store has not been initialized"),
		},
		{
			desc:  "Testing find images by name",
			store: NewImageStore(render.NewMockImageRender()),
			name:  "image_1",
			prepareAssertFunc: func(s *ImageStore) {

				s.imageNameVersionIndex = map[string]map[string]*image.Image{
					"image_1": {
						"version_1": &image.Image{
							Name:    "image_1",
							Version: "version_1",
						},
						"version_12": &image.Image{
							Name:    "image_1",
							Version: "version_12",
						},
					},
					"image_2": {
						"version_2": &image.Image{
							Name:    "image_2",
							Version: "version_2",
						},
					},
					"image_3": {
						"version_3": &image.Image{
							Name:    "image_3",
							Version: "version_3",
						},
					},
				}
			},
			assertFunc: func(t *testing.T, s *ImageStore, list []*image.Image) {
				expected := []*image.Image{
					{
						Name:    "image_1",
						Version: "version_1",
					},
					{
						Name:    "image_1",
						Version: "version_12",
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

			list, err := test.store.FindByName(test.name)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.store, list)
			}
		})
	}
}

func TestFind(t *testing.T) {

	errContext := "(store::Find)"

	tests := []struct {
		desc              string
		store             *ImageStore
		name              string
		version           string
		prepareAssertFunc func(*ImageStore)
		assertFunc        func(*testing.T, *ImageStore, *image.Image)
		err               error
	}{
		{
			desc:  "Testing error finding an image when store is not initialized",
			store: &ImageStore{},
			err:   errors.New(errContext, "Store has not been initialized"),
		},
		{
			desc:    "Testing find an image",
			store:   NewImageStore(render.NewMockImageRender()),
			name:    "image_1",
			version: "version_1",
			prepareAssertFunc: func(s *ImageStore) {
				s.imageNameVersionIndex = map[string]map[string]*image.Image{
					"image_1": {
						"version_1": &image.Image{
							Name:    "image_1",
							Version: "version_1",
						},
						"version_12": &image.Image{
							Name:    "image_1",
							Version: "version_12",
						},
					},
					"image_2": {
						"version_2": &image.Image{
							Name:    "image_2",
							Version: "version_2",
						},
					},
					"image_3": {
						"version_3": &image.Image{
							Name:    "image_3",
							Version: "version_3",
						},
					},
				}
			},
			assertFunc: func(t *testing.T, s *ImageStore, i *image.Image) {
				expected := &image.Image{
					Name:    "image_1",
					Version: "version_1",
				}

				assert.Equal(t, expected, i)
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing find the wildcard image",
			store:   NewImageStore(render.NewMockImageRender()),
			name:    "image",
			version: "*",
			prepareAssertFunc: func(s *ImageStore) {

				s.imageWildcardIndex = map[string]*image.Image{
					"image": {
						Name:    "image_wildcard",
						Version: "{{ .Version }}",
					},
				}

			},
			assertFunc: func(t *testing.T, s *ImageStore, i *image.Image) {
				expected := &image.Image{
					Name:    "image_wildcard",
					Version: "{{ .Version }}",
				}

				assert.Equal(t, expected, i)
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing find unexisting image",
			store:   NewImageStore(render.NewMockImageRender()),
			name:    "image",
			version: "unexisting",
			prepareAssertFunc: func(s *ImageStore) {
				s.imageWildcardIndex = map[string]*image.Image{}
			},
			assertFunc: func(t *testing.T, s *ImageStore, i *image.Image) {
				assert.Nil(t, i)
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

			i, err := test.store.Find(test.name, test.version)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.store, i)
			}
		})
	}
}

func TestFindGuaranteed(t *testing.T) {

	errContext := "(store::FindGuaranteed)"

	tests := []struct {
		desc              string
		store             *ImageStore
		findName          string
		findVersion       string
		imageName         string
		imageVersion      string
		prepareAssertFunc func(*ImageStore)
		assertFunc        func(*testing.T, *ImageStore, *image.Image)
		err               error
	}{
		{
			desc:  "Testing error finding an image when store is not initialized",
			store: &ImageStore{},
			err:   errors.New(errContext, "Store has not been initialized"),
		},
		{
			desc:         "Testing find an image",
			store:        NewImageStore(render.NewMockImageRender()),
			findName:     "image_1",
			findVersion:  "version_1",
			imageName:    "image_1",
			imageVersion: "version_1",
			prepareAssertFunc: func(s *ImageStore) {
				s.imageNameVersionIndex = map[string]map[string]*image.Image{
					"image_1": {
						"version_1": &image.Image{
							Name:    "image_1",
							Version: "version_1",
						},
						"version_12": &image.Image{
							Name:    "image_1",
							Version: "version_12",
						},
					},
					"image_2": {
						"version_2": &image.Image{
							Name:    "image_2",
							Version: "version_2",
						},
					},
					"image_3": {
						"version_3": &image.Image{
							Name:    "image_3",
							Version: "version_3",
						},
					},
				}
			},
			assertFunc: func(t *testing.T, s *ImageStore, i *image.Image) {
				expected := &image.Image{
					Name:    "image_1",
					Version: "version_1",
				}

				assert.Equal(t, expected, i)
			},
			err: &errors.Error{},
		},
		{
			desc:         "Testing find a wildcard image to be rendered",
			store:        NewImageStore(render.NewMockImageRender()),
			findName:     "image_wildcard",
			findVersion:  "wildcard",
			imageName:    "image_wildcard",
			imageVersion: "wildcard",
			prepareAssertFunc: func(s *ImageStore) {

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
					Name:    "image_parent",
					Version: "{{ .Parent.Version }}",
					Parent: &image.Image{
						Name:    "image_grandparent",
						Version: "version_grandparent",
					},
					Labels: map[string]string{},
					Vars:   map[string]interface{}{},
					Tags:   []string{},
				}).Return(
					&image.Image{
						Name:    "image_parent",
						Version: "version_grandparent",
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
						Labels: map[string]string{},
						Vars:   map[string]interface{}{},
						Tags:   []string{},
					},
					nil,
				)

				s.render.(*render.MockImageRender).On("Render", "image_wildcard", "wildcard", &image.Image{
					Name:    "image_wildcard",
					Version: "{{ .Version }}-{{ .Parent.Version }}",
					Parent: &image.Image{
						Name:    "image_parent",
						Version: "version_grandparent",
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
						Labels: map[string]string{},
						Vars:   map[string]interface{}{},
						Tags:   []string{},
					},
					Labels: map[string]string{},
					Vars:   map[string]interface{}{},
					Tags:   []string{},
				}).Return(
					&image.Image{
						Name:    "image_wildcard",
						Version: "wildcard-version_grandparent",
						Parent: &image.Image{
							Name:    "image_parent",
							Version: "version_grandparent",
							Parent: &image.Image{
								Name:    "image_grandparent",
								Version: "version_grandparent",
							},
							Labels: map[string]string{},
							Vars:   map[string]interface{}{},
							Tags:   []string{},
						},
						Labels: map[string]string{},
						Vars:   map[string]interface{}{},
						Tags:   []string{},
					},
					nil,
				)
			},
			assertFunc: func(t *testing.T, s *ImageStore, i *image.Image) {
				expected := &image.Image{
					Name:    "image_wildcard",
					Version: "wildcard-version_grandparent",
					Parent: &image.Image{
						Name:    "image_parent",
						Version: "version_grandparent",
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
						Labels: map[string]string{},
						Vars:   map[string]interface{}{},
						Tags:   []string{},
					},
					Labels: map[string]string{},
					Vars:   map[string]interface{}{},
					Tags:   []string{},
				}

				assert.Equal(t, expected, i)
			},
			err: &errors.Error{},
		},
		{
			desc:        "Testing find the wildcard image",
			store:       NewImageStore(render.NewMockImageRender()),
			findName:    "image",
			findVersion: "*",
			prepareAssertFunc: func(s *ImageStore) {

				s.imageWildcardIndex = map[string]*image.Image{
					"image": {
						Name:    "image_wildcard",
						Version: "{{ .Version }}",
					},
				}

			},
			assertFunc: func(t *testing.T, s *ImageStore, i *image.Image) {
				expected := &image.Image{
					Name:    "image_wildcard",
					Version: "{{ .Version }}",
				}

				assert.Equal(t, expected, i)
			},
			err: &errors.Error{},
		},
		{
			desc:         "Testing find unexisting image",
			store:        NewImageStore(render.NewMockImageRender()),
			findName:     "image",
			findVersion:  "unexisting",
			imageName:    "image",
			imageVersion: "unexisting",
			prepareAssertFunc: func(s *ImageStore) {
				s.imageWildcardIndex = map[string]*image.Image{}
			},
			assertFunc: func(t *testing.T, s *ImageStore, i *image.Image) {
				assert.Nil(t, i)
			},
			err: errors.New(errContext, "Image 'image:unexisting' does not exist on the store"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.store)
			}

			i, err := test.store.FindGuaranteed(test.findName, test.findVersion, test.imageName, test.imageVersion)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertFunc(t, test.store, i)
			}
		})
	}
}

func TestFindWildcardImage(t *testing.T) {
	errContext := "(store::FindWildcardImage)"

	tests := []struct {
		desc  string
		store *ImageStore
		name  string
		res   *image.Image
		err   error
	}{

		{
			desc:  "Testing find the wildcard image",
			store: &ImageStore{},

			err: errors.New(errContext, "Store has not been initialized"),
		},
		{
			desc: "Testing find the wildcard image",
			store: &ImageStore{
				store: []*image.Image{},
			},
			err: errors.New(errContext, "Wildcard index has not been initialized"),
		},
		{
			desc: "Testing find the wildcard image",
			store: &ImageStore{
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
	errContext := "(store::GenerateImageFromWildcard)"

	tests := []struct {
		desc              string
		store             *ImageStore
		image             *image.Image
		name              string
		version           string
		prepareAssertFunc func(*ImageStore)
		assertFunc        func(*testing.T, *ImageStore, *image.Image)
		err               error
	}{
		{
			desc:  "Testing error generating an image when wildcard image is nil",
			store: &ImageStore{},
			err:   errors.New(errContext, "Provided wildcard image is nil"),
		},
		{
			desc:  "Testing error generating an image when store is not initialized",
			store: &ImageStore{},
			image: &image.Image{
				Parent: &image.Image{},
			},
			err: errors.New(errContext, "Store has not been initialized"),
		},
		{
			desc:    "Testing  generating an image from wildcard image",
			store:   NewImageStore(render.NewMockImageRender()),
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
			prepareAssertFunc: func(s *ImageStore) {

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
					Name:    "image_parent",
					Version: "{{ .Parent.Version }}",
					Parent: &image.Image{
						Name:    "image_grandparent",
						Version: "version_grandparent",
					},
					Labels: map[string]string{},
					Vars:   map[string]interface{}{},
					Tags:   []string{},
				}).Return(
					&image.Image{
						Name:    "image_parent",
						Version: "version_grandparent",
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
						Labels: map[string]string{},
						Vars:   map[string]interface{}{},
						Tags:   []string{},
					},
					nil,
				)

				s.render.(*render.MockImageRender).On("Render", "image_wildcard", "wildcard", &image.Image{
					Name:    "image_wildcard",
					Version: "{{ .Version }}-{{ .Parent.Version }}",
					Parent: &image.Image{
						Name:    "image_parent",
						Version: "version_grandparent",
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
						Labels: map[string]string{},
						Vars:   map[string]interface{}{},
						Tags:   []string{},
					},
					Labels: map[string]string{},
					Vars:   map[string]interface{}{},
					Tags:   []string{},
				}).Return(
					&image.Image{
						Name:    "image_wildcard",
						Version: "wildcard-version_grandparent",
						Parent: &image.Image{
							Name:    "image_parent",
							Version: "version_grandparent",
							Parent: &image.Image{
								Name:    "image_grandparent",
								Version: "version_grandparent",
							},
							Labels: map[string]string{},
							Vars:   map[string]interface{}{},
							Tags:   []string{},
						},
						Labels: map[string]string{},
						Vars:   map[string]interface{}{},
						Tags:   []string{},
					},
					nil,
				)
			},
			assertFunc: func(t *testing.T, s *ImageStore, i *image.Image) {
				expected := &image.Image{
					Name:    "image_wildcard",
					Version: "wildcard-version_grandparent",
					Parent: &image.Image{
						Name:    "image_parent",
						Version: "version_grandparent",
						Parent: &image.Image{
							Name:    "image_grandparent",
							Version: "version_grandparent",
						},
						Labels: map[string]string{},
						Vars:   map[string]interface{}{},
						Tags:   []string{},
					},
					Labels: map[string]string{},
					Vars:   map[string]interface{}{},
					Tags:   []string{},
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
