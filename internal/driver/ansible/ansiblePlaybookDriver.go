package driver

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/varsmap"
	"github.com/gostevedore/stevedore/internal/driver"
	"github.com/gostevedore/stevedore/internal/images/image"
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
func (d *AnsiblePlaybookDriver) Build(ctx context.Context, i *image.Image, o *driver.BuildDriverOptions) error {

	builderName := "builder"
	errContext := "(ansibledriver::Build)"

	if d.driver == nil {
		return errors.New(errContext, "To build an image is required a driver")
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
	if inventory == "" {
		return errors.New(errContext, "Inventory has not been defined on build options")
	}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		Inventory: inventory,
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{}
	if o.ConnectionLocal {
		ansiblePlaybookConnectionOptions.Connection = "local"
	}

	if o.ImageName == "" {
		return errors.New(errContext, "Image has not been defined on build options")
	}

	ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageNameKey], o.ImageName)

	if o.RegistryNamespace != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingRegistryNamespaceKey], o.RegistryNamespace)
		builderName = strings.Join([]string{builderName, o.RegistryNamespace}, "_")
	}

	builderName = strings.Join([]string{builderName, o.ImageName}, "_")

	if o.ImageVersion != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageTagKey], o.ImageVersion)
		builderName = strings.Join([]string{builderName, o.ImageVersion}, "_")
	}

	if len(o.BuilderName) > 0 {
		builderName = o.BuilderName
	}
	ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageBuilderLabelKey], builderName)

	if o.RegistryHost != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingRegistryHostKey], o.RegistryHost)
	}

	// Persistent vars contains the variables defined by the user on execution time and has precedences over vars and the persistent vars defined on the image
	if len(o.PersistentVars) > 0 {
		for varName, varValue := range o.PersistentVars {
			ansiblePlaybookOptions.AddExtraVar(varName, varValue)
		}
	}

	// Vars contains the variables defined by the user on execution time and has precedences over the default values
	if len(o.Vars) > 0 {
		for varName, varValue := range o.Vars {
			ansiblePlaybookOptions.AddExtraVar(varName, varValue)
		}
	}

	if len(o.Tags) > 0 {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageExtraTagsKey], o.Tags)
	}

	if len(o.Labels) > 0 {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageExtraTagsKey], o.Labels)
	}

	if o.OutputPrefix == "" {
		o.OutputPrefix = o.ImageName
		if o.ImageVersion != "" {
			o.OutputPrefix = strings.Join([]string{o.OutputPrefix, o.ImageVersion}, ":")
		}
	}

	if o.ImageFromName != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageFromNameKey], o.ImageFromName)
	}

	if o.ImageFromRegistryNamespace != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryNamespaceKey], o.ImageFromRegistryNamespace)
	}

	if o.ImageFromRegistryHost != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageFromRegistryHostKey], o.ImageFromRegistryHost)
	}

	if o.ImageFromVersion != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageFromTagKey], o.ImageFromVersion)
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
		return errors.New(errContext, err.Error())
	}

	return nil
}
