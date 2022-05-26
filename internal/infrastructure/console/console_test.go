package console

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWrite tests function Write
func TestWrite(t *testing.T) {
	var w bytes.Buffer

	expected := `I'm a test`
	data := []byte(expected)
	c := &Console{
		write: io.Writer(&w),
	}
	c.Write(data)

	t.Log("Testing Write a message")
	assert.Equal(t, expected, w.String())

}

// TestColumnizeLine tests function TestColumnizeLine
func TestColumnizeLine(t *testing.T) {
	tests := []struct {
		desc  string
		items []string
		res   string
	}{
		{
			desc:  "Testing nil items input",
			items: nil,
			res:   "",
		},
		{
			desc:  "Testing one items input",
			items: []string{"one"},
			res:   "one",
		},
		{
			desc:  "Testing two items input",
			items: []string{"one", "two"},
			res:   "one" + columnSeparator + "two",
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		columnized := columnizeLine(test.items)
		assert.Equal(t, test.res, columnized)

	}
}
