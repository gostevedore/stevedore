package images

import (
	"fmt"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	graph "github.com/apenella/go-data-structures/extendedTree"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// TreeOutputOptionsFunc is a function used to configure the service
type TreeOutputOptionsFunc func(*TreeOutput)

type TreeOutput struct {
	writer        io.Writer
	graph         Grapher
	referenceName repository.ImageReferenceNamer
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

func WithReferenceName(ref repository.ImageReferenceNamer) TreeOutputOptionsFunc {
	return func(o *TreeOutput) {
		o.referenceName = ref
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

	//	if (parent != nil && i.Parent != nil) || (parent == nil && i.Parent == nil) {
	nodeName, err := o.nodeName(i)
	if err != nil {
		return errors.New(errContext, "", err)
	}

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
	//	}

	return nil
}

func (o *TreeOutput) nodeName(i *image.Image) (string, error) {
	errContext := "(output::images::tree::nodeName)"

	if o.referenceName == nil {
		return "", errors.New(errContext, "Images plain text output requires a reference name")
	}

	ref, err := o.referenceName.GenerateName(i)
	if err != nil {
		// return "", errors.New(errContext, "", err)
		// instead of returned the error fmt is used as a fallback
		ref = fmt.Sprintf("%s:%s", i.Name, i.Version)
	}

	return ref, nil
}
