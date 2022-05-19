package goansible

import (
	"context"
	"io"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	"github.com/stretchr/testify/mock"
)

// MockAnsibleDriver is a mock of AnsibleDriver interface
type MockAnsibleDriver struct {
	mock.Mock
}

// NewMockAnsibleDriver returns a new mock of AnsibleDriver interface
func NewMockAnsibleDriver() *MockAnsibleDriver {
	return &MockAnsibleDriver{}
}

// WithPlaybook returns a mock of AnsibleDriver interface
func (d *MockAnsibleDriver) WithPlaybook(playbook string) {
	d.Mock.Called(playbook)
}

// WithOptions returns a mock of AnsibleDriver interface
func (d *MockAnsibleDriver) WithOptions(opts *ansible.AnsiblePlaybookOptions) {
	d.Mock.Called(opts)
}

// WithConnectionOptions returns a mock of AnsibleDriver interface
func (d *MockAnsibleDriver) WithConnectionOptions(opts *options.AnsibleConnectionOptions) {
	d.Mock.Called(opts)
}

// WithPriviledgedEscalationOptions returns a mock of AnsibleDriver interface
func (d *MockAnsibleDriver) WithPriviledgedEscalationOptions(opts *options.AnsiblePrivilegeEscalationOptions) {
	d.Mock.Called(opts)
}

// WithStdoutCallback returns a mock of AnsibleDriver interface
func (d *MockAnsibleDriver) WithStdoutCallback(callback string) {
	d.Mock.Called(callback)
}

// WithExecutor returns a mock of AnsibleDriver interface
func (d *MockAnsibleDriver) WithExecutor(executor execute.Executor) {
	d.Mock.Called(executor)
}

// PrepareExecutor returns a mock of AnsibleDriver interface
func (d *MockAnsibleDriver) PrepareExecutor(writer io.Writer, prefix string) {
	d.Mock.Called(writer, prefix)
}

// Run returns a mock of AnsibleDriver interface
func (d *MockAnsibleDriver) Run(ctx context.Context) error {
	args := d.Mock.Called(ctx)

	return args.Error(0)
}
