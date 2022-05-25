package promote

import (
	"context"

	"github.com/gostevedore/stevedore/internal/application/promote"
)

// ServicePromoter
type ServicePromoter interface {
	Promote(ctx context.Context, options *promote.ServiceOptions) error
}
