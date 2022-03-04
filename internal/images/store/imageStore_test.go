package store

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
	"github.com/gostevedore/stevedore/internal/images/image/render"
	"github.com/stretchr/testify/assert"
)

func TestAddImage(t *testing.T) {

	errContext := "(store::AddImage)"

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

			err := test.store.AddImage(test.name, test.version, test.image)
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

}

func TestFindByName(t *testing.T) {

}

func TestFind(t *testing.T) {

}
