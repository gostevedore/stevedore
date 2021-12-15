package goansible

import (
	"context"
	"io"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
)

// GoAnsibleDriver is a driver for building docker images from ansible-playbooks
type GoAnsibleDriver struct {
	ansible *ansible.AnsiblePlaybookCmd
}

// NewGoAnsibleDriver creates a new GoAnsibleDriver
func NewGoAnsibleDriver() *GoAnsibleDriver {
	return &GoAnsibleDriver{
		ansible: &ansible.AnsiblePlaybookCmd{},
	}
}

// WithPlaybook sets the playbook to ansible command
func (d *GoAnsibleDriver) WithPlaybook(playbook string) {
	d.ansible.Playbooks = append(d.ansible.Playbooks, playbook)
}

// WithOptions sets the options to ansible command
func (d *GoAnsibleDriver) WithOptions(opts *ansible.AnsiblePlaybookOptions) {
	d.ansible.Options = opts
}

// WithConnectionOptions sets the connection options to ansible command
func (d *GoAnsibleDriver) WithConnectionOptions(opts *options.AnsibleConnectionOptions) {
	d.ansible.ConnectionOptions = opts
}

// WithStdoutCallback sets the callback to ansible command
func (d *GoAnsibleDriver) WithStdoutCallback(callback string) {
	d.ansible.StdoutCallback = callback
}

// WithExecutor sets the executor to ansible command
func (d *GoAnsibleDriver) WithExecutor(executor execute.Executor) {
	d.ansible.Exec = executor
}

// WithPriviledgedEscalationOptions sets the priviliged escalations options to ansible command
func (d *GoAnsibleDriver) WithPriviledgedEscalationOptions(opts *options.AnsiblePrivilegeEscalationOptions) {
	d.ansible.PrivilegeEscalationOptions = opts
}

// PrepareExecutor prepares the executor
func (d *GoAnsibleDriver) PrepareExecutor(writer io.Writer, prefix string) {
	executor := execute.NewDefaultExecute(
		execute.WithWrite(writer),
		execute.WithTransformers(
			results.Prepend(prefix),
		),
	)

	d.ansible.Exec = executor
}

// Run executes the ansible command
func (d *GoAnsibleDriver) Run(ctx context.Context) error {
	// TODO: only debug, replace with decorator
	// console.Blue(d.ansible.String())

	// Setup ansible options
	options.AnsibleForceColor()

	return d.ansible.Run(ctx)
}
