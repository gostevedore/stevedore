package credentials

import "github.com/gostevedore/stevedore/internal/core/domain/credentials"

type Outputter interface {
	Output(badge *credentials.Badge) (string, string, error)
}

type OutputWriter interface {
	PrintTable(content [][]string) error
}
