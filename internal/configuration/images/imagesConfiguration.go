package images

import (
	"fmt"
	"strings"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration/images/graph"
	"github.com/gostevedore/stevedore/internal/configuration/images/image"
	domainimage "github.com/gostevedore/stevedore/internal/images/image"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// const (
// 	ImageNodeNameSeparator = ":"
// )

// ImagesConfiguration
//
// Image structure
// image_tree:
// 	image_name:
//		image_tag1:
//			<Image>
//		image_tag2:
//			<Image>
type ImagesConfiguration struct {
	compatibility Compatibilitier
	graph         ImagesGraphTemplatesStorer
	fs            afero.Fs
	mutex         sync.RWMutex
	wg            sync.WaitGroup
	store         ImagesStorer

	// DEPRECATED_ImagesTree is replaced by Images
	DEPRECATED_ImagesTree map[string]map[string]*image.Image `yaml:"images_tree"`
	Images                map[string]map[string]*image.Image `yaml:"images"`
}

// NewImagesConfiguration method create a new ImagesConfiguration struct
func NewImagesConfiguration(fs afero.Fs, graph ImagesGraphTemplatesStorer, store ImagesStorer, compatibility Compatibilitier) *ImagesConfiguration {
	return &ImagesConfiguration{
		fs:            fs,
		compatibility: compatibility,
		graph:         graph,
		store:         store,

		DEPRECATED_ImagesTree: make(map[string]map[string]*image.Image),
		Images:                make(map[string]map[string]*image.Image),
	}
}

// CheckCompatibility method ensures that ImagesConfiguration is compatible with current version
func (t *ImagesConfiguration) CheckCompatibility() error {

	if t.DEPRECATED_ImagesTree != nil && len(t.DEPRECATED_ImagesTree) > 0 {
		t.compatibility.AddDeprecated("'images_tree' is deprecated and will be removed on v0.12.0, please use 'images' instead")
	}

	return nil
}

// LoadImagesToStore method loads images defined on configuration to images store
func (t *ImagesConfiguration) LoadImagesToStore(path string) error {

	var err error
	errContext := "(images::LoadImagesToStore)"
	var nodeDomainImage, copyDomainImage *domainimage.Image
	_ = copyDomainImage

	err = t.LoadImagesConfiguration(path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	// render image
	/*
		- get parent domain image
		- prepare render object with parent domain image, image name, image version and config image
		- render config image,
		- generate domain image from config image
		- add domain image as child of parent domain image
		- add domain image to store
		- initiate render of child config images
	*/

	storedNodes := map[string]struct{}{}
	for node := range t.graph.Iterate() {

		// skip node if already stored
		_, stored := storedNodes[node.Name()]
		if stored {
			continue
		}

		name, version, err := graph.ParseNodeName(node)
		if err != nil {
			return errors.New(errContext, err.Error())
		}

		nodeImage := node.Item().(*image.Image)
		nodeDomainImage, err = nodeImage.CreateDomainImage()
		if err != nil {
			return errors.New(errContext, err.Error())
		}

		if node.Parents() == nil || len(node.Parents()) <= 0 {
			err = t.store.AddImage(name, version, nodeDomainImage)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
		} else {
			for _, parent := range node.Parents() {

				parentName, parentVersion, err := graph.ParseNodeName(parent.(graph.GraphNoder))
				parentImage, err := t.store.Find(parentName, parentVersion)
				if err != nil {
					return errors.New(errContext, err.Error())
				}

				copyDomainImage, err = nodeDomainImage.Copy()
				if err != nil {
					return errors.New(errContext, err.Error())
				}

				copyDomainImage.Options(domainimage.WithParent(parentImage))
				parentImage.AddChild(copyDomainImage)
				err = t.store.AddImage(name, version, copyDomainImage)
				if err != nil {
					return errors.New(errContext, err.Error())
				}

				storedNodes[node.Name()] = struct{}{}
			}
		}
	}

	return nil
}

// LoadImagesConfiguration method generate and return an ImagesConfiguration struct from a file
func (t *ImagesConfiguration) LoadImagesConfiguration(path string) error {

	var err error
	var isDir bool

	errContext := "(images::LoadImagesConfiguration)"

	isDir, err = afero.IsDir(t.fs, path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	if isDir {
		return t.LoadImagesConfigurationFromDir(path)
	} else {
		return t.LoadImagesConfigurationFromFile(path)
	}
}

// LoadImagesConfigurationFromDir loads images tree from all files on directory
func (t *ImagesConfiguration) LoadImagesConfigurationFromDir(dir string) error {
	var err error
	errFuncs := []func() error{}
	errContext := "(images::LoadImagesConfigurationFromDir)"

	yamlFiles, err := afero.Glob(t.fs, dir+"/*.yaml")
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	ymlFiles, err := afero.Glob(t.fs, dir+"/*.yml")
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	files := append(yamlFiles, ymlFiles...)

	loadImagesConfigurationFromFile := func(path string) func() error {
		var err error

		c := make(chan struct{}, 1)
		go func() {
			defer close(c)
			err = t.LoadImagesConfigurationFromFile(path)
			t.wg.Done()
		}()

		return func() error {
			<-c
			return err
		}
	}

	for _, file := range files {
		t.wg.Add(1)
		f := loadImagesConfigurationFromFile(file)
		errFuncs = append(errFuncs, f)
	}

	t.wg.Wait()

	errMsg := ""
	for _, f := range errFuncs {
		err = f()
		if err != nil {
			errMsg = fmt.Sprintf("%s%s\n", errMsg, err.Error())
		}
	}
	if errMsg != "" {
		return errors.New(errContext, errMsg)
	}

	return nil
}

// LoadImagesConfigurationFromFile loads images tree from file
func (t *ImagesConfiguration) LoadImagesConfigurationFromFile(path string) error {

	var err error
	var fileData []byte

	errContext := "(images::LoadImagesConfigurationFromFile)"

	if t == nil {
		return errors.New(errContext, "Builders is nil")
	}

	imageTreeAux := NewImagesConfiguration(t.fs, t.graph, t.store, t.compatibility)

	fileData, err = afero.ReadFile(t.fs, path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	err = yaml.Unmarshal(fileData, imageTreeAux)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error loading images tree from file '%s'\nfound:\n%s", path, string(fileData)), err)
	}

	err = imageTreeAux.CheckCompatibility()
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	for name, images := range imageTreeAux.Images {
		if !isAValidName(name) {
			return errors.New(errContext, fmt.Sprintf("Found an invalid image name '%s' defined in file '%s'", name, path))
		}

		for version, image := range images {
			if !isAValidVersion(version) {
				return errors.New(errContext, fmt.Sprintf("Found an invalid image version '%s' defined in file '%s'", version, path))
			}

			err = t.graph.AddImage(name, version, image)
			// err = t.AddImage(name, version, image)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
		}
	}

	// TO BE REMOVE on v0.12: is kept just for compatibility concerns
	for name, images := range imageTreeAux.DEPRECATED_ImagesTree {
		if !isAValidName(name) {
			return errors.New(errContext, fmt.Sprintf("Found an invalid image name '%s' defined in file '%s'", name, path))
		}

		for version, image := range images {
			if !isAValidVersion(version) {
				return errors.New(errContext, fmt.Sprintf("Found an invalid image version '%s' defined in file '%s'", version, path))
			}

			err := t.graph.AddImage(name, version, image)
			if err != nil {
				return errors.New(errContext, err.Error())
			}
		}
	}

	return nil
}

// isValidName method checks if a string is a valid image name
func isAValidName(name string) bool {
	if strings.IndexRune(name, ':') != -1 {
		return false
	}
	return true
}

// isValidVersion method checks if a string is a valid image version
func isAValidVersion(name string) bool {
	if strings.IndexRune(name, ':') != -1 {
		return false
	}
	return true
}

// AddImage method add an image to the tree
// func (t *ImagesConfiguration) AddImage(name, version string, i *image.Image) error {

// 	errContext := "(images::AddImage)"

// 	if i == nil {
// 		return errors.New(errContext, "Image to add is null")
// 	}

// 	if t.Images == nil {
// 		t.Images = make(map[string]map[string]*image.Image)
// 	}

// 	t.mutex.Lock()
// 	defer t.mutex.Unlock()

// 	_, exist := t.Images[name]
// 	if !exist {
// 		t.Images[name] = make(map[string]*image.Image)
// 	}

// 	_, exist = t.Images[name][version]
// 	if exist {
// 		return errors.New(errContext, fmt.Sprintf("Image '%s:%s' already defined on image tree", name, version))
// 	}

// 	t.Images[name][version] = i

// 	return nil
// }

// // GenerateGraph method returns a graph having the images and its relationships and a index the improve its searches
// func (t *ImagesConfiguration) GenerateGraph() (*gdstree.Graph, *ImageIndex, error) {

// 	imagesTemplateGraph := &gdsexttree.Graph{}

// 	for imageName, imageVersions := range t.Images {
// 		for imageVersion, imageDef := range imageVersions {
// 			// root nodes has no parent then its argument is nil
// 			err := t.generateTemplateGraph(imageName, imageVersion, imageDef, imagesTemplateGraph, nil)
// 			if err != nil {
// 				return nil, nil, errors.New("(images::GenerateGraph)", "Error generating images graph", err)
// 			}
// 		}
// 	}

// 	imagesGraph, index, err := RenderizeGraph(imagesTemplateGraph)
// 	if err != nil {
// 		return nil, nil, errors.New("(images::GenerateGraph)", "Error renderizing images tree", err)
// 	}

// 	return imagesGraph, index, nil
// }

// // generateTemplateGraph method create the template graph which must be renderized to generate images graph
// func (t *ImagesConfiguration) generateTemplateGraph(imageName string, imageVersion string, nodeImage *Image, imagesGraph *gdsexttree.Graph, parent *gdsexttree.Node) error {

// 	if nodeImage == nil {
// 		return errors.New("(images::generateGraphRec)", "Node Image is null")
// 	}

// 	// enrich image date with a Name and a Version
// 	if nodeImage.Name == "" {
// 		nodeImage.Name = imageName
// 	}
// 	if nodeImage.Version == "" {
// 		nodeImage.Version = imageVersion
// 	}

// 	// validate compatibility
// 	nodeImage.CheckCompatibility()

// 	node := &gdsexttree.Node{
// 		Name: imageName + ImageNodeNameSeparator + imageVersion,
// 		Item: nodeImage,
// 	}

// 	if imagesGraph.Exist(node) {
// 		node, _ = imagesGraph.GetNode(node.Name)
// 	} else {
// 		err := imagesGraph.AddNode(node)
// 		if err != nil {
// 			return errors.New("(images::generateTemplateGraph)", fmt.Sprintf("Node '%s' could not be added to tree", node.Name), err)
// 		}
// 	}

// 	if parent != nil {
// 		//if parent != nil && !node.HasParent(parent) {

// 		if !node.HasParent(parent) {
// 			err := imagesGraph.AddRelationship(parent, node)
// 			if err != nil {
// 				return errors.New("(images::generateTemplateGraph)", fmt.Sprintf("Relationship from '%s' to '%s' could not be created", parent.Name, node.Name), err)
// 			}
// 		}
// 	}

// 	if imagesGraph.HasCycles() {
// 		return errors.New("(images::generateTemplateGraph)", "Cycle detected")
// 	}

// 	for childName, childVersions := range nodeImage.Children {
// 		for _, childVersion := range childVersions {
// 			childImage, exist := t.Images[childName][childVersion]

// 			if exist {
// 				err := t.generateTemplateGraph(childName, childVersion, childImage, imagesGraph, node)
// 				if err != nil {
// 					return errors.New("(images::generateTemplateGraph)", fmt.Sprintf("Error generating template tree from '%s' to '%s'", childName, node.Name), err)
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

// // GenerateNodeName
// func GenerateNodeName(i *Image) string {
// 	return i.Name + ImageNodeNameSeparator + i.Version
// }

// // RenderizeGraph method do the template graph renderization to generate an images graph
// func RenderizeGraph(g *gdsexttree.Graph) (*gdstree.Graph, *ImageIndex, error) {
// 	imagesGraph := &gdstree.Graph{}
// 	index := &ImageIndex{}

// 	for _, root := range g.Root {
// 		err := renderizeGraphRec(imagesGraph, index, nil, root)
// 		if err != nil {
// 			return nil, nil, errors.New("(images::RenderizeGraph)", "Error renderizing images graph", err)
// 		}
// 	}

// 	return imagesGraph, index, nil
// }

// // renderizeGraphRec method its the RenderizeGraph worker
// func renderizeGraphRec(imagesGraph *gdstree.Graph, index *ImageIndex, parent *gdstree.Node, node *gdsexttree.Node) error {

// 	var renderParent *Image
// 	if parent == nil {
// 		renderParent = &image.Image{}
// 	} else {
// 		renderParent = parent.Item.(*Image)
// 	}

// 	// copy image before to be processed
// 	originalImageNode := node.Item.(*Image)
// 	imageNode, err := originalImageNode.Copy()
// 	if err != nil {
// 		return errors.New("(images::renderizeGraphRec)", "Error coping image '"+originalImageNode.Name+"'", err)
// 	}

// 	imageDetail := strings.Split(node.Name, ImageNodeNameSeparator)
// 	if len(imageDetail) != 2 {
// 		return errors.New("(images::renderizeGraphRec)", "Node name '"+imageNode.Name+"' not valid")
// 	}
// 	imageName := imageDetail[0]
// 	imageVersion := imageDetail[1]

// 	renderImageData := &ImageRender{
// 		Name:    imageName,
// 		Version: imageVersion,
// 		Parent:  renderParent,
// 		Image:   imageNode,
// 	}

// 	err = RenderizeImage(renderImageData)
// 	if err != nil {
// 		return errors.New("(images::renderizeGraphRec)", "Error renderinzing image '"+imageName+"'", err)
// 	}

// 	if len(renderParent.PersistentVars) > 0 {
// 		for keyVar, keyValue := range renderParent.PersistentVars {
// 			// set all persistent vars defined on parent node an overwrite any matching node persistent var
// 			imageNode.PersistentVars[keyVar] = keyValue
// 		}
// 	}

// 	// generate node name for imagesGraph
// 	nodeName := GenerateNodeName(imageNode)
// 	if parent != nil {
// 		nodeName = nodeName + "@" + renderParent.Name + ":" + renderParent.Version
// 	}

// 	newImageNode := &gdstree.Node{
// 		Name: nodeName,
// 		Item: imageNode,
// 	}
// 	err = imagesGraph.AddNode(newImageNode)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	if parent != nil {
// 		imagesGraph.AddRelationship(parent, newImageNode)
// 	}

// 	// Include node to index.
// 	// Three entries are included:
// 	//  1 - from image tree definition
// 	index.AddNode(imageName, imageVersion, newImageNode)
// 	//  2 - from image rendered values
// 	if imageNode.Name != imageName || imageNode.Version != imageVersion {
// 		index.AddAlternativeIndexImage(imageNode.Name, imageNode.Version, newImageNode)
// 	}
// 	// 3 - include to wildcard index nodes
// 	if imageVersion == wildCardVersionSymbol {
// 		// imageVersion is used on find nodes by name
// 		index.AddWildcardIndexImage(imageName, imageVersion)
// 		// imageNode.Version is used when find nodes by name and version
// 		index.AddWildcardIndexImage(imageName, imageNode.Version)
// 	}

// 	for _, child := range node.Children {
// 		err := renderizeGraphRec(imagesGraph, index, newImageNode, child)
// 		if err != nil {
// 			return errors.New("(images::renderizeGraphRec)", "Error renderizing image graph", err)
// 		}
// 	}

// 	return nil
// }

// func GetNodeImage(node *gdstree.Node) (*Image, error) {
// 	if node == nil {
// 		return nil, errors.New("(images::GetNodeImage)", "Node is nil")
// 	}
// 	if node.Item == nil {
// 		return nil, errors.New("(images::GetNodeImage)", "Node item is nil")
// 	}
// 	i := node.Item.(*Image)

// 	return i, nil
// }

// // GenerateWilcardVersionNode generate a new node based wildcard version definition
// func (t *ImagesConfiguration) GenerateWilcardVersionNode(node *gdstree.Node, version string) (*gdstree.Node, error) {

// 	var err error
// 	var exist bool
// 	var imageAux *Image
// 	var imageAuxWildcard *Image
// 	var imageWildcard *Image
// 	var nodeAuxChilds []*gdstree.Node

// 	if t == nil {
// 		return nil, errors.New("(images::GenerateNodeWithWilcardVersion)", "Images tree is nil")
// 	}
// 	if node == nil {
// 		return nil, errors.New("(images::GenerateNodeWithWilcardVersion)", "Node is nil")
// 	}

// 	imageAux, err = GetNodeImage(node)
// 	if err != nil {
// 		return nil, errors.New("(images::GenerateNodeWithWilcardVersion)", "Error when achieve image from node '"+node.Name+"'")
// 	}
// 	nodeName := imageAux.Name

// 	imageAuxWildcard, exist = t.Images[nodeName][wildCardVersionSymbol]
// 	if !exist {
// 		return nil, errors.New("(images::GenerateNodeWithWilcardVersion)", "Node '"+nodeName+"' does not exists or not has not got a wildcard version")
// 	}

// 	imageWildcard, err = imageAuxWildcard.Copy()
// 	if err != nil {
// 		return nil, errors.New("(images::GenerateNodeWithWilcardVersion)", "Error coping image '"+node.Name+"'", err)
// 	}
// 	imageWildcard.Version = version

// 	nodeWildcardName := GenerateNodeName(imageWildcard)

// 	parent := &image.Image{}
// 	if node.Parent != nil && node.Parent.Item != nil {
// 		parent = node.Parent.Item.(*Image)
// 	}

// 	renderImageData := &ImageRender{
// 		Name:    nodeWildcardName,
// 		Version: version,
// 		Parent:  parent,
// 		Image:   imageWildcard,
// 	}

// 	err = RenderizeImage(renderImageData)
// 	if err != nil {
// 		return nil, errors.New("(images::GenerateNodeWithWilcardVersion)", "Error renderinzing image '"+nodeName+"'", err)
// 	}

// 	for _, aux := range node.Children {
// 		nodeChildAux, _ := t.GenerateWilcardVersionNode(aux, version)
// 		if nodeChildAux != nil {
// 			nodeAuxChilds = append(nodeAuxChilds, nodeChildAux)
// 		}
// 	}

// 	nodeAux := &gdstree.Node{
// 		Name:     nodeWildcardName,
// 		Item:     imageWildcard,
// 		Children: nodeAuxChilds,
// 		Parent:   node.Parent,
// 	}

// 	return nodeAux, nil
// }
