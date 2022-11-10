package images

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestParseFilterOpration(t *testing.T) {

	errContext := "(application::get:immages::parseFilterOpration)"
	errContextEqualFilter := "(application::get:immages::parseEqualFilterOpration)"

	tests := []struct {
		desc   string
		filter string
		res    filterOperation
		err    error
	}{
		{
			desc:   "Testing parse filter operation with an equality",
			filter: "a=b",
			res: filterOperation{
				attribute: "a",
				operation: EQ,
				item:      "b",
			},
			err: &errors.Error{},
		},
		{
			desc:   "Testing error parsing filter operation with an invalid equality",
			filter: "a=b=c",
			res:    filterOperation{},
			err: errors.New(
				errContext,
				"",
				errors.New(errContextEqualFilter, "Invalid filter 'a=b=c'"),
			),
		},
		{
			desc:   "Testing parse filter operation with unmanaged operation",
			filter: "a>b",
			res:    filterOperation{},
			err:    &errors.Error{},
		},
	}

	for _, test := range tests {
		res, err := ParseFilterOpration(test.filter)
		if err != nil {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.res, res)
		}
	}
}

func TestIsDefined(t *testing.T) {

	tests := []struct {
		desc string
		op   filterOperation
		res  bool
	}{
		{
			desc: "Testing whether get image filter operation is not defined",
			op:   filterOperation{},
			res:  false,
		},
		{
			desc: "Testing whether get image filter operation is defined",
			op: filterOperation{
				attribute: "name",
				operation: EQ,
				item:      "value",
			},
			res: true,
		},
	}

	for _, test := range tests {
		res := test.op.IsDefined()
		assert.Equal(t, test.res, res)
	}
}
