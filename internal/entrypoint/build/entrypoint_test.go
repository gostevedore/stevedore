package build

import (
	"context"
	"io/ioutil"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/credentials"
	"github.com/gostevedore/stevedore/internal/driver"
	ansibledriver "github.com/gostevedore/stevedore/internal/driver/ansible"
	defaultdriver "github.com/gostevedore/stevedore/internal/driver/default"
	dockerdriver "github.com/gostevedore/stevedore/internal/driver/docker"
	dryrundriver "github.com/gostevedore/stevedore/internal/driver/dryrun"
	"github.com/gostevedore/stevedore/internal/engine/build/plan"
	build "github.com/gostevedore/stevedore/internal/handler/build"
	"github.com/gostevedore/stevedore/internal/images/store"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {

}

func TestExecute(t *testing.T) {

	errContext := "(Entrypoint::Execute)"

	tests := []struct {
		desc              string
		entrypoint        *Entrypoint
		args              []string
		configuration     *configuration.Configuration
		entrypointOptions *EntrypointOptions
		handlerOptions    *build.HandlerOptions
		err               error
		assertions        func(*testing.T, *Entrypoint, []string, *EntrypointOptions, *build.HandlerOptions)
	}{
		{
			desc:       "Testing error when configuration is not provided",
			entrypoint: &Entrypoint{},
			err:        errors.New(errContext, "To execute the build entrypoint, configuration is required"),
		},
		{
			desc:          "Testing error when arguments are not provided",
			entrypoint:    &Entrypoint{},
			configuration: &configuration.Configuration{},
			err:           errors.New(errContext, "To execute the build entrypoint, arguments are required"),
		},
		{
			desc:          "Testing error when entrypoint options are not provided",
			entrypoint:    &Entrypoint{},
			configuration: &configuration.Configuration{},
			args:          []string{"image"},
			err:           errors.New(errContext, "To execute the build entrypoint, entrypoint options are required"),
		},
		{
			desc:              "Testing error when handler options are not provided",
			entrypoint:        &Entrypoint{},
			configuration:     &configuration.Configuration{},
			args:              []string{"image"},
			entrypointOptions: &EntrypointOptions{},
			err:               errors.New(errContext, "To execute the build entrypoint, handler options are required"),
		},
		{
			desc: "Testing execute entrypoint",
			entrypoint: &Entrypoint{
				writer: ioutil.Discard,
			},
			configuration:     &configuration.Configuration{},
			args:              []string{"image"},
			entrypointOptions: &EntrypointOptions{},
			handlerOptions:    &build.HandlerOptions{},
			err:               &errors.Error{},
			assertions:        func(*testing.T, *Entrypoint, []string, *EntrypointOptions, *build.HandlerOptions) {},
		},
		{
			desc: "Testing execute entrypoint overriding handler options with config",
			entrypoint: &Entrypoint{
				writer: ioutil.Discard,
			},
			args: []string{"image"},
			configuration: &configuration.Configuration{
				Concurrency:               5,
				PushImages:                true,
				EnableSemanticVersionTags: true,
				SemanticVersionTagsTemplates: []string{
					"template",
				},
			},
			entrypointOptions: &EntrypointOptions{},
			handlerOptions:    &build.HandlerOptions{},
			err:               &errors.Error{},
			assertions: func(t *testing.T, e *Entrypoint, args []string, entrypointOptions *EntrypointOptions, handlerOptions *build.HandlerOptions) {
				assert.Equal(t, 5, entrypointOptions.Concurrency, "Concurrency should be 5")
				assert.True(t, handlerOptions.PushImagesAfterBuild, "Push images after build should be true")
				assert.True(t, handlerOptions.EnableSemanticVersionTags, "Enable semantic version tags should be true")
				assert.Equal(t, []string{"template"}, handlerOptions.SemanticVersionTagsTemplates, "Semantic version tags templates is not as expected")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			err := test.entrypoint.Execute(context.TODO(), test.args, test.configuration, test.entrypointOptions, test.handlerOptions)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertions(t, test.entrypoint, test.args, test.entrypointOptions, test.handlerOptions)
			}
		})
	}
}

func TestCreateBuildDriverFactory(t *testing.T) {

	errContext := "(entrypoint::createBuildDriverFactory)"

	tests := []struct {
		desc        string
		entrypoint  *Entrypoint
		credentials *credentials.CredentialsStore
		options     *EntrypointOptions
		err         error
		assertions  func(t *testing.T, driverFactory driver.BuildDriverFactory)
	}{
		{
			desc:        "Testing create build driver factory with empty credentials",
			entrypoint:  NewEntrypoint(),
			credentials: nil,
			options:     nil,
			err:         errors.New(errContext, "Register drivers requires a credentials store"),
		},
		{
			desc:        "Testing create build driver factory with empty options",
			entrypoint:  NewEntrypoint(),
			credentials: credentials.NewCredentialsStore(),
			options:     nil,
			err:         errors.New(errContext, "Register drivers requires options"),
		},
		{
			desc:        "Testing create build driver factory with nil writer",
			entrypoint:  NewEntrypoint(),
			credentials: credentials.NewCredentialsStore(),
			options:     &EntrypointOptions{},
			err:         errors.New(errContext, "Register drivers requires a writer"),
		},
		{
			desc:        "Testing create build driver factory",
			entrypoint:  NewEntrypoint(WithWriter(ioutil.Discard)),
			credentials: credentials.NewCredentialsStore(),
			options:     &EntrypointOptions{},
			err:         &errors.Error{},
			assertions: func(t *testing.T, f driver.BuildDriverFactory) {
				dDocker, eDocker := f.Get("docker")
				assert.Nil(t, eDocker)
				assert.NotNil(t, dDocker)
				assert.IsType(t, &dockerdriver.DockerDriver{}, dDocker)

				dAnsible, eAnsible := f.Get("ansible-playbook")
				assert.Nil(t, eAnsible)
				assert.NotNil(t, dAnsible)
				assert.IsType(t, &ansibledriver.AnsiblePlaybookDriver{}, dAnsible)

				dDefault, eDefault := f.Get("default")
				assert.Nil(t, eDefault)
				assert.NotNil(t, dDefault)
				assert.IsType(t, &defaultdriver.DefaultDriver{}, dDefault)
			},
		},
		{
			desc:        "Testing create build driver factory with dry run",
			entrypoint:  NewEntrypoint(WithWriter(ioutil.Discard)),
			credentials: credentials.NewCredentialsStore(),
			options: &EntrypointOptions{
				DryRun: true,
			},
			err: &errors.Error{},
			assertions: func(t *testing.T, f driver.BuildDriverFactory) {
				dDocker, eDocker := f.Get("docker")
				assert.Nil(t, eDocker)
				assert.NotNil(t, dDocker)
				assert.IsType(t, &dryrundriver.DryRunDriver{}, dDocker)

				dAnsible, eAnsible := f.Get("ansible-playbook")
				assert.Nil(t, eAnsible)
				assert.NotNil(t, dAnsible)
				assert.IsType(t, &dryrundriver.DryRunDriver{}, dDocker)

				dDefault, eDefault := f.Get("default")
				assert.Nil(t, eDefault)
				assert.NotNil(t, dDefault)
				assert.IsType(t, &dryrundriver.DryRunDriver{}, dDocker)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			factory, err := test.entrypoint.createBuildDriverFactory(test.credentials, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.assertions(t, factory)
			}
		})
	}
}

func TestCreateDryRunDriver(t *testing.T) {
	desc := "Testing create dry-run driver"

	t.Run(desc, func(t *testing.T) {
		e := NewEntrypoint()

		dryRunDriver, err := e.createDryRunDriver()

		assert.Nil(t, err)
		assert.NotNil(t, dryRunDriver)
	})
}

func TestCreateDefaultDriver(t *testing.T) {

	errContext := "(entrypoint::createDefaultDriver)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		options    *EntrypointOptions
		res        driver.BuildDriverer
		err        error
	}{
		{
			desc:       "Testing error when creating default driver with nil options",
			entrypoint: NewEntrypoint(),
			options:    nil,
			err:        errors.New(errContext, "Entrypoint options are required to create default driver"),
		},
		{
			desc:       "Testing create default driver",
			entrypoint: NewEntrypoint(),
			options:    &EntrypointOptions{},
			res:        &defaultdriver.DefaultDriver{},
		},
		{
			desc:       "Testing create default driver with dry-run enabled",
			entrypoint: NewEntrypoint(),
			options: &EntrypointOptions{
				DryRun: true,
			},
			res: &dryrundriver.DryRunDriver{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			driver, err := test.entrypoint.createDefaultDriver(test.options)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}
}

func TestCreateAnsibleDriver(t *testing.T) {

	errContext := "(entrypoint::createAnsibleDriver)"

	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		options    *EntrypointOptions
		res        driver.BuildDriverer
		err        error
	}{
		{
			desc:       "Testing error when creating ansible driver with nil options",
			entrypoint: NewEntrypoint(),
			options:    nil,
			err:        errors.New(errContext, "Entrypoint options are required to create ansible driver"),
		},
		{
			desc:       "Testing create ansible driver",
			entrypoint: NewEntrypoint(),
			options:    &EntrypointOptions{},
			res:        &ansibledriver.AnsiblePlaybookDriver{},
		},
		{
			desc:       "Testing create ansible driver with dry-run enabled",
			entrypoint: NewEntrypoint(),
			options: &EntrypointOptions{
				DryRun: true,
			},
			res: &dryrundriver.DryRunDriver{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			driver, err := test.entrypoint.createAnsibleDriver(test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}
		})
	}
}

func TestCreateDockerDriver(t *testing.T) {
	errContext := "(entrypoint::createDockerDriver)"

	tests := []struct {
		desc        string
		entrypoint  *Entrypoint
		credentials *credentials.CredentialsStore
		options     *EntrypointOptions
		res         driver.BuildDriverer
		err         error
	}{
		{
			desc:        "Testing error when creating docker driver with empty credentials",
			entrypoint:  NewEntrypoint(),
			credentials: nil,
			err:         errors.New(errContext, "Docker driver requires a credentials store"),
		},
		{
			desc:        "Testing error when creating docker driver with empty options",
			entrypoint:  NewEntrypoint(),
			credentials: credentials.NewCredentialsStore(),
			options:     nil,
			err:         errors.New(errContext, "Entrypoint options are required to create docker driver"),
		},
		{
			desc:        "Testing create docker driver",
			entrypoint:  NewEntrypoint(),
			credentials: credentials.NewCredentialsStore(),
			options:     &EntrypointOptions{},
			res:         &dockerdriver.DockerDriver{},
			err:         &errors.Error{},
		},
		{
			desc:        "Testing create docker driver with dry-run enabled",
			entrypoint:  NewEntrypoint(),
			credentials: credentials.NewCredentialsStore(),
			options: &EntrypointOptions{
				DryRun: true,
			},
			res: &dryrundriver.DryRunDriver{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			driver, err := test.entrypoint.createDockerDriver(test.credentials, test.options)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, driver)
				assert.IsType(t, test.res, driver)
			}

		})
	}
}

func TestCreateDispatcher(t *testing.T) {
	desc := "Testing create dispatcher"

	t.Run(desc, func(t *testing.T) {
		e := NewEntrypoint()
		options := &EntrypointOptions{
			Concurrency: 5,
		}

		dispatch, err := e.createDispatcher(options)

		assert.Nil(t, err)
		assert.NotNil(t, dispatch)
		assert.NotNil(t, dispatch.WorkerPool)
		assert.Equal(t, dispatch.NumWorkers, 5)
	})
}

func TestCreatePlanFactory(t *testing.T) {
	desc := "Testing create build plan factory"

	t.Run(desc, func(t *testing.T) {
		e := NewEntrypoint()
		options := &EntrypointOptions{}

		imageStore := store.NewImageStore(nil)
		planFactory, err := e.createPlanFactory(imageStore, options)

		assert.Nil(t, err)
		assert.NotNil(t, planFactory)
		assert.IsType(t, plan.NewPlanFactory(imageStore), planFactory)
	})
}
