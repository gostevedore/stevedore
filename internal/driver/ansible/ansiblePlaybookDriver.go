package ansibledriver

import (
	"context"
	"stevedore/internal/build/varsmap"
	"stevedore/internal/driver/common"
	drivercommon "stevedore/internal/driver/common"
	"stevedore/internal/types"
	"stevedore/internal/ui/console"
	"strings"

	ansibler "github.com/apenella/go-ansible"
	errors "github.com/apenella/go-common-utils/error"
)

const (
	DriverName = "ansible-playbook"
)

func NewAnsiblePlaybookDriver(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {

	builderName := ""

	if o == nil {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "Build options are nil")
	}

	builderConfOptions := o.BuilderOptions

	playbook, ok := builderConfOptions["playbook"]
	if !ok {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "playbook has not been defined on build options")
	}

	if o.ImageName == "" {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "Image name is not set")
	}

	if o.RegistryNamespace == "" {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "Registry namespace is not set")
	}

	builderName = strings.Join([]string{"builder", o.RegistryNamespace, o.ImageName}, "_")

	inventory, ok := builderConfOptions["inventory"]
	if !ok {
		return nil, errors.New("(build::NewAnsiblePlaybookDriver)", "inventory has not been defined on build options")
	}

	ansiblePlaybookOptions := &ansibler.AnsiblePlaybookOptions{
		Inventory: inventory.(string),
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

	ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageNameKey], o.ImageName)
	ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingRegistryNamespaceKey], o.RegistryNamespace)

	if o.RegistryHost != "" {
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingRegistryHostKey], o.RegistryHost)
	}

	if o.ImageVersion != "" {
		o.ImageVersion = common.SanitizeTag(o.ImageVersion)
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageTagKey], o.ImageVersion)
		builderName = strings.Join([]string{builderName, o.ImageVersion}, "_")
	}

	if o.OutputPrefix == "" {
		o.OutputPrefix = o.ImageName
		if o.ImageVersion != "" {
			o.OutputPrefix = strings.Join([]string{o.OutputPrefix, o.ImageVersion}, ":")
		}
	}

	if len(o.BuilderName) > 0 {
		builderName = o.BuilderName
	}

	if len(o.Tags) > 0 {
		sanitizedTags := []string{}
		for _, tag := range o.Tags {
			tag = drivercommon.SanitizeTag(tag)
			sanitizedTags = append(sanitizedTags, tag)
		}
		ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageExtraTagsKey], sanitizedTags)
	}

	ansiblePlaybookOptions.AddExtraVar(o.BuilderVarMappings[varsmap.VarMappingImageBuilderLabelKey], builderName)

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

	ansiblePlaybookConnectionOptions := &ansibler.AnsiblePlaybookConnectionOptions{}
	if o.ConnectionLocal {
		ansiblePlaybookConnectionOptions.Connection = "local"
	}

	ansibler.AnsibleForceColor()

	ansiblePlaybook := &ansibler.AnsiblePlaybookCmd{
		Writer:            console.GetConsole(),
		Playbook:          playbook.(string),
		ExecPrefix:        o.OutputPrefix,
		Options:           ansiblePlaybookOptions,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
	}

	console.Blue(ansiblePlaybook.String())

	return ansiblePlaybook, nil
}
