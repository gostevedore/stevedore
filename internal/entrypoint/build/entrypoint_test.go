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

func TestCreateDefaultDriver(t *testing.T) {
	desc := "Testing create default driver"

	t.Run(desc, func(t *testing.T) {
		e := NewEntrypoint()

		defaultDriver, err := e.createDefaultDriver()

		assert.Nil(t, err)
		assert.NotNil(t, defaultDriver)
	})
}

func TestCreateAnsibleDriver(t *testing.T) {
	desc := "Testing create ansible driver"

	t.Run(desc, func(t *testing.T) {
		e := NewEntrypoint()

		ansibleDriver, err := e.createAnsibleDriver()

		assert.Nil(t, err)
		assert.NotNil(t, ansibleDriver)
	})
}

func TestCreateDockerDriver(t *testing.T) {

	tests := []struct {
		desc     string
		testFunc func(t *testing.T)
	}{
		{
			desc: "Testing error when creating docker driver with empty credentials",
			testFunc: func(t *testing.T) {
				e := NewEntrypoint()
				_, err := e.createDockerDriver(nil)

				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "Docker driver requires a credentials store")
			},
		},
		{
			desc: "Testing create docker driver",
			testFunc: func(t *testing.T) {
				e := NewEntrypoint()
				dockerDriver, err := e.createDockerDriver(credentials.NewCredentialsStore())

				assert.Nil(t, err)
				assert.NotNil(t, dockerDriver)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			test.testFunc(t)
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
