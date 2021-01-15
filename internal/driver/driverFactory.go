package driver

import (
	"context"

	ansibledriver "github.com/gostevedore/stevedore/internal/driver/ansible"
	defaultdriver "github.com/gostevedore/stevedore/internal/driver/default"
	dockerdriver "github.com/gostevedore/stevedore/internal/driver/docker"
	"github.com/gostevedore/stevedore/internal/types"

	errors "github.com/apenella/go-common-utils/error"
)

// DriverFactory type define functions that provides a Driverer
type DriverFactory func(ctx context.Context, o *types.BuildOptions) (types.Driverer, error)

//  driverFactories maps each driver to its builder factory
var driverFactories map[string]DriverFactory

// InitFactories initizalize the driverFactories data structure mapping each driver to its builder factory
func InitFactories() error {
	var err error

	driverFactories = map[string]DriverFactory{}

	err = RegisterDriverFactory(ansibledriver.DriverName, ansibledriver.NewAnsiblePlaybookDriver)
	if err != nil {
		return errors.New("(build::Init)", "Builder could not be registered", err)
	}

	err = RegisterDriverFactory(dockerdriver.DriverName, dockerdriver.NewDockerDriver)
	if err != nil {
		return errors.New("(build::Init)", "Builder could not be registered", err)
	}

	err = RegisterDriverFactory(defaultdriver.DriverName, defaultdriver.NewDefaultDriver)
	if err != nil {
		return errors.New("(build::Init)", "Builder could not be registered", err)
	}

	return nil
}

// Register registers a new factory
func RegisterDriverFactory(driver string, factory DriverFactory) error {
	if factory == nil {
		return errors.New("(builder::RegisterDriverFactory)", "Registring a nil factory")
	}

	_, registered := driverFactories[driver]
	if registered {
		return errors.New("(builder::RegisterDriverFactory)", "Driver factory '"+driver+"' already registered")
	}

	driverFactories[driver] = factory
	return nil
}

// GetFactory returns a factory associated to a driver
func GetDriverFactory(driver string) (DriverFactory, bool) {
	factory, exists := driverFactories[driver]
	return factory, exists
}

// ClearBDriverFactory clears factories defined on driverFactories
func ClearDriverFactory() {
	driverFactories = map[string]DriverFactory{}
}
