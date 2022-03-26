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
		options           []OptionFunc
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
			desc:              "Testing create image providing all parameters",
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
		{
			desc:              "Testing create image with options",
			name:              "image",
			version:           "version",
			registryHost:      "registry.test",
			registryNamesapce: "namespace",
			res: &Image{
				Builder:           "builder",
				Children:          []*Image{{Name: "child"}},
				Labels:            map[string]string{"label": "value"},
				Name:              "image",
				Parent:            &Image{Name: "parent"},
				PersistentLabels:  map[string]string{"plabel": "pvalue"},
				PersistentVars:    map[string]interface{}{"pvar": "pvalue"},
				RegistryHost:      "registry.test",
				RegistryNamespace: "namespace",
				Tags:              []string{"tag"},
				Vars:              map[string]interface{}{"var": "value"},
				Version:           "version",
			},
			options: []OptionFunc{
				WithBuilder("builder"),
				WithChildren([]*Image{{Name: "child"}}...),
				WithLabels(map[string]string{"label": "value"}),
				WithParent(&Image{Name: "parent"}),
				WithPersistentLabels(map[string]string{"plabel": "pvalue"}),
				WithPersistentVars(map[string]interface{}{"pvar": "pvalue"}),
				WithTags([]string{"tag"}...),
				WithVars(map[string]interface{}{"var": "value"}),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			image, err := NewImage(test.name, test.version, test.registryHost, test.registryNamesapce, test.options...)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, image)
			}
		})
	}
}

func TestAddChild(t *testing.T) {
	tests := []struct {
		desc  string
		image *Image
		child *Image
		res   *Image
	}{
		{
			desc: "Testing add child to image",
			image: &Image{
				Name: "image",
			},
			child: &Image{
				Name: "child",
			},
			res: &Image{
				Name:     "image",
				Children: []*Image{{Name: "child"}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			test.image.AddChild(test.child)

			assert.Equal(t, test.res, test.image)
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

func TestCopy(t *testing.T) {
	tests := []struct {
		desc  string
		image *Image
		res   *Image
		err   error
	}{
		{
			desc: "Testing copy of an image",
			image: &Image{
				Builder:  "builder",
				Children: []*Image{},
				Labels: map[string]string{
					"label": "value",
				},
				Name: "image",
				PersistentLabels: map[string]string{
					"plabel": "value",
				},
				PersistentVars: map[string]interface{}{
					"pvar": "value",
				},
				RegistryHost:      "registry.test",
				RegistryNamespace: "namespace",
				Tags: []string{
					"tag",
				},
				Vars: map[string]interface{}{
					"var": "value",
				},
				Version: "version",
			},
			res: &Image{
				Builder:  "builder",
				Children: []*Image{},
				Labels: map[string]string{
					"label": "value",
				},
				Name: "image",
				PersistentLabels: map[string]string{
					"plabel": "value",
				},
				PersistentVars: map[string]interface{}{
					"pvar": "value",
				},
				RegistryHost:      "registry.test",
				RegistryNamespace: "namespace",
				Tags: []string{
					"tag",
				},
				Vars: map[string]interface{}{
					"var": "value",
				},
				Version: "version",
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			image, err := test.image.Copy()

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, image)
			}
		})
	}
}

func TestIsWildcardImage(t *testing.T) {
	tests := []struct {
		desc  string
		image *Image
		res   bool
	}{
		{
			desc: "Testing wildcard image when is not a wildcard image",
			image: &Image{
				Name:    "image",
				Version: "version",
			},
			res: false,
		},
		{
			desc: "Testing wildcard image when is wildcard image",
			image: &Image{
				Name:    "image",
				Version: "*",
			},
			res: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res := test.image.IsWildcardImage()

			assert.Equal(t, test.res, res)
		})
	}
}
