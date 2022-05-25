package promote

import (
	"context"

	"github.com/gostevedore/stevedore/internal/application/promote"
)

// PromoteApplication
type PromoteApplication interface {
	Promote(ctx context.Context, options *promote.Options) error
}
