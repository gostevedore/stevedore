package common

// import (
// 	"testing"

// 	"gotest.tools/assert"
// )

// func TestSanitizeTag(t *testing.T) {
// 	tests := []struct {
// 		desc  string
// 		input string
// 		res   string
// 	}{
// 		{
// 			desc:  "Testing sanitize with no changes to apply",
// 			input: "no-changes",
// 			res:   "no-changes",
// 		},
// 		{
// 			desc:  "Testing sanitize /",
// 			input: "no/changes",
// 			res:   "no_changes",
// 		},
// 		{
// 			desc:  "Testing sanitize :",
// 			input: "no:changes",
// 			res:   "no_changes",
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			res := SanitizeTag(test.input)
// 			assert.Equal(t, test.res, res, "Unexpected value")
// 		})
// 	}
// }
