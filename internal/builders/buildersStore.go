package builders

import (
	"fmt"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/builders/builder"
)

// BuildersStore contains builders details
type BuildersStore struct {
	mutex    sync.RWMutex
	Builders map[string]*builder.Builder
}

// NewBuildersStore creates a new builders configuration
func NewBuildersStore() *BuildersStore {
	return &BuildersStore{
		Builders: make(map[string]*builder.Builder),
	}
}

// AddBuilder include a new builder to builders
func (b *BuildersStore) AddBuilder(builder *builder.Builder) error {

	errContext := "(builders::AddBuilder)"

	if b == nil {
		return errors.New(errContext, "Builders is nil")
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	_, exist := b.Builders[builder.Name]
	if exist {
		return errors.New(errContext, fmt.Sprintf("Builder '%s' already exist", builder.Name))
	}

	b.Builders[builder.Name] = builder

	return nil
}

// Find returns the builder registered with input name
func (b *BuildersStore) Find(name string) (*builder.Builder, error) {

	errContext := "(builders::GetBuilder)"

	if b == nil {
		return nil, errors.New(errContext, "Builders is nil")
	}

	b.mutex.RLock()
	builder, exists := b.Builders[name]
	if !exists {
		return nil, errors.New(errContext, fmt.Sprintf("Builder '%s' does not exists", name))
	}
	b.mutex.RUnlock()

	return builder, nil
}
