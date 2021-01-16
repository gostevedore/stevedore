package engine

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"testing"

	"github.com/gostevedore/stevedore/internal/build"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/tree"
	"github.com/gostevedore/stevedore/internal/types"

	factory "github.com/gostevedore/stevedore/internal/driver"
	mockdriver "github.com/gostevedore/stevedore/internal/driver/mock"
	"github.com/gostevedore/stevedore/internal/schedule"
	"github.com/gostevedore/stevedore/internal/ui/console"

	errors "github.com/apenella/go-common-utils/error"
	gdstree "github.com/apenella/go-data-structures/tree"
	"github.com/stretchr/testify/assert"
)

// TestNewImagesEngine
func TestNewImagesEngine(t *testing.T) {

	testBaseDir := "test"
	testImagesBaseDir := filepath.Join(testBaseDir)
	testBuilderBaseDir := filepath.Join(testBaseDir)

	ctx := context.TODO()

	tests := []struct {
		desc          string
		err           error
		engine        *ImagesEngine
		ctx           context.Context
		numWorkers    int
		imageTreePath string
		builderPath   string
	}{
		{
			desc: "Testing create a new engine with wrong image tree file path",
			err: errors.New("(engine::NewImagesEngine)", "Error loading image tree definition file",
				errors.New("(tree::LoadImagesTree)", "Error loading images tree configuration",
					errors.New("", "(LoadYAMLFile) Error loading file unexistent. open unexistent: no such file or directory"))),
			engine:        &ImagesEngine{},
			ctx:           ctx,
			numWorkers:    1,
			imageTreePath: "unexistent",
			builderPath:   filepath.Join(testBuilderBaseDir, "stevedore_multiple_images.yml"),
		},
		{
			desc: "Testing create a new engine with wrong builders file path",
			err: errors.New("(engine::NewImagesEngine)", "Error loading builders definition file",
				errors.New("", "Could not be load configuration builders file"),
				errors.New("", "(LoadYAMLFile) Error loading file unexistent. open unexistent: no such file or directory")),
			engine:        &ImagesEngine{},
			ctx:           ctx,
			numWorkers:    1,
			imageTreePath: filepath.Join(testImagesBaseDir, "stevedore_multiple_images.yml"),
			builderPath:   "unexistent",
		},
		{
			desc:          "Testing create a new engine with a wrong number of workers on job despatcher",
			err:           errors.New("(engine::NewImagesEngine)", "Error creating dispatcher", errors.New("(schedule::NewDispatch)", "Invalid value for number of workers, it must be greater than zero")),
			engine:        &ImagesEngine{},
			ctx:           ctx,
			numWorkers:    0,
			imageTreePath: filepath.Join(testImagesBaseDir, "stevedore_multiple_images.yml"),
			builderPath:   filepath.Join(testBuilderBaseDir, "stevedore_multiple_images.yml"),
		},
		{
			desc: "Testing create a new image engine",
			err:  nil,
			engine: &ImagesEngine{
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": &image.Image{
								Name:      "ubuntu",
								Registry:  "registry",
								Builder:   "mock-builder",
								Namespace: "namespace",
								Version:   "16.04",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
							Item: &image.Image{
								Name:      "ubuntu",
								Registry:  "registry",
								Builder:   "mock-builder",
								Namespace: "namespace",
								Version:   "16.04",
								Tags:      []string{},
								Children:  map[string][]string{},
								Vars:      map[string]interface{}{},
							},
						},
					},
				},
			},
			ctx:           ctx,
			numWorkers:    1,
			imageTreePath: filepath.Join(testImagesBaseDir, "image_tree_single_image.yml"),
			builderPath:   filepath.Join(testBuilderBaseDir, "image_tree_single_image.yml"),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		engine, err := NewImagesEngine(test.ctx, test.numWorkers, test.imageTreePath, test.builderPath)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, test.engine.ImagesTree, engine.ImagesTree, "Unexpected value")
		}
	}
}

func TestFindNodes(t *testing.T) {

	testNode := &gdstree.Node{
		Name: "node",
	}
	testNode2 := &gdstree.Node{
		Name: "node2",
	}
	testNode3 := &gdstree.Node{
		Name: "node3",
	}
	testNode4 := &gdstree.Node{
		Name: "node4",
	}
	testNode5 := &gdstree.Node{
		Name: "node5",
		Item: &image.Image{
			Name:    "node5",
			Version: "*",
		},
	}

	tests := []struct {
		desc     string
		name     string
		versions []string
		engine   *ImagesEngine
		err      error
		res      map[string]uint8
	}{
		{
			desc:     "Testing find an image with no version defined",
			name:     "imageName",
			versions: nil,
			engine: &ImagesEngine{
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"imageName": {"imageName:imageVersion"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"imageName:imageVersion": {testNode, testNode2},
						"imageName:*":            {testNode, testNode3},
					},
				},
			},
			res: map[string]uint8{
				"node":  0,
				"node2": 0,
			},
			err: nil,
		},
		{
			desc:     "Testing find an undefined image with no version defined",
			name:     "imageNameUnexisting",
			versions: nil,
			engine: &ImagesEngine{
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"imageName": {"imageName:imageVersion"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"imageName:imageVersion": {testNode, testNode2},
						"imageName:*":            {testNode, testNode3},
					},
				},
			},
			err: errors.New("(ImagesEngine::findNodes)", "No image 'imageNameUnexisting' found on images tree",
				errors.New("(tree::Find)", "Error when finding images by name 'imageNameUnexisting'",
					errors.New("(tree::FindByName) ", "Image name 'imageNameUnexisting' does not exists"))),
			res: nil,
		},
		{
			desc:     "Testing find an image with a version defined",
			name:     "imageName",
			versions: []string{"imageVersion"},
			engine: &ImagesEngine{
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"imageName": {"imageName:imageVersion"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"imageName:imageVersion": {testNode, testNode2},
						"imageName:*":            {testNode3, testNode4},
					},
				},
			},
			err: nil,
			res: map[string]uint8{
				"node":  0,
				"node2": 0,
			},
		},
		{
			desc:     "Testing find an image by alternative name and version",
			name:     "imageName",
			versions: []string{"alternative-imageVersion"},
			engine: &ImagesEngine{
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"imageName": {"imageName:imageVersion"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"imageName:imageVersion": {testNode},
					},
					NameVersionAlternativeIndex: map[string][]*gdstree.Node{
						"imageName:alternative-imageVersion": {testNode},
					},
				},
			},
			err: nil,
			res: map[string]uint8{
				"node": 0,
			},
		},
		{
			desc:     "Testing find an image with a wildcard version defined",
			name:     "node5",
			versions: []string{"wildcard"},
			engine: &ImagesEngine{
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"node5": {"node5:imageVersion", "node5:*"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"node5:imageVersion": {testNode, testNode2},
						"node5:*":            {testNode5},
					},
				},
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"node5": {
							"imageVersion": &image.Image{
								Name:    "imageName",
								Version: "imageVersion",
							},
							"*": &image.Image{
								Name:    "imageName",
								Version: "*",
							},
						},
					},
				},
			},
			err: nil,
			res: map[string]uint8{
				"imageName:wildcard": 0,
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		list, err := test.engine.findNodes(test.name, test.versions)

		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {

			for _, node := range list {
				_, exists := test.res[node.Name]
				assert.True(t, exists, "unexisting key "+node.Name)
				delete(test.res, node.Name)
			}
			assert.Equal(t, 0, len(test.res), "Not all expected nodes has been created", list)
		}
	}
}

// TestEngineBuild
func TestEngineBuild(t *testing.T) {

	ctx := context.TODO()

	dispatch, err := schedule.NewDispatch(ctx, 1)
	if err != nil {
		t.Fatalf("Error creating dispatcher: %s", err.Error())
	}

	err = dispatch.Start()
	if err != nil {
		t.Fatalf("Error starting dispatcher: %s", err.Error())
	}

	imageUbuntu1604 := &image.Image{
		Name:      "ubuntu",
		Version:   "16.04",
		Namespace: "namespace",
		Registry:  "registry",
		Builder:   "builder",
		PersistentVars: map[string]interface{}{
			"pvar1": "pvar1",
		},
		Vars: map[string]interface{}{
			"var1": "var1",
		},
	}
	nodeUbuntu1604 := &gdstree.Node{
		Name: "ubuntu:16.04",
		Item: imageUbuntu1604,
	}
	imageUbuntu1804 := &image.Image{
		Name:      "ubuntu",
		Version:   "18.04",
		Namespace: "namespace",
		Registry:  "registry",
		Builder:   "builder",
	}
	nodeUbuntu1804 := &gdstree.Node{
		Name: "ubuntu:16.04",
		Item: imageUbuntu1804,
	}

	tests := []struct {
		desc      string
		err       error
		engine    *ImagesEngine
		imageName string
		versions  []string
		options   *types.BuildOptions
		depth     int
		preFunc   func()
	}{
		{
			desc:      "Testing build an image with nil options",
			err:       errors.New("(ImagesEngine::Build)", "Build options is nil"),
			imageName: "image",
			versions:  []string{},
			depth:     -1,
			engine:    &ImagesEngine{},
			options:   nil,
			preFunc: func() {
				factory.ClearDriverFactory()
			},
		},
		{
			desc: "Testing build an image unexisting image",
			err: errors.New("(ImagesEngine::Build)", "Error finding image 'unexists' and versions '[]'",
				errors.New("(ImagesEngine::findNodes)", "No image 'unexists' found on images tree",
					errors.New("(tree::Find)", "Error when finding images by name 'unexists'",
						errors.New("(tree::FindByName)", "Image name 'unexists' does not exists")))),
			imageName: "unexists",
			versions:  []string{},
			depth:     -1,
			engine: &ImagesEngine{
				Dispatch: &schedule.Dispatch{},
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": imageUbuntu1604,
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:16.04": {
							Name: "ubuntu:16.04",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu": {"ubuntu:16.04"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:16.04": {nodeUbuntu1604},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"mock-driver": {
							Name:   "mock-driver",
							Driver: "ansible-playbook",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{},
			preFunc: func() {
				factory.ClearDriverFactory()
			},
		},
		{
			desc: "Testing build an image with an unexisting builder",
			err: errors.New("(ImagesEngine::Build)", "Image could not be built",
				errors.New("(ImagesEngine::buildWorker)", "Error building 'ubuntu:16.04'",
					errors.New("", "Error getting builder",
						errors.New("", "Image 'ubuntu' has not a builder defined")))),
			imageName: "ubuntu",
			versions:  []string{},
			depth:     -1,
			engine: &ImagesEngine{
				Dispatch: &schedule.Dispatch{},
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": &image.Image{
								Name:    "ubuntu",
								Version: "16.04",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:16.04": {
							Name: "ubuntu:16.04",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu": {"ubuntu:16.04"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:16.04": {
							&gdstree.Node{
								Name: "ubuntu:16.04",
								Item: &image.Image{
									Name:    "ubuntu",
									Version: "16.04",
								},
							},
						},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{},
			preFunc: func() {
				factory.ClearDriverFactory()
			},
		},
		{
			desc: "Testing build an image with an unexisting driver",
			err: errors.New("(ImagesEngine::Build)", "Image could not be built",
				errors.New("(ImagesEngine::buildWorker)", "Error building 'ubuntu:16.04'",
					errors.New("", "Unexisting driver for builder 'mock-builder' required to build image 'ubuntu'"))),
			imageName: "ubuntu",
			versions:  []string{},
			depth:     -1,
			engine: &ImagesEngine{
				Dispatch: &schedule.Dispatch{},
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": imageUbuntu1604,
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:16.04": {
							Name: "ubuntu:16.04",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu": {"ubuntu:16.04"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:16.04": {nodeUbuntu1804},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{},
			preFunc: func() {
				factory.ClearDriverFactory()
			},
		},
		{
			desc: "Testing build an image with an unkowm version",
			err: errors.New("(ImagesEngine::Build)", "Error finding image 'ubuntu' and versions '[unknown]'",
				errors.New("", "No matching images to be build found for 'ubuntu' (versions: [unknown])")),
			imageName: "ubuntu",
			versions:  []string{"unknown"},
			depth:     -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": imageUbuntu1604,
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:16.04": {
							Name: "ubuntu:16.04",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu": {"ubuntu:16.04"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:16.04": {nodeUbuntu1604},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"mock-builder": {
							Name:   "mock-builder",
							Driver: "ansible-playbook",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{
				ConnectionLocal: true,
			},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
		},
		{
			desc:      "Testing build an image",
			err:       &errors.Error{},
			imageName: "ubuntu",
			versions:  []string{"16.04"},
			depth:     -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": imageUbuntu1604,
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:16.04": {
							Name: "ubuntu:16.04",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu": {"ubuntu:16.04"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:16.04": {nodeUbuntu1604},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{
				ConnectionLocal: true,
				PersistentVars: map[string]interface{}{
					"pvar1": "pvar1",
				},
				Vars: map[string]interface{}{
					"var1": "var1",
				},
			},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
		},
		{
			desc:      "Testing build an image with no version defined",
			err:       &errors.Error{},
			imageName: "ubuntu",
			versions:  nil,
			depth:     -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": imageUbuntu1604,
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:16.04": {
							Name: "ubuntu:16.04",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu": {"ubuntu:16.04"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:16.04": {nodeUbuntu1604},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{
				ConnectionLocal: true,
			},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
		},
		{
			desc:      "Testing build an image with two versions defined",
			err:       &errors.Error{},
			imageName: "ubuntu",
			versions:  nil,
			depth:     -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": imageUbuntu1604,
							"18.04": imageUbuntu1804,
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
						},
						{
							Name: "ubuntu:18.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:16.04": {
							Name: "ubuntu:16.04",
						},
						"ubuntu:18.04": {
							Name: "ubuntu:18.04",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu": {"ubuntu:16.04", "ubuntu:18.04"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:16.04": {nodeUbuntu1604},
						"ubuntu:18.04": {nodeUbuntu1804},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{
				ConnectionLocal: true,
			},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
		},
		{
			desc:      "Testing build an image on cascade",
			err:       &errors.Error{},
			imageName: "ubuntu",
			versions:  nil,
			depth:     -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": &image.Image{
								Name:      "ubuntu",
								Version:   "16.04",
								Namespace: "namespace",
								Registry:  "registry",
								Builder:   "mock-builder",
								Children: map[string][]string{
									"php_fpm": {
										"7.1",
										"7.2",
									},
								},
							},
						},
						"php_fpm": {
							"7.1": &image.Image{
								Name:      "php_fpm",
								Version:   "7.1",
								Namespace: "namespace",
								Registry:  "registry",
								Builder:   "mock-builder",
							},
							"7.2": &image.Image{
								Name:      "php_fpm",
								Version:   "7.2",
								Namespace: "namespace",
								Registry:  "registry",
								Builder:   "mock-builder",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:16.04": {
							Name: "ubuntu:16.04",
						},
						"php-fpm:7.1": {
							Name: "php-fpm:7.1",
						},
						"php-fpm:7.2": {
							Name: "php-fpm:7.2",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu":  {"ubuntu:16.04", "ubuntu:18.04"},
						"php-fpm": {"php-fpm:7.1", "php-fpm:7.2"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:16.04": {nodeUbuntu1604},
						"ubuntu:18.04": {nodeUbuntu1604},
						"php-fpm:7.1":  {nodeUbuntu1604},
						"php-fpm:7.2":  {nodeUbuntu1604},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{
				Cascade:         true,
				ConnectionLocal: true,
			},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.preFunc != nil {
				test.preFunc()
			}

			err := test.engine.Build(test.imageName, test.versions, test.options, test.depth)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			}
		})
	}
}

// TestEngineBuildWorker tests buildWorker method
func TestEngineBuildWorker(t *testing.T) {

	ctx := context.TODO()

	imageUbuntu1604 := &image.Image{
		Name:      "ubuntu",
		Version:   "16.04",
		Namespace: "namespace",
		Builder:   "mock-builder",
	}
	nodeUbuntu1604 := &gdstree.Node{
		Name: "ubuntu:16.04",
		Item: imageUbuntu1604,
	}

	dispatch, err := schedule.NewDispatch(ctx, 1)
	if err != nil {
		t.Fatalf("Error creating dispatcher: %s", err.Error())
	}

	err = dispatch.Start()
	if err != nil {
		t.Fatalf("Error starting dispatcher: %s", err.Error())
	}

	tests := []struct {
		desc    string
		err     error
		node    *gdstree.Node
		engine  *ImagesEngine
		options *types.BuildOptions
		depth   int
		preFunc func()
	}{
		{
			desc:    "Testing buildWorker when is gave an undefined image",
			err:     errors.New("(ImagesEngine::buildWorker)", "Node is not defined"),
			node:    nil,
			depth:   -1,
			options: &types.BuildOptions{},
			engine:  &ImagesEngine{},
			preFunc: nil,
		},
		{
			desc:    "Testing buildWorker when is gave an undefined builder options",
			err:     errors.New("(ImagesEngine::buildWorker)", "Builder options is not defined"),
			node:    &gdstree.Node{},
			depth:   -1,
			options: nil,
			engine:  &ImagesEngine{},
			preFunc: nil,
		},
		{
			desc: "Testing buildWorker when builder configuration does not exist",
			err: errors.New("(ImagesEngine::buildWorker)", "Error getting builder",
				errors.New("", "Error getting 'ubuntu' builder",
					errors.New("", "Unexisting builder configuration for type 'mock-builder'"))),
			node: &gdstree.Node{
				Name: "ubuntu",
				Item: &image.Image{
					Name:    "ubuntu",
					Builder: "mock-builder",
				},
			},
			depth: -1,
			engine: &ImagesEngine{
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-driver",
						},
					},
				},
			},
			options: &types.BuildOptions{},
			preFunc: func() {
				factory.ClearDriverFactory()
			},
		},
		{
			desc: "Testing buildWorker when builder's driver does not exist",
			err:  errors.New("(ImagesEngine::buildWorker)", "Unexisting driver for builder 'mock-builder' required to build image 'ubuntu'"),
			node: &gdstree.Node{
				Name: "ubuntu",
				Item: &image.Image{
					Name:    "ubuntu",
					Builder: "builder",
				},
			},
			depth: -1,
			engine: &ImagesEngine{
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
						},
					},
				},
			},
			options: &types.BuildOptions{},
			preFunc: func() {
				factory.ClearDriverFactory()
			},
		},
		{
			desc: "Testing buildWorker error when is created a builder instance",
			err: errors.New("(ImagesEngine::buildWorker)", "Error creating builder instance",
				errors.New("", "Error")),
			node: &gdstree.Node{
				Name: "ubuntu",
				Item: &image.Image{
					Name:    "ubuntu",
					Version: "16.04",
					Builder: "builder",
				},
			},
			depth: -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriverErrOnNew)
			},
		},
		{
			desc: "Testing buildWorker when job returns an error",
			err:  errors.New("(ImagesEngine::buildWorker)", "Error building image 'ubuntu:16.04'", errors.New("(MockBuilderRunErr)", "Error")),
			node: &gdstree.Node{
				Name: "ubuntu",
				Item: &image.Image{
					Name:    "ubuntu",
					Version: "16.04",
					Builder: "builder",
				},
			},
			depth: -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriverErr)
			},
		},
		{
			desc: "Testing buildWorker in cascade mode",
			err:  &errors.Error{},
			node: &gdstree.Node{
				Name: "ubuntu",
				Item: &image.Image{
					Name:    "ubuntu",
					Builder: "builder",
					Children: map[string][]string{
						"nginx": {"1.10"},
					},
				},
			},
			depth: -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"16.04": &image.Image{
								Name:    "ubuntu",
								Builder: "builder",
							},
						},
						"nginx": {
							"1.10": &image.Image{
								Name:    "nginx",
								Builder: "builder",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:16.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:16.04": {
							Name: "ubuntu:16.04",
						},
						"nginx:1.10": {
							Name: "nginx:1.10",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu": {"ubuntu:16.04"},
						"nginx":  {"nginx:1.10"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:16.04": {nodeUbuntu1604},
						"nginx:1.10":   {nodeUbuntu1604},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{
				Cascade: true,
			},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
		},
		{
			desc: "Testing buildWorker and setting details from parent node to build",
			err:  &errors.Error{},
			node: &gdstree.Node{
				Name: "ubuntu",
				Item: &image.Image{
					Name:    "nginx",
					Builder: "builder",
					Children: map[string][]string{
						"nginx": {"1.19"},
					},
				},
				Parent: &gdstree.Node{
					Name: "ubuntu",
					Item: &image.Image{
						Name:      "ubuntu",
						Version:   "20.04",
						Namespace: "namespace",
						Registry:  "registry",
						Builder:   "builder",
						Children: map[string][]string{
							"nginx": {"20.04"},
						},
					},
				},
			},
			depth: -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"ubuntu": {
							"20.04": &image.Image{
								Name:    "ubuntu",
								Version: "20.04",
								Builder: "builder",
							},
						},
						"nginx": {
							"1.19": &image.Image{
								Name:    "nginx",
								Builder: "builder",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "ubuntu:20.04",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"ubuntu:20.04": {
							Name: "ubuntu:20.04",
						},
						"nginx:1.19": {
							Name: "nginx:1.19",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"ubuntu": {"ubuntu:20.04"},
						"nginx":  {"nginx:1.19"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:20.04": {nodeUbuntu1604},
						"nginx:1.19":   {nodeUbuntu1604},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{
				Cascade: true,
			},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
		},
		{
			desc: "Testing buildWorker and generate extra tags based on semver",
			err:  &errors.Error{},
			node: &gdstree.Node{
				Name: "semver-image",
				Item: &image.Image{
					Name:    "semver-image",
					Version: "1.2.3",
					Builder: "builder",
				},
			},
			depth: -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"semver-image": {
							"1.2.3": &image.Image{
								Name:    "semver-image",
								Version: "1.2.3",
								Builder: "builder",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "semver-image:1.2.3",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"semver-image:1.2.3": {
							Name: "semver-image:1.2.3",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"semver-image": {"semver-image:1.2.3"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:1.2.3": nil,
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{
				EnableSemanticVersionTags: true,
			},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
		},
		{
			desc: "Testing buildWorker error when generate extra tags based on semver",
			err:  &errors.Error{},
			node: &gdstree.Node{
				Name: "semver-image",
				Item: &image.Image{
					Name:    "semver-image",
					Version: "a.b.c",
					Builder: "builder",
				},
			},
			depth: -1,
			engine: &ImagesEngine{
				Dispatch: dispatch,
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"semver-image": {
							"1.2.3": &image.Image{
								Name:    "semver-image",
								Version: "1.2.3",
								Builder: "builder",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "semver-image:1.2.3",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"semver-image:1.2.3": {
							Name: "semver-image:1.2.3",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"semver-image": {"semver-image:1.2.3"},
					},
					NameVersionIndex: map[string][]*gdstree.Node{
						"ubuntu:1.2.3": nil,
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"builder": {
							Name:   "builder",
							Driver: "mock-builder",
							Options: map[string]interface{}{
								"inventory": "127.0.0.1,",
								"playbook":  "site.yml",
							},
						},
					},
				},
			},
			options: &types.BuildOptions{
				EnableSemanticVersionTags: true,
			},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.preFunc != nil {
				test.preFunc()
			}

			err := test.engine.buildWorker(test.node, test.options, test.depth)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			}
		})
	}
}

func TestGetBuilder(t *testing.T) {

	builder := &build.Builder{
		Name:    "mock-builder",
		Driver:  "mock-builder",
		Options: map[string]interface{}{},
	}

	tests := []struct {
		desc    string
		err     error
		image   *image.Image
		engine  *ImagesEngine
		options *types.BuildOptions
		res     *build.Builder
		preFunc func()
	}{
		{
			desc: "Testing get builder from builders configuration from type",
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
			image: &image.Image{
				Name:    "my-image",
				Version: "1.2.3",
				Type:    "mock-builder",
			},
			engine: &ImagesEngine{
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"my-image": {
							"1.2.3": &image.Image{
								Name:    "my-image",
								Version: "1.2.3",
								Builder: "mock-builder",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "my-image:1.2.3",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"my-image:1.2.3": {
							Name: "my-image:1.2.3",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"my-image": {"my-image:1.2.3"},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"mock-builder": builder,
					},
				},
			},
			res: builder,
		},
		{
			desc: "Testing get builder from builders configuration",
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
			image: &image.Image{
				Name:    "my-image",
				Version: "1.2.3",
				Builder: "mock-builder",
			},
			engine: &ImagesEngine{
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"my-image": {
							"1.2.3": &image.Image{
								Name:    "my-image",
								Version: "1.2.3",
								Builder: "mock-builder",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "my-image:1.2.3",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"my-image:1.2.3": {
							Name: "my-image:1.2.3",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"my-image": {"my-image:1.2.3"},
					},
				},
				Builders: &build.Builders{
					Builders: map[string]*build.Builder{
						"mock-builder": builder,
					},
				},
			},
			res: builder,
		},
		{
			desc: "Testing get in-line builder from image",
			err:  &errors.Error{},
			preFunc: func() {
				factory.ClearDriverFactory()
				factory.RegisterDriverFactory("mock-builder", mockdriver.NewMockDriver)
			},
			image: &image.Image{
				Name:    "my-image",
				Version: "1.2.3",
				Builder: &build.Builder{
					Name:    "mock-builder",
					Driver:  "mock-builder",
					Options: map[string]interface{}{},
				},
			},
			engine: &ImagesEngine{
				ImagesTree: &tree.ImagesTree{
					Images: map[string]map[string]*image.Image{
						"my-image": {
							"1.2.3": &image.Image{
								Name:    "my-image",
								Version: "1.2.3",
								Builder: "mock-builder",
							},
						},
					},
				},
				ImagesGraph: &gdstree.Graph{
					Root: []*gdstree.Node{
						{
							Name: "my-image:1.2.3",
						},
					},
					NodesIndex: map[string]*gdstree.Node{
						"my-image:1.2.3": {
							Name: "my-image:1.2.3",
						},
					},
				},
				ImageIndex: &tree.ImageIndex{
					NameIndex: map[string][]string{
						"my-image": {"my-image:1.2.3"},
					},
				},
			},
			res: builder,
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.preFunc != nil {
				test.preFunc()
			}

			builder, err := test.engine.getBuilder(test.image)

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, builder, "Unexpected builder")
			}
		})
	}

}

func TestDrawGraph(t *testing.T) {

	var w bytes.Buffer
	console.SetWriter(io.Writer(&w))

	ctx := context.TODO()

	phpFpmDev71 := &gdstree.Node{
		Name: "php-fpm-dev:7.1",
		Item: &image.Image{
			Name:    "php-fpm-dev",
			Version: "7.1",
		},
	}
	phpCliDev71 := &gdstree.Node{
		Name: "php-cli-dev:7.1",
		Item: &image.Image{
			Name:    "php-cli-dev",
			Version: "7.1",
		},
	}
	phpFpm71 := &gdstree.Node{
		Name: "php-fpm:7.1",
		Item: &image.Image{
			Name:    "php-fpm",
			Version: "7.1",
		},
		Children: []*gdstree.Node{
			phpFpmDev71,
		},
	}
	phpFpm72 := &gdstree.Node{
		Name: "php-fpm:7.2",
		Item: &image.Image{
			Name:    "php-fpm",
			Version: "7.2",
		},
	}
	phpCli71 := &gdstree.Node{
		Name: "php-cli:7.1",
		Item: &image.Image{
			Name:    "php-cli",
			Version: "7.1",
		},
		Children: []*gdstree.Node{
			phpFpmDev71,
		},
	}
	phpCli72 := &gdstree.Node{
		Name: "php-cli:7.2",
		Item: &image.Image{
			Name:    "php-cli",
			Version: "7.2",
		},
	}
	phpBuilder71 := &gdstree.Node{
		Name: "php-builder:7.1",
		Item: &image.Image{
			Name:    "php-builder",
			Version: "7.1",
		},
	}
	ubuntu := &gdstree.Node{
		Name: "ubuntu:16.04",
		Children: []*gdstree.Node{
			phpFpm71,
			phpFpm72,
			phpCli71,
			phpCli72,
			phpBuilder71,
		},
	}

	imagesGraph := &gdstree.Graph{
		Root: []*gdstree.Node{
			ubuntu,
		},
		NodesIndex: map[string]*gdstree.Node{
			"ubuntu":          ubuntu,
			"php-fpm:7.1":     phpFpm71,
			"php-fpm:7.2":     phpFpm72,
			"php-cli:7.1":     phpCli71,
			"php-cli:7.2":     phpCli72,
			"php-builder:7.1": phpBuilder71,
			"php-fpm-dev:7.1": phpFpmDev71,
			"php-cli-dev:7.1": phpCliDev71,
		},
	}

	engine := &ImagesEngine{
		ImagesGraph: imagesGraph,
	}

	res := " \u251C\u2500\u2500\u2500 ubuntu:16.04\n"
	res = res + " \u2502  \u251C\u2500\u2500\u2500 php-fpm:7.1\n"
	res = res + " \u2502  \u2502  \u251C\u2500\u2500\u2500 php-fpm-dev:7.1\n"
	res = res + " \u2502  \u251C\u2500\u2500\u2500 php-fpm:7.2\n"
	res = res + " \u2502  \u251C\u2500\u2500\u2500 php-cli:7.1\n"
	res = res + " \u2502  \u2502  \u251C\u2500\u2500\u2500 php-fpm-dev:7.1\n"
	res = res + " \u2502  \u251C\u2500\u2500\u2500 php-cli:7.2\n"
	res = res + " \u2502  \u251C\u2500\u2500\u2500 php-builder:7.1\n"

	engine.DrawGraph(ctx)
	assert.Equal(t, res, w.String(), "Output not equal")
}

func TestListImages(t *testing.T) {

	nginx := &gdstree.Node{
		Name: "nginx:1.15",
		Item: &image.Image{
			Name:    "nginx",
			Version: "1.15-ubuntu16.04",
			Builder: "infrastructure",
		},
	}
	ubuntu := &gdstree.Node{
		Name: "ubuntu:16.04",
		Item: &image.Image{
			Name:    "ubuntu",
			Version: "16.04",
			Builder: "infrastructure",
		},
		Children: []*gdstree.Node{
			nginx,
		},
	}
	nginx.Parent = ubuntu

	imagesGraph := &gdstree.Graph{
		Root: []*gdstree.Node{
			ubuntu,
		},
		NodesIndex: map[string]*gdstree.Node{
			"ubuntu": ubuntu,
			"nginx":  nginx,
		},
	}

	engine := &ImagesEngine{
		ImagesGraph: imagesGraph,
	}
	expected := [][]string{
		{"ubuntu", "16.04", "infrastructure", "", "", "-"},
		//Parent not correct
		{"nginx", "1.15-ubuntu16.04", "infrastructure", "", "", "ubuntu:16.04"},
	}
	list, _ := engine.ListImages()
	assert.Equal(t, expected, list)
}

func TestListImageHeader(t *testing.T) {

	t.Log("Testing list Builders header")
	expected := []string{"NAME", "VERSION", "BUILDER", "NAMESPACE", "REGISTRY", "PARENT"}
	res := ListImageHeader()

	assert.Equal(t, expected, res)
}
