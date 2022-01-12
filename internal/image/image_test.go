package image

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestTesting(t *testing.T) {

	errContext := "(image::NewImage)"

	tests := []struct {
		desc              string
		name              string
		version           string
		registryHost      string
		registryNamesapce string
		res               *Image
		err               error
	}{
		{
			desc: "Testing error no name provides",
			err:  errors.New(errContext, "Image could not be created\n\tinvalid reference format"),
		},
		{
			desc:         "Testing error when invalid registy host is provided",
			name:         "image",
			registryHost: "registry",
			err:          errors.New(errContext, "Registry host name must by a FQDN"),
		},
		{
			desc: "Testing create image providing only a name",
			name: "image",
			res: &Image{
				Name:              "image",
				Version:           "latest",
				RegistryHost:      "docker.io",
				RegistryNamespace: "library",
			},
		},
		{
			desc:              "Testing create image providing all the parameters",
			name:              "image",
			version:           "version",
			registryHost:      "registry.test",
			registryNamesapce: "namespace",
			res: &Image{
				Name:              "image",
				Version:           "version",
				RegistryHost:      "registry.test",
				RegistryNamespace: "namespace",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			image, err := NewImage(test.name, test.version, test.registryHost, test.registryNamesapce)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, image)
			}

		})
	}

}
