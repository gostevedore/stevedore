package promote

import (
	"context"

	"github.com/gostevedore/stevedore/internal/engine/promote"
)

// ServicePromoter
type ServicePromoter interface {
	Promote(ctx context.Context, options *promote.ServiceOptions, promoteType string) error
}
