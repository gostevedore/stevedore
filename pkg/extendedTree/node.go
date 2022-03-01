package gdsexttree

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
)

// Node is the extended tree graph node
type Node struct {
	Name     string
	Item     interface{}
	Parents  []*Node
	Children []*Node
}

// AddParent method update node's parents list adding a new one. It also update parent's childs list
func (n *Node) AddParent(parent *Node) error {
	if n == nil {
		return errors.New("(graph::AddParent)", "Adding parent to a nil node")
	}

	if parent == nil {
		return errors.New("(graph::AddParent)", "Adding nil parent to node")
	}

	if n.Parents == nil || len(n.Parents) == 0 {
		n.Parents = []*Node{}
	}

	if !n.HasParent(parent) {
		n.Parents = append(n.Parents, parent)
	} else {
		return errors.New("(graph::AddParent)", fmt.Sprintf("Parent '%s' already exists to '%s'", parent.Name, n.Name))
	}

	// node node a parent childe
	err := parent.AddChild(n)
	if err != nil {
		return errors.New("(graph:AddParent)", fmt.Sprintf("Child could not be add to '%s'", parent.Name), err)
	}

	return nil
}

// AddChild method update node's childs list adding a new one
func (n *Node) AddChild(child *Node) error {
	if n == nil {
		return errors.New("(graph::AddChild)", "Adding child to a nil node")
	}

	if child == nil {
		return errors.New("(graph::AddChild)", "Adding nil parent to node")
	}

	if n.Children == nil || len(n.Children) == 0 {
		n.Children = []*Node{}
	}

	if !n.HasChild(child) {
		n.Children = append(n.Children, child)
	}

	return nil
}

// HasChild method validate whether a child node already exists in node's child list. Two nodes are equal when they have the same node name
func (n *Node) HasChild(child *Node) bool {
	hasChild := false
	for _, c := range n.Children {
		if c.Name == child.Name {
			return true
		}
	}

	return hasChild
}

// HasParent method validate whether a parent node already exists in node's parent list. Two nodes are equal when they have the same node name
func (n *Node) HasParent(parent *Node) bool {

	for _, p := range n.Parents {
		if p.Name == parent.Name {
			return true
		}
	}

	return false
}
