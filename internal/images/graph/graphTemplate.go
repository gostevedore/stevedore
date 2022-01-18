package graph

import (
	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
)

// GraphTemplate holds the graph template for images
type GraphTemplate struct {
	*gdsexttree.Graph
}

// NewGraphTemplate creates a new graph template for images
func NewGraphTemplate() *GraphTemplate {
	return &GraphTemplate{
		&gdsexttree.Graph{},
	}
}

// GetNode gets a node from the graph template
func (m *GraphTemplate) GetNode(name string) GraphTemplateNoder {
	node, _ := m.Graph.GetNode(name)
	return &GraphTemplateNode{node}
}

// AddNode adds a node to the graph template
func (m *GraphTemplate) AddNode(node GraphTemplateNoder) error {
	return m.Graph.AddNode(node.(*GraphTemplateNode).Node)
}

// AddRelationship adds a relationship to the graph template
func (m *GraphTemplate) AddRelationship(parent, child GraphTemplateNoder) error {

	return m.Graph.AddRelationship(parent.(*GraphTemplateNode).Node, child.(*GraphTemplateNode).Node)
}

// Exists checks if a node exists in the graph template
func (m *GraphTemplate) Exists(name string) bool {
	node, _ := m.Graph.GetNode(name)
	return node != nil
}

// HasCycles checks if the graph template has cycles
func (m *GraphTemplate) HasCycles() bool {
	return m.Graph.HasCycles()
}
