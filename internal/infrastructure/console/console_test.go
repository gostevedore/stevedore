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

// TestRead tests function Read
func TestRead(t *testing.T) {

	tests := []struct {
		desc string
		buff *bytes.Buffer
		res  string
	}{
		{
			desc: "Testing read empty line",
			buff: &bytes.Buffer{},
			res:  "",
		},
		{
			desc: "Testing read a line",
			buff: bytes.NewBuffer([]byte("word")),
			res:  "word",
		},
	}

	for _, test := range tests {
		t.Log(test.desc)
		c := &Console{
			read: io.Reader(test.buff),
		}
		res := c.Read()
		assert.Equal(t, test.res, res)
	}

}

// TestReadPassword tests function Read
// func TestReadPassword(t *testing.T) {
// 	var wbuff bytes.Buffer
// 	var rbuff bytes.Buffer

// 	tests := []struct {
// 		desc   string
// 		prompt string
// 		res    string
// 		err    error
// 	}{
// 		{
// 			desc:   "Testing read password",
// 			prompt: "Password: ",

// 			res: "password",
// 			err: &errors.Error{},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Log(test.desc)

// 		rbuff.WriteString("password\r")
// 		console := Console{
// 			read:  io.Reader(&rbuff),
// 			write: io.Writer(&wbuff),
// 		}

// 		res, err := console.ReadPassword(test.prompt)
// 		if err != nil {
// 			assert.Equal(t, test.err.Error(), err.Error())
// 		} else {
// 			assert.Equal(t, test.res, res)
// 		}
// 	}
// }

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
