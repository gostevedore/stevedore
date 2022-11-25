package builder

import (
	"bytes"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/varsmap"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	tests := []struct {
		desc   string
		name   string
		driver string
		res    *Builder
	}{
		{
			desc:   "Testing create a new builder",
			name:   "builder",
			driver: "ansible-playbook",
			res: &Builder{
				Name:       "builder",
				Driver:     "ansible-playbook",
				Options:    &BuilderOptions{},
				VarMapping: varsmap.New(),
			},
		},
	}
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			builder := NewBuilder(test.name, test.driver, nil, nil)
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
				VarMapping: varsmap.New(),
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
				VarMapping: varsmap.New(),
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

	errContext := "(core::domain::builder::NewBuilderFromIOReader)"

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
				VarMapping: varsmap.New(),
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
				VarMapping: varsmap.New(),
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
