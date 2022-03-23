package gdsexttree

import (
	"bytes"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

// TestLoadImage tests
func TestAddNode(t *testing.T) {

	node := &Node{
		Name:     "node",
		Parents:  nil,
		Item:     "string item",
		Children: nil,
	}

	node2 := &Node{
		Name: "node2",
		Parents: []*Node{
			node,
		},
		Item:     "string item",
		Children: nil,
	}

	node3 := &Node{
		Name:     "node3",
		Parents:  nil,
		Item:     "string item",
		Children: nil,
	}

	tests := []struct {
		desc  string
		graph *Graph
		node  *Node
		err   error
		res   *Graph
	}{
		{
			desc:  "Add Node to nil graph",
			graph: nil,
			node:  node,
			err:   errors.New("(graph::AddNode)", "Adding a node to a nil graph"),
			res:   nil,
		},
		{
			desc:  "Add Node to an empty graph",
			graph: &Graph{},
			node:  node,
			err:   nil,
			res: &Graph{
				Root: []*Node{
					node,
				},
				NodesIndex: map[string]*Node{
					"node": node,
				},
			},
		},
		{
			desc: "Add Node to an empty graph with and empty index",
			graph: &Graph{
				NodesIndex: map[string]*Node{},
			},
			node: node,
			err:  nil,
			res: &Graph{
				Root: []*Node{
					node,
				},
				NodesIndex: map[string]*Node{
					"node": node,
				},
			},
		},
		{
			desc: "Add an existing node to a graph",
			graph: &Graph{
				Root: []*Node{
					{
						Name: "root",
					},
				},
				NodesIndex: map[string]*Node{
					"root": {
						Name: "root",
					},
				},
			},
			node: &Node{
				Name: "root",
			},
			err: errors.New("(graph::AddNode)", "Node 'root' already exists on the graph"),
			res: nil,
		},
		{
			desc: "Add Node as a child to another node",
			graph: &Graph{
				Root: []*Node{
					node,
				},
				NodesIndex: map[string]*Node{
					"node": node,
				},
			},
			node: node2,
			err:  nil,
			res: &Graph{
				Root: []*Node{
					{
						Name:    "node",
						Parents: nil,
						Item:    "string item",
						Children: []*Node{
							{
								Name: "node2",
								Parents: []*Node{
									node,
								},
								Item:     "string item",
								Children: nil,
							},
						},
					},
				},
				NodesIndex: map[string]*Node{
					"node":  node,
					"node2": node2,
				},
			},
		},
		{
			desc: "Add Node without parent to a graph with one element",
			graph: &Graph{
				Root: []*Node{
					node,
				},
				NodesIndex: map[string]*Node{
					"node": node,
				},
			},
			node: node3,
			err:  nil,
			res: &Graph{
				Root: []*Node{
					{
						Name:     "node",
						Parents:  nil,
						Item:     "string item",
						Children: nil,
					},
					{
						Name:     "node3",
						Parents:  nil,
						Item:     "string item",
						Children: nil,
					},
				},
				NodesIndex: map[string]*Node{
					"node":  node,
					"node2": node3,
				},
			},
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.graph.AddNode(test.node)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, len(test.graph.Root), len(test.res.Root), "Root size not equal")
				assert.Equal(t, len(test.graph.NodesIndex), len(test.res.NodesIndex), "Root size not equal")
			}
		})
	}
}

// TestDrawGraph
func TestDrawGraph(t *testing.T) {
	tests := []struct {
		desc  string
		graph *Graph
		err   error
		res   string
	}{
		{
			desc: "Print one root",
			res: `|->  root
  |->  level1-sibling1
    |->  level2-sibling1
  |->  level1-sibling2
    |->  level2-sibling1
    |->  level2-sibling2
`,
			graph: &Graph{
				Root: []*Node{
					{
						Name: "root",
						Children: []*Node{
							{
								Name: "level1-sibling1",
								Children: []*Node{
									{
										Name:     "level2-sibling1",
										Children: nil,
									},
								},
							},
							{
								Name: "level1-sibling2",
								Children: []*Node{
									{
										Name:     "level2-sibling1",
										Children: nil,
									},
									{
										Name:     "level2-sibling2",
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	var w bytes.Buffer
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			w.Reset()
			test.graph.DrawGraph(&w)
			assert.Equal(t, test.res, w.String(), "Output not equal")
		})
	}
}

// TestAddRelationship
func TestAddRelationship(t *testing.T) {

	errContext := "(graph::AddRelationship)"

	tests := []struct {
		desc   string
		graph  *Graph
		parent *Node
		node   *Node
		err    error
		res    *Graph
	}{
		{
			desc:   "Add parent to node into a nil graph",
			graph:  nil,
			parent: nil,
			node:   nil,
			err:    errors.New(errContext, "Graph is null"),
			res:    nil,
		},
		{
			desc:   "Add node to a nil parent",
			graph:  &Graph{},
			parent: nil,
			node:   nil,
			err:    errors.New(errContext, "Parent is null"),
			res:    nil,
		},
		{
			desc: "Add nil node to a parent",
			graph: &Graph{
				Root: []*Node{
					{
						Name: "A",
					},
				},
				NodesIndex: map[string]*Node{
					"A": {
						Name: "A",
					},
				},
			},
			parent: &Node{
				Name: "A",
			},
			node: nil,
			err:  errors.New(errContext, "Child is null"),
			res:  nil,
		},
		{
			desc: "Add parent to orphan node",
			graph: &Graph{
				Root: []*Node{
					{
						Name: "A",
					},
					{
						Name:    "orphan",
						Parents: nil,
					},
				},
				NodesIndex: map[string]*Node{
					"A": {
						Name: "A",
					},
					"orphan": {
						Name:    "orphan",
						Parents: nil,
					},
				},
			},
			parent: &Node{
				Name: "A",
			},
			node: &Node{
				Name:    "orphan",
				Parents: nil,
			},
			err: nil,
			res: &Graph{
				Root: []*Node{
					{
						Name: "A",
					},
				},
				NodesIndex: map[string]*Node{
					"A": {
						Name: "A",
					},
					"orphan": {
						Name:    "orphan",
						Parents: nil,
					},
				},
			},
		},
		{
			desc: "Add parent unexistent parent to a node",
			graph: &Graph{
				Root: []*Node{
					{
						Name: "A",
					},
					{
						Name: "orphan",
					},
				},
				NodesIndex: map[string]*Node{
					"A": {
						Name: "A",
					},
					"orphan": {
						Name: "orphan",
					},
				},
			},
			parent: &Node{
				Name: "unexistent",
			},
			node: &Node{
				Name: "orphan",
			},
			err: errors.New(errContext, "Parent does not exist"),
			res: nil,
		},
		{
			desc: "Add parent parent to an unexistent node",
			graph: &Graph{
				Root: []*Node{
					{
						Name: "A",
					},
				},
				NodesIndex: map[string]*Node{
					"A": {
						Name: "A",
					},
				},
			},
			parent: &Node{
				Name: "A",
			},
			node: &Node{
				Name: "unexistent",
			},
			err: errors.New(errContext, "Child does not exist"),
			res: nil,
		},
		{
			desc: "Add relationship already defined",
			graph: &Graph{
				Root: []*Node{
					{
						Name: "A",
					},
					{
						Name: "B",
						Parents: []*Node{
							{
								Name: "A",
							},
						},
					},
				},
				NodesIndex: map[string]*Node{
					"A": {
						Name: "A",
					},
					"B": {
						Name: "B",
						Parents: []*Node{
							{
								Name: "A",
							},
						},
					},
				},
			},
			parent: &Node{
				Name: "A",
			},
			node: &Node{
				Name: "B",
				Parents: []*Node{
					{
						Name: "A",
					},
				},
			},
			err: errors.New(errContext, "Parent can not be added to 'B'", errors.New("(graph::AddParent)", "Parent 'A' already exists to 'B'")),
			res: nil,
		},
		{
			desc: "Add cyclic relationship",
			graph: &Graph{
				Root: []*Node{
					{
						Name: "A",
						Children: []*Node{
							{
								Name: "B",
								Parents: []*Node{
									{
										Name: "A",
									},
								},
							},
						},
					},
				},
				NodesIndex: map[string]*Node{
					"A": {
						Name: "A",
						Children: []*Node{
							{
								Name: "B",
								Parents: []*Node{
									{
										Name: "B",
									},
								},
							},
						},
					},
					"B": {
						Name: "B",
						Parents: []*Node{
							{
								Name: "A",
							},
						},
					},
				},
			},
			parent: &Node{
				Name: "B",
				Parents: []*Node{
					{
						Name: "A",
						Children: []*Node{
							{Name: "B"},
						},
					},
				},
			},
			node: &Node{
				Name: "A",
				Children: []*Node{
					{
						Name: "B",
						Children: []*Node{
							{Name: "A"},
						},
					},
				},
			},
			err: errors.New(errContext, "Cycle detected adding relationship from 'B' to 'A'"),
			res: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.graph.AddRelationship(test.parent, test.node)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, len(test.graph.Root), len(test.res.Root), "Root size not equal")
				assert.Equal(t, len(test.graph.NodesIndex), len(test.res.NodesIndex), "Root size not equal")
			}
		})
	}
}

func TestHasCycles(t *testing.T) {

	tests := []struct {
		desc              string
		graph             *Graph
		prepareAssertFunc func(*Graph)
		res               bool
	}{
		{
			desc:  "Testing cyclic graph",
			graph: &Graph{},
			prepareAssertFunc: func(graph *Graph) {
				nodeA := &Node{Name: "A"}
				nodeB := &Node{Name: "B"}
				nodeC := &Node{Name: "C"}
				nodeD := &Node{Name: "D"}

				graph.AddNode(nodeA)
				graph.AddNode(nodeB)
				graph.AddNode(nodeC)
				graph.AddNode(nodeD)

				graph.AddRelationship(nodeA, nodeB)
				graph.AddRelationship(nodeB, nodeC)
				graph.AddRelationship(nodeC, nodeD)

				nodeA.AddParent(nodeC)
			},
			res: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.graph)
			}

			res := test.graph.HasCycles()
			assert.Equal(t, test.res, res)
		})
	}

}

func TestExist(t *testing.T) {
	tests := []struct {
		desc  string
		graph *Graph
		node  *Node
		res   bool
	}{
		{
			desc: "Search an existing node",
			graph: &Graph{
				NodesIndex: map[string]*Node{
					"node": {
						Name: "node",
					},
				},
			},
			node: &Node{
				Name: "node",
			},
			res: true,
		},
		{
			desc: "Search an unexisting node",
			graph: &Graph{
				NodesIndex: map[string]*Node{
					"node": {
						Name: "node",
					},
				},
			},
			node: &Node{
				Name: "node2",
			},
			res: false,
		},
		{
			desc:  "Search on a nil graph",
			graph: nil,
			node: &Node{
				Name: "node",
			},
			res: false,
		},
		{
			desc: "Search on a nil nodesindex",
			graph: &Graph{
				NodesIndex: nil,
			},
			node: &Node{
				Name: "node",
			},
			res: false,
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			exists := test.graph.Exist(test.node)
			assert.Equal(t, exists, test.res, "Unexpected return of existence")
		})
	}
}

func TestGetNode(t *testing.T) {
	tests := []struct {
		desc  string
		graph *Graph
		node  string
		res   *Node
		err   error
	}{
		{
			desc: "Get an existing node",
			graph: &Graph{
				NodesIndex: map[string]*Node{
					"node": {
						Name: "node",
					},
				},
			},
			node: "node",
			res: &Node{
				Name: "node",
			},
			err: nil,
		},
		{
			desc: "Get an unexisting node",
			graph: &Graph{
				NodesIndex: map[string]*Node{
					"node": {
						Name: "node",
					},
				},
			},
			node: "node2",
			res:  nil,
			err:  errors.New("(graph::GetNode)", "Node 'node2' does not exists on the graph"),
		},
		{
			desc:  "Search on a nil graph",
			graph: nil,
			node:  "node",
			res:   nil,
			err:   errors.New("(graph::GetNode)", "Graph is nil"),
		},
		{
			desc: "Search on a nil nodesindex",
			graph: &Graph{
				NodesIndex: nil,
			},
			node: "node",
			res:  nil,
			err:  errors.New("(graph::GetNode)", "NodesIndex is nil"),
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			node, err := test.graph.GetNode(test.node)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, node, "Unexpected node")
			}
		})
	}
}
