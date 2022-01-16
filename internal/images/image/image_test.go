package image

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {

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
			err:  errors.New(errContext, "Image could not be parsed\n\tinvalid reference format"),
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

func TestDockerNormalizedNamed(t *testing.T) {
	errContext := "(image::DockerNormalizedNamed)"

	tests := []struct {
		desc  string
		res   string
		image *Image
		err   error
	}{
		{
			desc:  "Testing error no name provided",
			err:   errors.New(errContext, "Image name is empty"),
			image: &Image{},
		},
		{
			desc: "Testing error no version is provided",
			err:  errors.New(errContext, "Image version is empty"),
			image: &Image{
				Name: "image",
			},
		},
		{
			desc: "Testing error no registry host is provided",
			err:  errors.New(errContext, "Registry host is empty"),
			image: &Image{
				Name:    "image",
				Version: "version",
			},
		},
		{
			desc: "Testing error no registry namespace is provided",
			err:  errors.New(errContext, "Registry namespace is empty"),
			image: &Image{
				Name:         "image",
				Version:      "version",
				RegistryHost: "registry.test",
			},
		},
		{
			desc: "Testing docekr normalized name",
			err:  &errors.Error{},
			image: &Image{
				Name:              "image",
				Version:           "version",
				RegistryHost:      "registry.test",
				RegistryNamespace: "namespace",
			},
			res: "registry.test/namespace/image:version",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			name, err := test.image.DockerNormalizedNamed()

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, name)
			}
		})
	}

}
