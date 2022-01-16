package graph

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration/images/image"
	"github.com/gostevedore/stevedore/internal/images/graph"
	"github.com/stretchr/testify/assert"
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
			desc:    "Testing add an image with parents and children neither added to graph nor pending",
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
			graph: NewImagesGraphTemplate(*graph.NewGraphTemplateFactory(true)),
			prepareAssertFunc: func(g *ImagesGraphTemplate, i *image.Image) {
				// node
				g.graph.(*graph.MockGraphTemplate).On("Exists", "image_name:image_version").Return(false, nil)
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("image_name", "image_version")).Return(nil)
				node := g.graphFactory.NewGraphTemplateNode(generateNodeName("image_name", "image_version"))
				node.AddItem(i)
				g.graph.(*graph.MockGraphTemplate).On("AddNode", node).Return(nil)

				// parents
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("parent_name", "parent_version")).Return(nil)

				// children
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("child_name", "child_version")).Return(nil)

			},
			assertFunc: func(t *testing.T, g *ImagesGraphTemplate, i *image.Image) {
				g.graph.(*graph.MockGraphTemplate).AssertExpectations(t)

				node := g.graphFactory.NewGraphTemplateNode(generateNodeName("image_name", "image_version"))
				node.AddItem(i)

				parent := g.graphFactory.NewGraphTemplateNode(generateNodeName("parent_name", "parent_version"))
				node.AddParent(parent)

				child := g.graphFactory.NewGraphTemplateNode(generateNodeName("child_name", "child_version"))
				node.AddChild(child)

				res := &ImagesGraphTemplate{
					pendingNodes: map[string]map[string]GraphNoder{
						"parent_name": {
							"parent_version": parent,
						},
						"child_name": {
							"child_version": child,
						},
					},
				}

				assert.Equal(t, res.pendingNodes, g.pendingNodes)
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
			graph: NewImagesGraphTemplate(*graph.NewGraphTemplateFactory(true)),
			prepareAssertFunc: func(g *ImagesGraphTemplate, i *image.Image) {

				// node
				g.graph.(*graph.MockGraphTemplate).On("Exists", "image_name:image_version").Return(false, nil)
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("image_name", "image_version")).Return(nil)
				node := g.graphFactory.NewGraphTemplateNode(generateNodeName("image_name", "image_version"))
				node.AddItem(i)
				g.graph.(*graph.MockGraphTemplate).On("AddNode", node).Return(nil)

				// parents
				parent := g.graphFactory.NewGraphTemplateNode(generateNodeName("parent_name", "parent_version"))
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("parent_name", "parent_version")).Return(parent)

				// children
				child := g.graphFactory.NewGraphTemplateNode(generateNodeName("child_name", "child_version"))
				g.graph.(*graph.MockGraphTemplate).On("GetNode", generateNodeName("child_name", "child_version")).Return(child)
			},

			assertFunc: func(t *testing.T, g *ImagesGraphTemplate, i *image.Image) {
				g.graph.(*graph.MockGraphTemplate).AssertExpectations(t)

				res := &ImagesGraphTemplate{
					pendingNodes: map[string]map[string]GraphNoder{},
				}

				assert.Equal(t, res.pendingNodes, g.pendingNodes)
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing add an image defined on pending nodes",
			name:    "image_name",
			version: "image_version",
			image: &image.Image{
				Name:    "image_name",
				Version: "image_version",
			},
			graph: NewImagesGraphTemplate(*graph.NewGraphTemplateFactory(true)),
			prepareAssertFunc: func(g *ImagesGraphTemplate, i *image.Image) {

				node := g.graphFactory.NewGraphTemplateNode(generateNodeName("image_name", "image_version"))
				g.addNodeToPendingNodes("image_name", "image_version", node)

				// node
				g.graph.(*graph.MockGraphTemplate).On("Exists", "image_name:image_version").Return(false, nil)

			},
			assertFunc: func(t *testing.T, g *ImagesGraphTemplate, i *image.Image) {
				g.graph.(*graph.MockGraphTemplate).AssertExpectations(t)

				res := &ImagesGraphTemplate{
					pendingNodes: map[string]map[string]GraphNoder{},
				}

				assert.Equal(t, res.pendingNodes, g.pendingNodes)
			},
			err: &errors.Error{},
		},

		{
			desc:    "Testing add an existing image",
			name:    "image_name",
			version: "image_version",
			image: &image.Image{
				Name:    "image_name",
				Version: "image_version",
			},
			graph: NewImagesGraphTemplate(*graph.NewGraphTemplateFactory(true)),
			prepareAssertFunc: func(g *ImagesGraphTemplate, i *image.Image) {
				g.graph.(*graph.MockGraphTemplate).On("Exists", "image_name:image_version").Return(true)
			},
			err: errors.New("", "Image 'image_name:image_version' already added to images graph template"),
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