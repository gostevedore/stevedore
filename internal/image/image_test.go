package image

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	build "github.com/gostevedore/stevedore/internal/build"
	"github.com/gostevedore/stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

// TestLoadImage tests
func TestLoadImage(t *testing.T) {

	testBaseDir := "test"

	tests := []struct {
		desc  string
		file  string
		err   error
		image *Image
	}{
		{
			desc: "testing an unexistent file",
			file: "nofile",
			err: errors.New("(images::LoadImage)", "Images file could not be load",
				errors.New("", "(LoadYAMLFile) Error loading file nofile. open nofile: no such file or directory")),
			image: &Image{},
		},
		{
			desc: "Testing a simple image",
			file: filepath.Join(testBaseDir, "image.yml"),
			err:  &errors.Error{},
			image: &Image{
				Name:         "ubuntu",
				RegistryHost: "registry",
				Builder:      "infrastructure",
				Tags: []string{
					"16.04",
					"xenial",
				},
				Vars: map[string]interface{}{
					"container_name":   "ubuntu",
					"source_image_tag": "16.04",
				},
				Children: map[string][]string{
					"php-builder": {
						"7.1",
						"7.2",
					},
					"php-fpm": {
						"7.1",
					},
					"php-cli": {
						"7.1",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		image, err := LoadImage(test.file)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, image, test.image, "Unexpected value")
		}
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
			desc:  "Testing a nil image copy",
			image: nil,
			res:   nil,
			err:   errors.New("(image::Image::Copy)", "Image is nil"),
		},
		{
			desc: "Testing an image copy",
			image: &Image{
				Name:         "ubuntu",
				RegistryHost: "registry",
				Builder:      "infrastructure",
				Tags: []string{
					"16.04",
				},
				PersistentVars: map[string]interface{}{
					"ubuntu_version": "16.04",
				},
				Vars: map[string]interface{}{
					"container_name":   "ubuntu",
					"source_image_tag": "16.04",
				},
			},
			res: &Image{
				Name:         "ubuntu",
				RegistryHost: "registry",
				Builder:      "infrastructure",
				Tags: []string{
					"16.04",
				},
				PersistentVars: map[string]interface{}{
					"ubuntu_version": "16.04",
				},
				Vars: map[string]interface{}{
					"container_name":   "ubuntu",
					"source_image_tag": "16.04",
				},
			},
			err: nil,
		},
		{
			desc: "Testing an image copy with childs",
			image: &Image{
				Name:         "ubuntu",
				RegistryHost: "registry",
				Builder:      "infrastructure",
				Tags: []string{
					"16.04",
				},
				PersistentVars: map[string]interface{}{
					"ubuntu_version": "16.04",
				},
				Vars: map[string]interface{}{
					"container_name":   "ubuntu",
					"source_image_tag": "16.04",
				},
				Children: map[string][]string{
					"php-fpm": {
						"7.1",
						"7.3",
					},
				},
			},
			res: &Image{
				Name:         "ubuntu",
				RegistryHost: "registry",
				Builder:      "infrastructure",
				Tags: []string{
					"16.04",
				},
				PersistentVars: map[string]interface{}{
					"ubuntu_version": "16.04",
				},
				Vars: map[string]interface{}{
					"container_name":   "ubuntu",
					"source_image_tag": "16.04",
				},
				Children: map[string][]string{
					"php-fpm": {
						"7.1",
						"7.3",
					},
				},
			},
			err: errors.New("(image::Image::Copy)", "Image is nil"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.image.Copy()

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})

	}

}

func TestImageToArray(t *testing.T) {
	tests := []struct {
		desc  string
		image *Image
		res   []string
		err   error
	}{
		{
			desc:  "Testing array generation from a nil image",
			image: nil,
			res:   nil,
			err:   errors.New("(image::Image::ToArray)", "Image is nil"),
		},

		{
			desc: "Testing array generation from an image",
			image: &Image{
				Name:              "name",
				Version:           "version",
				Builder:           "type",
				RegistryNamespace: "namespace",
				RegistryHost:      "registry",
			},
			res: []string{"name", "version", "type", "namespace", "registry"},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.image.ToArray()

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestGetBuilderType(t *testing.T) {
	tests := []struct {
		desc  string
		image *Image
		res   string
	}{
		{
			desc: "Testing get builder type when builder is an string",
			image: &Image{
				Name:              "name",
				Version:           "version",
				Builder:           "builder",
				RegistryNamespace: "namespace",
				RegistryHost:      "registry",
			},
			res: "builder",
		},
		{
			desc: "Testing get builder type when builder is an string",
			image: &Image{
				Name:    "name",
				Version: "version",
				Builder: &build.Builder{
					Name:   "builder",
					Driver: "driver",
					Options: map[string]interface{}{
						"option1": "option1",
					},
				},
				RegistryNamespace: "namespace",
				RegistryHost:      "registry",
			},
			res: InlineBuilder,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res := test.image.getBuilderType()
			assert.Equal(t, test.res, res)
		})
	}
}

func TestCheckCompatibility(t *testing.T) {
	tests := []struct {
		desc  string
		image *Image
		res   *Image
	}{
		{
			desc: "Testing childs compatibility",
			image: &Image{
				Name: "image",
				Childs: map[string][]string{
					"child": {"v1", "v2"},
				},
			},
			res: &Image{
				Name: "image",
				Childs: map[string][]string{
					"child": {"v1", "v2"},
				},
				Children: map[string][]string{
					"child": {"v1", "v2"},
				},
			},
		},
	}

	console.Init(ioutil.Discard)

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			test.image.CheckCompatibility()
			assert.Equal(t, test.res, test.image)
		})
	}
}
