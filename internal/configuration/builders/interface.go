package builders

import "github.com/gostevedore/stevedore/internal/builders/builder"

// BuildersStorer interface
type BuildersStorer interface {
	AddBuilder(builder *builder.Builder) error
}
