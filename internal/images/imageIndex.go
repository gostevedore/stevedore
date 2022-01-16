package tree

import (
	errors "github.com/apenella/go-common-utils/error"
	gdstree "github.com/apenella/go-data-structures/tree"
)

const (
	wildCardVersionSymbol = "*"
)

// ImageIndex contains indexes to images defined on the images grapth
type ImageIndex struct {
	NameIndex                   map[string][]string
	NameVersionIndex            map[string][]*gdstree.Node
	NameVersionAlternativeIndex map[string][]*gdstree.Node
	WildcardIndex               map[string]uint8
}

// AddNode method includes a node to index
func (ii *ImageIndex) AddNode(name, version string, node *gdstree.Node) error {
	if ii == nil {
		return errors.New("(tree::AddImage)", "Image index is null")
	}

	if ii.NameIndex == nil {
		ii.NameIndex = map[string][]string{}
	}

	if ii.NameVersionIndex == nil {
		ii.NameVersionIndex = map[string][]*gdstree.Node{}
	}
	imageNameVersion := name + ImageNodeNameSeparator + version

	imageNameVersionList, exist := ii.NameIndex[name]
	if !exist {
		ii.NameIndex[name] = []string{imageNameVersion}
	}

	nodeList, exist := ii.NameVersionIndex[imageNameVersion]
	if !exist {
		ii.NameVersionIndex[imageNameVersion] = []*gdstree.Node{node}
		ii.NameIndex[name] = append(imageNameVersionList, imageNameVersion)
	} else {
		ii.NameVersionIndex[imageNameVersion] = append(nodeList, node)
	}

	return nil
}

// AddAlternativeIndexImage method includes an node to alternative searches
func (ii *ImageIndex) AddAlternativeIndexImage(name, version string, node *gdstree.Node) error {
	if ii == nil {
		return errors.New("(tree::AddAlternativeIndexImage)", "Image index is null")
	}

	if ii.NameVersionAlternativeIndex == nil {
		ii.NameVersionAlternativeIndex = map[string][]*gdstree.Node{}
	}

	imageNameVersion := name + ImageNodeNameSeparator + version
	imageList, exist := ii.NameVersionAlternativeIndex[imageNameVersion]
	if !exist {
		ii.NameVersionAlternativeIndex[imageNameVersion] = []*gdstree.Node{node}
	} else {
		ii.NameVersionAlternativeIndex[imageNameVersion] = append(imageList, node)
	}

	return nil
}

// AddWildcardIndexImage add an image identified by a wildcard version
func (ii *ImageIndex) AddWildcardIndexImage(name, version string) error {

	if ii.WildcardIndex == nil {
		ii.WildcardIndex = map[string]uint8{}
	}

	imageNameWildcardVersion := name + ImageNodeNameSeparator + version

	_, exists := ii.WildcardIndex[imageNameWildcardVersion]
	if !exists {
		ii.WildcardIndex[imageNameWildcardVersion] = uint8(0)
	}

	return nil
}

// Find returns a node list with node matching to imageName and version. The list does not return wildcarded versions.
func (ii *ImageIndex) Find(imageName string, version string) ([]*gdstree.Node, error) {
	var err error

	nodesList := []*gdstree.Node{}

	imageNameWildcardVersion := imageName + ImageNodeNameSeparator + version
	_, isWildcardVersion := ii.WildcardIndex[imageNameWildcardVersion]

	// return an empty node list and a nil error when a wildcard version is found
	// it is required because imagesEngine walks through image tree and find wildcards --> 0.8.2 may solve it
	if isWildcardVersion {
		return nodesList, nil
	}

	if version == "" {
		nodesList, err = ii.findByName(imageName)
		if err != nil {
			return nil, errors.New("(tree::Find)", "Error when finding images by name '"+imageName+"'", err)
		}
	} else {

		nodesList, err = ii.findByNameAndVersion(imageName, version)
		if err != nil {
			nodesList, err = ii.findAlternativeByNameAndVersion(imageName, version)
			if err != nil {
				return nil, errors.New("(tree::Find)", "Error when finding image '"+imageName+":"+version+"'")
			}
		}
	}

	return nodesList, nil
}

// findByName method returns a node list with al nodes matching to node name. List does not contain wildcarded versions.
func (ii *ImageIndex) findByName(name string) ([]*gdstree.Node, error) {

	nodes := []*gdstree.Node{}

	imageNameversions, exists := ii.NameIndex[name]
	if !exists {
		return nil, errors.New("(tree::FindByName)", "Image name '"+name+"' does not exists")
	}

	for _, imageNameVersion := range imageNameversions {
		// append node when it has not have a wildcard version
		_, exists := ii.WildcardIndex[imageNameVersion]

		if !exists {
			nodeList, exists := ii.NameVersionIndex[imageNameVersion]
			if !exists {
				return nil, errors.New("(tree::FindByName)", "Image name-version '"+imageNameVersion+"' does not exists")
			}
			for _, image := range nodeList {
				nodes = append(nodes, image)
			}
		}
	}

	return nodes, nil
}

// findByNameAndVersion method returns a node list with nodes matching to version parameters. When version is registered as wildcard version it returns an empty list.
func (ii *ImageIndex) findByNameAndVersion(name, version string) ([]*gdstree.Node, error) {

	nodes := []*gdstree.Node{}

	imageNameVersion := name + ImageNodeNameSeparator + version
	// append node when it has not have a wildcard version
	_, exists := ii.WildcardIndex[imageNameVersion]
	if !exists {
		nodeList, exists := ii.NameVersionIndex[imageNameVersion]
		if !exists {
			return nil, errors.New("(tree::FindByNameAndVersion)", "Image name-version '"+imageNameVersion+"' does not exists")
		}
		for _, node := range nodeList {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

// findAlternativeByNameAndVersion return a node list with nodes matching to the alternative version parameters. When version is registered as wildcard version it returns an empty list.
func (ii *ImageIndex) findAlternativeByNameAndVersion(name, version string) ([]*gdstree.Node, error) {

	nodes := []*gdstree.Node{}

	imageNameVersion := name + ImageNodeNameSeparator + version

	// append node when it has not have a wildcard version
	_, exists := ii.WildcardIndex[imageNameVersion]
	if !exists {
		nodeList, exists := ii.NameVersionAlternativeIndex[imageNameVersion]
		if !exists {
			return nil, errors.New("(tree::FindAlternativeByNameAndVersion)", "Image name-version '"+imageNameVersion+"' does not exists")
		}
		for _, node := range nodeList {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

// FindWildcardVersion return a wildcarded image when it exists. Otherwise returns an error.
func (ii *ImageIndex) FindWildcardVersion(name string) (*gdstree.Node, error) {

	imageNameVersion := name + ImageNodeNameSeparator + wildCardVersionSymbol
	nodeList, exists := ii.NameVersionIndex[imageNameVersion]
	if !exists {
		return nil, errors.New("(tree::FindWildcardVersion)", "Image '"+name+"' does not have a wildcard version definition")
	}

	return nodeList[0], nil
}

// IsWildcardVersion return if an image is a wilcarded image
func (ii *ImageIndex) IsWildcardVersion(name, version string) bool {
	imageNameWildcardVersion := name + ImageNodeNameSeparator + version
	_, isWildcardVersion := ii.WildcardIndex[imageNameWildcardVersion]

	return isWildcardVersion
}
