package graph

import "github.com/stretchr/testify/mock"

// MockGraphTemplate mocks a graph template for images
type MockGraphTemplate struct {
	mock.Mock
}

// NewMockGraphTemplate creates a new graph template for images
func NewMockGraphTemplate() *MockGraphTemplate {
	return &MockGraphTemplate{}
}

// GetNode gets a node from the graph template
func (m *MockGraphTemplate) GetNode(name string) GraphTemplateNoder {
	args := m.Called(name)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(GraphTemplateNoder)
}

// AddNode adds a node to the graph template
func (m *MockGraphTemplate) AddNode(node GraphTemplateNoder) error {
	args := m.Called(node)

	return args.Error(0)
}

// AddRelationship adds a relationship to the graph template
func (m *MockGraphTemplate) AddRelationship(parent, child GraphTemplateNoder) error {
	args := m.Called(parent, child)

	return args.Error(0)
}

// Exists checks if a node exists in the graph template
func (m *MockGraphTemplate) Exists(name string) bool {
	args := m.Called(name)

	return args.Bool(0)
}

// HasCycles checks if the graph template has cycles
func (m *MockGraphTemplate) HasCycles() bool {
	args := m.Called()

	return args.Bool(0)
}
