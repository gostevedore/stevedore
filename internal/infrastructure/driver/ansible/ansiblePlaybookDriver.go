package ansible

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/domain/varsmap"
)

const (
	// DriverName is the name for the driver
	DriverName = "ansible-playbook"
)

// AnsiblePlaybookDriver drives the build through ansible
type AnsiblePlaybookDriver struct {
	driver AnsibleDriverer
	writer io.Writer
}

// NewAnsiblePlaybookDriver returns an AnsiblePlaybookDriver. In case driver is null, it returns an error
func NewAnsiblePlaybookDriver(driver AnsibleDriverer, writer io.Writer) (*AnsiblePlaybookDriver, error) {

	errContext := "(ansibledriver::NewAnsiblePlaybookDriver)"

	if driver == nil {
		return nil, errors.New(errContext, "To create an AnsiblePlaybookDriver is required a driver")
	}

	if writer == nil {
		writer = os.Stdout
	}

	return &AnsiblePlaybookDriver{
		driver: driver,
		writer: writer,
	}, nil
}

// Build performs the build. In case the build could not performed it returns an error
func (d *AnsiblePlaybookDriver) Build(ctx context.Context, i *image.Image, o *image.BuildDriverOptions) error {

	builderName := "builder"
	errContext := "(ansibledriver::Build)"

	if d.driver == nil {
		return errors.New(errContext, "To build an image is required a driver")
	}

	if i == nil {
		return errors.New(errContext, "To build an image is required a image")
	}

	if o == nil {
		return errors.New(errContext, "To build an image is required a build options")
	}

	if ctx == nil {
		return errors.New(errContext, "To build an image is required a golang context")
	}

	if o.BuilderOptions == nil {
		return errors.New(errContext, "To build an image are required the options from the builder")
	}

	playbook := o.BuilderOptions.Playbook
	if playbook == "" {
		return errors.New(errContext, "Playbook has not been defined on build options")
	}

	inventory := o.BuilderOptions.Inventory
	if o.AnsibleInventoryPath != "" {
		inventory = o.AnsibleInventoryPath
	}

	if inventory == "" {
		return errors.New(errContext, "Inventory has not been defined on build options")
	}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		Inventory: inventory,
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{}
	if o.AnsibleConnectionLocal {
		ansiblePlaybookConnectionOptions.Connection = "local"
	}

	if i.Name == "" {
		return errors.New(errContext, "Image name is not defined")
	}

	ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageNameKey], i.Name)

	if i.RegistryNamespace != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingRegistryNamespaceKey], i.RegistryNamespace)
		builderName = strings.Join([]string{builderName, i.RegistryNamespace}, "_")
	}

	builderName = strings.Join([]string{builderName, i.Name}, "_")
	if i.Version != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageTagKey], i.Version)
		builderName = strings.Join([]string{builderName, i.Version}, "_")
	}

	if len(o.AnsibleIntermediateContainerName) > 0 {
		builderName = o.AnsibleIntermediateContainerName
	}

	if o.AnsibleLimit != "" {
		ansiblePlaybookOptions.Limit = o.AnsibleLimit
	}

	ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageBuilderLabelKey], builderName)

	if i.RegistryHost != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingRegistryHostKey], i.RegistryHost)
	}

	// Persistent vars contains the variables defined by the user on execution time and has precedences over vars and the persistent vars defined on the image
	if len(i.PersistentVars) > 0 {
		for varName, varValue := range i.PersistentVars {
			ansiblePlaybookOptions.AddExtraVar(varName, varValue)
		}
	}

	// Vars contains the variables defined by the user on execution time and has precedences over the default values
	if len(i.Vars) > 0 {
		for varName, varValue := range i.Vars {
			ansiblePlaybookOptions.AddExtraVar(varName, varValue)
		}
	}

	if len(i.Tags) > 0 {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageExtraTagsKey], i.Tags)
	}

	// Persistent labels contains the variables defined by the user on execution time and has precedences over labels and the persistent vars defined on the image
	if len(i.PersistentLabels) > 0 {
		for varName, varValue := range i.PersistentLabels {
			ansiblePlaybookOptions.AddExtraVar(varName, varValue)
		}
	}

	if len(i.Labels) > 0 {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageExtraTagsKey], i.Labels)
	}

	if o.OutputPrefix == "" {
		o.OutputPrefix = i.Name
		if i.Version != "" {
			o.OutputPrefix = strings.Join([]string{o.OutputPrefix, i.Version}, ":")
		}
	}

	if i.Parent != nil && i.Parent.Name != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageFromNameKey], i.Parent.Name)

	}

	if i.Parent != nil && i.Parent.RegistryNamespace != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryNamespaceKey], i.Parent.RegistryNamespace)
	}

	if i.Parent != nil && i.Parent.RegistryHost != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryHostKey], i.Parent.RegistryHost)
	}

	if i.Parent != nil && i.Parent.Version != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageFromTagKey], i.Parent.Version)
	}

	if !o.PushImageAfterBuild {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingPushImagetKey], false)
	}

	// TODO:
	// go-ansible library is not able to pass secrets, auth values won't be passed to ansible playbook till this is allowed

	d.driver.WithPlaybook(playbook)
	d.driver.WithOptions(ansiblePlaybookOptions)
	d.driver.WithConnectionOptions(ansiblePlaybookConnectionOptions)
	d.driver.PrepareExecutor(d.writer, o.OutputPrefix)

	err := d.driver.Run(ctx)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}
