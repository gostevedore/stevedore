package ansibler

import (
	"context"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	ansible "github.com/apenella/go-ansible/pkg/playbook"
	"github.com/stretchr/testify/mock"
)

type MockAnsibleDriver struct {
	mock.Mock
}

func NewMockAnsibleDriver() *MockAnsibleDriver {
	return &MockAnsibleDriver{}
}

func (d *MockAnsibleDriver) WithPlaybook(playbook string) {
	d.Mock.Called(playbook)
}

func (d *MockAnsibleDriver) WithOptions(opts *ansible.AnsiblePlaybookOptions) {
	d.Mock.Called(opts)
}

func (d *MockAnsibleDriver) WithConnectionOptions(opts *options.AnsibleConnectionOptions) {
	d.Mock.Called(opts)
}

func (d *MockAnsibleDriver) WithPriviledgedEscalationOptions(opts *options.AnsiblePrivilegeEscalationOptions) {
	d.Mock.Called(opts)
}

func (d *MockAnsibleDriver) WithStdoutCallback(callback string) {
	d.Mock.Called(callback)
}

func (d *MockAnsibleDriver) WithExecutor(executor execute.Executor) {
	d.Mock.Called(executor)
}

func (d *MockAnsibleDriver) Run(ctx context.Context) error {
	args := d.Mock.Called(ctx)

	return args.Error(0)
}
