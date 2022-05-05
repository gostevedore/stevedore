package main

import (
	"context"
	"os"

	"github.com/gostevedore/stevedore/internal/cli/stevedore"
	"github.com/gostevedore/stevedore/internal/compatibility"
	"github.com/gostevedore/stevedore/internal/configuration"
	"github.com/gostevedore/stevedore/internal/ui/console"
	"github.com/spf13/afero"
)

func main() {

	fs := afero.NewOsFs()
	cons := console.NewConsole(os.Stdout)
	compatibility := compatibility.NewCompatibility(cons)
	conf, err := configuration.New(fs, compatibility)
	if err != nil {
		cons.Error(err.Error())
		os.Exit(1)
	}

	stevedore := stevedore.NewCommand(context.Background(), fs, compatibility, compatibility, cons, conf)
	err = stevedore.Execute()
	if err != nil {
		os.Exit(1)
	}
}
