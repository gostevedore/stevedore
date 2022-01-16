package graph

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	"github.com/stretchr/testify/assert"
)

func TestNewGraphTemplateAddNode(t *testing.T) {
	tests := []struct {
		desc  string
		graph *GraphTemplate
		node  GraphTemplateNoder
		res   *GraphTemplate
		err   error
	}{
		{
			desc:  "Testing add node to GraphTemplate",
			graph: NewGraphTemplate(),
			node: &GraphTemplateNode{&gdsexttree.Node{
				Name: "node",
			}},
			res: &GraphTemplate{&gdsexttree.Graph{
				Root: []*gdsexttree.Node{
					{
						Name: "node",
					},
				},
				NodesIndex: map[string]*gdsexttree.Node{
					"node": {
						Name: "node",
					},
				},
			}},
			err: &errors.Error{},
		},
		{
			desc: "Testing error adding existing node to GraphTemplate",
			graph: &GraphTemplate{&gdsexttree.Graph{
				Root: []*gdsexttree.Node{
					{
						Name: "node",
					},
				},
				NodesIndex: map[string]*gdsexttree.Node{
					"node": {
						Name: "node",
					},
				},
			}},
			node: &GraphTemplateNode{&gdsexttree.Node{
				Name: "node",
			}},
			err: errors.New("", "Node 'node' already exists on the graph"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.graph.AddNode(test.node)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.graph, test.res)
			}
		})
	}
}

func TestNewGraphTemplateAddRelationship(t *testing.T) {
	tests := []struct {
		desc        string
		graph       *GraphTemplate
		parent      *GraphTemplateNode
		child       *GraphTemplateNode
		prepareTest func(*GraphTemplate, *GraphTemplateNode, *GraphTemplateNode)
		err         error
	}{
		{
			desc:   "Testing add relationship to GraphTemplate node",
			graph:  NewGraphTemplate(),
			parent: NewGraphTemplateNode("parent"),
			child:  NewGraphTemplateNode("child"),
			prepareTest: func(graph *GraphTemplate, parent, child *GraphTemplateNode) {
				graph.AddNode(parent)
				graph.AddNode(child)
			},
			err: &errors.Error{},
		},
		{
			desc:   "Testing error adding relationship to an unexisting parent",
			graph:  NewGraphTemplate(),
			parent: NewGraphTemplateNode("parent"),
			child:  NewGraphTemplateNode("child"),
			prepareTest: func(graph *GraphTemplate, parent, child *GraphTemplateNode) {
				graph.AddNode(child)
			},
			err: errors.New("", "Parent does not exist"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareTest != nil {
				test.prepareTest(test.graph, test.parent, test.child)
			}

			err := test.graph.AddRelationship(test.parent, test.child)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, len(test.graph.Root), 1)
				assert.Equal(t, len(test.parent.Children), 1)
				assert.Equal(t, len(test.child.Parents), 1)
			}
		})
	}
}

func TestNewGraphTemplateExist(t *testing.T) {
	tests := []struct {
		desc        string
		graph       *GraphTemplate
		name        string
		prepareTest func(*GraphTemplate)
		res         bool
		err         error
	}{
		{
			desc:  "Testing whether a node exists GraphTemplate, and exists",
			graph: NewGraphTemplate(),
			name:  "node",
			prepareTest: func(graph *GraphTemplate) {
				graph.AddNode(NewGraphTemplateNode("node"))
			},
			res: true,
			err: &errors.Error{},
		},
		{
			desc:        "Testing whether a node exists GraphTemplate, and not exists",
			graph:       NewGraphTemplate(),
			name:        "node",
			prepareTest: func(graph *GraphTemplate) {},
			res:         false,
			err:         &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareTest != nil {
				test.prepareTest(test.graph)
			}

			exists := test.graph.Exists(test.name)
			assert.Equal(t, exists, test.res)
		})
	}
}
