package image

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/compatibility"
	domainimage "github.com/gostevedore/stevedore/internal/images/image"
	"github.com/stretchr/testify/assert"
)

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
				Builder: "builder",
				Labels:  map[string]string{"label": "value"},
				Name:    "ubuntu",
				Parents: map[string][]string{"parent": {"parent_version"}},
				PersistentLabels: map[string]string{
					"plabel": "pvalue",
				},
				PersistentVars: map[string]interface{}{
					"ubuntu_version": "16.04",
				},
				Tags: []string{
					"16.04",
				},
				RegistryHost: "registry",
				Vars: map[string]interface{}{
					"container_name":   "ubuntu",
					"source_image_tag": "16.04",
				},
			},
			res: &Image{
				Builder: "builder",
				Labels:  map[string]string{"label": "value"},
				Name:    "ubuntu",
				Parents: map[string][]string{"parent": {"parent_version"}},
				PersistentLabels: map[string]string{
					"plabel": "pvalue",
				},
				PersistentVars: map[string]interface{}{
					"ubuntu_version": "16.04",
				},
				Tags: []string{
					"16.04",
				},
				RegistryHost: "registry",
				Vars: map[string]interface{}{
					"container_name":   "ubuntu",
					"source_image_tag": "16.04",
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing an image copy with childs",
			image: &Image{
				Builder: "builder",
				Children: map[string][]string{
					"php-fpm": {
						"7.1",
						"7.3",
					},
				},
				Name:             "ubuntu",
				RegistryHost:     "registry",
				PersistentLabels: map[string]string{},
				PersistentVars: map[string]interface{}{
					"ubuntu_version": "16.04",
				},
				Tags: []string{
					"16.04",
				},
				Vars: map[string]interface{}{
					"container_name":   "ubuntu",
					"source_image_tag": "16.04",
				},
			},
			res: &Image{
				Builder: "builder",
				Children: map[string][]string{
					"php-fpm": {
						"7.1",
						"7.3",
					},
				},
				Labels:           map[string]string{},
				Name:             "ubuntu",
				RegistryHost:     "registry",
				PersistentLabels: map[string]string{},
				PersistentVars: map[string]interface{}{
					"ubuntu_version": "16.04",
				},
				Tags: []string{
					"16.04",
				},
				Vars: map[string]interface{}{
					"container_name":   "ubuntu",
					"source_image_tag": "16.04",
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

func TestCreateDomainImage(t *testing.T) {

	tests := []struct {
		desc  string
		image *Image
		res   *domainimage.Image
		err   error
	}{
		{
			desc: "Testing create a domain image",
			image: &Image{
				Builder: "builder",
				Labels: map[string]string{
					"label": "value",
				},
				Name: "image",
				PersistentLabels: map[string]string{
					"plabel": "pvalue",
				},
				PersistentVars: map[string]interface{}{
					"pvar": "pvalue",
				},
				RegistryHost:      "registry.test",
				RegistryNamespace: "namespace",
				Tags: []string{
					"tag",
				},
				Vars: map[string]interface{}{
					"var": "value",
				},
				Version: "1.0.0",
			},
			res: &domainimage.Image{
				Builder: "builder",
				Name:    "image",
				Labels: map[string]string{
					"label": "value",
				},
				PersistentLabels: map[string]string{
					"plabel": "pvalue",
				},
				PersistentVars: map[string]interface{}{
					"pvar": "pvalue",
				},
				RegistryHost:      "registry.test",
				RegistryNamespace: "namespace",
				Tags: []string{
					"tag",
				},
				Vars: map[string]interface{}{
					"var": "value",
				},
				Version: "1.0.0",
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.image.CreateDomainImage()

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestImageCheckCompatibility(t *testing.T) {
	tests := []struct {
		desc              string
		image             *Image
		compatibility     Compatibilitier
		prepareAssertFunc func(Compatibilitier)
	}{
		// {
		// 	desc:          "Testing childs compatibility",
		// 	compatibility: compatibility.NewMockCompatibility(),
		// 	image: &Image{
		// 		Name: "image",
		// 		Childs: map[string][]string{
		// 			"child": {"v1", "v2"},
		// 		},
		// 	},
		// 	prepareAssertFunc: func(c Compatibilitier) {
		// 		c.(*compatibility.MockCompatibility).On("AddDeprecated", []string{"On 'image', 'childs' attribute must be replaced by 'children' before 0.11.0"}).Return(nil)
		// 	},
		// },
	}

	//console.Init(ioutil.Discard)

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.compatibility)
			}

			test.image.CheckCompatibility(test.compatibility)
			test.compatibility.(*compatibility.MockCompatibility).AssertExpectations(t)

		})
	}
}
