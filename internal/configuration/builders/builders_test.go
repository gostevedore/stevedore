package builders

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

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
		desc              string
		path              string
		builders          *Builders
		prepareAssertFunc func(*Builders)
		err               error
	}{
		{
			desc:     "Testing loading builders from file",
			path:     "/builders/file1.yml",
			err:      &errors.Error{},
			builders: NewBuilders(testFs, builders.NewMockBuildersStore()),
			prepareAssertFunc: func(b *Builders) {
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "first",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/context"}},
						},
					},
				).Return(nil)
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "second",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/another/context"}},
						},
					},
				).Return(nil)
			},
		},
		{
			desc:     "Testing error loading builders from file",
			path:     "/builders/tab_error_file.yml",
			builders: NewBuilders(testFs, builders.NewMockBuildersStore()),
			err:      errors.New(errContext, "Error loading builders from file '/builders/tab_error_file.yml'\nfound:\n\nbuilders:\n  first:\n    driver: docker\n    options:\n\t  dockerfile: Dockerfile.test\n      context:\n      path: /path/to/context\n\n\tyaml: line 6: found character that cannot start any token"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.builders)
			}

			err := test.builders.LoadBuildersFromFile(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.builders.store.(*builders.MockBuildersStore).AssertExpectations(t)
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
		desc              string
		path              string
		builders          *Builders
		prepareAssertFunc func(*Builders)
		err               error
	}{
		{
			desc:     "Testing loading builders from directory",
			path:     "/builders",
			err:      &errors.Error{},
			builders: NewBuilders(testFs, builders.NewMockBuildersStore()),
			prepareAssertFunc: func(b *Builders) {
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "first",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/context"}},
						},
					},
				).Return(nil)
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "second",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/another/context"}},
						},
					},
				).Return(nil)
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "third",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/even/another/path/to/context"}},
						},
					},
				).Return(nil)
			},
		},
		{
			desc:     "Testing error loading builders from directory",
			path:     "/builders_error",
			err:      errors.New(errContext, "Error loading builders from file '/builders_error/file1.yml'\nfound:\n\nbuilders:\n  third:\n    driver: docker\n    options:\n\t  dockerfile: Dockerfile.test\n      context:\n        - path: /even/another/path/to/context\n\n\tyaml: line 6: found character that cannot start any token\nError loading builders from file '/builders_error/file2.yml'\nfound:\n\nbuilders:\n  fourth:\n    driver: docker\n    options:\n\t  dockerfile: Dockerfile.test\n      context:\n        - path: /even/another/path/to/context\n\n\tyaml: line 6: found character that cannot start any token\n"),
			builders: NewBuilders(testFs, builders.NewMockBuildersStore()),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.builders)
			}

			err := test.builders.LoadBuildersFromDir(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.builders.store.(*builders.MockBuildersStore).AssertExpectations(t)
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
		desc              string
		path              string
		builders          *Builders
		prepareAssertFunc func(*Builders)
		err               error
	}{
		{
			desc:     "Testing loading builders from file",
			path:     "/builders/file1.yml",
			builders: NewBuilders(testFs, builders.NewMockBuildersStore()),
			prepareAssertFunc: func(b *Builders) {
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "first",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/context"}},
						},
					},
				).Return(nil)
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "second",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/another/context"}},
						},
					},
				).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc:     "Testing loading builders from directory",
			path:     "/builders",
			builders: NewBuilders(testFs, builders.NewMockBuildersStore()),
			prepareAssertFunc: func(b *Builders) {
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "first",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/context"}},
						},
					},
				).Return(nil)
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "second",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/path/to/another/context"}},
						},
					},
				).Return(nil)
				b.store.(*builders.MockBuildersStore).On("Store",
					&builder.Builder{
						Name:   "third",
						Driver: "docker",
						Options: &builder.BuilderOptions{
							Dockerfile: "Dockerfile.test",
							Context:    []*builder.DockerDriverContextOptions{{Path: "/even/another/path/to/context"}},
						},
					},
				).Return(nil)
			},
			err: &errors.Error{},
		},

		{
			desc:     "Testing error loading builders from unexisting directory",
			path:     "/builders_unexisting",
			builders: NewBuilders(testFs, builders.NewMockBuildersStore()),
			err:      errors.New(errContext, "open /builders_unexisting: file does not exist"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.builders)
			}

			err := test.builders.LoadBuilders(test.path)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.builders.store.(*builders.MockBuildersStore).AssertExpectations(t)
			}
		})
	}
}
