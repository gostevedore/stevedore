package repository

import "github.com/gostevedore/stevedore/internal/core/domain/builder"

// BuildersStorer interface
type BuildersStorer interface {
	BuildersStoreWriter
	BuildersStorerReader
}

// BuildersStoreWriter interface defines which methods are needed to save images to an images store
type BuildersStoreWriter interface {
	Store(builder *builder.Builder) error
}

// BuildersStorerReader interface defines which methods are needed to read images from an images store
type BuildersStorerReader interface {
	Find(name string) (*builder.Builder, error)
	List() ([]*builder.Builder, error)
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

// BuildersSelector interface defines which methods are needed to select builders from a list
type BuildersSelector interface {
	Select(builders []*builder.Builder, operation string, item string) ([]*builder.Builder, error)
}

// BuildersOutputter interface defines which methods are needed to return buidlers definitions
type BuildersOutputter interface {
	Output(list []*builder.Builder) error
}
