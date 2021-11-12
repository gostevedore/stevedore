package ansibledriver

import (
	"context"
	"io"

	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
)

// Ansibler is an interface that describes the methods that must implement any objects for building images from ansible
type Ansibler interface {
	WithPlaybook(playbook string)
	WithOptions(opts *ansible.AnsiblePlaybookOptions)
	WithConnectionOptions(opts *options.AnsibleConnectionOptions)
	WithPriviledgedEscalationOptions(opts *options.AnsiblePrivilegeEscalationOptions)
	WithStdoutCallback(callback string)
	PrepareExecutor(writer io.Writer, prefix string)
	Run(ctx context.Context) error
}
