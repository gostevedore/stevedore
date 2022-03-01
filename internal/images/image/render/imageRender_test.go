package render

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	domainimage "github.com/gostevedore/stevedore/internal/images/image"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {

	errContext := "(render::Render)"
	_ = errContext

	tests := []struct {
		desc    string
		render  *ImageRender
		name    string
		version string
		image   *domainimage.Image
		res     *domainimage.Image
		err     error
	}{
		{
			desc:    "Testing render domain image",
			render:  &ImageRender{},
			name:    "image_name",
			version: "Image_version",
			image: &domainimage.Image{
				Name:              "{{.Name}}-{{.Parent.Name}}",
				Version:           "{{.Version}}-{{.Parent.Version}}",
				RegistryNamespace: "{{.Parent.RegistryNamespace}}",
				Parent: &domainimage.Image{
					Name:              "parent_name",
					Version:           "parent_version",
					RegistryNamespace: "parent_registry_namespace",
				},
			},
			res: &domainimage.Image{
				Name:              "image_name-parent_name",
				Version:           "Image_version-parent_version",
				RegistryNamespace: "parent_registry_namespace",
				Labels:            map[string]string{},
				PersistentVars:    map[string]interface{}{},
				Tags:              []string{},
				Vars:              map[string]interface{}{},
				Parent: &domainimage.Image{
					Name:              "parent_name",
					Version:           "parent_version",
					RegistryNamespace: "parent_registry_namespace",
				},
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing render domain image using grand parent details",
			render:  &ImageRender{},
			name:    "image_name",
			version: "Image_version",
			image: &domainimage.Image{
				Name:              "{{.Name}}-{{.Parent.Name}}",
				Version:           "{{.Version}}-{{.Parent.Version}}",
				RegistryNamespace: "{{.Parent.Parent.RegistryNamespace}}",
				Parent: &domainimage.Image{
					Name:    "parent_name",
					Version: "parent_version",
					Parent: &domainimage.Image{
						Name:              "parent_parent_name",
						Version:           "parent_parent_version",
						RegistryNamespace: "parent_parent_registry_namespace",
					},
				},
			},
			res: &domainimage.Image{
				Name:              "image_name-parent_name",
				Version:           "Image_version-parent_version",
				RegistryNamespace: "parent_parent_registry_namespace",
				Labels:            map[string]string{},
				PersistentVars:    map[string]interface{}{},
				Tags:              []string{},
				Vars:              map[string]interface{}{},
				Parent: &domainimage.Image{
					Name:    "parent_name",
					Version: "parent_version",
					Parent: &domainimage.Image{
						Name:              "parent_parent_name",
						Version:           "parent_parent_version",
						RegistryNamespace: "parent_parent_registry_namespace",
					},
				},
			},
			err: &errors.Error{},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.render.Render(test.name, test.version, test.image)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.image)
			}

		})
	}
}
