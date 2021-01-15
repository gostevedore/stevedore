package tree

import (
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	gdstree "github.com/apenella/go-data-structures/tree"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/stretchr/testify/assert"
)

// LoadImagesTree test
func TestLoadImagesTree(t *testing.T) {

	testBaseDir := "test"

	tests := []struct {
		desc       string
		file       string
		err        error
		imagesTree *ImagesTree
	}{
		{
			desc:       "Testing an unexistent file",
			file:       "nofile",
			err:        errors.New("(tree::LoadImagesTree)", "Error loading images tree configuration", errors.New("", "(LoadYAMLFile) Error loading file nofile. open nofile: no such file or directory")),
			imagesTree: &ImagesTree{},
		},
		{
			desc: "Testing a simple tree",
			file: filepath.Join(testBaseDir, "stevedore_multiple_images.yml"),
			err:  nil,
			imagesTree: &ImagesTree{
				Images: map[string]map[string]*image.Image{
					"php-fpm": {
						"7.1": &image.Image{
							Builder: "mock-builder",
							Tags: []string{
								"7.1",
							},
							Vars: map[string]interface{}{
								"container_name":   "php-fpm",
								"source_image_tag": "16.04",
							},
							Children: map[string][]string{
								"php-fpm-dev": {
									"7.1",
								},
							},
						},
						"7.2": &image.Image{
							Builder: "mock-builder",
							Tags: []string{
								"7.2",
							},
							Vars: map[string]interface{}{
								"container_name":   "php-fpm",
								"source_image_tag": "16.04",
							},
							Children: map[string][]string{
								"php-fpm-dev": {
									"7.2",
								},
							},
						},
					},
					"php-fpm-dev": {
						"7.1": &image.Image{
							Builder: "mock-builder",
							Tags: []string{
								"7.1",
							},
							Vars: map[string]interface{}{
								"container_name":   "php-fpm-dev",
								"source_image_tag": "16.04",
							},
						},
						"7.2": &image.Image{
							Builder: "mock-builder",
							Tags: []string{
								"7.2",
							},
							Vars: map[string]interface{}{
								"container_name":   "php-fpm-dev",
								"source_image_tag": "16.04",
							},
						},
					},
					"ubuntu": {
						"16.04": &image.Image{
							Builder: "mock-builder",
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
								},
								"php-fpm": {
									"7.1",
									"7.2",
								},
							},
						},
					},
				},
			},
		},
		{
			desc:       "Testing a simple tree",
			file:       filepath.Join(testBaseDir, "stevedore_nil.yml"),
			err:        errors.New("(tree::LoadImagesTree)", "Image tree is not defined properly on "+filepath.Join(testBaseDir, "stevedore_nil.yml")),
			imagesTree: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		imagesTree, err := LoadImagesTree(test.file)

		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.imagesTree, imagesTree, "Unexpected value")
		}
	}
}

func TestGenerateGraph(t *testing.T) {
	tests := []struct {
		desc       string
		file       string
		err        error
		imagesTree *ImagesTree
	}{
		{
			desc: "Testing a simple tree",
			file: "../test/images/simpleImagesTree.yml",
			err:  nil,
			imagesTree: &ImagesTree{
				Images: map[string]map[string]*image.Image{
					"php-fpm": {
						"7.1": &image.Image{
							Builder: "infrastructure",
							Tags: []string{
								"7.1",
							},
							PersistentVars: map[string]interface{}{
								"php_version": "7.1",
							},
							Vars: map[string]interface{}{
								"container_name":   "php-fpm",
								"source_image_tag": "16.04",
							},
							Children: map[string][]string{
								"php-fpm-dev": {
									"7.1",
								},
							},
						},
						"7.2": &image.Image{
							Builder: "infrastructure",
							Tags: []string{
								"7.2",
							},
							PersistentVars: map[string]interface{}{
								"php_version": "7.1",
							},
							Vars: map[string]interface{}{
								"container_name":   "php-fpm",
								"source_image_tag": "16.04",
							},
							Children: map[string][]string{
								"php-fpm-dev": {
									"7.2",
								},
							},
						},
					},
					"php-fpm-dev": {
						"7.1": &image.Image{
							Builder: "infrastructure",
							Tags: []string{
								"7.1",
							},
							Vars: map[string]interface{}{
								"container_name":   "php-fpm-dev",
								"source_image_tag": "16.04",
							},
						},
						"7.2": &image.Image{
							Builder: "infrastructure",
							Tags: []string{
								"7.2",
							},
							Vars: map[string]interface{}{
								"container_name":   "php-fpm-dev",
								"source_image_tag": "16.04",
							},
						},
					},
					"ubuntu": {
						"16.04": &image.Image{
							Builder: "infrastructure",
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
								},
								"php-fpm": {
									"7.1",
									"7.2",
								},
								"php-cli": {
									"7.1",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		_, _, err := test.imagesTree.GenerateGraph()
		if err != nil {
			t.Fatal(err.Error())
		}
	}

}

func TestGenerateTemplateGraph(t *testing.T) {

	imagesTree := &ImagesTree{
		Images: map[string]map[string]*image.Image{
			"php-fpm": {
				"7.1": &image.Image{
					Builder: "infrastructure",
					Tags: []string{
						"7.1",
					},
					Vars: map[string]interface{}{
						"container_name":   "php-fpm",
						"source_image_tag": "16.04",
					},
					Children: map[string][]string{
						"php-fpm-dev": {
							"7.1",
						},
					},
				},
				"7.2": &image.Image{
					Builder: "infrastructure",
					Tags: []string{
						"7.2",
					},
					Vars: map[string]interface{}{
						"container_name":   "php-fpm",
						"source_image_tag": "16.04",
					},
					Children: map[string][]string{
						"php-fpm-dev": {
							"7.2",
						},
					},
				},
			},
			"php-fpm-dev": {
				"7.1": &image.Image{
					Builder: "infrastructure",
					Tags: []string{
						"7.1",
					},
					Vars: map[string]interface{}{
						"container_name":   "php-fpm-dev",
						"source_image_tag": "16.04",
					},
				},
				"7.2": &image.Image{
					Builder: "infrastructure",
					Tags: []string{
						"7.2",
					},
					Vars: map[string]interface{}{
						"container_name":   "php-fpm-dev",
						"source_image_tag": "16.04",
					},
				},
			},
			"ubuntu": {
				"16.04": &image.Image{
					Builder: "infrastructure",
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
						},
						"php-fpm": {
							"7.1",
							"7.2",
						},
						"php-cli": {
							"7.1",
						},
					},
				},
			},
		},
	}

	tests := []struct {
		desc         string
		nodeImage    *image.Image
		imageGraph   *gdsexttree.Graph
		imageName    string
		imageVersion string
		parent       *gdsexttree.Node
		res          *gdstree.Graph
		err          error
	}{
		{
			desc:         "Generate graph with an empty node image",
			nodeImage:    nil,
			imageGraph:   &gdsexttree.Graph{},
			imageName:    "",
			imageVersion: "",
			parent:       nil,
			res:          nil,
			err:          errors.New("(tree::generateGraphRec)", "Node Image is null"),
		},
		{
			desc: "Adding Image to an existing graph",
			nodeImage: &image.Image{
				Name:    "nginx",
				Version: "1.10",
			},
			imageName:    "nginx",
			imageVersion: "1.10",
			parent: &gdsexttree.Node{
				Name: "ubuntu:16.04",
				Item: &image.Image{},
			},
			imageGraph: &gdsexttree.Graph{
				Root: []*gdsexttree.Node{
					{
						Name: "ubuntu:16.04",
					},
				},
				NodesIndex: map[string]*gdsexttree.Node{
					"ubuntu:16.04": {
						Name: "ubuntu:16.04",
					},
					"php-fpm:7.1": {
						Name: "php-fpm:7.1",
						Parents: []*gdsexttree.Node{
							{
								Name: "ubuntu:16.04",
							},
						},
					},
				},
			},
			res: &gdstree.Graph{
				Root: []*gdstree.Node{
					{
						Name: "ubuntu:16.04",
					},
				},
				NodesIndex: map[string]*gdstree.Node{
					"ubuntu:16.04": {
						Name: "ubuntu:16.04",
					},
					"php-fpm:7.1": {
						Name: "php-fpm:7.1",
						Parent: &gdstree.Node{
							Name: "ubuntu:16.04",
						},
					},
					"nginx:1.10": {
						Name: "nginx:1.10",
						Parent: &gdstree.Node{
							Name: "ubuntu:16.04",
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		err := imagesTree.generateTemplateGraph(test.imageName, test.imageVersion, test.nodeImage, test.imageGraph, test.parent)

		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, len(test.res.Root), len(test.imageGraph.Root), "Unexpected lenght of root elements")
			assert.Equal(t, len(test.res.NodesIndex), len(test.imageGraph.NodesIndex), "Unexpected lenght of nodes index elements")
		}
	}
}

func TestGenerateNodeName(t *testing.T) {
	t.Log("Testing node name generation")

	name := "name"
	version := "version"
	res := name + ImageNodeNameSeparator + version

	i := &image.Image{
		Name:    name,
		Version: version,
	}
	nodename := GenerateNodeName(i)

	assert.Equal(t, nodename, res, "Nodename is not valid")
}

func TestRenderizeGraph(t *testing.T) {

	phpFpmDev71 := &gdsexttree.Node{
		Name: "php-fpm-dev:7.1",
		Item: &image.Image{
			Name:    "php-fpm-dev",
			Version: "{{ .Parent.Version }}",
		},
	}
	phpCliDev71 := &gdsexttree.Node{
		Name: "php-cli-dev:7.1",
		Item: &image.Image{
			Name:    "php-cli-dev",
			Version: "{{ .Parent.Version }}",
		},
	}
	phpFpm71 := &gdsexttree.Node{
		Name: "php-fpm:7.1",
		Item: &image.Image{
			Name:    "php-fpm",
			Version: "7.1-ubuntu{{ .Parent.Version }}",
		},
		Children: []*gdsexttree.Node{
			phpFpmDev71,
		},
	}
	phpFpm72 := &gdsexttree.Node{
		Name: "php-fpm:7.2",
		Item: &image.Image{
			Name:    "php-fpm",
			Version: "7.2-ubuntu{{ .Parent.Version }}",
		},
	}
	phpCli71 := &gdsexttree.Node{
		Name: "php-cli:7.1",
		Item: &image.Image{
			Name:    "php-cli",
			Version: "7.1-ubuntu{{ .Parent.Version }}",
		},
		Children: []*gdsexttree.Node{
			phpCliDev71,
		},
	}
	phpCli72 := &gdsexttree.Node{
		Name: "php-cli:7.2",
		Item: &image.Image{
			Name:    "php-cli",
			Version: "7.2-ubuntu{{ .Parent.Version }}",
		},
	}
	phpBuilder71 := &gdsexttree.Node{
		Name: "php-builder:7.1",
		Item: &image.Image{
			Name:    "php-builder",
			Version: "7.1-ubuntu{{ .Parent.Version }}",
		},
	}
	ubuntu16 := &gdsexttree.Node{
		Name: "ubuntu:16.04",
		Item: &image.Image{
			Name:    "ubuntu",
			Version: "16.04",
		},
		Children: []*gdsexttree.Node{
			phpFpm71,
			phpFpm72,
			phpCli71,
			phpCli72,
			phpBuilder71,
		},
	}
	ubuntu18 := &gdsexttree.Node{
		Name: "ubuntu:18.04",
		Item: &image.Image{
			Name:    "ubuntu",
			Version: "18.04",
		},
		Children: []*gdsexttree.Node{
			phpFpm71,
			phpFpm72,
			phpCli71,
			phpCli72,
			phpBuilder71,
		},
	}

	phpBuilder71.AddParent(ubuntu16)
	phpBuilder71.AddParent(ubuntu18)
	phpCli71.AddParent(ubuntu16)
	phpCli71.AddParent(ubuntu18)
	phpCli72.AddParent(ubuntu16)
	phpCli72.AddParent(ubuntu18)
	phpFpm71.AddParent(ubuntu16)
	phpFpm71.AddParent(ubuntu18)
	phpFpm72.AddParent(ubuntu16)
	phpFpm72.AddParent(ubuntu18)
	phpFpmDev71.AddParent(phpFpm71)
	phpCliDev71.AddParent(phpCli71)

	imagesGraph := &gdsexttree.Graph{
		Root: []*gdsexttree.Node{
			ubuntu16,
			ubuntu18,
		},
		NodesIndex: map[string]*gdsexttree.Node{
			"ubuntu:16.04":    ubuntu16,
			"ubuntu:18.04":    ubuntu18,
			"php-fpm:7.1":     phpFpm71,
			"php-fpm:7.2":     phpFpm72,
			"php-cli:7.1":     phpCli71,
			"php-cli:7.2":     phpCli72,
			"php-builder:7.1": phpBuilder71,
			"php-fpm-dev:7.1": phpFpmDev71,
			"php-cli-dev:7.1": phpCliDev71,
		},
	}

	ResPhpFpmDev7116 := &gdstree.Node{
		Name: "php-fpm-dev:7.1-ubuntu16.04@php-fpm:7.1-ubuntu16.04",
		Item: &image.Image{
			Name:    "php-fpm-dev",
			Version: "7.1-ubuntu16.04",
		},
	}
	ResPhpFpmDev7118 := &gdstree.Node{
		Name: "php-fpm-dev:7.1-ubuntu18.04@php-fpm:7.1-ubuntu18.04",
		Item: &image.Image{
			Name:    "php-fpm-dev",
			Version: "7.1-ubuntu18.04",
		},
	}
	ResPhpCliDev7116 := &gdstree.Node{
		Name: "php-cli-dev:7.1-ubuntu16.04@php-fpm:7.1-ubuntu16.04",
		Item: &image.Image{
			Name:    "php-cli-dev",
			Version: "7.1-ubuntu16.04",
		},
	}
	ResPhpCliDev7118 := &gdstree.Node{
		Name: "php-cli-dev:7.1-ubuntu18.04@php-fpm:7.1-ubuntu18.04",
		Item: &image.Image{
			Name:    "php-cli-dev",
			Version: "7.1-ubuntu18.04",
		},
	}
	ResPhpFpm7116 := &gdstree.Node{
		Name: "php-fpm:7.1-ubuntu16.04@ubuntu:16.04",
		Item: &image.Image{
			Name:    "php-fpm",
			Version: "7.1-ubuntu16.04",
		},
	}
	ResPhpFpm7118 := &gdstree.Node{
		Name: "php-fpm:7.1-ubuntu18.04@ubuntu:18.04",
		Item: &image.Image{
			Name:    "php-fpm",
			Version: "7.1-ubuntu18.04",
		},
	}
	ResPhpFpm7216 := &gdstree.Node{
		Name: "php-fpm:7.2-ubuntu16.04@ubuntu:16.04",
		Item: &image.Image{
			Name:    "php-fpm",
			Version: "7.2-ubuntu16.04",
		},
	}
	ResPhpFpm7218 := &gdstree.Node{
		Name: "php-fpm:7.2-ubuntu18.04@ubuntu:18.04",
		Item: &image.Image{
			Name:    "php-fpm",
			Version: "7.2-ubuntu18.04",
		},
	}
	ResPhpCli7116 := &gdstree.Node{
		Name: "php-cli:7.1-ubuntu16.04@ubuntu:16.04",
		Item: &image.Image{
			Name:    "php-cli",
			Version: "7.1-ubuntu16.04",
		},
	}
	ResPhpCli7118 := &gdstree.Node{
		Name: "php-cli:7.1-ubuntu18.04@ubuntu:18.04",
		Item: &image.Image{
			Name:    "php-cli",
			Version: "7.1-ubuntu18.04",
		},
	}
	ResPhpCli7216 := &gdstree.Node{
		Name: "php-cli:7.2-ubuntu16.04@ubuntu:16.04",
		Item: &image.Image{
			Name:    "php-cli",
			Version: "7.2-ubuntu16.04",
		},
	}
	ResPhpCli7218 := &gdstree.Node{
		Name: "php-cli:7.2-ubuntu18.04@ubuntu:18.04",
		Item: &image.Image{
			Name:    "php-cli",
			Version: "7.2-ubuntu18.04",
		},
	}
	ResPhpBuilder7116 := &gdstree.Node{
		Name: "php-builder:7.1-ubuntu16.04@ubuntu:16.04",
		Item: &image.Image{
			Name:    "php-builder",
			Version: "7.1-ubuntu16.04",
		},
	}
	ResPhpBuilder7118 := &gdstree.Node{
		Name: "php-builder:7.1-ubuntu18.04@ubuntu:16.04",
		Item: &image.Image{
			Name:    "php-builder",
			Version: "7.1-ubuntu18.04",
		},
	}
	ResUbuntu16 := &gdstree.Node{
		Name: "ubuntu:16.04",
		Item: &image.Image{
			Name:    "ubuntu",
			Version: "16.04",
		},
	}
	ResUbuntu18 := &gdstree.Node{
		Name: "ubuntu:18.04",
		Item: &image.Image{
			Name:    "ubuntu",
			Version: "18.04",
		},
	}

	graphRes := &gdstree.Graph{}
	// bases
	graphRes.AddNode(ResUbuntu16)
	graphRes.AddNode(ResUbuntu18)
	// ubuntu16 base
	graphRes.AddNode(ResPhpBuilder7116)
	graphRes.AddRelationship(ResUbuntu16, ResPhpBuilder7116)
	graphRes.AddNode(ResPhpFpm7116)
	graphRes.AddRelationship(ResUbuntu16, ResPhpFpm7116)
	graphRes.AddNode(ResPhpFpm7216)
	graphRes.AddRelationship(ResUbuntu16, ResPhpFpm7216)
	graphRes.AddNode(ResPhpCli7116)
	graphRes.AddRelationship(ResUbuntu16, ResPhpCli7116)
	graphRes.AddNode(ResPhpCli7216)
	graphRes.AddRelationship(ResUbuntu16, ResPhpCli7216)
	graphRes.AddNode(ResPhpFpmDev7116)
	graphRes.AddRelationship(ResPhpFpm7116, ResPhpFpmDev7116)
	graphRes.AddNode(ResPhpCliDev7116)
	graphRes.AddRelationship(ResPhpCli7116, ResPhpCliDev7116)
	// ubuntu18 base
	graphRes.AddNode(ResPhpBuilder7118)
	graphRes.AddRelationship(ResUbuntu16, ResPhpBuilder7118)
	graphRes.AddNode(ResPhpFpm7118)
	graphRes.AddRelationship(ResUbuntu18, ResPhpFpm7118)
	graphRes.AddNode(ResPhpFpm7218)
	graphRes.AddRelationship(ResUbuntu18, ResPhpFpm7218)
	graphRes.AddNode(ResPhpCli7118)
	graphRes.AddRelationship(ResUbuntu18, ResPhpCli7118)
	graphRes.AddNode(ResPhpCli7218)
	graphRes.AddRelationship(ResUbuntu18, ResPhpCli7218)
	graphRes.AddNode(ResPhpFpmDev7118)
	graphRes.AddRelationship(ResPhpFpm7118, ResPhpFpmDev7118)
	graphRes.AddNode(ResPhpCliDev7118)
	graphRes.AddRelationship(ResPhpCli7118, ResPhpCliDev7118)

	indexhRes := &ImageIndex{
		NameIndex: map[string][]string{
			"ubuntu":      {"ubuntu:16.04", "ubuntu:18.04"},
			"php-fpm":     {""},
			"php-fpm-dev": {""},
			"php-cli":     {""},
			"php-cli-dev": {""},
			"php-builder": {""},
		},
		NameVersionIndex: map[string][]*gdstree.Node{
			"ubuntu:16.04":    nil,
			"ubuntu:18.04":    nil,
			"php-fpm:7.1":     nil,
			"php-fpm-dev:7.1": nil,
			"php-cli:7.1":     nil,
			"php-cli-dev:7.1": nil,
			"php-builder:7.1": nil,
			"php-fpm:7.2":     nil,
			"php-cli:7.2":     nil,
		},
		NameVersionAlternativeIndex: map[string][]*gdstree.Node{
			"php-fpm:7.1-ubuntu16.04":     nil,
			"php-fpm-dev:7.1-ubuntu16.04": nil,
			"php-cli:7.1-ubuntu16.04":     nil,
			"php-cli-dev:7.1-ubuntu16.04": nil,
			"php-builder:7.1-ubuntu16.04": nil,
			"php-fpm:7.2-ubuntu16.04":     nil,
			"php-cli:7.2-ubuntu16.04":     nil,
			"php-fpm:7.1-ubuntu18.04":     nil,
			"php-fpm-dev:7.1-ubuntu18.04": nil,
			"php-cli:7.1-ubuntu18.04":     nil,
			"php-cli-dev:7.1-ubuntu18.04": nil,
			"php-builder:7.1-ubuntu18.04": nil,
			"php-fpm:7.2-ubuntu18.04":     nil,
			"php-cli:7.2-ubuntu18.04":     nil,
		},
	}

	g, i, _ := RenderizeGraph(imagesGraph)
	t.Log("Testing graph elements")
	assert.Equal(t, len(graphRes.Root), len(g.Root), "Unexpected lenght of root elements")
	assert.Equal(t, len(graphRes.NodesIndex), len(g.NodesIndex), "Unexpected lenght of nodes index elements")
	t.Log("Testing images index elements")
	assert.Equal(t, len(indexhRes.NameIndex), len(i.NameIndex), "Unexpected lenght of index names")
	assert.Equal(t, len(indexhRes.NameVersionIndex), len(i.NameVersionIndex), "Unexpected lenght of name-version elements")
	assert.Equal(t, len(indexhRes.NameVersionAlternativeIndex), len(i.NameVersionAlternativeIndex), "Unexpected lenght of alternatives elements")
}

func TestRenderizeGraphRec(t *testing.T) {
	//TODO
}

func TestGetNodeImage(t *testing.T) {

	imageNode := &image.Image{}

	tests := []struct {
		desc string
		node *gdstree.Node
		res  *image.Image
		err  error
	}{
		{
			desc: "Testing get node image from a nil node",
			node: nil,
			res:  nil,
			err:  errors.New("(tree::GetNodeImage)", "Node is nil"),
		},
		{
			desc: "Testing get node image from a node with nil image",
			node: &gdstree.Node{
				Item: nil,
			},
			res: nil,
			err: errors.New("(tree::GetNodeImage)", "Node item is nil"),
		},
		{
			desc: "Testing get node image from a node",
			node: &gdstree.Node{
				Item: imageNode,
			},
			res: imageNode,
			err: errors.New("(tree::GetNodeImage)", "Node item is nil"),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		res, err := GetNodeImage(test.node)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.res, res, "Unexpected Image")
		}
	}
}

func TestGenerateWilcardVersionNode(t *testing.T) {

	phpFpm71 := &image.Image{
		Name:    "php-fpm",
		Builder: "infrastructure",
		Tags: []string{
			"7.1",
		},
		PersistentVars: map[string]interface{}{
			"php_version": "7.1",
		},
		Vars: map[string]interface{}{
			"container_name": "php-fpm",
		},
		Children: map[string][]string{
			"php-fpm-dev": {
				"7.1",
			},
		},
	}

	phpFpmWilcard := &image.Image{
		Name:    "php-fpm",
		Builder: "infrastructure",
		Tags: []string{
			"{{ .Version }}",
		},
		PersistentVars: map[string]interface{}{
			"php_version": "{{ .Version }}",
		},
		Vars: map[string]interface{}{
			"container_name": "php-fpm",
		},
		Children: map[string][]string{
			"php-fpm-dev": {
				"{{ .Version }}",
			},
		},
	}

	phpFpmDev71 := &image.Image{
		Name:    "php-fpm-dev",
		Builder: "infrastructure",
		Tags: []string{
			"7.1",
		},
		Vars: map[string]interface{}{
			"container_name": "php-fpm-dev",
		},
	}

	phpFpmDevWilcard := &image.Image{
		Name:    "php-fpm-dev",
		Builder: "infrastructure",
		Tags: []string{
			"{{ .Version }}",
		},
		Vars: map[string]interface{}{
			"container_name": "php-fpm-dev",
		},
	}

	ubuntu16 := &image.Image{
		Name:    "ubuntu",
		Builder: "infrastructure",
		Tags: []string{
			"16.04",
			"xenial",
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
				"*",
			},
		},
	}

	tree := &ImagesTree{
		Images: map[string]map[string]*image.Image{
			"php-fpm": {
				"7.1": phpFpm71,
				"*":   phpFpmWilcard,
			},
			"php-fpm-dev": {
				"7.1": phpFpmDev71,
				"*":   phpFpmDevWilcard,
			},
			"ubuntu": {
				"16.04": ubuntu16,
			},
		},
	}

	tests := []struct {
		desc     string
		tree     *ImagesTree
		version  string
		nodeBase *gdstree.Node
		res      *gdstree.Node
		err      error
	}{
		{
			desc:     "Testing generate wildcard version gave a nil images tree",
			tree:     nil,
			version:  "",
			nodeBase: nil,
			res:      nil,
			err:      errors.New("(tree::GenerateNodeWithWilcardVersion)", "Images tree is nil"),
		},
		{
			desc:     "Testing generate wildcard version gave a nil node",
			tree:     tree,
			version:  "",
			nodeBase: nil,
			res:      nil,
			err:      errors.New("(tree::GenerateNodeWithWilcardVersion)", "Node is nil"),
		},
		{
			desc:    "Testing generate wildcard version node",
			version: "version",
			nodeBase: &gdstree.Node{
				Name:   "php-fpm:*",
				Item:   phpFpmWilcard,
				Parent: nil,
			},
			tree: tree,
			res: &gdstree.Node{
				Name: "php-fpm:version",
				Item: &image.Image{
					Name:      "php-fpm",
					Namespace: "",
					Version:   "version",
					Builder:   "infrastructure",
					Tags: []string{
						"version",
					},
					PersistentVars: map[string]interface{}{
						"php_version": "version",
					},
					Vars: map[string]interface{}{
						"container_name": "php-fpm",
					},
					Children: map[string][]string{
						"php-fpm-dev": {
							"version",
						},
					},
					Childs: map[string][]string{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		res, err := test.tree.GenerateWilcardVersionNode(test.nodeBase, test.version)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.res, res, "Unexpected Node")
		}
	}
}
