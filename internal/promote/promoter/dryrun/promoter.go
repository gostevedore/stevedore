package dryrunpromoter

import (
	"context"
	"fmt"
	"stevedore/internal/types"
	"stevedore/internal/ui/console"
)

func Promote(ctx context.Context, options *types.PromoteOptions) error {
	fmt.Fprintln(console.GetConsole(), fmt.Sprintf("%+v", *options))
	return nil
}
