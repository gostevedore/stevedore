package builders

import (
	"fmt"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/builder"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Builders list of builders to create images
type Builders struct {
	fs       afero.Fs
	mutex    sync.RWMutex
	wg       sync.WaitGroup
	store    BuildersStorer
	Builders map[string]*builder.Builder `yaml:"builders"`
}

// NewBuilders creates a new builders configuration
func NewBuilders(fs afero.Fs, store BuildersStorer) *Builders {
	return &Builders{
		fs:       fs,
		store:    store,
		Builders: make(map[string]*builder.Builder),
	}
}

// LoadBuilders loads builders from path to Builders map
func (b *Builders) LoadBuilders(path string) error {
	var err error
	var isDir bool

	errContext := "(builders::LoadBuilders)"

	isDir, err = afero.IsDir(b.fs, path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	if isDir {
		return b.LoadBuildersFromDir(path)
	} else {
		return b.LoadBuildersFromFile(path)
	}

}

// LoadBuildersFromFile loads builders from file
func (b *Builders) LoadBuildersFromFile(path string) error {

	var err error
	var fileData []byte

	errContext := "(builders::LoadBuildersFromFile)"

	if b == nil {
		return errors.New(errContext, "Builders is nil")
	}

	buildersAux := NewBuilders(b.fs, b.store)

	fileData, err = afero.ReadFile(b.fs, path)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	err = yaml.Unmarshal(fileData, buildersAux)
	if err != nil {
		return errors.New(errContext, fmt.Sprintf("Error loading builders from file '%s'\nfound:\n%s", path, string(fileData)), err)
	}

	for name, builder := range buildersAux.Builders {
		if builder.Name == "" {
			builder.Name = name
		}

		err = b.store.Store(builder)
		if err != nil {
			return errors.New(errContext, fmt.Sprintf("Error loading builders from file '%s'", path), err)
		}
	}

	return nil
}

// LoadBuildersFromDir loads builders from all files on directory
func (b *Builders) LoadBuildersFromDir(dir string) error {
	var err error
	errFuncs := []func() error{}
	errContext := "(builders::loadBuildersFromDir)"

	if b == nil {
		return errors.New(errContext, "Builders is nil")
	}

	yamlFiles, err := afero.Glob(b.fs, dir+"/*.yaml")
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	ymlFiles, err := afero.Glob(b.fs, dir+"/*.yml")
	if err != nil {
		return errors.New(errContext, err.Error())
	}
	files := append(yamlFiles, ymlFiles...)

	// promise function to load builders from file
	loadBuildersFromFile := func(path string) func() error {
		var err error

		c := make(chan struct{}, 1)
		go func() {
			defer close(c)
			err = b.LoadBuildersFromFile(path)
			b.wg.Done()
		}()

		return func() error {
			<-c
			return err
		}
	}

	for _, file := range files {
		b.wg.Add(1)
		f := loadBuildersFromFile(file)
		errFuncs = append(errFuncs, f)
	}
	b.wg.Wait()

	errMsg := ""
	for _, f := range errFuncs {
		err = f()
		if err != nil {
			errMsg = fmt.Sprintf("%s%s\n", errMsg, err.Error())
		}
	}
	if errMsg != "" {
		return errors.New(errContext, errMsg)
	}

	return nil
}
