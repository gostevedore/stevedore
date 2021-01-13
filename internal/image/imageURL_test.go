package image

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		desc string
		name string
		err  error
		res  *ImageURL
	}{
		{
			desc: "Testing parse image URL error when multiple coloms",
			name: "name:name:name",
			err:  errors.New("(ImageURL::Parse)", "Invalid image name"),
			res:  nil,
		},
		{
			desc: "Testing parse image URL error when multiple slashes",
			name: "name/name/name/name",
			err:  errors.New("(ImageURL::Parse)", "Invalid image name"),
			res:  nil,
		},
		{
			desc: "Testing parse image URL with name",
			name: "name",
			err:  nil,
			res: &ImageURL{
				Name: "name",
			},
		},
		{
			desc: "Testing parse image URL with name and tag",
			name: "name:tag",
			err:  nil,
			res: &ImageURL{
				Name: "name",
				Tag:  "tag",
			},
		},
		{
			desc: "Testing parse image URL with namespace, name and tag",
			name: "namespace/name:tag",
			err:  nil,
			res: &ImageURL{
				Namespace: "namespace",
				Name:      "name",
				Tag:       "tag",
			},
		},
		{
			desc: "Testing parse image URL with regitry, namespace, name and tag",
			name: "registry/namespace/name:tag",
			err:  nil,
			res: &ImageURL{
				Registry:  "registry",
				Namespace: "namespace",
				Name:      "name",
				Tag:       "tag",
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		parsed, err := Parse(test.name)

		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.res, parsed, "Unexpected value")
		}
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		desc  string
		image *ImageURL
		err   error
		res   string
	}{
		{
			desc: "Testing parse image URL with regitry, namespace, name and tag",
			err:  nil,
			image: &ImageURL{
				Registry:  "registry",
				Namespace: "namespace",
				Name:      "name",
				Tag:       "tag",
			},
			res: "registry/namespace/name:tag",
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		url, err := test.image.URL()

		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.res, url, "Unexpected value")
		}
	}
}

func TestSanitizeTag(t *testing.T) {
	tests := []struct {
		desc  string
		input string
		res   string
	}{
		{
			desc:  "Testing sanitize with no changes to apply",
			input: "no-changes",
			res:   "no-changes",
		},
		{
			desc:  "Testing sanitize /",
			input: "no/changes",
			res:   "no_changes",
		},
		{
			desc:  "Testing sanitize :",
			input: "no:changes",
			res:   "no_changes",
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		res := sanitizeTag(test.input)
		assert.Equal(t, test.res, res, "Unexpected value")
	}
}
