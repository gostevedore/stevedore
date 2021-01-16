package dryrunpromoter

import (
	"context"
	"fmt"

	"github.com/gostevedore/stevedore/internal/types"
	"github.com/gostevedore/stevedore/internal/ui/console"
)

func Promote(ctx context.Context, options *types.PromoteOptions) error {
	fmt.Fprintln(console.GetConsole(), fmt.Sprintf("%+v", *options))
	return nil
}
