package build

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gostevedore/stevedore/internal/build/varsmap"
	defaultbuilder "github.com/gostevedore/stevedore/internal/driver/default"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	tests := []struct {
		desc    string
		name    string
		driver  string
		options map[string]interface{}
		varsmap map[string]string
		res     *Builder
		err     error
	}{
		{
			desc:   "Testing create a new builder",
			name:   "builder",
			driver: "driver",
			res: &Builder{
				Name:       "builder",
				Driver:     "driver",
				Options:    map[string]interface{}{},
				VarMapping: varsmap.New(),
			},
		},
		{
			desc:   "Testing error creating a new builder with blank name",
			name:   "",
			driver: "driver",
			err:    errors.New("(builder::NewBuilder", "Name must be provided to create a builder"),
		},
		{
			desc:   "Testing error creating a new builder with blank name",
			name:   "name",
			driver: "",
			err:    errors.New("(builder::NewBuilder", "Driver must be provided to create a builder"),
		},
	}
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			builder, err := NewBuilder(test.name, test.driver, test.options, test.varsmap)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, builder)
			}

		})
	}

}

func TestToArray(t *testing.T) {
	tests := []struct {
		desc    string
		builder *Builder
		res     []string
		err     error
	}{
		{
			desc: "Testing array generation from a builder conf with map",
			builder: &Builder{
				Name:   "builder",
				Driver: "driver",
				Options: map[string]interface{}{
					"option1": "option1",
					"option2": 2,
					"option3": map[string]interface{}{
						"suboption3.2": 3,
						"suboption3.1": 3,
					},
				},
			},
			err: nil,
			res: []string{"builder", "driver", "option1=option1", "option2=2", "option3=map[suboption3.1:3 suboption3.2:3]"},
		},
		{
			desc: "Testing array generation from a builder conf with array",
			builder: &Builder{
				Name:   "builder",
				Driver: "driver",
				Options: map[string]interface{}{
					"option1": "option1",
					"option2": 2,
					"option3": []string{
						"suboption3.1", "suboption3.2",
					},
				},
			},
			err: nil,
			res: []string{"builder", "driver", "option1=option1", "option2=2", "option3=[suboption3.1 suboption3.2]"},
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			res, err := test.builder.ToArray()
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				sort.Strings(test.res)
				sort.Strings(res)

				assert.True(t, reflect.DeepEqual(test.res, res), "Unexpected response\n", res, test.res)
			}
		})
	}
}

func TestSanitizeBuilder(t *testing.T) {
	tests := []struct {
		desc    string
		name    string
		driver  string
		builder *Builder
		res     *Builder
		err     error
	}{
		{
			desc:    "Testing sanetize a nil builder",
			name:    "",
			builder: nil,
			res:     nil,
			err:     errors.New("(builder::SanetizeBuilder)", "Builder is nil"),
		},
		{
			desc:    "Testing sanetize builder with no name defined",
			name:    "name",
			builder: &Builder{},
			res: &Builder{
				Name:       "name",
				Driver:     defaultbuilder.DriverName,
				VarMapping: varsmap.New(),
			},
			err: errors.New("(builder::SanetizeBuilder)", "Builder is nil"),
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			err := test.builder.SanetizeBuilder(test.name)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, test.builder)
			}
		})
	}
}

// func TestInitVariablesMapping(t *testing.T) {

// 	t.Log("Testing update var mappings")

// 	builder := &Builder{
// 		Name:   "builder",
// 		Driver: "driver",
// 	}

// 	res := &Builder{
// 		Name:       "builder",
// 		Driver:     "driver",
// 		VarMapping: varsmap.New(),
// 	}

// 	assert.Equal(t, res, builder)
// }
