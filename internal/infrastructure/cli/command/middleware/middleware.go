package middleware

import (
	"context"

	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command/middleware/audit"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command/middleware/compatibility"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command/middleware/error"
	"github.com/gostevedore/stevedore/internal/infrastructure/cli/command/middleware/interruption"
)

func Command(ctx context.Context, c *command.StevedoreCommand, compatibilityReport CompatibilityReporter, log Logger, cons Consoler) *command.StevedoreCommand {

	cmd := compatibility.NewCommand(c, compatibilityReport)
	cmd = audit.NewCommand(cmd, log)
	cmd = error.NewCommand(cmd, cons, log)
	cmd = interruption.NewCommand(ctx, cmd, cons, log)

	return cmd
}
