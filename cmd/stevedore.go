package main

import (
	"context"
	"os"
	"stevedore/internal/command/stevedore"
)

func main() {

	stevedore := stevedore.NewCommand(context.TODO(), nil)
	err := stevedore.Execute()
	if err != nil {
		os.Exit(1)
	}
}
