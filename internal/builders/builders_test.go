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

			res, err := test.builders.Find(test.builder)
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
