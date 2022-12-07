package name

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/stretchr/testify/assert"
)

func TestGenerateName(t *testing.T) {

	errContext := "(name::images::DefaultReferenceName::GenerateName)"

	tests := []struct {
		desc  string
		name  *DefaultReferenceName
		image *image.Image
		res   string
		err   error
	}{
		{
			desc: "Testing generate image name with default named",
			name: NewDefaultReferenceName(),
			image: &image.Image{
				Name:              "name",
				Version:           "version",
				RegistryNamespace: "namespace",
				RegistryHost:      "registry.test",
			},
			res: "registry.test/namespace/name:version",
		},
		{
			desc:  "Testing error on generating and image name with default named when image name is not defined",
			name:  NewDefaultReferenceName(),
			image: &image.Image{},
			err:   errors.New(errContext, "Image reference name can not be generated because image name is undefined"),
		},
		{
			desc: "Testing error on generating and image name with default named when image version is not defined",
			name: NewDefaultReferenceName(),
			image: &image.Image{
				Name: "name",
			},
			err: errors.New(errContext, "Image reference name can not be generated because image version is undefined"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			name, err := test.name.GenerateName(test.image)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, name)
			}
		})
	}
}
