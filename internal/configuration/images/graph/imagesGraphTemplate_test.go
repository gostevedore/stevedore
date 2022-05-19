package graph

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration/images/image"
	"github.com/gostevedore/stevedore/internal/images/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddImage(t *testing.T) {
	errContext := "(graph::AddImage)"

	tests := []struct {
		desc              string
		name              string
		version           string
		image             *image.Image
		graph             *ImagesGraphTemplate
		res               *ImagesGraphTemplate
		prepareAssertFunc func(*ImagesGraphTemplate, *image.Image)
		assertFunc        func(*testing.T, *ImagesGraphTemplate, *image.Image)
		err               error
	}{
		{
			desc:  "Testing error adding an image when name is empty",
			name:  "",
			graph: &ImagesGraphTemplate{},
			err:   errors.New(errContext, "To add an image, the name must be specified"),
		},
		{
			desc:    "Testing error adding an image when version is empty",
			name:    "image_name",
			version: "",
			graph:   &ImagesGraphTemplate{},
			err:     errors.New(errContext, "To add an image, the version must be specified"),
		},
		{
			desc:    "Testing error adding an image when image is empty",
			name:    "image_name",
			version: "image_version",
			image:   nil,
			graph:   &ImagesGraphTemplate{},
			err:     errors.New(errContext, "To add and image, image must be provided"),
		},
		{
			desc:    "Testing add an image with parents and children that does not exists on graph nor pending nodes",
			name:    "image_name",
			version: "image_version",
			image: &image.Image{
				Name:    "image_name",
				Version: "image_version",
				Parents: map[string][]string{
					"parent_name": {
						"parent_version",
					},
				},
				Children: map[string][]string{
					"child_name": {
						"child_version",
					},
				},
			},
			graph: NewImagesGraphTemplate(graph.NewGraphTemplateFactory(true)),
			prepareAssertFunc: func(g *ImagesGraphTemplate, i *image.Image) {
				// node
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("image_name", "image_version")).Return(nil)
				node := g.graphFactory.NewGraphTemplateNode(generateNodeName("image_name", "image_version"))
				g.graph.(*graph.MockGraphTemplate).On("AddNode", node).Return(nil)

				// parents
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("parent_name", "parent_version")).Return(nil)
				parent := g.graphFactory.NewGraphTemplateNode(generateNodeName("parent_name", "parent_version"))
				g.graph.(*graph.MockGraphTemplate).On("AddNode", parent).Return(nil)
				g.graph.(*graph.MockGraphTemplate).On("AddRelationship", parent, mock.Anything).Return(nil)

				// children
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("child_name", "child_version")).Return(nil)
				child := g.graphFactory.NewGraphTemplateNode(generateNodeName("child_name", "child_version"))
				g.graph.(*graph.MockGraphTemplate).On("AddNode", child).Return(nil)
				g.graph.(*graph.MockGraphTemplate).On("AddRelationship", mock.Anything, child).Return(nil)

				g.graph.(*graph.MockGraphTemplate).On("HasCycles").Return(false)
			},
			assertFunc: func(t *testing.T, g *ImagesGraphTemplate, i *image.Image) {
				g.graph.(*graph.MockGraphTemplate).AssertExpectations(t)
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing add an image with parents and children already added to graph",
			name:    "image_name",
			version: "image_version",
			image: &image.Image{
				Name:    "image_name",
				Version: "image_version",
				Parents: map[string][]string{
					"parent_name": {
						"parent_version",
					},
				},
				Children: map[string][]string{
					"child_name": {
						"child_version",
					},
				},
			},
			graph: NewImagesGraphTemplate(graph.NewGraphTemplateFactory(true)),
			prepareAssertFunc: func(g *ImagesGraphTemplate, i *image.Image) {
				// node
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("image_name", "image_version")).Return(nil)
				node := g.graphFactory.NewGraphTemplateNode(generateNodeName("image_name", "image_version"))
				g.graph.(*graph.MockGraphTemplate).On("AddNode", node).Return(nil)

				// parents
				parent := g.graphFactory.NewGraphTemplateNode(generateNodeName("parent_name", "parent_version"))
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("parent_name", "parent_version")).Return(parent)
				g.graph.(*graph.MockGraphTemplate).On("AddRelationship", parent, mock.Anything).Return(nil)

				// children
				child := g.graphFactory.NewGraphTemplateNode(generateNodeName("child_name", "child_version"))
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("child_name", "child_version")).Return(child)
				g.graph.(*graph.MockGraphTemplate).On("AddRelationship", mock.Anything, child).Return(nil)

				g.graph.(*graph.MockGraphTemplate).On("HasCycles").Return(false)
			},

			assertFunc: func(t *testing.T, g *ImagesGraphTemplate, i *image.Image) {
				g.graph.(*graph.MockGraphTemplate).AssertExpectations(t)
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing error adding node that creates a cycles",
			name:    "image_name",
			version: "image_version",
			image: &image.Image{
				Name:    "image_name",
				Version: "image_version",
			},
			graph: NewImagesGraphTemplate(graph.NewGraphTemplateFactory(true)),
			prepareAssertFunc: func(g *ImagesGraphTemplate, i *image.Image) {
				// node
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("image_name", "image_version")).Return(nil)
				node := g.graphFactory.NewGraphTemplateNode(generateNodeName("image_name", "image_version"))
				g.graph.(*graph.MockGraphTemplate).On("AddNode", node).Return(nil)
				g.graph.(*graph.MockGraphTemplate).On("HasCycles").Return(true)
			},
			err: errors.New("", "Detected a cycle in the graph template after adding node 'image_name:image_version'"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.graph, test.image)
			}

			err := test.graph.AddImage(test.name, test.version, test.image)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {

				if test.assertFunc != nil {
					test.assertFunc(t, test.graph, test.image)
				} else {
					test.graph.graph.(*graph.MockGraphTemplate).AssertExpectations(t)
				}
			}
		})
	}
}

func TestIterate(t *testing.T) {
	t.Log("Testing Iterate")

	i01 := &image.Image{}
	i011 := &image.Image{}
	i0111 := &image.Image{}
	i0112 := &image.Image{}
	i012 := &image.Image{}
	i02 := &image.Image{}
	i021 := &image.Image{}
	i0211 := &image.Image{}
	i0212 := &image.Image{}
	i022 := &image.Image{}

	graph := NewImagesGraphTemplate(graph.NewGraphTemplateFactory(false))
	graph.AddImage("i01", "1.0.0", i01)
	graph.AddImage("i011", "1.0.1", i011)
	graph.AddImage("i0111", "1.0.11", i0111)
	graph.AddImage("i0112", "1.0.12", i0112)
	graph.AddImage("i012", "1.0.2", i012)
	graph.AddImage("i02", "2.0.0", i02)
	graph.AddImage("i021", "2.0.1", i021)
	graph.AddImage("i0211", "2.0.11", i0211)
	graph.AddImage("i0212", "2.0.12", i0212)
	graph.AddImage("i022", "2.0.2", i022)

	numNodes := 0
	for range graph.Iterate() {
		numNodes++
	}

	assert.Equal(t, numNodes, 10)
}

func TestParseNodeName(t *testing.T) {

	f := graph.NewGraphTemplateFactory(false)

	tests := []struct {
		desc    string
		node    GraphNoder
		name    string
		version string
		err     error
	}{
		{
			desc:    "Testing parsing node name",
			node:    f.NewGraphTemplateNode("image_name:image_version"),
			name:    "image_name",
			version: "image_version",
			err:     &errors.Error{},
		},
		{
			desc:    "Testing parsing node name with invalid node",
			node:    f.NewGraphTemplateNode("image_name"),
			name:    "",
			version: "",
			err:     errors.New("", "Node name 'image_name' is not valid"),
		},
		{
			desc:    "Testing parsing node with undefined name",
			node:    f.NewGraphTemplateNode(""),
			name:    "",
			version: "",
			err:     errors.New("", "Node name is undefined"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			name, version, err := ParseNodeName(test.node)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.name, name)
				assert.Equal(t, test.version, version)
			}
		})
	}

}
