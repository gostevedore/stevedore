package main

import (
	"context"
	"os"

	"github.com/gostevedore/stevedore/internal/infrastructure/cli/stevedore"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	"github.com/gostevedore/stevedore/internal/infrastructure/logger"
	"github.com/spf13/afero"
)

func main() {
	log := logger.New()
	defer log.Sync()

	fs := afero.NewOsFs()
	cons := console.NewConsole(os.Stdout)
	compatibility := compatibility.NewCompatibility(cons)
	conf, err := configuration.New(fs, compatibility)
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
