package driver

import (
	"context"
	"fmt"
	ansibledriver "stevedore/internal/driver/ansible"
	defaultdriver "stevedore/internal/driver/default"
	dockerdriver "stevedore/internal/driver/docker"
	mockdriver "stevedore/internal/driver/mock"
	"stevedore/internal/types"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	expected := map[string]DriverFactory{
		ansibledriver.DriverName: ansibledriver.NewAnsiblePlaybookDriver,
		dockerdriver.DriverName:  dockerdriver.NewDockerDriver,
		defaultdriver.DriverName: defaultdriver.NewDefaultDriver,
	}

	InitFactories()
	for driver, _ := range expected {
		_, exists := driverFactories[driver]
		assert.True(t, exists, fmt.Sprintf("Unregistered driver '%s'", driver))
	}

}

func TestRegisterFactory(t *testing.T) {

	tests := []struct {
		desc    string
		name    string
		factory DriverFactory
		preFunc func()
		res     types.Driverer
		err     error
	}{
		{
			desc:    "Testing nil factory",
			name:    "test",
			factory: nil,
			preFunc: nil,
			res:     nil,
			err:     errors.New("(builder::RegisterDriverFactory)", "Registring a nil factory"),
		},
		{
			desc: "Testing register factory",
			name: "test",
			factory: func(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {
				return &mockdriver.MockDriver{}, nil
			},
			preFunc: func() {
				driverFactories = map[string]DriverFactory{}
			},
			res: &mockdriver.MockDriver{},
			err: nil,
		},
		{
			desc: "Testing register registerd factory",
			name: "repeat",
			factory: func(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {
				return &mockdriver.MockDriver{}, nil
			},
			preFunc: func() {
				driverFactories = map[string]DriverFactory{}
				RegisterDriverFactory("repeat", func(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {
					return &mockdriver.MockDriver{}, nil
				})
			},
			res: nil,
			err: errors.New("(builder::RegisterDriverFactory)", "Driver factory 'repeat' already registered"),
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		if test.preFunc != nil {
			test.preFunc()
		}

		err := RegisterDriverFactory(test.name, test.factory)
		if err != nil && assert.Error(t, err) {
			assert.Equal(t, test.err.Error(), err.Error())
		} else {

			f, _ := driverFactories[test.name]
			builderer, _ := f(nil, nil)
			assert.Equal(t, builderer, test.res, "Unexpected value")
		}
	}
}

func TestGetDriverFactory(t *testing.T) {

	tests := []struct {
		desc    string
		name    string
		preFunc func()
		exists  bool
	}{
		{
			desc: "Testing existing builder",
			name: "builder",
			preFunc: func() {
				driverFactories = map[string]DriverFactory{}
				RegisterDriverFactory("builder", func(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {
					return &mockdriver.MockDriver{}, nil
				})
			},
			exists: true,
		},
		{
			desc: "Testing unexisting builder",
			name: "builder",
			preFunc: func() {
				driverFactories = map[string]DriverFactory{}
			},
			exists: false,
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		if test.preFunc != nil {
			test.preFunc()
		}

		_, exists := GetDriverFactory(test.name)
		assert.Equal(t, exists, test.exists, "Unexpected value")

	}
}

func TestClearBDriverFactory(t *testing.T) {
	t.Log("Testing clear DriverFactory data structure")

	// Clear and register a new factory
	driverFactories = map[string]DriverFactory{}
	RegisterDriverFactory("builder", func(ctx context.Context, o *types.BuildOptions) (types.Driverer, error) {
		return &mockdriver.MockDriver{}, nil
	})

	// Clear factories
	ClearDriverFactory()

	assert.Equal(t, driverFactories, map[string]DriverFactory{}, "driverFactories has not been cleared")
}
