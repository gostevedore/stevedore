package engine

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	gdstree "github.com/apenella/go-data-structures/tree"
	"github.com/gostevedore/stevedore/internal/build"
	factory "github.com/gostevedore/stevedore/internal/driver"
	defaultbuilder "github.com/gostevedore/stevedore/internal/driver/default"
	"github.com/gostevedore/stevedore/internal/image"

	//images "github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/promote"
	"github.com/gostevedore/stevedore/internal/schedule"
	"github.com/gostevedore/stevedore/internal/semver"
	"github.com/gostevedore/stevedore/internal/tree"
	"github.com/gostevedore/stevedore/internal/types"
	"github.com/gostevedore/stevedore/internal/ui/console"
	"gopkg.in/yaml.v2"
)

type ImagesEngine struct {
	Dispatch    types.Dispatcher
	ImagesTree  *tree.ImagesTree
	ImagesGraph *gdstree.Graph
	ImageIndex  *tree.ImageIndex
	Builders    *build.Builders
	context     context.Context
}

// NewImagesEngine
func NewImagesEngine(ctx context.Context, numWorkers int, imageTreePath string, buildersPath string) (*ImagesEngine, error) {

	var err error
	var imagesTree *tree.ImagesTree
	var imagesGraph *gdstree.Graph
	var imageIndex *tree.ImageIndex
	var builders *build.Builders
	var dispatch types.Dispatcher

	err = factory.InitFactories()
	if err != nil {
		return nil, errors.New("(engine::NewImagesEngine)", "Error initializing builder factories", err)

	}

	imagesTree, err = tree.LoadImagesTree(imageTreePath)
	if err != nil {
		msg := "Error loading image tree definition file"
		console.Print(msg)
		return nil, errors.New("(engine::NewImagesEngine)", msg, err)
	}

	imagesGraph, imageIndex, err = imagesTree.GenerateGraph()
	if err != nil {
		msg := "Error creating images tree"
		console.Print(msg)
		return nil, errors.New("(engine::NewImagesEngine)", msg, err)
	}

	builders, err = build.LoadBuilders(buildersPath)
	if err != nil {
		msg := "Error loading builders definition file"
		console.Print(msg)
		return nil, errors.New("(engine::NewImagesEngine)", msg, err)
	}

	dispatch, err = schedule.NewDispatch(ctx, numWorkers)
	if err != nil {
		msg := "Error creating dispatcher"
		console.Print(msg)
		return nil, errors.New("(engine::NewImagesEngine)", msg, err)
	}

	err = dispatch.Start()
	if err != nil {
		msg := "Error starting dispatcher"
		console.Print(msg)
		return nil, errors.New("(engine::NewImagesEngine)", msg, err)
	}

	engine := &ImagesEngine{
		Dispatch:    dispatch,
		ImagesTree:  imagesTree,
		ImagesGraph: imagesGraph,
		ImageIndex:  imageIndex,
		Builders:    builders,
		context:     ctx,
	}

	return engine, nil
}

func (e *ImagesEngine) findNodes(imageName string, versions []string) ([]*gdstree.Node, error) {
	var err error
	var nodeList []*gdstree.Node
	var nodeListAux []*gdstree.Node

	if versions == nil || len(versions) == 0 {
		nodeList, err = e.ImageIndex.Find(imageName, "")
		if err != nil {
			return nil, errors.New("(ImagesEngine::findNodes)", fmt.Sprintf("No image '%s' found on images tree", imageName), err)
		}
		if len(nodeList) == 0 {
			return nil, errors.New("(ImagesEngine::findNodes)", fmt.Sprintf("No image '%s' found on images tree", imageName))
		}
	} else {
		for _, version := range versions {
			nodeListAux, err = e.ImageIndex.Find(imageName, version)
			if err != nil {
				nodeAux, err := e.ImageIndex.FindWildcardVersion(imageName)
				if err != nil {
					continue
				}
				nodeAux, err = e.ImagesTree.GenerateWilcardVersionNode(nodeAux, version)
				if err != nil {
					continue
				}
				nodeListAux = append(nodeListAux, nodeAux)
			}

			for _, node := range nodeListAux {
				nodeList = append(nodeList, node)
			}
			nodeListAux = nil
		}

		if len(nodeList) == 0 {
			// That conditiona must be review when engine does not manages indexes directly

			// when no nodes found, one version is defined and
			if (len(versions) == 1 && !e.ImageIndex.IsWildcardVersion(imageName, versions[0])) ||
				// when no nodes found and more than one version is defined
				len(versions) > 1 {
				return nil, errors.New("(ImagesEngine::findNodes) ", fmt.Sprintf("No matching images to be build found for '%s' (versions: %s)", imageName, versions))
			}
		}
	}

	return nodeList, nil
}

// Build find the images to be build based on the imageName and version and runs each build
func (e *ImagesEngine) Build(imageName string, versions []string, options *types.BuildOptions, depth int) error {

	if options == nil {
		return errors.New("(ImagesEngine::Build)", "Build options is nil")
	}

	nodesList, err := e.findNodes(imageName, versions)
	if err != nil {
		return errors.New("(ImagesEngine::Build)", fmt.Sprintf("Error finding image '%s' and versions '%v'", imageName, versions), err)
	}

	endChildBuild := make(chan bool)
	endChildBuildErr := make(chan error)
	// start a goroutine for each image to built
	for _, node := range nodesList {
		go func(node *gdstree.Node, options types.BuildOptions) {
			// the build is done by buildworker method
			err := e.buildWorker(node, &options, depth)
			if err != nil {
				endChildBuildErr <- errors.New("(ImagesEngine::Build::goroutine)", "Error building '"+node.Name+"'", err)
			} else {
				endChildBuild <- true
			}
		}(node, *options)
	}

	// wait for all images to be built
	var imageBuildErrs error
	for range nodesList {
		select {
		case <-endChildBuild:
		case imageBuildErr := <-endChildBuildErr:
			if imageBuildErrs == nil {
				imageBuildErrs = imageBuildErr
			} else {
				if imageBuildErrs == nil {
					imageBuildErrs = fmt.Errorf("")
				}
				imageBuildErrs = fmt.Errorf("\n%s\n%s", imageBuildErrs.Error(), imageBuildErr.Error())
			}
		}
	}

	if imageBuildErrs != nil {
		return errors.New("(ImagesEngine::Build)", "Image could not be built", imageBuildErrs)
	}

	return nil
}

func (e *ImagesEngine) buildWorker(node *gdstree.Node, options *types.BuildOptions, depth int) error {
	var err error
	var imageBuilder *build.Builder
	var factoryFunc factory.DriverFactory
	var builder types.Driverer
	var exists bool

	if node == nil {
		return errors.New("(ImagesEngine::buildWorker)", "Node is not defined")
	}

	if options == nil {
		return errors.New("(ImagesEngine::buildWorker)", "Builder options is not defined")
	}

	image, err := tree.GetNodeImage(node)
	if err != nil {
		return errors.New("(ImagesEngine::buildWorker)", "Error getting image from node", err)
	}
	// An originalOptions' copy is kept because it will be passed to children build on cascade mode.
	originalOptions := &types.BuildOptions{}
	*originalOptions = *options

	// Set the image name when it has not been defined
	if options.ImageName == "" {
		options.ImageName = image.Name
	}
	options.ImageVersion = image.Version

	if options.EnableSemanticVersionTags {
		sv, err := semver.NewSemVer(options.ImageVersion)
		if err != nil {
			console.Warn(errors.New("(ImagesEngine::buildWorker)", "Version does not match to a semver expression", err))
		} else {
			svtree, err := sv.VersionTree(options.SemanticVersionTagsTemplate)
			if err != nil {
				console.Warn(errors.New("(ImagesEngine::buildWorker)", "Error generating version tree", err))
			} else {
				for _, v := range svtree {
					options.Tags = append(options.Tags, v)
				}
			}
		}
	}

	// copy the original map to avoid overwrite the original values
	imagePersistentVars := make(map[string]interface{})
	for k, v := range originalOptions.PersistentVars {
		imagePersistentVars[k] = v
	}
	// add the image persistent vars to the options build struct
	// a definition comming from options has precedence over the image one
	for varKey, varValue := range image.PersistentVars {
		_, exist := imagePersistentVars[varKey]
		if !exist {
			imagePersistentVars[varKey] = varValue
		}
	}
	options.PersistentVars = imagePersistentVars

	// copy the original map to avoid overwrite the original values
	imageVars := make(map[string]interface{})
	for k, v := range originalOptions.Vars {
		imageVars[k] = v
	}
	// add the image vars to the options build struct
	// a definition comming from options has precedence over the image one
	for varKey, varValue := range image.Vars {
		_, exist := imageVars[varKey]
		if !exist {
			imageVars[varKey] = varValue
		}
	}
	options.Vars = imageVars

	if options.RegistryHost == "" && len(image.Registry) > 0 {
		options.RegistryHost = image.Registry
	}
	if options.RegistryNamespace == "" && len(image.Namespace) > 0 {
		options.RegistryNamespace = image.Namespace
	}

	imageBuilder, err = e.getBuilder(image)
	if err != nil {
		return errors.New("(ImagesEngine::buildWorker)", "Error getting builder", err)
	}

	// add data comming from builder into build options
	options.BuilderOptions = imageBuilder.Options
	options.BuilderVarMappings = imageBuilder.VarMapping

	options.BuilderName = strings.Join([]string{"builder", options.RegistryNamespace, options.ImageName, options.ImageVersion}, "_")

	driver := imageBuilder.Driver
	// when dry-run is enabled is used the default driver which prints build options
	if options.DryRun {
		driver = defaultbuilder.DriverName
	}

	// Get builder's factory for the specified driver. It return a func that is executed later on with the specific options
	factoryFunc, exists = factory.GetDriverFactory(driver)
	if !exists {
		return errors.New("(ImagesEngine::buildWorker)", "Unexisting driver for builder '"+imageBuilder.Driver+"' required to build image '"+image.Name+"'")
	}

	// Fill options comming paranet node
	if node.Parent != nil {
		parent, err := tree.GetNodeImage(node.Parent)
		if err != nil {
			console.Warn(errors.New("(ImagesEngine::buildWorker)", "Parent image could not be achieved", err))
		} else {
			// when image from name is not defined image parent's name is used on build options
			if options.ImageFromName == "" && parent.Name != "" {
				options.ImageFromName = parent.Name
			}
			// when image from version/tag is not defined image parent's name is used on build options
			if options.ImageFromVersion == "" && parent.Version != "" {
				options.ImageFromVersion = parent.Version
			}
			// when image from namespace is not defined image parent's name is used on build options
			if options.ImageFromRegistryNamespace == "" && parent.Namespace != "" {
				options.ImageFromRegistryNamespace = parent.Namespace
			}
			// when image from registry host is not defined image parent's name is used on build options
			if options.ImageFromRegistryHost == "" && parent.Registry != "" {
				options.ImageFromRegistryHost = parent.Registry
			}
		}
	}

	// factory creates a new intances for the builder such NewAnsiblePlaybookBuilder
	builder, err = factoryFunc(e.context, options)
	if err != nil {
		return errors.New("(ImagesEngine::buildWorker)", "Error creating builder instance", err)
	}

	job := &Job{
		Driver: builder,
		Done:   make(chan bool),
		Err:    make(chan error),
	}

	e.Dispatch.Enqueue(job)

	select {
	case <-job.Done:
		// check when a build on cascade is requested
		if options.Cascade && depth != 0 {
			endChildBuild := make(chan bool)
			endChildBuildErr := make(chan error)

			for _, childNode := range node.Children {
				// for each image start the whole building process
				go func(node *gdstree.Node) {
					image, err := tree.GetNodeImage(node)
					if err != nil {
						endChildBuildErr <- errors.New("(goroutine::ImagesEngine::buildWorker)", "Error getting image from node", err)
						return
					}
					err = e.Build(image.Name, []string{image.Version}, originalOptions, depth-1)
					if err != nil {
						endChildBuildErr <- errors.New("(goroutine::ImagesEngine::build)", "Error building '"+image.Name+":"+image.Version+"'", err)
					} else {
						endChildBuild <- true
					}
				}(childNode)
			}

			var imageBuildErrs error
			for range node.Children {
				select {
				case <-endChildBuild:
				case imageBuildErr := <-endChildBuildErr:
					if imageBuildErrs == nil {
						imageBuildErrs = fmt.Errorf("")
					}
					imageBuildErrs = fmt.Errorf("%s\n%s", imageBuildErrs.Error(), imageBuildErr.Error())
				}
			}

			if imageBuildErrs != nil {
				return errors.New("(ImagesEngine::buildWorker)", "Error building '"+tree.GenerateNodeName(image)+"' children:\n"+imageBuildErrs.Error())
			}

		}
	case err = <-job.Err:
		return errors.New("(ImagesEngine::buildWorker)", "Error building image '"+image.Name+":"+image.Version+"'", err)
	}

	return nil
}

// getBuilder returns the builder required to build the image
func (e *ImagesEngine) getBuilder(i *image.Image) (*build.Builder, error) {
	var builder *build.Builder
	var err error

	// this validation is for compatibuility to v0.10.0
	if i.Type != "" && i.Builder == nil {
		i.Builder = i.Type
	}

	if i.Builder == nil {
		return nil, errors.New("(ImagesEngine::getBuilder)", fmt.Sprintf("Image '%s' has not a builder defined", i.Name))
	}

	switch i.Builder.(type) {
	case string:
		// Getting builder from global builders definition
		builder, err = e.Builders.GetBuilder(i.Builder.(string))
		if err != nil {
			return nil, errors.New("(ImagesEngine::getBuilder)", fmt.Sprintf("Error getting '%s' builder", i.Name), err)
		}
	case *build.Builder:
		builder = i.Builder.(*build.Builder)

	case interface{}:
		// In-line builder definition
		var data []byte
		builder = &build.Builder{}

		// builder should be defined as a yaml structure and is needed to marshal to a []byte before unmarshal it to the build.Builder struct
		data, err = yaml.Marshal(i.Builder)
		if err != nil {
			return nil, errors.New("(ImagesEngine::getBuilder)", fmt.Sprintf("Error marshaling '%s' builder's definition", i.Name), err)
		}

		err = yaml.Unmarshal(data, builder)
		if err != nil {
			return nil, errors.New("(ImagesEngine::getBuilder)", fmt.Sprintf("Error unmarshaling '%s' builder", i.Name), err)
		}

		// sanetize builder to ensure that all items are properly defined and initialized
		builder.SanetizeBuilder(i.Name)
	default:
		return nil, errors.New("(ImagesEngine::getBuilder)", "Builder definition type is unknown")
	}

	return builder, nil

}

func (e *ImagesEngine) DrawGraph(ctx context.Context) {
	g := e.ImagesGraph
	prefix := "\u251C\u2500\u2500\u2500"
	for _, rootNode := range g.Root {
		e.drawGraphRec(ctx, rootNode, prefix)
	}
}

func (e *ImagesEngine) drawGraphRec(ctx context.Context, nodeImage *gdstree.Node, prefix string) {

	msg := ""
	if nodeImage.Item == nil {
		msg = fmt.Sprintf(" %s %s", prefix, nodeImage.Name)
	} else {
		i := nodeImage.Item.(*image.Image)
		msg = fmt.Sprintf(" %s %s:%s", prefix, i.Name, i.Version)
	}
	console.Print(msg)

	prefix = "\u2502  " + prefix
	for _, child := range nodeImage.Children {
		e.drawGraphRec(ctx, child, prefix)
	}
}

func (e *ImagesEngine) ListImages() ([][]string, error) {
	var err error
	graph := e.ImagesGraph
	images := [][]string{}

	for _, root := range graph.Root {
		images, err = e.listImagesRec(root, images)
		if err != nil {

			return nil, errors.New("(ImagesEngine::ListImages)", "Error listing images", err)
		}
	}

	return images, nil
}

// listImagesRec
func (e *ImagesEngine) listImagesRec(nodeImage *gdstree.Node, listImages [][]string) ([][]string, error) {
	var err error
	var array []string

	img := nodeImage.Item.(*image.Image)
	array, err = img.ToArray()
	if err != nil {
		return nil, errors.New("(ImagesEngine::listImagesRec)", "Error listing images", err)
	}

	if nodeImage.Parent != nil {
		if nodeImage.Parent.Item != nil {
			p := nodeImage.Parent.Item.(*image.Image)
			array = append(array, p.Name+":"+p.Version)
		}
	} else {
		array = append(array, "-")
	}

	listImages = append(listImages, array)
	for _, child := range nodeImage.Children {
		listImages, err = e.listImagesRec(child, listImages)
		if err != nil {
			continue
		}
	}

	return listImages, nil
}

// ListImageHeader
func ListImageHeader() []string {
	h := []string{
		"NAME",
		"VERSION",
		"BUILDER",
		"NAMESPACE",
		"REGISTRY",
		"PARENT",
	}

	return h
}

// Promote an image
func (e *ImagesEngine) Promote(options *types.PromoteOptions) error {

	err := promote.Promote(e.context, options)
	if err != nil {
		return errors.New("(ImagesEngine::Promote)", "Error promoting image", err)
	}

	return nil
}
