package builders

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestAddBuilder(t *testing.T) {
	errContext := "(builders::AddBuilder)"

	tests := []struct {
		desc     string
		err      error
		builders *Builders
		builder  *builder.Builder
		res      map[string]*builder.Builder
	}{
		{
			desc:     "Testing add a builder",
			builders: NewBuilders(afero.NewMemMapFs()),
			err:      &errors.Error{},
			builder: &builder.Builder{
				Name: "first",
			},
			res: map[string]*builder.Builder{
				"first": {Name: "first"},
			},
		},

		{
			desc: "Testing error adding already existing builder",
			builders: &Builders{
				fs: afero.NewMemMapFs(),
				Builders: map[string]*builder.Builder{
					"first": {Name: "first"},
				},
			},
			err: errors.New(errContext, "Builder 'first' already exist"),
			builder: &builder.Builder{
				Name: "first",
			},
			res: map[string]*builder.Builder{
				"first": {Name: "first"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.builders.AddBuilder(test.builder)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.builders.Builders)
			}
		})
	}
}

func TestGetBuilder(t *testing.T) {
	errContext := "(builders::GetBuilder)"

	tests := []struct {
		desc     string
		err      error
		builders *Builders
		builder  string
		res      *builder.Builder
	}{
		{
			desc: "Testing get a builder",
			builders: &Builders{
				fs: afero.NewMemMapFs(),
				Builders: map[string]*builder.Builder{
					"first": {Name: "first"},
				},
			},
			err:     &errors.Error{},
			builder: "first",
			res: &builder.Builder{
				Name: "first",
			},
		},

		{
			desc: "Testing error getting an unexisting",
			builders: &Builders{
				fs:       afero.NewMemMapFs(),
				Builders: map[string]*builder.Builder{},
			},
			err:     errors.New(errContext, "Builder 'first' does not exists"),
			builder: "first",

			res: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := test.builders.GetBuilder(test.builder)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestLoadBuildersFromFile(t *testing.T) {
	var err error

	errContext := "(builders::loadBuilderFile)"

	testFs := afero.NewMemMapFs()
	testFs.MkdirAll("/builders", 0755)

	err = afero.WriteFile(testFs, "/builders/file1.yml", []byte(`
builders:
  first:
    driver: docker
    options:
      dockerfile: Dockerfile.test
      context:
        - path: /path/to/context
  second:
    driver: docker
    options:
      dockerfile: Dockerfile.test
      context:
        - path: /path/to/another/context
`), 0666)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, "/builders/tab_error_file.yml", []byte(`
builders:
  first:
    driver: docker
    options:
	  dockerfile: Dockerfile.test
      context:
      path: /path/to/context
`), 0666)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc     string
		path     string
		builders *Builders
		res      map[string]*builder.Builder
		err      error
	}{
		{
			desc:     "Testing loading builders from file",
			path:     "/builders/file1.yml",
			err:      &errors.Error{},
			builders: NewBuilders(testFs),
			res: map[string]*builder.Builder{
				"first": {
					Name:   "first",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/context"}},
					},
				},
				"second": {
					Name:   "second",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/another/context"}},
					},
				},
			},
		},
		{
			desc:     "Testing error loading builders from file",
			path:     "/builders/tab_error_file.yml",
			builders: NewBuilders(testFs),
			res:      nil,
			err:      errors.New(errContext, "Error loading builders from file '/builders/tab_error_file.yml'\nfound:\n\nbuilders:\n  first:\n    driver: docker\n    options:\n\t  dockerfile: Dockerfile.test\n      context:\n      path: /path/to/context\n\n\tyaml: line 6: found character that cannot start any token"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.builders.LoadBuildersFromFile(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.builders.Builders)
			}
		})
	}

}

func TestLoadBuildersFromDir(t *testing.T) {
	var err error

	errContext := "(builders::loadBuildersFromDir)"

	testFs := afero.NewMemMapFs()
	testFs.MkdirAll("/builders", 0755)
	testFs.MkdirAll("/builders_error", 0755)

	err = afero.WriteFile(testFs, "/builders/file1.yml", []byte(`
builders:
  first:
    driver: docker
    options:
      dockerfile: Dockerfile.test
      context:
        - path: /path/to/context
  second:
    driver: docker
    options:
      dockerfile: Dockerfile.test
      context:
        - path: /path/to/another/context
`), 0666)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, "/builders/file2.yml", []byte(`
builders:
  third:
    driver: docker
    options:
      dockerfile: Dockerfile.test
      context:
        - path: /even/another/path/to/context
`), 0666)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, "/builders_error/file1.yml", []byte(`
builders:
  third:
    driver: docker
    options:
	  dockerfile: Dockerfile.test
      context:
        - path: /even/another/path/to/context
`), 0666)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, "/builders_error/file2.yml", []byte(`
builders:
  fourth:
    driver: docker
    options:
	  dockerfile: Dockerfile.test
      context:
        - path: /even/another/path/to/context
`), 0666)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc     string
		path     string
		builders *Builders
		res      map[string]*builder.Builder
		err      error
	}{
		{
			desc:     "Testing loading builders from directory",
			path:     "/builders",
			err:      &errors.Error{},
			builders: NewBuilders(testFs),
			res: map[string]*builder.Builder{
				"first": {
					Name:   "first",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/context"}},
					},
				},
				"second": {
					Name:   "second",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/another/context"}},
					},
				},
				"third": {
					Name:   "third",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/even/another/path/to/context"}},
					},
				},
			},
		},
		{
			desc:     "Testing error loading builders from directory",
			path:     "/builders_error",
			err:      errors.New(errContext, "Error loading builders from file '/builders_error/file1.yml'\nfound:\n\nbuilders:\n  third:\n    driver: docker\n    options:\n\t  dockerfile: Dockerfile.test\n      context:\n        - path: /even/another/path/to/context\n\n\tyaml: line 6: found character that cannot start any token\nError loading builders from file '/builders_error/file2.yml'\nfound:\n\nbuilders:\n  fourth:\n    driver: docker\n    options:\n\t  dockerfile: Dockerfile.test\n      context:\n        - path: /even/another/path/to/context\n\n\tyaml: line 6: found character that cannot start any token\n"),
			builders: NewBuilders(testFs),
			res:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.builders.LoadBuildersFromDir(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.builders.Builders)
			}
		})
	}
}

func TestLoadBuilders(t *testing.T) {
	var err error

	errContext := "(builders::loadBuildersFromDir)"

	testFs := afero.NewMemMapFs()
	testFs.MkdirAll("/builders", 0755)
	testFs.MkdirAll("/builders_error", 0755)

	err = afero.WriteFile(testFs, "/builders/file1.yml", []byte(`
builders:
  first:
    driver: docker
    options:
      dockerfile: Dockerfile.test
      context:
        - path: /path/to/context
  second:
    driver: docker
    options:
      dockerfile: Dockerfile.test
      context:
        - path: /path/to/another/context
`), 0666)
	if err != nil {
		t.Log(err)
	}

	err = afero.WriteFile(testFs, "/builders/file2.yml", []byte(`
builders:
  third:
    driver: docker
    options:
      dockerfile: Dockerfile.test
      context:
        - path: /even/another/path/to/context
`), 0666)

	err = afero.WriteFile(testFs, "/builders_error/file1.yml", []byte(`
builders:
  third:
    driver: docker
    options:
	  dockerfile: Dockerfile.test
      context:
        - path: /even/another/path/to/context
`), 0666)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc     string
		path     string
		builders *Builders
		res      map[string]*builder.Builder
		err      error
	}{
		{
			desc:     "Testing loading builders from file",
			path:     "/builders/file1.yml",
			builders: NewBuilders(testFs),
			res: map[string]*builder.Builder{
				"first": {
					Name:   "first",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/context"}},
					},
				},
				"second": {
					Name:   "second",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/another/context"}},
					},
				},
			},
			err: &errors.Error{},
		},
		{
			desc:     "Testing loading builders from directory",
			path:     "/builders",
			builders: NewBuilders(testFs),
			res: map[string]*builder.Builder{
				"first": {
					Name:   "first",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/context"}},
					},
				},
				"second": {
					Name:   "second",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/another/context"}},
					},
				},
				"third": {
					Name:   "third",
					Driver: "docker",
					Options: &builder.BuilderOptions{
						Dockerfile: "Dockerfile.test",
						Context:    []*builder.DockerDriverContextOptions{{Path: "/even/another/path/to/context"}},
					},
				},
			},
			err: &errors.Error{},
		},

		{
			desc:     "Testing error loading builders from unexisting directory",
			path:     "/builders_unexisting",
			builders: NewBuilders(testFs),
			res:      nil,
			err:      errors.New(errContext, "open /builders_unexisting: file does not exist"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.builders.LoadBuilders(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.builders.Builders)
			}
		})
	}
}

// // TestLoadImage tests
// func TestLoadBuilders(t *testing.T) {

// 	testBaseDir := "test"

// 	tests := []struct {
// 		desc     string
// 		file     string
// 		err      error
// 		builders *Builders
// 	}{
// 		{
// 			desc:     "testing an unexistent file",
// 			file:     "nofile",
// 			err:      errors.New("(builder::LoadBuilders)", "Could not be load configuration builders file", errors.New("", "(LoadYAMLFile) Error loading file nofile. open nofile: no such file or directory")),
// 			builders: &Builders{},
// 		},
// 		{
// 			desc:     "testing nil drivers",
// 			file:     filepath.Join(testBaseDir, "nil_builders.yml"),
// 			err:      &errors.Error{},
// 			builders: &Builders{map[string]*build.Builder{}},
// 		},
// 		{
// 			desc: "testing builders",
// 			file: filepath.Join(testBaseDir, "builders.yml"),
// 			err:  &errors.Error{},
// 			builders: &Builders{
// 				Builders: map[string]*build.Builder{
// 					"infrastructure": {
// 						Name:   "infrastructure",
// 						Driver: "ansible-playbook",
// 						Options: map[string]interface{}{
// 							"inventory": "inventory/all",
// 							"playbook":  "site.yml",
// 						},
// 						VarMapping: varsmap.Varsmap{
// 							varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
// 							varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
// 							varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
// 							varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
// 							varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
// 							varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
// 							varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
// 							varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
// 							varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
// 							varsmap.VarMappingImageNameKey:                     "image",
// 							varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
// 							varsmap.VarMappingImageExtraTagsKey:                varsmap.VarMappingImageExtraTagsDefaultValue,
// 							varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
// 							varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
// 							varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
// 						},
// 					},
// 					"php-code": {
// 						Name:   "php-code",
// 						Driver: "ansible-playbook",
// 						Options: map[string]interface{}{
// 							"inventory": "inventory/all",
// 							"playbook":  "code_builder.yml",
// 						},
// 						VarMapping: varsmap.New(),
// 					},
// 					"dummy": {
// 						Name:       "dummy",
// 						Driver:     "default",
// 						Options:    map[string]interface{}{},
// 						VarMapping: varsmap.New(),
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Log(test.desc)

// 		builders, err := LoadBuilders(test.file)
// 		if err != nil && assert.Error(t, err) {
// 			assert.Equal(t, test.err.Error(), err.Error())
// 		} else {
// 			assert.Equal(t, builders, test.builders, "Unexpected value")
// 		}
// 	}
// }

// func TestGetBuilders(t *testing.T) {

// 	b1 := &build.Builder{
// 		Name:       "builder1",
// 		Driver:     "ansible-playbook",
// 		VarMapping: varsmap.New(),
// 	}
// 	b2 := &build.Builder{
// 		Name:       "builder2",
// 		Driver:     "docker",
// 		VarMapping: varsmap.New(),
// 	}

// 	builders := &Builders{
// 		Builders: map[string]*build.Builder{
// 			"builder1": b1,
// 			"builder2": b2,
// 		},
// 	}

// 	tests := []struct {
// 		desc     string
// 		builders *Builders
// 		builder  string
// 		res      *build.Builder
// 		err      error
// 	}{
// 		{
// 			desc:     "Testing error getting builder to a nil builders",
// 			builders: nil,
// 			res:      nil,
// 			err:      errors.New("(images::GetBuilder)", "Builders is nil"),
// 		},
// 		{
// 			desc:     "Testing get builder",
// 			builders: builders,
// 			builder:  "builder1",
// 			res:      b1,
// 			err:      &errors.Error{},
// 		},
// 		{
// 			desc:     "Testing get builder an unexisting builder",
// 			builders: builders,
// 			builder:  "unexisting",
// 			res:      nil,
// 			err:      errors.New("(images::GetBuilder)", "Unexisting builder configuration for type 'unexisting'"),
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {

// 			t.Log(test.desc)

// 			res, err := test.builders.GetBuilder(test.builder)
// 			if err != nil && assert.Error(t, err) {
// 				assert.Equal(t, test.err.Error(), err.Error())
// 			} else {
// 				assert.Equal(t, test.res, res, "Unexpected value")
// 			}

// 		})
// 	}

// }

// func TestListBuilders(t *testing.T) {
// 	tests := []struct {
// 		desc     string
// 		builders *Builders
// 		res      [][]string
// 	}{
// 		{
// 			desc:     "Testing to list an empty buildersConf defined",
// 			builders: &Builders{},
// 			res:      [][]string{},
// 		},
// 		{
// 			desc: "Testing to list one builder on buildersConf",
// 			builders: &Builders{
// 				Builders: map[string]*build.Builder{
// 					"one": {
// 						Name:   "builder",
// 						Driver: "driver",
// 						Options: map[string]interface{}{
// 							"option1": "option1",
// 						},
// 					},
// 				},
// 			},
// 			res: [][]string{
// 				{"builder", "driver", "option1=option1"},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Log(test.desc)
// 		res, _ := test.builders.ListBuilders()
// 		assert.Equal(t, test.res, res)
// 	}
// }

// func TestAddBuilder(t *testing.T) {

// 	b1 := &build.Builder{
// 		Name:       "builder1",
// 		Driver:     "ansible-playbook",
// 		VarMapping: varsmap.New(),
// 	}
// 	b2 := &Builder{
// 		Name:       "builder2",
// 		Driver:     "docker",
// 		VarMapping: varsmap.New(),
// 	}

// 	tests := []struct {
// 		desc     string
// 		builders *Builders
// 		builder  *build.Builder
// 		res      *Builders
// 		err      error
// 	}{
// 		{
// 			desc: "Testing add new builder",
// 			builders: &Builders{
// 				Builders: map[string]*build.Builder{
// 					"builder1": b1,
// 				},
// 			},
// 			builder: b2,
// 			res: &Builders{
// 				Builders: map[string]*build.Builder{
// 					"builder1": b1,
// 					"builder2": b2,
// 				},
// 			},
// 			err: &errors.Error{},
// 		},
// 		{
// 			desc: "Testing add an existing builder",
// 			builders: &Builders{
// 				Builders: map[string]*build.Builder{
// 					"builder1": b1,
// 				},
// 			},
// 			builder: b1,
// 			res:     nil,
// 			err:     errors.New("(builder::AddBuilder)", "Builder 'builder1' already exist"),
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			err := test.builders.AddBuilder(test.builder)
// 			if err != nil && assert.Error(t, err) {
// 				assert.Equal(t, test.err.Error(), err.Error())
// 			} else {
// 				assert.Equal(t, test.res, test.builders)
// 			}
// 		})
// 	}
// }

// func TestListBuildersHeader(t *testing.T) {

// 	t.Log("Testing list Builders header")
// 	expected := []string{"BUILDER", "DRIVER", "OPTIONS"}
// 	res := ListBuildersHeader()

// 	assert.Equal(t, expected, res)
// }
