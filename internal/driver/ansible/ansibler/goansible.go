package ansibler

import (
	"context"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	"github.com/gostevedore/stevedore/internal/ui/console"
)

type GoAnsibleDriver struct {
	ansible *ansible.AnsiblePlaybookCmd
}

func NewGoAnsibleDriver() *GoAnsibleDriver {
	return &GoAnsibleDriver{
		ansible: &ansible.AnsiblePlaybookCmd{},
	}
}

func (d *GoAnsibleDriver) WithPlaybook(playbook string) {
	d.ansible.Playbooks = append(d.ansible.Playbooks, playbook)
}

func (d *GoAnsibleDriver) WithOptions(opts *ansible.AnsiblePlaybookOptions) {
	d.ansible.Options = opts
}

func (d *GoAnsibleDriver) WithConnectionOptions(opts *options.AnsibleConnectionOptions) {
	d.ansible.ConnectionOptions = opts
}

func (d *GoAnsibleDriver) WithStdoutCallback(callback string) {
	d.ansible.StdoutCallback = callback
}

func (d *GoAnsibleDriver) WithExecutor(executor execute.Executor) {
	d.ansible.Exec = executor
}

func (d *GoAnsibleDriver) WithPriviledgedEscalationOptions(opts *options.AnsiblePrivilegeEscalationOptions) {
	d.ansible.PrivilegeEscalationOptions = opts
}

func (d *GoAnsibleDriver) Run(ctx context.Context) error {
	// TODO: only debug, replace with decorator
	console.Blue(d.ansible.String())

	// Setup ansible options
	options.AnsibleForceColor()

	return d.ansible.Run(ctx)
}
