package ansibledriver

import (
	"context"
	"strings"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/build/varsmap"
	"github.com/gostevedore/stevedore/internal/driver/common"
	drivercommon "github.com/gostevedore/stevedore/internal/driver/common"
	"github.com/gostevedore/stevedore/internal/types"
	"github.com/gostevedore/stevedore/internal/ui/console"
)

const (
	DriverName = "ansible-playbook"
)

func NewAnsiblePlaybookDriver(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {

	builderName := "builder"

	if o == nil {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "Build options are nil")
	}

	if ctx == nil {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "Context is nil")
	}

	builderConfOptions := o.BuilderOptions

	playbook, ok := builderConfOptions["playbook"]
	if !ok {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "playbook has not been defined on build options")
	}

	inventory, ok := builderConfOptions["inventory"]
	if !ok {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "inventory has not been defined on build options")
	}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		Inventory: inventory.(string),
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{}
	if o.ConnectionLocal {
		ansiblePlaybookConnectionOptions.Connection = "local"
	}

	if o.ImageName == "" {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "Image name is not set")
	}
	ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageNameKey], o.ImageName)

	if o.RegistryNamespace != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingRegistryNamespaceKey], o.RegistryNamespace)
		builderName = strings.Join([]string{builderName, o.RegistryNamespace}, "_")
	}

	builderName = strings.Join([]string{builderName, o.ImageName}, "_")

	if o.ImageVersion != "" {
		o.ImageVersion = common.SanitizeTag(o.ImageVersion)
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
		sanitizedTags := []string{}
		for _, tag := range o.Tags {
			tag = drivercommon.SanitizeTag(tag)
			sanitizedTags = append(sanitizedTags, tag)
		}
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageExtraTagsKey], sanitizedTags)
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

	options.AnsibleForceColor()

	ansiblePlaybook := &ansible.AnsiblePlaybookCmd{
		Playbooks:         []string{playbook.(string)},
		Options:           ansiblePlaybookOptions,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithWrite(console.GetConsole()),
			execute.WithTransformers(
				results.Prepend(o.OutputPrefix),
			),
		),
	}

	console.Blue(ansiblePlaybook.String())

	return ansiblePlaybook, nil
}
