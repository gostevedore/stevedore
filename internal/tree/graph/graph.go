package graph

import (
	"errors"
	"fmt"
	"stevedore/internal/image"
	"stevedore/internal/tree"
	"strings"

	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	gdstree "github.com/apenella/go-data-structures/tree"
)

type Graph struct {
	Graph *gdstree.Graph
	Index *tree.ImageIndex
}

func GenerateGraph(template *GraphTemplate) (*Graph, error) {

	graph := &gdstree.Graph{}
	index := &tree.ImageIndex{}

	for _, root := range template.Graph.Root {
		err := generateGraphRec(graph, index, nil, root)
		if err != nil {
			return nil, errors.New("(graph::GenerateGraph) Error renderizing images graph. " + err.Error())
		}
	}

	g := &Graph{
		Graph: graph,
		Index: index,
	}

	return g, nil
}

func generateGraphRec(graph *gdstree.Graph, index *tree.ImageIndex, parent *gdstree.Node, node *gdsexttree.Node) error {

	var renderParent *image.Image
	if parent == nil {
		renderParent = &image.Image{}
	} else {
		renderParent = parent.Item.(*image.Image)
	}

	// copy image before to be processed
	originalImageNode := node.Item.(*image.Image)
	imageNode, err := originalImageNode.Copy()
	if err != nil {
		return errors.New("(tree::renderizeGraphRec) Error coping image '" + originalImageNode.Name + "' -> " + err.Error())
	}

	imageDetail := strings.Split(node.Name, ImageNodeNameSeparator)
	if len(imageDetail) != 2 {
		return errors.New("(tree::renderizeGraphRec) Node name '" + imageNode.Name + "' not valid")
	}
	imageName := imageDetail[0]
	imageVersion := imageDetail[1]

	renderImageData := &tree.ImageRender{
		Name:    imageName,
		Version: imageVersion,
		Parent:  renderParent,
		Image:   imageNode,
	}

	err = tree.RenderizeImage(renderImageData)
	if err != nil {
		return errors.New("(tree::renderizeGraphRec) Error renderinzing image '" + imageName + "'.\n " + err.Error())
	}

	if len(renderParent.PersistentVars) > 0 {
		for keyVar, keyValue := range renderParent.PersistentVars {
			// set all persistent vars defined on parent node an overwrite any matching node persistent var
			imageNode.PersistentVars[keyVar] = keyValue
		}
	}

	// generate node name for imagesGraph
	nodeName := tree.GenerateNodeName(imageNode)
	if parent != nil {
		nodeName = nodeName + "@" + renderParent.Name + ":" + renderParent.Version
	}

	newImageNode := &gdstree.Node{
		Name: nodeName,
		Item: imageNode,
	}
	err = graph.AddNode(newImageNode)
	if err != nil {
		fmt.Println(err.Error())
	}
	if parent != nil {
		graph.AddRelationship(parent, newImageNode)
	}

	// Include node to index.
	// Two entries are included:
	//  1 - from image tree definition
	index.AddNode(imageName, imageVersion, newImageNode)
	//  2 - from image rendered values
	if imageNode.Name != imageName || imageNode.Version != imageVersion {
		index.AddAlternativeIndexImage(imageNode.Name, imageNode.Version, newImageNode)
	}
	// 3 - include to wildcard index nodes
	if imageVersion == WildCardVersionSymbol {
		// imageVersion is used on find nodes by name
		index.AddWildcardIndexImage(imageName, imageVersion)
		// imageNode.Version is used when find nodes by name and version
		index.AddWildcardIndexImage(imageName, imageNode.Version)
	}

	for _, child := range node.Children {
		err := generateGraphRec(graph, index, newImageNode, child)
		if err != nil {
			return errors.New("(tree::renderizeGraphRec) Error renderizing image graph -> " + err.Error())
		}
	}

	return nil
}
