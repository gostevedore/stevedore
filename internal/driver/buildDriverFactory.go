package driver

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
)

// BuildDriverFactory type define a map of BuildDriverer
type BuildDriverFactory map[string]BuildDriverer

// NewBuildDriverFactory returns a new BuildDriverFactory
func NewBuildDriverFactory() BuildDriverFactory {
	return make(BuildDriverFactory)
}

// Get returns a BuildDriverer
func (f BuildDriverFactory) Get(id string) (BuildDriverer, error) {
	errContext := "(BuildDriverFactory::Get)"

	driver, exist := f[id]
	if !exist {
		return nil, errors.New(errContext, fmt.Sprintf("Driver '%s' has not been registered", id))
	}

	return driver, nil
}

// Register registers a BuildDriverer
func (f BuildDriverFactory) Register(id string, driver BuildDriverer) error {

	errContext := "(BuildDriverFactory::Register)"

	_, exist := f[id]
	if exist {
		return errors.New(errContext, fmt.Sprintf("Driver '%s' already registered", id))
	}

	f[id] = driver

	return nil
}

// DriverFactory type define functions that provides a Driverer
// type DriverFactory func(ctx context.Context, o *types.BuildOptions) (Driverer, error)

// //  driverFactories maps each driver to its builder factory
// var driverFactories map[string]DriverFactory

// InitFactories initizalize the driverFactories data structure mapping each driver to its builder factory
// func InitFactories() error {
// 	var err error

// 	driverFactories = map[string]DriverFactory{}

// 	err = RegisterDriverFactory(ansibledriver.DriverName, ansibledriver.NewAnsiblePlaybookDriver)
// 	if err != nil {
// 		return errors.New("(build::Init)", "Builder could not be registered", err)
// 	}

// 	err = RegisterDriverFactory(dockerdriver.DriverName, dockerdriver.NewDockerDriver)
// 	if err != nil {
// 		return errors.New("(build::Init)", "Builder could not be registered", err)
// 	}

// 	err = RegisterDriverFactory(defaultdriver.DriverName, defaultdriver.NewDefaultDriver)
// 	if err != nil {
// 		return errors.New("(build::Init)", "Builder could not be registered", err)
// 	}

// 	return nil
// }

// // Register registers a new factory
// func RegisterDriverFactory(driver string, factory DriverFactory) error {
// 	if factory == nil {
// 		return errors.New("(builder::RegisterDriverFactory)", "Registring a nil factory")
// 	}

// 	_, registered := driverFactories[driver]
// 	if registered {
// 		return errors.New("(builder::RegisterDriverFactory)", "Driver factory '"+driver+"' already registered")
// 	}

// 	driverFactories[driver] = factory
// 	return nil
// }

// // GetFactory returns a factory associated to a driver
// func GetDriverFactory(driver string) (DriverFactory, bool) {
// 	factory, exists := driverFactories[driver]
// 	return factory, exists
// }

// // ClearBDriverFactory clears factories defined on driverFactories
// func ClearDriverFactory() {
// 	driverFactories = map[string]DriverFactory{}
// }
