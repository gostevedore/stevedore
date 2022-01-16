package graph

import (
	"github.com/gostevedore/stevedore/internal/images/graph"
)

// Grapher is a graph template for images
type Grapher interface {
	GetNode(string) graph.GraphTemplateNoder
	AddNode(graph.GraphTemplateNoder) error
	AddRelationship(graph.GraphTemplateNoder, graph.GraphTemplateNoder) error
	Exists(string) bool
}

// GraphNoder is a node for the graph template
type GraphNoder interface {
	AddChild(graph.GraphTemplateNoder) error
	AddParent(graph.GraphTemplateNoder) error
	AddItem(interface{})
}
