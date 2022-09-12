package mock

import (
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/mock"
)

// MockStore is a mocked store for credentials
type MockStore struct {
	mock.Mock
}

// NewMockStore creates a new mocked store for credentials
func NewMockStore() *MockStore {
	return &MockStore{}
}

// Store stores a badge
func (m *MockStore) Store(id string, badge *credentials.Badge) error {
	args := m.Mock.Called(id, badge)
	return args.Error(0)
}

// Get returns a auth for the badge id
func (m *MockStore) Get(id string) (*credentials.Badge, error) {
	args := m.Mock.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*credentials.Badge), args.Error(1)
	}
}

// All returns all badges
func (m *MockStore) All() []*credentials.Badge {
	args := m.Mock.Called()
	return args.Get(0).([]*credentials.Badge)
}
