package builder

import (
	"bytes"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/varsmap"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	tests := []struct {
		desc string
		res  *Builder
	}{
		{
			desc: "Testing create a new builder",
			res: &Builder{
				Options:    &BuilderOptions{},
				VarMapping: varsmap.New(),
			},
		},
	}
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			builder := NewBuilder()
			assert.Equal(t, test.res, builder)
		})
	}
}

func TestNewBuilderFromByteArray(t *testing.T) {

	tests := []struct {
		desc string
		data []byte
		res  *Builder
		err  error
	}{
		{
			desc: "Testing create builder from byte array with ansible-playbook data",
			data: []byte(`
driver: ansible-playbook
options:
  playbook: playbook
  inventory: inventory
`),
			res: &Builder{
				Name:   "",
				Driver: "ansible-playbook",
				Options: &BuilderOptions{
					Playbook:  "playbook",
					Inventory: "inventory",
				},
				VarMapping: nil,
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing create builder from byte array with docker data",
			data: []byte(`
driver: docker
options:
  dockerfile: Dockerfile.test
  context:
    - git:
        path: path
        repository: repository
        reference: reference
        auth:
          username: username
          password: password
`),
			res: &Builder{
				Name:   "",
				Driver: "docker",
				Options: &BuilderOptions{
					Context: []*DockerDriverContextOptions{
						{
							Git: &DockerDriverGitContextOptions{
								Path:       "path",
								Repository: "repository",
								Reference:  "reference",
								Auth: &DockerDriverGitContextAuthOptions{
									Username: "username",
									Password: "password",
								},
							},
						},
					},
					Dockerfile: "Dockerfile.test",
				},
				VarMapping: nil,
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			builder, err := NewBuilderFromByteArray(test.data)

			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, builder)
			}
		})
	}

}

func TestNewBuilderFromIOReader(t *testing.T) {

	errContext := "(builder::NewBuilderFromIOReader)"

	tests := []struct {
		desc string
		data string
		res  *Builder
		err  error
	}{
		{
			desc: "Testing create builder from byte array with ansible-playbook data",
			data: `
driver: ansible-playbook
options:
  playbook: playbook
  inventory: inventory
`,
			res: &Builder{
				Name:   "",
				Driver: "ansible-playbook",
				Options: &BuilderOptions{
					Playbook:  "playbook",
					Inventory: "inventory",
				},
				VarMapping: nil,
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing create builder from byte array with docker data",
			data: `
driver: docker
options:
  dockerfile: Dockerfile.test
  context:
    - git:
        path: path
        repository: repository
        reference: reference
        auth:
          username: username
          password: password
`,
			res: &Builder{
				Name:   "",
				Driver: "docker",
				Options: &BuilderOptions{
					Context: []*DockerDriverContextOptions{
						{
							Git: &DockerDriverGitContextOptions{
								Path:       "path",
								Repository: "repository",
								Reference:  "reference",
								Auth: &DockerDriverGitContextAuthOptions{
									Username: "username",
									Password: "password",
								},
							},
						},
					},
					Dockerfile: "Dockerfile.test",
				},
				VarMapping: nil,
			},
			err: &errors.Error{},
		},

		{
			desc: "Testing create builder from byte array with docker data",
			data: `
	driver: docker
`,
			res: &Builder{},
			err: errors.New(errContext, "Builder could not be created.\nfound:\n'\n\tdriver: docker\n'\n\n\tyaml: line 2: found character that cannot start any token"),
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			buff := new(bytes.Buffer)
			buff.WriteString(test.data)

			builder, err := NewBuilderFromIOReader(buff)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, builder)
			}
		})
	}

}

func TestWithName(t *testing.T) {
	builder := &Builder{}
	builder.WithName("name")
	assert.Equal(t, "name", builder.Name)
}

func TestWithDriver(t *testing.T) {
	builder := &Builder{}
	builder.WithDriver("driver")
	assert.Equal(t, "driver", builder.Driver)
}

func TestWithOptions(t *testing.T) {
	builder := &Builder{}
	options := &BuilderOptions{
		Playbook: "playbook",
	}
	builder.WithOptions(options)
	assert.Equal(t, options, builder.Options)
}

func TestWithVarMapping(t *testing.T) {
	builder := &Builder{}
	varsmap := varsmap.New()
	builder.WithVarMapping(varsmap)
	assert.Equal(t, varsmap, builder.VarMapping)
}

// func TestToArray(t *testing.T) {
// 	tests := []struct {
// 		desc    string
// 		builder *Builder
// 		res     []string
// 		err     error
// 	}{
// 		{
// 			desc: "Testing array generation from a builder conf with map",
// 			builder: &Builder{
// 				Name:   "builder",
// 				Driver: "driver",
// 				Options: map[string]interface{}{
// 					"option1": "option1",
// 					"option2": 2,
// 					"option3": map[string]interface{}{
// 						"suboption3.2": 3,
// 						"suboption3.1": 3,
// 					},
// 				},
// 			},
// 			err: nil,
// 			res: []string{"builder", "driver", "option1=option1", "option2=2", "option3=map[suboption3.1:3 suboption3.2:3]"},
// 		},
// 		{
// 			desc: "Testing array generation from a builder conf with array",
// 			builder: &Builder{
// 				Name:   "builder",
// 				Driver: "driver",
// 				Options: map[string]interface{}{
// 					"option1": "option1",
// 					"option2": 2,
// 					"option3": []string{
// 						"suboption3.1", "suboption3.2",
// 					},
// 				},
// 			},
// 			err: nil,
// 			res: []string{"builder", "driver", "option1=option1", "option2=2", "option3=[suboption3.1 suboption3.2]"},
// 		},
// 	}

// 	for _, test := range tests {

// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)
// 			res, err := test.builder.ToArray()
// 			if err != nil && assert.Error(t, err) {
// 				assert.Equal(t, test.err, err)
// 			} else {
// 				sort.Strings(test.res)
// 				sort.Strings(res)

// 				assert.True(t, reflect.DeepEqual(test.res, res), "Unexpected response\n", res, test.res)
// 			}
// 		})
// 	}
// }

// func TestSanitizeBuilder(t *testing.T) {
// 	tests := []struct {
// 		desc    string
// 		name    string
// 		driver  string
// 		builder *Builder
// 		res     *Builder
// 		err     error
// 	}{
// 		{
// 			desc:    "Testing sanetize a nil builder",
// 			name:    "",
// 			builder: nil,
// 			res:     nil,
// 			err:     errors.New("(builder::SanetizeBuilder)", "Builder is nil"),
// 		},
// 		{
// 			desc:    "Testing sanetize builder with no name defined",
// 			name:    "name",
// 			builder: &Builder{},
// 			res: &Builder{
// 				Name:       "name",
// 				Driver:     defaultbuilder.DriverName,
// 				VarMapping: varsmap.New(),
// 			},
// 			err: errors.New("(builder::SanetizeBuilder)", "Builder is nil"),
// 		},
// 	}

// 	for _, test := range tests {

// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)
// 			err := test.builder.SanetizeBuilder(test.name)
// 			if err != nil && assert.Error(t, err) {
// 				assert.Equal(t, test.err, err)
// 			} else {
// 				assert.Equal(t, test.res, test.builder)
// 			}
// 		})
// 	}
// }
