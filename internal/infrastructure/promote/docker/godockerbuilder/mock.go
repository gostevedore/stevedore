package godockerbuilder

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type PromoteMock struct {
	mock.Mock
}

// NewPromoteMock creates a new mock for docker client
func NewPromoteMock() *PromoteMock {
	return &PromoteMock{}
}

func (m *PromoteMock) AddTag(tags ...string) {
	m.Mock.Called(tags)
}

// Run
func (m *PromoteMock) Run(ctx context.Context) error {
	args := m.Mock.Called(ctx)
	return args.Error(0)
}

// AddAuth
func (m *PromoteMock) AddAuth(user string, pass string) error {
	args := m.Mock.Called(user, pass)
	return args.Error(0)
}

// AddAuthPull
func (m *PromoteMock) AddPullAuth(user string, pass string) error {
	args := m.Mock.Called(user, pass)
	return args.Error(0)
}

// AddAuthPush
func (m *PromoteMock) AddPushAuth(user string, pass string) error {
	args := m.Mock.Called(user, pass)
	return args.Error(0)
}

// WithSourceImage
func (m *PromoteMock) WithSourceImage(name string) {
	m.Mock.Called(name)
}

// WithSourceImage
func (m *PromoteMock) WithTargetImage(name string) {
	m.Mock.Called(name)
}

// WithTags
func (m *PromoteMock) WithTags(tags []string) {
	m.Mock.Called(tags)
}

// WithRemoteSource
func (m *PromoteMock) WithRemoteSource() {
	m.Mock.Called()
}

// WithRemoveAfterPush
func (m *PromoteMock) WithRemoveAfterPush() {
	m.Mock.Called()
}

// WithResponse
func (m *PromoteMock) WithResponse(w io.Writer, prefix string) {
	m.Mock.Called(w, prefix)
}

// WithUseNormalizedNamed
func (m *PromoteMock) WithUseNormalizedNamed() {
	m.Mock.Called()
}
