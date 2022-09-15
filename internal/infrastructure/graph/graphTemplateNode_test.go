package graph

import (
	"testing"

	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	"github.com/stretchr/testify/assert"
)

func TestGetNode(t *testing.T) {
	t.Log("Testing GraphTemplateNode node")

	node := NewGraphTemplateNode("node")
	res := &gdsexttree.Node{
		Name: "node",
	}

	assert.Equal(t, res, node.getNode())

}

func TestAddChild(t *testing.T) {
	t.Log("Testing add child to GraphTemplateNode")

	parent := NewGraphTemplateNode("parent")
	child := NewGraphTemplateNode("child")

	parent.AddChild(child)
	assert.Equal(t, 1, len(parent.Children()))
}

func TestAddParent(t *testing.T) {
	t.Log("Testing add parent to GraphTemplateNode")

	parent := NewGraphTemplateNode("parent")
	child := NewGraphTemplateNode("child")

	child.AddParent(parent)
	assert.Equal(t, 1, len(child.Parents()))
}

func TestAddItem(t *testing.T) {
	t.Log("Testing add item to GraphTemplateNode")

	node := NewGraphTemplateNode("node")
	node.AddItem("item")
	assert.Equal(t, "item", node.Item())

}
