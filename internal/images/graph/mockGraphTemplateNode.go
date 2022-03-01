package graph

import (
	// gdsexttree "github.com/apenella/go-data-structures/extendedTree"
	gdsexttree "github.com/gostevedore/stevedore/pkg/extendedTree"
	"github.com/stretchr/testify/mock"
)

// MockGraphTemplateNode is a mock of GraphTemplateNode
type MockGraphTemplateNode struct {
	mock.Mock
}

// NewMockGraphTemplateNode creates a new MockGraphTemplateNode
func NewMockGraphTemplateNode() *MockGraphTemplateNode {
	return &MockGraphTemplateNode{}
}

func (m *MockGraphTemplateNode) getNode() *gdsexttree.Node {
	args := m.Called()

	return args.Get(0).(*gdsexttree.Node)
}

// AddChild adds a child to the graph template node
func (m *MockGraphTemplateNode) AddChild(child GraphTemplateNoder) error {
	args := m.Called(child)

	return args.Error(0)
}

// AddParent adds a parent to the graph template node
func (m *MockGraphTemplateNode) AddParent(parent GraphTemplateNoder) error {
	args := m.Called(parent)

	return args.Error(0)
}

// AddItem adds a item to the graph template node
func (m *MockGraphTemplateNode) AddItem(item interface{}) {
	m.Called(item)
}
