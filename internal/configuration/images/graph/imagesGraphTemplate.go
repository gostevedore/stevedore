package graph

import (
	"fmt"
	"strings"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration/images/image"
	"github.com/gostevedore/stevedore/internal/images/graph"
)

// ImageGraphTemplate is a graph template for images
type ImagesGraphTemplate struct {
	graph        Grapher
	graphFactory graph.GraphTemplateFactory

	mutex sync.RWMutex
	// pendingNodes map[string]map[string]GraphNoder
}

// NewImagesGraphTemplate creates a new graph template for images
func NewImagesGraphTemplate(factory graph.GraphTemplateFactory) *ImagesGraphTemplate {
	return &ImagesGraphTemplate{
		graph:        factory.NewGraphTemplate(),
		graphFactory: factory,
		// pendingNodes: make(map[string]map[string]GraphNoder),
	}
}

// generateNodeName generates a node name
func generateNodeName(name, version string) string {
	return fmt.Sprintf("%s:%s", name, version)
}

func ParseNodeName(node GraphNoder) (string, string, error) {

	errContext := "(graph::ParseNodeName)"

	name := node.Name()
	if name == "" {
		return "", "", errors.New(errContext, "Node name is undefined")
	}

	idx := strings.IndexRune(name, ':')
	if idx == -1 {
		return "", "", errors.New(errContext, fmt.Sprintf("Node name '%s' is not valid", name))
	}

	return name[:idx], name[idx+1:], nil
}

// AddImage is a mock implementation of the AddImage method
func (m *ImagesGraphTemplate) AddImage(name, version string, image *image.Image) error {

	var err error
	var node GraphNoder

	errContext := "(graph::AddImage)"

	if name == "" {
		return errors.New(errContext, "To add an image, the name must be specified")
	}

	if version == "" {
		return errors.New(errContext, "To add an image, the version must be specified")
	}

	if image == nil {
		return errors.New(errContext, "To add and image, image must be provided")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	node = m.graph.GetNode(generateNodeName(name, version))
	if node == nil {
		node = m.graphFactory.NewGraphTemplateNode(generateNodeName(name, version))
		err = m.graph.AddNode(node)
		if err != nil {
			return errors.New(errContext, err.Error())
		}
	}
	node.AddItem(image)

	if len(image.Parents) > 0 {
		for parentName, versions := range image.Parents {
			for _, parentVersion := range versions {

				parentNode := m.graph.GetNode(generateNodeName(parentName, parentVersion))
				if parentNode == nil {
					parentNode = m.graphFactory.NewGraphTemplateNode(generateNodeName(parentName, parentVersion))
					err = m.graph.AddNode(parentNode)
					if err != nil {
						return errors.New(errContext, err.Error())
					}
				}

				err = m.graph.AddRelationship(parentNode, node)
				if err != nil {
					return errors.New(errContext, err.Error())
				}
			}
		}
	}

	if len(image.Children) > 0 {
		for childName, versions := range image.Children {
			for _, childVersion := range versions {

				childNode := m.graph.GetNode(generateNodeName(childName, childVersion))
				if childNode == nil {
					childNode = m.graphFactory.NewGraphTemplateNode(generateNodeName(childName, childVersion))
					err = m.graph.AddNode(childNode)
					if err != nil {
						return errors.New(errContext, err.Error())
					}
				}

				err = m.graph.AddRelationship(node, childNode)
				if err != nil {
					return errors.New(errContext, err.Error())
				}
			}
		}
	}

	if m.graph.HasCycles() {
		return errors.New(errContext, fmt.Sprintf("Detected a cycle in the graph template after adding node '%s'", generateNodeName(name, version)))
	}

	return nil
}

// Validate validates the graph template
func (m *ImagesGraphTemplate) Validate() error {
	errContext := "(graph::Validate)"

	if m.graph.HasCycles() {
		return errors.New(errContext, "Graph template has cycles")
	}
	return nil
}

// Iterate iterates over the graph template
func (m *ImagesGraphTemplate) Iterate() <-chan GraphNoder {
	it := make(chan GraphNoder)

	go func() {
		defer close(it)
		for node := range m.graph.Iterate() {
			it <- node
		}
	}()

	return it
}
