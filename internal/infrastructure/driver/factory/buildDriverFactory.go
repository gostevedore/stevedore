package factory

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

type BuildDriverFactoryFunc func() (repository.BuildDriverer, error)

// BuildDriverFactory type define a map of BuildDriverer
type BuildDriverFactory map[string]BuildDriverFactoryFunc

// NewBuildDriverFactory returns a new BuildDriverFactory
func NewBuildDriverFactory() BuildDriverFactory {
	return make(BuildDriverFactory)
}

// Get returns a BuildDriverer
func (f BuildDriverFactory) Get(id string) (BuildDriverFactoryFunc, error) {
	errContext := "(BuildDriverFactory::Get)"

	driver, exist := f[id]
	if !exist {
		return nil, errors.New(errContext, fmt.Sprintf("Driver '%s' has not been registered", id))
	}

	return driver, nil
}

// Register registers a BuildDriverer
func (f BuildDriverFactory) Register(id string, driver BuildDriverFactoryFunc) error {

	errContext := "(BuildDriverFactory::Register)"

	_, exist := f[id]
	if exist {
		return errors.New(errContext, fmt.Sprintf("Driver '%s' already registered", id))
	}

	f[id] = driver

	return nil
}
