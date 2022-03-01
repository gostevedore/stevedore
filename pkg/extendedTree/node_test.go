package gdsexttree

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

// TestAddParent
func TestAddParent(t *testing.T) {

	nodeChild := &Node{
		Name:     "node",
		Parents:  nil,
		Children: nil,
		Item:     nil,
	}

	nodeParent := &Node{
		Name:     "parent",
		Parents:  nil,
		Children: nil,
		Item:     nil,
	}

	tests := []struct {
		desc   string
		node   *Node
		parent *Node
		err    error
		res    *Node
	}{
		{
			desc:   "Add parent node",
			node:   nodeChild,
			parent: nodeParent,
			err:    nil,
			res: &Node{
				Name:     "node",
				Children: nil,
				Item:     nil,
				Parents: []*Node{
					{
						Name:    "parent",
						Parents: nil,
						Children: []*Node{
							{
								Name: "node",
							},
						},
						Item: nil,
					},
				},
			},
		},
		{
			desc:   "Add parent to nil node",
			node:   nil,
			parent: nodeParent,
			err:    errors.New("(graph::AddParent)", "Adding parent to a nil node"),
			res:    nil,
		},
		{
			desc:   "Add nil parent to node",
			parent: nil,
			node:   nodeChild,
			err:    errors.New("(graph::AddParent)", "Adding nil parent to node"),
			res:    nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		err := test.node.AddParent(test.parent)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.res.Name, test.node.Name, "Name not equal")
			assert.Equal(t, len(test.res.Parents), len(test.node.Parents), "Parent name not equal")
		}
	}
}

// TestAddChild
func TestAddChild(t *testing.T) {

	nodeChild := &Node{
		Name:     "node",
		Parents:  nil,
		Children: nil,
		Item:     nil,
	}

	nodeChild2 := &Node{
		Name:     "node2",
		Parents:  nil,
		Children: nil,
		Item:     nil,
	}

	nodeParent := &Node{
		Name:     "parent",
		Parents:  nil,
		Children: nil,
		Item:     nil,
	}

	nodeParent2 := &Node{
		Name:    "parent2",
		Parents: nil,
		Children: []*Node{
			nodeChild,
		},
		Item: nil,
	}

	tests := []struct {
		desc   string
		node   *Node
		parent *Node
		err    error
		res    *Node
	}{
		{
			desc:   "Add child to node",
			node:   nodeChild,
			parent: nodeParent,
			err:    nil,
			res: &Node{
				Name: "parent",
				Children: []*Node{
					{
						Name:     "node",
						Parents:  nil,
						Children: nil,
						Item:     nil,
					},
				},
				Item:    nil,
				Parents: nil,
			},
		},
		{
			desc:   "Add second child to node",
			node:   nodeChild2,
			parent: nodeParent2,
			err:    nil,
			res: &Node{
				Name: "parent2",
				Children: []*Node{
					nodeChild,
					nodeChild2,
				},
				Item:    nil,
				Parents: nil,
			},
		},
		{
			desc:   "Add child to nil parent",
			parent: nil,
			node:   nodeChild,
			err:    errors.New("(graph::AddChild)", "Adding child to a nil node"),
			res:    nil,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		err := test.parent.AddChild(test.node)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.res, test.parent, "Nodes not equal")
		}
	}
}

func TestHasChild(t *testing.T) {

	nodeChild := &Node{
		Name:     "node",
		Parents:  nil,
		Children: nil,
		Item:     nil,
	}

	nodeParent := &Node{
		Name:     "parent",
		Parents:  nil,
		Children: nil,
		Item:     nil,
	}

	tests := []struct {
		desc   string
		node   *Node
		parent *Node
		err    error
		res    bool
	}{
		{
			desc:   "Node is not a child",
			node:   nodeChild,
			parent: nodeParent,
			err:    nil,
			res:    false,
		},

		{
			desc: "Node is not a child",
			node: nodeChild,
			parent: &Node{
				Name:    "parent",
				Parents: nil,
				Children: []*Node{
					nodeChild,
				},
				Item: nil,
			},
			err: nil,
			res: true,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		has := test.parent.HasChild(test.node)
		assert.Equal(t, test.res, has, "Nodes not equal")

	}
}

func TestHasParent(t *testing.T) {

	nodeParent := &Node{
		Name:     "parent",
		Parents:  nil,
		Children: nil,
		Item:     nil,
	}

	nodeParent2 := &Node{
		Name:     "parent2",
		Parents:  nil,
		Children: nil,
		Item:     nil,
	}

	nodeChild := &Node{
		Name: "node",
		Parents: []*Node{
			nodeParent,
		},
		Children: nil,
		Item:     nil,
	}

	tests := []struct {
		desc   string
		node   *Node
		parent *Node
		err    error
		res    bool
	}{
		{
			desc:   "Testing node which has a gived parent",
			node:   nodeChild,
			parent: nodeParent,
			err:    nil,
			res:    true,
		},
		{
			desc:   "Testing node which does not have a gived parent",
			node:   nodeChild,
			parent: nodeParent2,
			err:    nil,
			res:    false,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		has := test.node.HasParent(test.parent)
		assert.Equal(t, test.res, has, "Nodes not equal")

	}
}
