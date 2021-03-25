package dryrunpromoter

import (
	"context"
	"fmt"

	"github.com/gostevedore/stevedore/internal/types"
	"github.com/gostevedore/stevedore/internal/ui/console"
)

func Promote(ctx context.Context, options *types.PromoteOptions) error {
	fmt.Fprint(console.GetConsole(), options.String())
	return nil
}
