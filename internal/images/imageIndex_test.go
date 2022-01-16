package tree

import (
	"reflect"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	gdstree "github.com/apenella/go-data-structures/tree"
	"github.com/stretchr/testify/assert"
)

func TestAddNode(t *testing.T) {

	testNode := &gdstree.Node{}
	testNode2 := &gdstree.Node{}

	tests := []struct {
		desc    string
		name    string
		version string
		node    *gdstree.Node
		index   *ImageIndex
		res     *ImageIndex
		err     error
	}{
		{
			desc:    "Testing add an image name to a nil image index",
			name:    "imageName",
			version: "imageVersion",
			node:    testNode,
			index:   nil,
			res:     nil,
			err:     errors.New("(tree::AddImage)", "Image index is null"),
		},
		{
			desc:    "Testing add an image name to an empty index",
			name:    "imageName",
			version: "imageVersion",
			node:    testNode,
			index:   &ImageIndex{},
			res: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
			},
			err: errors.New("(tree::AddImage)", "Image index is null"),
		},
		{
			desc:    "Testing add a new image version into an existing image name",
			name:    "imageName",
			version: "imageVersion2",
			node:    testNode,
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
			},
			res: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion", "imageName:imageVersion2"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion":  {testNode},
					"imageName:imageVersion2": {testNode},
				},
			},
			err: errors.New("(tree::AddImage)", "Image index is null"),
		},
		{
			desc:    "Testing add a second image name into names index",
			name:    "imageName2",
			version: "imageVersion",
			node:    testNode,
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
			},
			res: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName":  {"imageName:imageVersion"},
					"imageName2": {"imageName2:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion":  {testNode},
					"imageName2:imageVersion": {testNode},
				},
			},
			err: errors.New("(tree::AddImage)", "Image index is null"),
		},
		{
			desc:    "Testing add a second image into name-versions index",
			name:    "imageName",
			version: "imageVersion",
			node:    testNode2,
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
			},
			res: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode, testNode2},
				},
			},
			err: errors.New("(tree::AddImage)", "Image index is null"),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		err := test.index.AddNode(test.name, test.version, test.node)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			assert.True(t, reflect.DeepEqual(test.index, test.res))
		}
	}

}

func TestAddWildcardIndexImage(t *testing.T) {
	tests := []struct {
		desc    string
		name    string
		version string
		index   *ImageIndex
		res     map[string]uint8
		err     error
	}{
		{
			desc:    "Testing add a wildcard version to a nil wildcard index",
			name:    "ImageName",
			version: "*",
			index:   &ImageIndex{},
			res: map[string]uint8{
				"ImageName:*": uint8(0),
			},
			err: nil,
		},
		{
			desc:    "Testing add a wildcard version to wildcard index",
			name:    "ImageName2",
			version: "*",
			index: &ImageIndex{
				WildcardIndex: map[string]uint8{
					"ImageName:*": uint8(0),
				},
			},
			res: map[string]uint8{
				"ImageName:*":  uint8(0),
				"ImageName2:*": uint8(0),
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		err := test.index.AddWildcardIndexImage(test.name, test.version)

		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.True(t, reflect.DeepEqual(test.index.WildcardIndex, test.res))
		}
	}

}

func TestFind(t *testing.T) {

	testNode := &gdstree.Node{}
	testNode2 := &gdstree.Node{}
	testNode3 := &gdstree.Node{}
	testNode4 := &gdstree.Node{}

	tests := []struct {
		desc    string
		name    string
		version string
		index   *ImageIndex
		res     []*gdstree.Node
		err     error
	}{
		{
			desc:    "Testing find an image with no version defined",
			name:    "imageName",
			version: "",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode, testNode2},
					"imageName:*":            {testNode, testNode3},
				},
			},
			res: []*gdstree.Node{testNode, testNode2},
		},
		{
			desc:    "Testing find an undefined image with no version defined",
			name:    "imageNameUnexisting",
			version: "",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode, testNode2},
				},
			},
			res: nil,
			err: errors.New("(tree::Find)", "Error when finding images by name 'imageNameUnexisting'", errors.New("(tree::FindByName) ", "Image name 'imageNameUnexisting' does not exists")),
		},
		{
			desc:    "Testing find an image with a version defined",
			name:    "imageName",
			version: "imageVersion",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode, testNode2},
					"imageName:*":            {testNode3, testNode4},
				},
			},
			res: []*gdstree.Node{testNode, testNode2},
		},
		{
			desc:    "Testing find an inmage by alternative name and version",
			name:    "imageName",
			version: "alternative-imageVersion",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
				NameVersionAlternativeIndex: map[string][]*gdstree.Node{
					"imageName:alternative-imageVersion": {testNode},
				},
			},
			res: []*gdstree.Node{testNode},
		},
		{
			desc:    "Testing find an undefined inmage by alternative name and version",
			name:    "imageName",
			version: "alternative-imageVersion-unexisting",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
				NameVersionAlternativeIndex: map[string][]*gdstree.Node{
					"imageName:alternative-imageVersion": {testNode},
				},
			},
			res: nil,
			err: errors.New("(tree::Find)", "Error when finding image 'imageName:alternative-imageVersion-unexisting'"),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		list, err := test.index.Find(test.name, test.version)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.True(t, reflect.DeepEqual(list, test.res))
		}
	}

}

func TestFindByName(t *testing.T) {

	testNode := &gdstree.Node{}
	testNode2 := &gdstree.Node{}
	testNodeWildcard := &gdstree.Node{}

	tests := []struct {
		desc  string
		name  string
		index *ImageIndex
		res   []*gdstree.Node
		err   error
	}{
		{
			desc: "Testing find by name an unexisting image name",
			name: "imageName2",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
			},
			res: []*gdstree.Node{testNode},
			err: errors.New("(tree::FindByName)", "Image name 'imageName2' does not exists"),
		},
		{
			desc: "Testing find by name an unexisting name-version",
			name: "imageName",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion2"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
			},
			res: []*gdstree.Node{testNode},
			err: errors.New("(tree::FindByName)", "Image name-version 'imageName:imageVersion2' does not exists"),
		},
		{
			desc: "Testing find by name one image",
			name: "imageName",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
			},
			res: []*gdstree.Node{testNode},
			err: nil,
		},
		{
			desc: "Testing find by name multiple images",
			name: "imageName",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode, testNode2},
				},
			},
			res: []*gdstree.Node{testNode, testNode2},
			err: nil,
		},
		{
			desc: "Testing find by name multiple images with wildcard",
			name: "imageName",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion", "imageName:*"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode, testNode2},
					"imageName:*":            {testNodeWildcard},
				},
				WildcardIndex: map[string]uint8{
					"imageName:*": uint8(0),
				},
			},
			res: []*gdstree.Node{testNode, testNode2},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		list, err := test.index.findByName(test.name)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.True(t, reflect.DeepEqual(list, test.res))
		}
	}

}

func TestFindByNameAndVersion(t *testing.T) {

	testNode := &gdstree.Node{}

	tests := []struct {
		desc    string
		name    string
		version string
		index   *ImageIndex
		res     []*gdstree.Node
		err     error
	}{
		{
			desc:    "Testing find by name and version for an unexisting image name-version",
			name:    "imageName",
			version: "imageVersion2",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
			},
			res: nil,
			err: errors.New("(tree::FindByNameAndVersion)", "Image name-version 'imageName:imageVersion2' does not exists"),
		},
		{
			desc:    "Testing find by name and version one image",
			name:    "imageName",
			version: "imageVersion",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
			},
			res: []*gdstree.Node{testNode},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		list, err := test.index.findByNameAndVersion(test.name, test.version)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.True(t, reflect.DeepEqual(list, test.res))
		}
	}
}

func TestFindAlternativeByNameAndVersion(t *testing.T) {
	testNode := &gdstree.Node{}

	tests := []struct {
		desc    string
		name    string
		version string
		index   *ImageIndex
		res     []*gdstree.Node
		err     error
	}{
		{
			desc:    "Testing find alternative by name and version for an unexisting image name-version",
			name:    "imageName",
			version: "alternative-imageVersion2",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
				NameVersionAlternativeIndex: map[string][]*gdstree.Node{
					"imageName:alternative-imageVersion": {testNode},
				},
			},
			res: nil,
			err: errors.New("(tree::FindAlternativeByNameAndVersion)", "Image name-version 'imageName:alternative-imageVersion2' does not exists"),
		},
		{
			desc:    "Testing find alternative by name and version one image",
			name:    "imageName",
			version: "alternative-imageVersion",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
				NameVersionAlternativeIndex: map[string][]*gdstree.Node{
					"imageName:alternative-imageVersion": {testNode},
				},
			},
			res: []*gdstree.Node{testNode},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		list, err := test.index.findAlternativeByNameAndVersion(test.name, test.version)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.True(t, reflect.DeepEqual(list, test.res))
		}
	}
}

func TestFindWildcardVersion(t *testing.T) {

	testNode := &gdstree.Node{}
	testNodeWildcard := &gdstree.Node{}

	tests := []struct {
		desc    string
		name    string
		version string
		index   *ImageIndex
		res     *gdstree.Node
		err     error
	}{
		{
			desc:    "Testing find a node wildcard version",
			name:    "imageName",
			version: "imageVersion",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion", "imageName:*"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
					"imageName:*":            {testNodeWildcard},
				},
				NameVersionAlternativeIndex: map[string][]*gdstree.Node{
					"imageName:alternative-imageVersion": {testNode},
				},
			},
			res: testNodeWildcard,
			err: nil,
		},
		{
			desc:    "Testing find a node wildcard version when no wildcard version is defined",
			name:    "imageName",
			version: "imageVersion",
			index: &ImageIndex{
				NameIndex: map[string][]string{
					"imageName": {"imageName:imageVersion"},
				},
				NameVersionIndex: map[string][]*gdstree.Node{
					"imageName:imageVersion": {testNode},
				},
				NameVersionAlternativeIndex: map[string][]*gdstree.Node{
					"imageName:alternative-imageVersion": {testNode},
				},
			},
			res: nil,
			err: errors.New("(tree::FindWildcardVersion)", "Image 'imageName' does not have a wildcard version definition"),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		node, err := test.index.FindWildcardVersion(test.name)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.True(t, reflect.DeepEqual(node, test.res))
		}
	}
}
