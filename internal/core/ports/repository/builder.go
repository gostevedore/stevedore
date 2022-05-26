package repository

import "github.com/gostevedore/stevedore/internal/core/domain/builder"

// BuildersStorer interface
type BuildersStorer interface {
	Store(builder *builder.Builder) error
	Find(name string) (*builder.Builder, error)
}

// BuildersFilterer is an interface for filtering builders content output
type BuildersFilterer interface {
	All() []*builder.Builder
	FilterByName(string) *builder.Builder
	FilterByDriver(string) []*builder.Builder
}

// BuildersPrinter is an interface for printing builders content output.
type BuildersPrinter interface {
	PrintTable(content [][]string) error
}
