package promote

import (
	"context"
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	dockerpromoter "github.com/gostevedore/stevedore/internal/promote/promoter/docker"
	dryrunpromoter "github.com/gostevedore/stevedore/internal/promote/promoter/dryrun"
	"github.com/gostevedore/stevedore/internal/types"
)

func Promote(ctx context.Context, options *types.PromoteOptions) error {
	var err error

	if options.DryRun {
		err = dryrunpromoter.Promote(ctx, options)
	} else {
		err = dockerpromoter.Promote(ctx, options)
	}

	if err != nil {
		return errors.New("(promote::Promote)", fmt.Sprintf("Error promoting '%s' image", options.ImageName), err)
	}

	return nil
}
