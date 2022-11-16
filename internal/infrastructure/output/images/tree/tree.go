package images

import (
	"fmt"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	graph "github.com/apenella/go-data-structures/extendedTree"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// TreeOutputOptionsFunc is a function used to configure the service
type TreeOutputOptionsFunc func(*TreeOutput)

type TreeOutput struct {
	writer io.Writer
	graph  Grapher
}

func NewTreeOutput(options ...TreeOutputOptionsFunc) *TreeOutput {
	output := &TreeOutput{}
	output.Options(options...)
	return output
}

func WithWriter(w io.Writer) TreeOutputOptionsFunc {
	return func(o *TreeOutput) {
		o.writer = w
	}
}

func WithGraph(g Grapher) TreeOutputOptionsFunc {
	return func(o *TreeOutput) {
		o.graph = g
	}
}

// Options configure the service
func (o *TreeOutput) Options(opts ...TreeOutputOptionsFunc) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *TreeOutput) Output(list []*image.Image) error {
	var err error
	errContext := "(output::images::tree::Output)"

	if o.graph == nil {
		return errors.New(errContext, " Tree output requireds that graph must be initialized")
	}

	for _, i := range list {
		err = o.addNodeToGraph(i)
		if err != nil {
			return errors.New(errContext, "", err)
		}
	}

	o.graph.DrawGraph(o.writer)

	return nil
}

func (o *TreeOutput) addNodeToGraph(i *image.Image) error {
	var err error

	errContext := "(output::images::tree::addNodeToGraph)"

	if o.graph == nil {
		return errors.New(errContext, " Tree output requireds that graph must be initialized")
	}

	err = o.addNodeToGraphRec(i, nil)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil

}

func (o *TreeOutput) addNodeToGraphRec(i *image.Image, parent *graph.Node) error {
	var node *graph.Node
	errContext := "(output::images::tree::addNodeToGraphRec)"

	if o.graph == nil {
		return errors.New(errContext, " Tree output requireds that graph must be initialized")
	}

	nodeName := nodeName(i)

	node = &graph.Node{
		Name: nodeName,
		Item: i,
	}

	if !o.graph.Exist(node) {
		err := o.graph.AddNode(node)
		if err != nil {
			return errors.New(errContext, "", err)
		}
	}

	if parent != nil {
		err := o.graph.AddRelationship(parent, node)
		if err != nil {
			return errors.New(errContext, "", err)
		}
	}

	if len(i.Children) > 0 {
		for _, child := range i.Children {
			err := o.addNodeToGraphRec(child, node)
			if err != nil {
				return errors.New(errContext, "", err)
			}
		}
	}

	return nil
}

func nodeName(i *image.Image) string {
	nodeName := fmt.Sprintf("%s:%s", i.Name, i.Version)

	if i.RegistryNamespace != "" {
		nodeName = fmt.Sprintf("%s/%s", i.RegistryNamespace, nodeName)
	}

	if i.RegistryHost != "" {
		nodeName = fmt.Sprintf("%s/%s", i.RegistryHost, nodeName)
	}

	return nodeName
}
