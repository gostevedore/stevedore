package builder

import (
	"bytes"
	"fmt"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/varsmap"
	"gopkg.in/yaml.v3"
)

const (
	arrayOptionAssignment = "="
	// NameFilterAttribute is the attribute's filter value to filter by name
	NameFilterAttribute = "name"
	// DriverFilterAttribute is the attribute's filter value to filter by driver
	DriverFilterAttribute = "driver"
)

// Builder serializes each builder defined on user configuration
type Builder struct {
	Name       string          `yaml:"name"`
	Driver     string          `yaml:"driver"`
	Options    *BuilderOptions `yaml:"options"`
	VarMapping varsmap.Varsmap `yaml:"variables_mapping"`
}

// NewBuilder creates a new builder
func NewBuilder(name, driver string, options *BuilderOptions, varmap varsmap.Varsmap) *Builder {

	if options == nil {
		options = &BuilderOptions{}
	}

	if varmap != nil {
		// Combine existing values in varmap with those comming from a new varsmap
		varmap.Combine(varsmap.New())
	} else {
		varmap = varsmap.New()
	}

	return &Builder{
		Name:       name,
		Driver:     driver,
		Options:    options,
		VarMapping: varmap,
	}
}

// NewBuilderFromByteArray creates a new builder from a byte array
func NewBuilderFromByteArray(data []byte) (*Builder, error) {
	var builder *Builder

	errContext := "(core::domain::builder::NewBuilderFromByteArray)"

	err := yaml.Unmarshal(data, &builder)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Builder could not be created.\nfound:\n'%s'\n", string(data)), err)
	}

	if builder.VarMapping == nil {
		builder.VarMapping = varsmap.New()
	}

	return builder, nil
}

// NewBuilderFromIOReader creates a new builder from an io reader
func NewBuilderFromIOReader(reader io.Reader) (*Builder, error) {
	var builder *Builder
	var buff bytes.Buffer
	var err error

	errContext := "(core::domain::builder::NewBuilderFromIOReader)"

	_, err = buff.ReadFrom(reader)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	err = yaml.Unmarshal(buff.Bytes(), &builder)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Builder could not be created.\nfound:\n'%s'\n", buff.String()), err)
	}

	if builder.VarMapping == nil {
		builder.VarMapping = varsmap.New()
	}

	return builder, nil
}

// WithName sets the name of the builder
func (b *Builder) WithName(name string) {
	b.Name = name
}

// WithDriver sets the driver of the builder
func (b *Builder) WithDriver(driver string) {
	b.Driver = driver
}

// WithOptions sets the options of the builder
func (b *Builder) WithOptions(options *BuilderOptions) {
	b.Options = options
}

// WithVarMapping sets the variable mapping of the builder
func (b *Builder) WithVarMapping(mapping varsmap.Varsmap) {
	b.VarMapping = mapping
}
