package operation

import (
	"fmt"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
)

const (
	EQ = "="
)

// FilterOperation hold the filter operation tokens
type FilterOperation struct {
	attribute string
	operation string
	item      interface{}
}

func NewFilterOperation() *FilterOperation {
	fOp := new(FilterOperation)
	return fOp
}

func (f *FilterOperation) IsDefined() bool {
	if f == nil {
		return false
	}

	return *f != FilterOperation{}
}

func (f *FilterOperation) Attribute() string {
	return f.attribute
}

func (f *FilterOperation) Operation() string {
	return f.operation
}

func (f *FilterOperation) Item() interface{} {
	return f.item
}

func (f *FilterOperation) ParseFilterOpration(filter string) error {
	var fOp *FilterOperation
	var err error

	errContext := "(filters::operation::ParseFilterOpration)"

	if f == nil {
		return errors.New(errContext, "Filter operations is not not initialized")
	}

	fOp, err = parseEqualFilterOpration(filter)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	if fOp != nil {
		f.attribute = fOp.attribute
		f.item = fOp.item
		f.operation = fOp.operation

		return nil
	}

	return nil
}

func parseEqualFilterOpration(filter string) (*FilterOperation, error) {

	errContext := "(filters::operation::parseEqualFilterOpration)"
	tokens := strings.Split(filter, EQ)

	if len(tokens) > 2 {
		return nil, errors.New(errContext, fmt.Sprintf("Invalid filter '%s'", filter))
	}

	if len(tokens) < 2 {
		return nil, nil
	}

	return &FilterOperation{
		attribute: tokens[0],
		operation: EQ,
		item:      tokens[1],
	}, nil
}
