package repository

import "github.com/gostevedore/stevedore/internal/core/domain/release"

// ReleasePrinter is an interface for printing release information
type ReleasePrinter interface {
	Print(release *release.Release) error
}
