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
	//addedNode    map[string]map[string]struct{}
	pendingNodes map[string]map[string]GraphNoder
}

// NewImagesGraphTemplate creates a new graph template for images
func NewImagesGraphTemplate(factory graph.GraphTemplateFactory) *ImagesGraphTemplate {
	return &ImagesGraphTemplate{
		graph:        factory.NewGraphTemplate(),
		graphFactory: factory,
		//addedNode:    make(map[string]map[string]struct{}),
		pendingNodes: make(map[string]map[string]GraphNoder),
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
	var isPendingNode bool

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

	if m.pendingNodes == nil {
		m.pendingNodes = make(map[string]map[string]GraphNoder)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.graph.Exists(generateNodeName(name, version)) {
		return errors.New(errContext, fmt.Sprintf("Image '%s:%s' already added to images graph template", name, version))
	}

	node, isPendingNode = m.pendingNodes[name][version]
	if !isPendingNode {
		node = m.graphFactory.NewGraphTemplateNode(generateNodeName(name, version))
	} else {
		delete(m.pendingNodes[name], version)
		if len(m.pendingNodes[name]) <= 0 {
			delete(m.pendingNodes, name)
		}
	}
	node.AddItem(image)

	err = m.graph.AddNode(node)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	if m.graph.HasCycles() {
		return errors.New(errContext, fmt.Sprintf("Detected a cycle in the graph template after adding node '%s'", generateNodeName(name, version)))
	}

	if len(image.Parents) > 0 {
		for parentName, versions := range image.Parents {
			for _, parentVersion := range versions {
				parentNode := m.achieveGraphNode(parentName, parentVersion)
				err = node.AddParent(parentNode)
				if err != nil {
					return errors.New(errContext, err.Error())
				}
			}
		}
	}

	if len(image.Children) > 0 {
		for childName, versions := range image.Children {
			for _, childVersion := range versions {
				childNode := m.achieveGraphNode(childName, childVersion)
				node.AddChild(childNode)
				if err != nil {
					return errors.New(errContext, err.Error())
				}
			}
		}
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

// achieveGraphNode creates a new node
func (m *ImagesGraphTemplate) achieveGraphNode(name, version string) GraphNoder {
	var node GraphNoder
	var exists bool

	// if node is on pending nodes, use that node
	// if node is already defined on graph, use that node
	// otherwise, create a new node and add it to pending nodes, and use that node
	node, exists = m.pendingNodes[name][version]
	if !exists {
		node = m.graph.GetNode(generateNodeName(name, version))
		// when node does not exist
		if node == nil {
			node = m.graphFactory.NewGraphTemplateNode(generateNodeName(name, version))
			m.addNodeToPendingNodes(name, version, node)
		}
	}

	return node
}

// addNodeToPendingNodes adds a node to the pending nodes
func (m *ImagesGraphTemplate) addNodeToPendingNodes(name, version string, node GraphNoder) {

	if m.pendingNodes[name] == nil {
		m.pendingNodes[name] = make(map[string]GraphNoder)
	}

	m.pendingNodes[name][version] = node
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
