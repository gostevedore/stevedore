package main

import (
	"context"
	"os"

	"github.com/gostevedore/stevedore/internal/infrastructure/cli/stevedore"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration/loader"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	"github.com/gostevedore/stevedore/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func main() {
	log := logger.New()
	defer log.Sync()

	fs := afero.NewOsFs()
	cons := console.NewConsole(os.Stdout, os.Stdin)
	compatibility := compatibility.NewCompatibility(cons)
	configLoader := loader.NewConfigurationLoader(viper.New())
	conf, err := configuration.New(fs, configLoader, compatibility)
	if err != nil {
		cons.Error(err.Error())
		os.Exit(1)
	}

	stevedore := stevedore.NewCommand(context.Background(), fs, compatibility, compatibility, cons, log, conf)
	err = stevedore.Execute()
	if err != nil {
		os.Exit(1)
	}
}
