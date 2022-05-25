package graph

import (
	// gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	gdsexttree "github.com/gostevedore/stevedore/pkg/extendedTree"
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

func (n *GraphTemplateNode) getNode() *gdsexttree.Node {
	return n.Node
}

// AddChild adds a child to the graph template node
func (n *GraphTemplateNode) AddChild(child GraphTemplateNoder) error {
	return n.Node.AddChild(child.(*GraphTemplateNode).Node)
}

// AddParent adds a parent to the graph template node
func (n *GraphTemplateNode) AddParent(parent GraphTemplateNoder) error {
	return n.Node.AddParent(parent.(*GraphTemplateNode).Node)
}

// AddItem adds a item to the graph template node
func (n *GraphTemplateNode) AddItem(item interface{}) {
	n.Node.Item = item
}

// Name returns the name of the graph template node
func (n *GraphTemplateNode) Name() string {
	return n.Node.Name
}

// Item return node's item
func (n *GraphTemplateNode) Item() interface{} {
	return n.Node.Item
}

// Parents returns the parents of the graph template node
func (n *GraphTemplateNode) Parents() []GraphTemplateNoder {
	var parents []GraphTemplateNoder
	for _, parent := range n.Node.Parents {
		parents = append(parents, &GraphTemplateNode{parent})
	}
	return parents
}

// Children returns the children of the graph template node
func (n *GraphTemplateNode) Children() []GraphTemplateNoder {
	children := make([]GraphTemplateNoder, len(n.Node.Children))
	for i, child := range n.Node.Children {
		children[i] = &GraphTemplateNode{child}
	}
	return children
}
