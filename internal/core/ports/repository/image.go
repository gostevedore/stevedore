package repository

import (
	"context"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// BuildDriverer interface defines which methods are used to build a docker image
type BuildDriverer interface {
	Build(context.Context, *image.Image, *image.BuildDriverOptions) error
}

// Promoter
type Promoter interface {
	Promote(context.Context, *image.PromoteOptions) error
}
