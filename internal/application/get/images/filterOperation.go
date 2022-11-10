package images

import (
	"fmt"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
)

const (
	EQ = "="
)

// filterOperation hold the filter operation tokens
type filterOperation struct {
	attribute string
	operation string
	item      interface{}
}

func NewFilterOperation(filter string) filterOperation {
	fOp, _ := ParseFilterOpration(filter)
	return fOp
}

func (f filterOperation) IsDefined() bool {
	return f != filterOperation{}
}

func ParseFilterOpration(filter string) (filterOperation, error) {
	var fOp filterOperation
	var err error

	errContext := "(application::get:immages::parseFilterOpration)"

	fOp, err = parseEqualFilterOpration(filter)
	if err != nil {
		return filterOperation{}, errors.New(errContext, "", err)
	}
	if fOp.operation == EQ {
		return fOp, nil
	}

	return filterOperation{}, nil
}

func parseEqualFilterOpration(filter string) (filterOperation, error) {

	errContext := "(application::get:immages::parseEqualFilterOpration)"
	tokens := strings.Split(filter, EQ)

	if len(tokens) > 2 {
		return filterOperation{}, errors.New(errContext, fmt.Sprintf("Invalid filter '%s'", filter))
	}

	if len(tokens) < 2 {
		return filterOperation{}, nil
	}

	return filterOperation{
		attribute: tokens[0],
		operation: EQ,
		item:      tokens[1],
	}, nil
}
