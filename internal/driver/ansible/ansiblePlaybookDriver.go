package ansibledriver

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/build/varsmap"
	"github.com/gostevedore/stevedore/internal/types"
)

const (
	// DriverName is the name for the driver
	DriverName = "ansible-playbook"
	// BuilderConfOptionsPlaybookKey string that defines the playbook on builder options defintions
	BuilderConfOptionsPlaybookKey = "playbook"
	// BuilderConfOptionsInventoryKey string that defines the invetory on builder options defintions
	BuilderConfOptionsInventoryKey = "inventory"
)

// AnsiblePlaybookDriver drives the build through ansible
type AnsiblePlaybookDriver struct {
	driver Ansibler
	writer io.Writer
}

// NewAnsiblePlaybookDriver returns an AnsiblePlaybookDriver. In case driver is null, it returns an error
func NewAnsiblePlaybookDriver(driver Ansibler, writer io.Writer) (*AnsiblePlaybookDriver, error) {

	errContext := "(ansibledriver::NewAnsiblePlaybookDriver)"

	if driver == nil {
		return nil, errors.New(errContext, "To create an AnsiblePlaybookDriver is expected a driver")
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
func (d *AnsiblePlaybookDriver) Build(ctx context.Context, o *types.BuildOptions) error {

	builderName := "builder"
	errContext := "(ansibledriver::Build)"

	if d.driver == nil {
		return errors.New(errContext, "Build driver is missing")
	}

	if o == nil {
		return errors.New(errContext, "Build options are nil")
	}

	if ctx == nil {
		return errors.New(errContext, "Context is nil")
	}

	builderConfOptions := o.BuilderOptions

	playbook, ok := builderConfOptions[BuilderConfOptionsPlaybookKey]
	if !ok {
		return errors.New(errContext, "Playbook has not been defined on build options")
	}

	inventory, ok := builderConfOptions[BuilderConfOptionsInventoryKey]
	if !ok {
		return errors.New(errContext, "Inventory has not been defined on build options")
	}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		Inventory: inventory.(string),
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
		// Removed: stevedore does not sanitize image version
		// o.ImageVersion = common.SanitizeTag(o.ImageVersion)
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
		// Removed: stevedore does not sanitize image version
		// sanitizedTags := []string{}
		// for _, tag := range o.Tags {
		// 	tag = drivercommon.SanitizeTag(tag)
		// 	sanitizedTags = append(sanitizedTags, tag)
		// }
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageExtraTagsKey], o.Tags)
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

	if !o.PushImages {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingPushImagetKey], false)
	}

	d.driver.WithPlaybook(playbook.(string))
	d.driver.WithOptions(ansiblePlaybookOptions)
	d.driver.WithConnectionOptions(ansiblePlaybookConnectionOptions)
	d.driver.PrepareExecutor(d.writer, o.OutputPrefix)

	// ansiblePlaybook := &ansible.AnsiblePlaybookCmd{
	// 	Playbooks:         []string{playbook.(string)},
	// 	Options:           ansiblePlaybookOptions,
	// 	ConnectionOptions: ansiblePlaybookConnectionOptions,
	// 	Exec: execute.NewDefaultExecute(
	// 		execute.WithWrite(console.GetConsole()),
	// 		execute.WithTransformers(
	// 			results.Prepend(o.OutputPrefix),
	// 		),
	// 	),
	// }

	err := d.driver.Run(ctx)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	return nil
}
