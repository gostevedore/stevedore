package builders

import "github.com/gostevedore/stevedore/internal/core/domain/builder"

// BuildersPrinter is an interface for printing builders content output.
type BuildersPrinter interface {
	PrintTable(content [][]string) error
}

// BuildersFilterer is an interface for filtering builders content output
type BuildersFilterer interface {
	All() []*builder.Builder
	FilterByName(string) *builder.Builder
	FilterByDriver(string) []*builder.Builder
}
