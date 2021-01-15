package build

import (
	"fmt"
	"sort"

	"github.com/gostevedore/stevedore/internal/build/varsmap"
	defaultbuilder "github.com/gostevedore/stevedore/internal/driver/default"

	errors "github.com/apenella/go-common-utils/error"
)

const (
	arrayOptionAssignment = "="
)

// Builder serializes each builder defined on user configuration
type Builder struct {
	Name       string                 `yaml:"name"`
	Driver     string                 `yaml:"driver"`
	Options    map[string]interface{} `yaml:"options"`
	VarMapping varsmap.Varsmap        `yaml:"variables_mapping"`
}

// New
func NewBuilder(name, driver string, options map[string]interface{}, vmap varsmap.Varsmap) (*Builder, error) {

	if name == "" {
		return nil, errors.New("(builder::NewBuilder", "Name must be provided to create a builder")
	}

	if driver == "" {
		return nil, errors.New("(builder::NewBuilder", "Driver must be provided to create a builder")
	}

	if options == nil {
		options = map[string]interface{}{}
	}

	if vmap == nil {
		vmap = varsmap.New()
	} else {
		vmap.Combine(varsmap.New())
	}

	b := &Builder{
		Name:       name,
		Driver:     driver,
		Options:    options,
		VarMapping: vmap,
	}
	return b, nil
}

// SanetizeBuilder ensures that a builders has been created with all required parameters
func (b *Builder) SanetizeBuilder(name string) error {

	if b == nil {
		return errors.New("(builder::SanetizeBuilder)", "Builder is nil")
	}
	if len(b.Name) <= 0 {
		b.Name = name
	}

	if len(b.Driver) <= 0 {
		b.Driver = defaultbuilder.DriverName
	}

	if b.VarMapping == nil {
		b.VarMapping = varsmap.New()
	} else {
		b.VarMapping.Combine(varsmap.New())
	}

	return nil
}

// ToArray
func (b *Builder) ToArray() ([]string, error) {
	arrayBuilder := []string{}
	arrayBuilder = append(arrayBuilder, b.Name)
	arrayBuilder = append(arrayBuilder, b.Driver)
	arrayBuilder = append(arrayBuilder, b.listArrayOptions()...)

	return arrayBuilder, nil
}

// listArrayOptions
func (b *Builder) listArrayOptions() []string {
	options := []string{}
	for option, value := range b.Options {
		switch value.(type) {

		case string:
			options = append(options, fmt.Sprintf("%s%s%s", option, arrayOptionAssignment, value.(string)))
		case int:
			options = append(options, fmt.Sprintf("%s%s%v", option, arrayOptionAssignment, fmt.Sprintf("%v", value.(int))))
		default:
			options = append(options, fmt.Sprintf("%s%s%v", option, arrayOptionAssignment, fmt.Sprintf("%v", value)))
		}
	}

	sort.Strings(options)
	return options
}

// // InitVariablesMapping update builder
// func (b *Builder) initVariablesMapping() {
// 	var exists bool

// 	if b.VarMapping == nil {
// 		b.VarMapping = map[string]string{}
// 	}

// 	_, exists = b.VarMapping[VarMappingImageBuilderNameKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageBuilderNameKey] = VarMappingImageBuilderNameDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageBuilderTagKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageBuilderTagKey] = VarMappingImageBuilderTagDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageBuilderRegistryNamespaceKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageBuilderRegistryNamespaceKey] = VarMappingImageBuilderRegistryNamespaceDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageBuilderRegistryHostKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageBuilderRegistryHostKey] = VarMappingImageBuilderRegistryHostDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageBuilderLabelKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageBuilderLabelKey] = VarMappingImageBuilderLabelDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageFromNameKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageFromNameKey] = VarMappingImageFromNameDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageFromTagKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageFromTagKey] = VarMappingImageFromTagDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageFromRegistryNamespaceKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageFromRegistryNamespaceKey] = VarMappingImageFromRegistryNamespaceDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageFromRegistryHostKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageFromRegistryHostKey] = VarMappingImageFromRegistryHostDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageNameKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageNameKey] = VarMappingImageNameDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingImageTagKey]
// 	if !exists {
// 		b.VarMapping[VarMappingImageTagKey] = VarMappingImageTagDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingRegistryNamespaceKey]
// 	if !exists {
// 		b.VarMapping[VarMappingRegistryNamespaceKey] = VarMappingRegistryNamespaceDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingRegistryHostKey]
// 	if !exists {
// 		b.VarMapping[VarMappingRegistryHostKey] = VarMappingRegistryHostDefaultValue
// 	}

// 	_, exists = b.VarMapping[VarMappingPushImagetKey]
// 	if !exists {
// 		b.VarMapping[VarMappingPushImagetKey] = VarMappingPushImagetDefaultValue
// 	}
// }
