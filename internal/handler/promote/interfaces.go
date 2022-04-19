package handler

import (
	"context"

	"github.com/gostevedore/stevedore/internal/service/promote"
)

// ServicePromoter
type ServicePromoter interface {
	Promote(ctx context.Context, options *promote.ServiceOptions) error
}
