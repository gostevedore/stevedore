package middleware

import (
	"context"

	"github.com/gostevedore/stevedore/internal/cli/command"
	"github.com/gostevedore/stevedore/internal/cli/command/middleware/audit"
	"github.com/gostevedore/stevedore/internal/cli/command/middleware/compatibility"
	"github.com/gostevedore/stevedore/internal/cli/command/middleware/error"
	"github.com/gostevedore/stevedore/internal/cli/command/middleware/interruption"
	"github.com/gostevedore/stevedore/internal/cli/command/middleware/panic"
)

func Command(ctx context.Context, c *command.StevedoreCommand, compatibilityReport CompatibilityReporter, log Logger, cons Consoler) *command.StevedoreCommand {

	cmd := compatibility.NewCommand(c, compatibilityReport)
	cmd = audit.NewCommand(cmd, log)
	cmd = panic.NewCommand(cmd, cons, log)
	cmd = error.NewCommand(cmd, cons, log)
	cmd = interruption.NewCommand(ctx, cmd, cons, log)

	return cmd
}
