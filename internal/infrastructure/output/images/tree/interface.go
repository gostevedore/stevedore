package images

import (
	"io"

	graph "github.com/apenella/go-data-structures/extendedTree"
)

type Grapher interface {
	AddNode(node *graph.Node) error
	Exist(node *graph.Node) bool
	AddRelationship(parent, child *graph.Node) error
	DrawGraph(w io.Writer)
}
