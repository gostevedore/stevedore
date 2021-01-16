package main

import (
	"context"
	"os"

	"github.com/gostevedore/stevedore/internal/command/stevedore"
)

func main() {

	stevedore := stevedore.NewCommand(context.TODO(), nil)
	err := stevedore.Execute()
	if err != nil {
		os.Exit(1)
	}
}
