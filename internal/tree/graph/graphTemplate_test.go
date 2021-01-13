package graph

import (
	"errors"
	"stevedore/internal/image"
	"stevedore/internal/tree"
	"testing"

	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTemplateGraph(t *testing.T) {

	imagesTree := &tree.ImagesTree{
		Images: map[string]map[string]*image.Image{
			"php-fpm": {
				"7.1": &image.Image{
					Type: "infrastructure",
					Children: map[string][]string{
						"php-fpm-dev": {
							"7.1",
						},
					},
				},
				"7.2": &image.Image{
					Type: "infrastructure",
					Children: map[string][]string{
						"php-fpm-dev": {
							"7.2",
						},
					},
				},
			},
			"php-fpm-dev": {
				"7.1": &image.Image{
					Type: "infrastructure",
				},
				"7.2": &image.Image{
					Type: "infrastructure",
				},
			},
			"ubuntu": {
				"16.04": &image.Image{
					Type: "infrastructure",
					Children: map[string][]string{
						"php-fpm": {
							"7.1",
							"7.2",
						},
					},
				},
				"18.04": &image.Image{
					Type: "infrastructure",
					Children: map[string][]string{
						"php-fpm": {
							"7.1",
							"7.2",
						},
					},
				},
			},
		},
	}

	tests := []struct {
		desc string
		tree *tree.ImagesTree
		res  *gdsexttree.Graph
		err  error
	}{
		{
			desc: "Testing generate graph with a nil node image",
			tree: nil,
			res:  nil,
			err:  errors.New("(graph::GenerateGraphTemplate) Tree is null"),
		},
		{
			desc: "Testing generate subgraph from an image",
			tree: imagesTree,
			res: &gdsexttree.Graph{
				Root: []*gdsexttree.Node{
					{
						Name: "ubuntu:16.04",
					},
					{
						Name: "ubuntu:18.04",
					},
				},
				NodesIndex: map[string]*gdsexttree.Node{
					"ubuntu:16.04": {
						Name: "ubuntu:16.04",
					},
					"ubuntu:18.04": {
						Name: "ubuntu:18.04",
					},
					"php-fpm:7.1": {
						Name: "php-fpm:7.1",
						Parents: []*gdsexttree.Node{
							{
								Name: "ubuntu:16.04",
							},
						},
					},
					"php-fpm:7.2": {
						Name: "php-fpm:7.2",
						Parents: []*gdsexttree.Node{
							{
								Name: "ubuntu:16.04",
							},
						},
					},
					"php-fpm-dev:7.1": {
						Name: "php-fpm-dev:7.1",
						Parents: []*gdsexttree.Node{
							{
								Name: "ubuntu:16.04",
							},
						},
					},
					"php-fpm-dev:7.2": {
						Name: "php-fpm:7.2",
						Parents: []*gdsexttree.Node{
							{
								Name: "ubuntu:16.04",
							},
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		graph, err := GenerateGraphTemplate(test.tree)

		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, len(test.res.Root), len(graph.Graph.Root), "Unexpected lenght of root elements")
			assert.Equal(t, len(test.res.NodesIndex), len(graph.Graph.NodesIndex), "Unexpected lenght of nodes index elements")
		}
	}
}

func TestGenerateTemplateGraphRec(t *testing.T) {

	imagesTree := &tree.ImagesTree{
		Images: map[string]map[string]*image.Image{
			"php-fpm": {
				"7.1": &image.Image{
					Type: "infrastructure",
					Children: map[string][]string{
						"php-fpm-dev": {
							"7.1",
						},
					},
				},
				"7.2": &image.Image{
					Type: "infrastructure",
					Children: map[string][]string{
						"php-fpm-dev": {
							"7.2",
						},
					},
				},
			},
			"php-fpm-dev": {
				"7.1": &image.Image{
					Type: "infrastructure",
				},
				"7.2": &image.Image{
					Type: "infrastructure",
				},
			},
			"ubuntu": {
				"16.04": &image.Image{
					Type: "infrastructure",
					Children: map[string][]string{
						"php-fpm": {
							"7.1",
							"7.2",
						},
					},
				},
				"18.04": &image.Image{
					Type: "infrastructure",
					Children: map[string][]string{
						"php-fpm": {
							"7.1",
							"7.2",
						},
					},
				},
			},
		},
	}

	tests := []struct {
		desc         string
		imageName    string
		imageVersion string
		image        *image.Image
		parent       *gdsexttree.Node
		graph        *gdsexttree.Graph
		tree         *tree.ImagesTree
		res          *gdsexttree.Graph
		err          error
	}{
		{
			desc:         "Testing generate graph with a nil node image",
			imageName:    "",
			imageVersion: "",
			image:        nil,
			parent:       nil,
			graph:        &gdsexttree.Graph{},
			tree:         imagesTree,
			res:          nil,
			err:          errors.New("(graph::generateGraphTemplateRec) Node Image is null"),
		},
		{
			desc:         "Testing generate graph with a nil graph",
			imageName:    "",
			imageVersion: "",
			image:        &image.Image{},
			parent:       nil,
			graph:        nil,
			tree:         imagesTree,
			res:          nil,
			err:          errors.New("(graph::generateGraphTemplateRec) Graph is null"),
		},
		{
			desc:         "Testing generate graph with a nil tree",
			imageName:    "",
			imageVersion: "",
			image:        &image.Image{},
			parent:       nil,
			graph:        &gdsexttree.Graph{},
			tree:         nil,
			res:          nil,
			err:          errors.New("(graph::generateGraphTemplateRec) Tree is null"),
		},
		{
			desc:         "Testing generate subgraph from an image",
			imageName:    "ubuntu",
			imageVersion: "16.04",
			image: &image.Image{
				Type: "infrastructure",
				Children: map[string][]string{
					"php-fpm": {
						"7.1",
						"7.2",
					},
				},
			},
			parent: nil,
			graph:  &gdsexttree.Graph{},
			tree:   imagesTree,
			res: &gdsexttree.Graph{
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
					"php-fpm:7.2": {
						Name: "php-fpm:7.2",
						Parents: []*gdsexttree.Node{
							{
								Name: "ubuntu:16.04",
							},
						},
					},
					"php-fpm-dev:7.1": {
						Name: "php-fpm-dev:7.1",
						Parents: []*gdsexttree.Node{
							{
								Name: "ubuntu:16.04",
							},
						},
					},
					"php-fpm-dev:7.2": {
						Name: "php-fpm:7.2",
						Parents: []*gdsexttree.Node{
							{
								Name: "ubuntu:16.04",
							},
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		err := generateGraphTemplateRec(test.imageName, test.imageVersion, test.image, test.parent, test.graph, test.tree)

		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, len(test.res.Root), len(test.graph.Root), "Unexpected lenght of root elements")
			assert.Equal(t, len(test.res.NodesIndex), len(test.graph.NodesIndex), "Unexpected lenght of nodes index elements")
		}
	}
}
