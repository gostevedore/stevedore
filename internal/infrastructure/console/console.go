package console

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/ryanuber/columnize"
	"golang.org/x/term"
)

const (
	columnSeparator = "|"
	columnGlue      = " "
	columnPrefix    = ""

	resetColor  = "\033[0m"
	redColor    = "\033[31m"
	greenColor  = "\033[32m"
	yellowColor = "\033[33m"
	blueColor   = "\033[34m"
	purpleColor = "\033[35m"
	cyanColor   = "\033[36m"
	whiteColor  = "\033[37m"
)

// Console
type Console struct {
	write io.Writer
	read  io.Reader
}

// NewConsole creates a new console
func NewConsole(w io.Writer, r io.Reader) *Console {
	return &Console{
		write: w,
		read:  r,
	}
}

// Write
func (c *Console) Write(data []byte) (int, error) {

	if c.write == nil {
		c.write = os.Stdout
	}

	size := len(data)
	// if size > 0 && data[size-1] == '\n' {
	// data = data[:size-1]
	// }
	fmt.Fprint(c.write, string(data))

	return size, nil
}

// Print
func (c *Console) Print(msg ...interface{}) {
	if c.write == nil {
		c.write = os.Stdout
	}

	fmt.Fprintln(c.write, msg...)
}

// PrintTable prints a table
func (c *Console) PrintTable(content [][]string) error {

	if c.write == nil {
		c.write = os.Stdout
	}

	table := []string{}
	config := columnize.DefaultConfig()
	config.Delim = columnSeparator
	config.Glue = columnGlue
	config.Prefix = columnPrefix

	//	table = append(table, columnizeLine(header))

	for _, row := range content {
		table = append(table, columnizeLine(row))
	}

	fmt.Fprintf(c.write, "%s\n", columnize.Format(table, config))

	return nil
}

// message prints a message using the default color
func (c *Console) message(msg ...interface{}) {
	c.Print(resetColor, fmt.Sprint(msg...), resetColor)
}

// blue prints a message in blue color
func (c *Console) blue(msg ...interface{}) {
	c.Print(blueColor, fmt.Sprint(msg...), resetColor)
}

// green prints a message in green color
func (c *Console) green(msg ...interface{}) {
	c.Print(greenColor, fmt.Sprint(msg...), resetColor)
}

// purple prints a message in purple color
func (c *Console) purple(msg ...interface{}) {
	c.Print(purpleColor, fmt.Sprint(msg...), resetColor)
}

// red prints a message in red color
func (c *Console) red(msg ...interface{}) {
	c.Print(redColor, fmt.Sprint(msg...), resetColor)
}

// Info prints a info message
func (c *Console) Info(msg ...interface{}) {
	c.message(msg...)
}

// Warn prints a warning message
func (c *Console) Warn(msg ...interface{}) {
	c.purple(msg...)
}

// Error prints an error message
func (c *Console) Error(msg ...interface{}) {
	c.red(msg...)
}

// Debug prints a debug message
func (c *Console) Debug(msg ...interface{}) {
	c.blue(msg...)
}

func columnizeLine(items []string) string {
	var line string

	if len(items) > 0 {
		line = items[0]
		items = items[1:]
	}

	for _, item := range items {
		line = line + columnSeparator + item
	}
	return line
}

// Read read a line from console reader
func (c *Console) Read() string {
	var input string

	fmt.Fscanln(c.read, &input)
	//fmt.Fscanf(c.read, "%s", &input)
	return input
}

// ReadPassword read a line from console reader without echo
func (c *Console) ReadPassword(prompt string) (string, error) {
	stdin := int(syscall.Stdin)
	oldState, err := term.GetState(stdin)
	if err != nil {
		return "", err
	}
	defer term.Restore(stdin, oldState)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		for range sigch {
			c.Warn("Interrupted")
			term.Restore(stdin, oldState)
			os.Exit(0)
		}
	}()

	c.Write([]byte(prompt))
	password, err := term.ReadPassword(stdin)
	if err != nil {
		return "", err
	}
	return string(password), nil
}
