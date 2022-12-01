package name

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/stretchr/testify/assert"
)

func TestGenerateName(t *testing.T) {

	errContext := "(name::images::DockerNormalizedReferenceName::GenerateName)"

	tests := []struct {
		desc  string
		name  *DockerNormalizedReferenceName
		image *image.Image
		res   string
		err   error
	}{
		{
			desc: "Testing generate image name with docker normalized named providing name and version",
			name: NewDockerNormalizedReferenceName(),
			image: &image.Image{
				Name:    "name",
				Version: "version",
			},
			res: "docker.io/library/name:version",
		},
		{
			desc: "Testing generate image name with docker normalized named providing registry host, name and version",
			name: NewDockerNormalizedReferenceName(),
			image: &image.Image{
				Name:         "name",
				Version:      "version",
				RegistryHost: "registry.test",
			},
			res: "registry.test/name:version",
		},
		{
			desc: "Testing generate image name with docker normalized named providing registry namespace, name and version",
			name: NewDockerNormalizedReferenceName(),
			image: &image.Image{
				Name:              "name",
				Version:           "version",
				RegistryNamespace: "stable",
			},
			res: "docker.io/stable/name:version",
		},
		{
			desc: "Testing generate image name with docker normalized named",
			name: NewDockerNormalizedReferenceName(),
			image: &image.Image{
				Name:              "name",
				Version:           "version",
				RegistryHost:      "registry.test",
				RegistryNamespace: "stable",
			},
			res: "registry.test/stable/name:version",
		},
		{
			desc:  "Testing error on generating and image name with docker normalized named when image name is not defined",
			name:  NewDockerNormalizedReferenceName(),
			image: &image.Image{},
			err:   errors.New(errContext, "Name could not be generated because image name is undefined"),
		},
		{
			desc: "Testing error on generating and image name with docker normalized named when image version is not defined",
			name: NewDockerNormalizedReferenceName(),
			image: &image.Image{
				Name: "name",
			},
			err: errors.New(errContext, "Name could not be generated because image version is undefined"),
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
