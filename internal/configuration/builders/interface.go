package builders

import "github.com/gostevedore/stevedore/internal/core/domain/builder"

// BuildersStorer interface
type BuildersStorer interface {
	Store(builder *builder.Builder) error
}
