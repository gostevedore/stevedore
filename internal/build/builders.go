package build

import (
	"fmt"

	data "github.com/apenella/go-common-utils/data"
	errors "github.com/apenella/go-common-utils/error"
	defaultbuilder "github.com/gostevedore/stevedore/internal/driver/default"
)

// Builders serializes the builders defined on user configuration
type Builders struct {
	Builders map[string]*Builder `yaml:"builders"`
}

// LoadImagesTree
func LoadBuilders(file string) (*Builders, error) {

	builders := &Builders{}
	err := data.LoadYAMLFile(file, builders)
	if err != nil {
		return nil, errors.New("(builder::LoadBuilders)", "Could not be load configuration builders file", err)
	}

	if builders.Builders == nil {
		builders.Builders = map[string]*Builder{}
	}

	err = builders.sanetizeBuilders()
	if err != nil {
		return nil, errors.New("(builder::LoadBuilders)", "Builders configuration could not be prepared", err)
	}

	return builders, nil
}

// AddBuilder include a new builder to builders
func (b *Builders) AddBuilder(builder *Builder) error {
	_, exist := b.Builders[builder.Name]
	if exist {
		return errors.New("(builder::AddBuilder)", fmt.Sprintf("Builder '%s' already exist", builder.Name))
	}

	b.Builders[builder.Name] = builder

	return nil
}

// sanetizeBuilders ensures that all builders has been created with all required parameters
func (b *Builders) sanetizeBuilders() error {

	var err error

	if b == nil {
		return errors.New("(builder::sanetizeBuilders)", "Builders configuration is nil")
	}

	for builderName, builder := range b.Builders {
		if builder == nil {
			builder, err = NewBuilder(builderName, defaultbuilder.DriverName, nil, nil)
			if err != nil {
				return errors.New("(builder::sanetizeBuilders)", fmt.Sprintf("Builder '%s' could not created", builderName), err)
			}

			b.Builders[builderName] = builder
		}
		err = builder.SanetizeBuilder(builderName)
		if err != nil {
			return errors.New("(builder::sanetizeBuilders)", fmt.Sprintf("Error sanitizing builder '%s'", builderName), err)
		}

	}

	return nil
}

// GetBuilder returns the builder registered with input name
func (c *Builders) GetBuilder(name string) (*Builder, error) {
	if c == nil {
		return nil, errors.New("(images::GetBuilder)", "Builders is nil")
	}

	builder, exists := c.Builders[name]
	if !exists {
		return nil, errors.New("(images::GetBuilder)", "Unexisting builder configuration for type '"+name+"'")
	}

	return builder, nil
}

// ListBuilders
func (c *Builders) ListBuilders() ([][]string, error) {
	builders := [][]string{}

	for _, builder := range c.Builders {

		b, err := builder.ToArray()
		if err != nil {
			return nil, errors.New("(images::ListBuilders)", "Builders could not be listed", err)
		}
		builders = append(builders, b)
	}

	return builders, nil
}

// ListBuildersHeader
func ListBuildersHeader() []string {
	h := []string{
		"BUILDER",
		"DRIVER",
		"OPTIONS",
	}

	return h
}
