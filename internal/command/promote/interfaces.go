package promote

import (
	"context"

	handler "github.com/gostevedore/stevedore/internal/command/promote/handler"
)

// HandlerPromoter
type HandlerPromoter interface {
	Handler(ctx context.Context, options *handler.HandlerOptions) error
}
