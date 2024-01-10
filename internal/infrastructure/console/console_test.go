package console

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/creack/pty"
	"github.com/go-errors/errors"
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
func TestReadPassword(t *testing.T) {
	var err error
	var rbuff bytes.Buffer

	// Create a new pseudo-terminal.
	terminal, tty, err := pty.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer terminal.Close()
	defer tty.Close()

	_, err = terminal.Write([]byte("password\n"))
	if err != nil {
		t.Error(fmt.Sprintf("[%s] error writing password to test terminal.", t.Name()))
	}

	tests := []struct {
		desc   string
		reader io.Reader
		writer io.Writer
		prompt string
		res    string
		err    error
	}{
		{
			desc:   "Testing error when providing inappropriate ioctl for device to read password",
			prompt: "Password: ",
			res:    "password",
			reader: io.Reader(&rbuff),
			writer: io.Discard,
			err:    fmt.Errorf("inappropriate ioctl for device."),
		},
		{
			desc:   "Testing read password",
			prompt: "Password: ",
			res:    "password",
			reader: tty,
			writer: io.Discard,
			err:    &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Log(test.desc)

		c := &Console{
			read:  test.reader,
			write: test.writer,
		}

		res, err := c.ReadPassword(test.prompt)
		if err != nil {
			assert.Equal(t, test.err, err)
		} else {
			t.Log(res)
			assert.Equal(t, test.res, res)
		}
	}
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
