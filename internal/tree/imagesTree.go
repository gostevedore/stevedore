package tree

import (
	"fmt"
	"strings"

	common "github.com/apenella/go-common-utils/data"
	errors "github.com/apenella/go-common-utils/error"
	gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	gdstree "github.com/apenella/go-data-structures/tree"
	"github.com/gostevedore/stevedore/internal/image"
)

const (
	ImageNodeNameSeparator = ":"
)

// ImagesTree
//
// Image structure
// image_tree:
// 	image_name:
//		image_tag1:
//			<Image>
//		image_tag2:
//			<Image>
type ImagesTree struct {
	Images map[string]map[string]*image.Image `yaml:"images_tree"`
}

// LoadImagesTree method generate and return an ImagesTree struct from a file
func LoadImagesTree(file string) (*ImagesTree, error) {

	imagesTree := &ImagesTree{}
	err := common.LoadYAMLFile(file, imagesTree)
	if err != nil {
		return nil, errors.New("(tree::LoadImagesTree)", "Error loading images tree configuration", err)
	}

	if imagesTree.Images == nil {
		return nil, errors.New("(tree::LoadImagesTree)", "Image tree is not defined properly on "+file)
	}

	return imagesTree, nil
}

// GenerateGraph method returns a graph having the images and its relationships and a index the improve its searches
func (t *ImagesTree) GenerateGraph() (*gdstree.Graph, *ImageIndex, error) {

	imagesTemplateGraph := &gdsexttree.Graph{}

	for imageName, imageVersions := range t.Images {
		for imageVersion, imageDef := range imageVersions {
			// root nodes has no parent then its argument is nil
			err := t.generateTemplateGraph(imageName, imageVersion, imageDef, imagesTemplateGraph, nil)
			if err != nil {
				return nil, nil, errors.New("(tree::GenerateGraph)", "Error generating images graph", err)
			}
		}
	}

	imagesGraph, index, err := RenderizeGraph(imagesTemplateGraph)
	if err != nil {
		return nil, nil, errors.New("(tree::GenerateGraph)", "Error renderizing images tree", err)
	}

	return imagesGraph, index, nil
}

// generateTemplateGraph method create the template graph which must be renderized to generate images graph
func (t *ImagesTree) generateTemplateGraph(imageName string, imageVersion string, nodeImage *image.Image, imagesGraph *gdsexttree.Graph, parent *gdsexttree.Node) error {

	if nodeImage == nil {
		return errors.New("(tree::generateGraphRec)", "Node Image is null")
	}

	// enrich image date with a Name and a Version
	if nodeImage.Name == "" {
		nodeImage.Name = imageName
	}
	if nodeImage.Version == "" {
		nodeImage.Version = imageVersion
	}

	// validate compatibility
	nodeImage.CheckCompatibility()

	node := &gdsexttree.Node{
		Name: imageName + ImageNodeNameSeparator + imageVersion,
		Item: nodeImage,
	}

	if imagesGraph.Exist(node) {
		node, _ = imagesGraph.GetNode(node.Name)
	} else {
		err := imagesGraph.AddNode(node)
		if err != nil {
			return errors.New("(tree::generateTemplateGraph)", fmt.Sprintf("Node '%s' could not be added to tree", node.Name), err)
		}
	}

	if parent != nil {
		//if parent != nil && !node.HasParent(parent) {

		if !node.HasParent(parent) {
			err := imagesGraph.AddRelationship(parent, node)
			if err != nil {
				return errors.New("(tree::generateTemplateGraph)", fmt.Sprintf("Relationship from '%s' to '%s' could not be created", parent.Name, node.Name), err)
			}
		}
	}

	if imagesGraph.HasCycles() {
		return errors.New("(tree::generateTemplateGraph)", "Cycle detected")
	}

	for childName, childVersions := range nodeImage.Children {
		for _, childVersion := range childVersions {
			childImage, exist := t.Images[childName][childVersion]

			if exist {
				err := t.generateTemplateGraph(childName, childVersion, childImage, imagesGraph, node)
				if err != nil {
					return errors.New("(tree::generateTemplateGraph)", fmt.Sprintf("Error generating template tree from '%s' to '%s'", childName, node.Name), err)
				}
			}
		}
	}

	return nil
}

// GenerateNodeName
func GenerateNodeName(i *image.Image) string {
	return i.Name + ImageNodeNameSeparator + i.Version
}

// RenderizeGraph method do the template graph renderization to generate an images graph
func RenderizeGraph(g *gdsexttree.Graph) (*gdstree.Graph, *ImageIndex, error) {
	imagesGraph := &gdstree.Graph{}
	index := &ImageIndex{}

	for _, root := range g.Root {
		err := renderizeGraphRec(imagesGraph, index, nil, root)
		if err != nil {
			return nil, nil, errors.New("(tree::RenderizeGraph)", "Error renderizing images graph", err)
		}
	}

	return imagesGraph, index, nil
}

// renderizeGraphRec method its the RenderizeGraph worker
func renderizeGraphRec(imagesGraph *gdstree.Graph, index *ImageIndex, parent *gdstree.Node, node *gdsexttree.Node) error {

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
		return errors.New("(tree::renderizeGraphRec)", "Error coping image '"+originalImageNode.Name+"'", err)
	}

	imageDetail := strings.Split(node.Name, ImageNodeNameSeparator)
	if len(imageDetail) != 2 {
		return errors.New("(tree::renderizeGraphRec)", "Node name '"+imageNode.Name+"' not valid")
	}
	imageName := imageDetail[0]
	imageVersion := imageDetail[1]

	renderImageData := &ImageRender{
		Name:    imageName,
		Version: imageVersion,
		Parent:  renderParent,
		Image:   imageNode,
	}

	err = RenderizeImage(renderImageData)
	if err != nil {
		return errors.New("(tree::renderizeGraphRec)", "Error renderinzing image '"+imageName+"'", err)
	}

	if len(renderParent.PersistentVars) > 0 {
		for keyVar, keyValue := range renderParent.PersistentVars {
			// set all persistent vars defined on parent node an overwrite any matching node persistent var
			imageNode.PersistentVars[keyVar] = keyValue
		}
	}

	// generate node name for imagesGraph
	nodeName := GenerateNodeName(imageNode)
	if parent != nil {
		nodeName = nodeName + "@" + renderParent.Name + ":" + renderParent.Version
	}

	newImageNode := &gdstree.Node{
		Name: nodeName,
		Item: imageNode,
	}
	err = imagesGraph.AddNode(newImageNode)
	if err != nil {
		fmt.Println(err.Error())
	}
	if parent != nil {
		imagesGraph.AddRelationship(parent, newImageNode)
	}

	// Include node to index.
	// Three entries are included:
	//  1 - from image tree definition
	index.AddNode(imageName, imageVersion, newImageNode)
	//  2 - from image rendered values
	if imageNode.Name != imageName || imageNode.Version != imageVersion {
		index.AddAlternativeIndexImage(imageNode.Name, imageNode.Version, newImageNode)
	}
	// 3 - include to wildcard index nodes
	if imageVersion == wildCardVersionSymbol {
		// imageVersion is used on find nodes by name
		index.AddWildcardIndexImage(imageName, imageVersion)
		// imageNode.Version is used when find nodes by name and version
		index.AddWildcardIndexImage(imageName, imageNode.Version)
	}

	for _, child := range node.Children {
		err := renderizeGraphRec(imagesGraph, index, newImageNode, child)
		if err != nil {
			return errors.New("(tree::renderizeGraphRec)", "Error renderizing image graph", err)
		}
	}

	return nil
}

func GetNodeImage(node *gdstree.Node) (*image.Image, error) {
	if node == nil {
		return nil, errors.New("(tree::GetNodeImage)", "Node is nil")
	}
	if node.Item == nil {
		return nil, errors.New("(tree::GetNodeImage)", "Node item is nil")
	}
	i := node.Item.(*image.Image)

	return i, nil
}

// GenerateWilcardVersionNode generate a new node based wildcard version definition
func (t *ImagesTree) GenerateWilcardVersionNode(node *gdstree.Node, version string) (*gdstree.Node, error) {

	var err error
	var exist bool
	var imageAux *image.Image
	var imageAuxWildcard *image.Image
	var imageWildcard *image.Image
	var nodeAuxChilds []*gdstree.Node

	if t == nil {
		return nil, errors.New("(tree::GenerateNodeWithWilcardVersion)", "Images tree is nil")
	}
	if node == nil {
		return nil, errors.New("(tree::GenerateNodeWithWilcardVersion)", "Node is nil")
	}

	imageAux, err = GetNodeImage(node)
	if err != nil {
		return nil, errors.New("(tree::GenerateNodeWithWilcardVersion)", "Error when achieve image from node '"+node.Name+"'")
	}
	nodeName := imageAux.Name

	imageAuxWildcard, exist = t.Images[nodeName][wildCardVersionSymbol]
	if !exist {
		return nil, errors.New("(tree::GenerateNodeWithWilcardVersion)", "Node '"+nodeName+"' does not exists or not has not got a wildcard version")
	}

	imageWildcard, err = imageAuxWildcard.Copy()
	if err != nil {
		return nil, errors.New("(tree::GenerateNodeWithWilcardVersion)", "Error coping image '"+node.Name+"'", err)
	}
	imageWildcard.Version = version

	nodeWildcardName := GenerateNodeName(imageWildcard)

	parent := &image.Image{}
	if node.Parent != nil && node.Parent.Item != nil {
		parent = node.Parent.Item.(*image.Image)
	}

	renderImageData := &ImageRender{
		Name:    nodeWildcardName,
		Version: version,
		Parent:  parent,
		Image:   imageWildcard,
	}

	err = RenderizeImage(renderImageData)
	if err != nil {
		return nil, errors.New("(tree::GenerateNodeWithWilcardVersion)", "Error renderinzing image '"+nodeName+"'", err)
	}

	for _, aux := range node.Children {
		nodeChildAux, _ := t.GenerateWilcardVersionNode(aux, version)
		if nodeChildAux != nil {
			nodeAuxChilds = append(nodeAuxChilds, nodeChildAux)
		}
	}

	nodeAux := &gdstree.Node{
		Name:     nodeWildcardName,
		Item:     imageWildcard,
		Children: nodeAuxChilds,
		Parent:   node.Parent,
	}

	return nodeAux, nil
}
