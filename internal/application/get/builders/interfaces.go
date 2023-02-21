package builders

import (
	operationfilter "github.com/gostevedore/stevedore/internal/infrastructure/filters/operation"
)

type FilterFactorier interface {
	FilterOperation() *operationfilter.FilterOperation
}

type FilterOperationer interface {
	ParseFilterOpration(filter string) error
	Attribute() string
	Operation() string
	Item() interface{}
	IsDefined() bool
}
