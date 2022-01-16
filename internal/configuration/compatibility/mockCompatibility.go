package compatibility

import (
	"github.com/stretchr/testify/mock"
)

// Compatibility holds compatibility details for configuration
type MockCompatibility struct {
	mock.Mock
}

// NewCompatibility creates a new compatibility checker
func NewMockCompatibility() *MockCompatibility {
	return &MockCompatibility{}
}

// AddDeprecated adds a deprecated field to the compatibility list
func (c *MockCompatibility) AddDeprecated(deprecated ...string) {
	c.Called(deprecated)
}

// AddRemoved adds a removed field to the compatibility list
func (c *MockCompatibility) AddRemoved(removed ...string) {
	c.Called(removed)
}

// AddChanged adds a changed field to the compatibility list
func (c *MockCompatibility) AddChanged(changed ...string) {
	c.Called(changed)
}

// Report return the compatibility issues
func (c *MockCompatibility) Report() {
	c.Called()
}
