package builders

import "github.com/gostevedore/stevedore/internal/builders/builder"

// BuildersStorer interface
type BuildersStorer interface {
	Store(builder *builder.Builder) error
}
