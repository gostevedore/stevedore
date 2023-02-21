package stevedore

import (
	"context"
	"testing"

	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/gostevedore/stevedore/internal/infrastructure/console"
	"github.com/gostevedore/stevedore/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestStevedoreRun(t *testing.T) {
	log := logger.New()
	cons := console.NewMockConsole()
	config := configuration.DefaultConfig()
	compatibility := compatibility.NewMockCompatibility()
	fs := afero.NewMemMapFs()
	args := []string{}

	t.Run("Testing stevedore Run func", func(t *testing.T) {
		cmd := NewCommand(context.TODO(), fs, compatibility, compatibility, cons, log, config)
		cmd.Command.ParseFlags(args)
		cmd.Command.Run(cmd.Command, args)
	})
}

func TestStevedorePersistentPreRunE(t *testing.T) {
	log := logger.New()
	cons := console.NewMockConsole()
	config := configuration.DefaultConfig()
	compatibility := compatibility.NewMockCompatibility()
	fs := afero.NewMemMapFs()

	args := []string{}

	t.Run("Testing stevedore PersistentPreRunE function", func(t *testing.T) {
		cmd := NewCommand(context.TODO(), fs, compatibility, compatibility, cons, log, config)
		cmd.Command.ParseFlags(args)
		err := cmd.Command.PersistentPreRunE(cmd.Command, args)
		assert.NoError(t, err)
	})
}
