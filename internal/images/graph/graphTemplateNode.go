package graph

import (
	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
)

// GraphTemplateNode is a node for the graph template
type GraphTemplateNode struct {
	*gdsexttree.Node
}

// NewGraphTemplateNode creates a new graph template node
func NewGraphTemplateNode(name string) *GraphTemplateNode {
	return &GraphTemplateNode{&gdsexttree.Node{
		Name: name,
	}}
}

func (m *GraphTemplateNode) getNode() *gdsexttree.Node {
	return m.Node
}

// AddChild adds a child to the graph template node
func (m *GraphTemplateNode) AddChild(child GraphTemplateNoder) error {
	return m.Node.AddChild(child.(*GraphTemplateNode).Node)
}

// AddParent adds a parent to the graph template node
func (m *GraphTemplateNode) AddParent(parent GraphTemplateNoder) error {
	return m.Node.AddParent(parent.(*GraphTemplateNode).Node)
}

// AddItem adds a item to the graph template node
func (m *GraphTemplateNode) AddItem(item interface{}) {
	m.Node.Item = item
}
