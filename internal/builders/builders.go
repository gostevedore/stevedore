package builders

import (
	"fmt"
	"os"
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
	Builders map[string]*builder.Builder `yaml:"builders"`
}

// NewBuilders creates a new builders configuration
func NewBuilders(fs afero.Fs) *Builders {
	return &Builders{
		fs: fs,

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

// AddBuilder include a new builder to builders
func (b *Builders) AddBuilder(builder *builder.Builder) error {

	errContext := "(builders::AddBuilder)"

	if b == nil {
		return errors.New(errContext, "Builders is nil")
	}

	b.mutex.Lock()
	_, exist := b.Builders[builder.Name]
	if exist {
		b.mutex.Unlock()
		return errors.New(errContext, fmt.Sprintf("Builder '%s' already exist", builder.Name))
	}

	b.Builders[builder.Name] = builder
	b.mutex.Unlock()

	return nil
}

// GetBuilder returns the builder registered with input name
func (b *Builders) GetBuilder(name string) (*builder.Builder, error) {

	errContext := "(builders::GetBuilder)"

	if b == nil {
		return nil, errors.New(errContext, "Builders is nil")
	}

	b.mutex.RLock()
	builder, exists := b.Builders[name]
	if !exists {
		return nil, errors.New(errContext, fmt.Sprintf("Builder '%s' does not exists", name))
	}
	b.mutex.RUnlock()

	return builder, nil
}

// LoadBuildersFromFile loads builders from file
func (b *Builders) LoadBuildersFromFile(path string) error {

	var err error
	var fileData []byte

	errContext := "(builders::loadBuilderFile)"

	if b == nil {
		return errors.New(errContext, "Builders is nil")
	}

	buildersAux := NewBuilders(b.fs)

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

		err = b.AddBuilder(builder)
		if err != nil {
			return errors.New(errContext, fmt.Sprintf("Error loading builders from file '%s'", path), err)
		}
	}

	return nil
}

// LoadBuildersFromDir loads builders from all files on directory
func (b *Builders) LoadBuildersFromDir(path string) error {
	var err error
	errFuncs := []func() error{}
	errContext := "(builders::loadBuildersFromDir)"

	if b == nil {
		return errors.New(errContext, "Builders is nil")
	}

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

	err = afero.Walk(b.fs, path, func(path string, info os.FileInfo, err error) error {
		var isDir bool

		isDir, err = afero.IsDir(b.fs, path)
		if err != nil {
			return errors.New(errContext, err.Error())
		}

		if isDir {
			return nil
		}

		b.wg.Add(1)
		f := loadBuildersFromFile(path)
		errFuncs = append(errFuncs, f)

		return nil
	})

	if err != nil {
		return errors.New(errContext, err.Error())
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

// // ListBuilders
// func (c *Builders) ListBuilders() ([][]string, error) {
// 	builders := [][]string{}

// 	for _, builder := range c.Builders {

// 		b, err := builder.ToArray()
// 		if err != nil {
// 			return nil, errors.New("(images::ListBuilders)", "Builders could not be listed", err)
// 		}
// 		builders = append(builders, b)
// 	}

// 	return builders, nil
// }

// // ListBuildersHeader
// func ListBuildersHeader() []string {
// 	h := []string{
// 		"BUILDER",
// 		"DRIVER",
// 		"OPTIONS",
// 	}

// 	return h
// }
