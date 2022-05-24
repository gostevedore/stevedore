package builders

import (
	"fmt"
	"sync"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
)

// Store contains builders details
type Store struct {
	mutex    sync.RWMutex
	Builders map[string]*builder.Builder
}

// NewStore creates a new builders configuration
func NewStore() *Store {
	return &Store{
		Builders: make(map[string]*builder.Builder),
	}
}

// Store include a new builder to builders
func (b *Store) Store(builder *builder.Builder) error {

	errContext := "(builders::Store)"

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
func (b *Store) Find(name string) (*builder.Builder, error) {

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
