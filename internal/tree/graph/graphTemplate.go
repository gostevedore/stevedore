package graph

import (
	"errors"

	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/tree"
)

type GraphTemplate struct {
	Graph *gdsexttree.Graph
}

// GenerateGraphTemplate generates a graph based on image tree definition. The graph returned is a templete which will be used to generate the images graph.
func GenerateGraphTemplate(tree *tree.ImagesTree) (*GraphTemplate, error) {

	if tree == nil {
		return nil, errors.New("(graph::GenerateGraphTemplate) Tree is null")
	}

	graph := &gdsexttree.Graph{}

	for imageName, imageVersions := range tree.Images {
		for imageVersion, imageDetail := range imageVersions {
			err := generateGraphTemplateRec(imageName, imageVersion, imageDetail, nil, graph, tree)
			if err != nil {
				return nil, errors.New("(graph::GenerateGraph) Error generating template graph. " + err.Error())
			}
		}
	}

	template := &GraphTemplate{
		Graph: graph,
	}

	return template, nil
}

// generateGraphTemplateRec method create the template graph which must be renderized to generate images graph
func generateGraphTemplateRec(imageName string, imageVersion string, image *image.Image, parent *gdsexttree.Node, graph *gdsexttree.Graph, tree *tree.ImagesTree) error {

	if image == nil {
		return errors.New("(graph::generateGraphTemplateRec) Node Image is null")
	}

	if graph == nil {
		return errors.New("(graph::generateGraphTemplateRec) Graph is null")
	}

	if tree == nil {
		return errors.New("(graph::generateGraphTemplateRec) Tree is null")
	}

	// enrich image date with a Name and a Version
	if image.Name == "" {
		image.Name = imageName
	}
	if image.Version == "" {
		image.Version = imageVersion
	}

	node := &gdsexttree.Node{
		Name: imageName + ImageNodeNameSeparator + imageVersion,
		Item: image,
	}

	if graph.Exist(node) {
		node, _ = graph.GetNode(node.Name)
	} else {
		graph.AddNode(node)
	}

	graph.AddRelationship(parent, node)

	for childName, childVersions := range image.Children {
		for _, childVersion := range childVersions {
			childImage, exist := tree.Images[childName][childVersion]

			if exist {
				// the current node is the next node to explore parent
				generateGraphTemplateRec(childName, childVersion, childImage, node, graph, tree)
			}
		}
	}

	return nil
}
