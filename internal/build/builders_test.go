package build

import (
	"path/filepath"
	"stevedore/internal/build/varsmap"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

// TestLoadImage tests
func TestLoadBuilders(t *testing.T) {

	testBaseDir := "test"

	tests := []struct {
		desc     string
		file     string
		err      error
		builders *Builders
	}{
		{
			desc:     "testing an unexistent file",
			file:     "nofile",
			err:      errors.New("(builder::LoadBuilders)", "Could not be load configuration builders file", errors.New("", "(LoadYAMLFile) Error loading file nofile. open nofile: no such file or directory")),
			builders: &Builders{},
		},
		{
			desc:     "testing nil drivers",
			file:     filepath.Join(testBaseDir, "nil_builders.yml"),
			err:      &errors.Error{},
			builders: &Builders{map[string]*Builder{}},
		},
		{
			desc: "testing builders",
			file: filepath.Join(testBaseDir, "builders.yml"),
			err:  &errors.Error{},
			builders: &Builders{
				Builders: map[string]*Builder{
					"infrastructure": {
						Name:   "infrastructure",
						Driver: "ansible-playbook",
						Options: map[string]interface{}{
							"inventory": "inventory/all",
							"playbook":  "site.yml",
						},
						VarMapping: varsmap.Varsmap{
							varsmap.VarMappingImageBuilderNameKey:              varsmap.VarMappingImageBuilderNameDefaultValue,
							varsmap.VarMappingImageBuilderTagKey:               varsmap.VarMappingImageBuilderTagDefaultValue,
							varsmap.VarMappingImageBuilderRegistryNamespaceKey: varsmap.VarMappingImageBuilderRegistryNamespaceDefaultValue,
							varsmap.VarMappingImageBuilderRegistryHostKey:      varsmap.VarMappingImageBuilderRegistryHostDefaultValue,
							varsmap.VarMappingImageBuilderLabelKey:             varsmap.VarMappingImageBuilderLabelDefaultValue,
							varsmap.VarMappingImageFromNameKey:                 varsmap.VarMappingImageFromNameDefaultValue,
							varsmap.VarMappingImageFromTagKey:                  varsmap.VarMappingImageFromTagDefaultValue,
							varsmap.VarMappingImageFromRegistryNamespaceKey:    varsmap.VarMappingImageFromRegistryNamespaceDefaultValue,
							varsmap.VarMappingImageFromRegistryHostKey:         varsmap.VarMappingImageFromRegistryHostDefaultValue,
							varsmap.VarMappingImageNameKey:                     "image",
							varsmap.VarMappingImageTagKey:                      varsmap.VarMappingImageTagDefaultValue,
							varsmap.VarMappingImageExtraTagsKey:                varsmap.VarMappingImageExtraTagsDefaultValue,
							varsmap.VarMappingRegistryNamespaceKey:             varsmap.VarMappingRegistryNamespaceDefaultValue,
							varsmap.VarMappingRegistryHostKey:                  varsmap.VarMappingRegistryHostDefaultValue,
							varsmap.VarMappingPushImagetKey:                    varsmap.VarMappingPushImagetDefaultValue,
						},
					},
					"php-code": {
						Name:   "php-code",
						Driver: "ansible-playbook",
						Options: map[string]interface{}{
							"inventory": "inventory/all",
							"playbook":  "code_builder.yml",
						},
						VarMapping: varsmap.New(),
					},
					"dummy": {
						Name:       "dummy",
						Driver:     "default",
						Options:    map[string]interface{}{},
						VarMapping: varsmap.New(),
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		builders, err := LoadBuilders(test.file)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {
			assert.Equal(t, builders, test.builders, "Unexpected value")
		}
	}
}

func TestGetBuilders(t *testing.T) {

	b1 := &Builder{
		Name:       "builder1",
		Driver:     "ansible-playbook",
		VarMapping: varsmap.New(),
	}
	b2 := &Builder{
		Name:       "builder2",
		Driver:     "docker",
		VarMapping: varsmap.New(),
	}

	builders := &Builders{
		Builders: map[string]*Builder{
			"builder1": b1,
			"builder2": b2,
		},
	}

	tests := []struct {
		desc     string
		builders *Builders
		builder  string
		res      *Builder
		err      error
	}{
		{
			desc:     "Testing error getting builder to a nil builders",
			builders: nil,
			res:      nil,
			err:      errors.New("(images::GetBuilder)", "Builders is nil"),
		},
		{
			desc:     "Testing get builder",
			builders: builders,
			builder:  "builder1",
			res:      b1,
			err:      &errors.Error{},
		},
		{
			desc:     "Testing get builder an unexisting builder",
			builders: builders,
			builder:  "unexisting",
			res:      nil,
			err:      errors.New("(images::GetBuilder)", "Unexisting builder configuration for type 'unexisting'"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			t.Log(test.desc)

			res, err := test.builders.GetBuilder(test.builder)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res, "Unexpected value")
			}

		})
	}

}

func TestListBuilders(t *testing.T) {
	tests := []struct {
		desc     string
		builders *Builders
		res      [][]string
	}{
		{
			desc:     "Testing to list an empty buildersConf defined",
			builders: &Builders{},
			res:      [][]string{},
		},
		{
			desc: "Testing to list one builder on buildersConf",
			builders: &Builders{
				Builders: map[string]*Builder{
					"one": {
						Name:   "builder",
						Driver: "driver",
						Options: map[string]interface{}{
							"option1": "option1",
						},
					},
				},
			},
			res: [][]string{
				{"builder", "driver", "option1=option1"},
			},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		res, _ := test.builders.ListBuilders()
		assert.Equal(t, test.res, res)
	}
}

func TestAddBuilder(t *testing.T) {

	b1 := &Builder{
		Name:       "builder1",
		Driver:     "ansible-playbook",
		VarMapping: varsmap.New(),
	}
	b2 := &Builder{
		Name:       "builder2",
		Driver:     "docker",
		VarMapping: varsmap.New(),
	}

	tests := []struct {
		desc     string
		builders *Builders
		builder  *Builder
		res      *Builders
		err      error
	}{
		{
			desc: "Testing add new builder",
			builders: &Builders{
				Builders: map[string]*Builder{
					"builder1": b1,
				},
			},
			builder: b2,
			res: &Builders{
				Builders: map[string]*Builder{
					"builder1": b1,
					"builder2": b2,
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing add an existing builder",
			builders: &Builders{
				Builders: map[string]*Builder{
					"builder1": b1,
				},
			},
			builder: b1,
			res:     nil,
			err:     errors.New("(builder::AddBuilder)", "Builder 'builder1' already exist"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.builders.AddBuilder(test.builder)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.builders)
			}
		})
	}
}

func TestListBuildersHeader(t *testing.T) {

	t.Log("Testing list Builders header")
	expected := []string{"BUILDER", "DRIVER", "OPTIONS"}
	res := ListBuildersHeader()

	assert.Equal(t, expected, res)
}
