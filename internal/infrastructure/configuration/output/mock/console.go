package console

import (
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/mock"
)

type ConfigurationMockOutput struct {
	mock.Mock
}

func NewConfigurationMockOutput() *ConfigurationMockOutput {
	return &ConfigurationMockOutput{}
}

func (o *ConfigurationMockOutput) Write(conf *configuration.Configuration) error {

	args := o.Called(conf)
	return args.Error(0)
}
