package operation

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

func TestParseFilterOpration(t *testing.T) {

	errContext := "(filters::operation::ParseFilterOpration)"
	errContextEqualFilter := "(filters::operation::parseEqualFilterOpration)"

	tests := []struct {
		desc      string
		filter    string
		operation *FilterOperation
		res       *FilterOperation
		err       error
	}{
		{
			desc:      "Testing error on parse filter operation when filter operation is nil",
			operation: nil,
			err:       errors.New(errContext, "Filter operations is not not initialized"),
		},
		{
			desc:      "Testing parse filter operation with an equality",
			filter:    "a=b",
			operation: NewFilterOperation(),
			res: &FilterOperation{
				attribute: "a",
				operation: EQ,
				item:      "b",
			},
			err: &errors.Error{},
		},
		{
			desc:      "Testing error parsing filter operation with an invalid equality",
			filter:    "a=b=c",
			operation: NewFilterOperation(),
			res:       NewFilterOperation(),
			err: errors.New(
				errContext,
				"",
				errors.New(errContextEqualFilter, "Invalid filter 'a=b=c'"),
			),
		},
		{
			desc:      "Testing parse filter operation with unmanaged operation",
			filter:    "a>b",
			operation: NewFilterOperation(),
			res:       NewFilterOperation(),
			err:       &errors.Error{},
		},
	}

	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			err := test.operation.ParseFilterOpration(test.filter)

			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, test.operation)
			}
		})
	}
}

func TestIsDefined(t *testing.T) {

	tests := []struct {
		desc string
		op   *FilterOperation
		res  bool
	}{
		{
			desc: "Testing whether filter operation is not defined with nil filter operation",
			op:   nil,
			res:  false,
		},
		{
			desc: "Testing whether filter operation is not defined",
			op:   NewFilterOperation(),
			res:  false,
		},
		{
			desc: "Testing whether filter operation is defined",
			op: &FilterOperation{
				attribute: "name",
				operation: EQ,
				item:      "value",
			},
			res: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			res := test.op.IsDefined()
			assert.Equal(t, test.res, res)
		})
	}
}
