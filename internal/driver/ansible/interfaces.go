package ansibledriver

import (
	"context"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
)

type Ansibler interface {
	WithPlaybook(playbook string)
	WithOptions(opts *ansible.AnsiblePlaybookOptions)
	WithConnectionOptions(opts *options.AnsibleConnectionOptions)
	WithPriviledgedEscalationOptions(opts *options.AnsiblePrivilegeEscalationOptions)
	WithStdoutCallback(callback string)
	WithExecutor(executor execute.Executor)
	Run(ctx context.Context) error
}
