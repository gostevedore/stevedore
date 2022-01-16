package console

import (
	"fmt"
	"io"
	"os"

	"github.com/ryanuber/columnize"
)

const (
	columnSeparator = "|"
	columnGlue      = " "
	columnPrefix    = ""

	resetColor = "\033[0m"
	red        = "\033[31m"
	green      = "\033[32m"
	yellow     = "\033[33m"
	blue       = "\033[34m"
	purple     = "\033[35m"
	cyan       = "\033[36m"
	white      = "\033[37m"
)

var ui *Console

// Console
type Console struct {
	write io.Writer
}

// NewConsole creates a new console
func NewConsole(w io.Writer) *Console {
	return &Console{
		write: w,
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

// Init initialzes the ui console
func Init(w io.Writer) {
	if ui == nil {
		ui = &Console{w}
	}
}

// GetConsole
func GetConsole() io.Writer {
	if ui == nil {
		ui = &Console{
			write: os.Stdout,
		}
	}

	return ui
}

// SetWriter defines a writer to ui console
func SetWriter(w io.Writer) {
	ui = &Console{w}
}

// Print
func Print(msg ...interface{}) {
	if ui == nil {
		Init(os.Stdout)
	}

	fmt.Fprintln(ui.write, msg...)
}

// Blue prints a message in blue color
func Blue(msg ...interface{}) {
	Print(blue, msg, resetColor)
}

// Green prints a message in green color
func Green(msg ...interface{}) {
	Print(green, msg, resetColor)
}

// Purple prints a message in purple color
func Purple(msg ...interface{}) {
	Print(purple, msg, resetColor)
}

// Red prints a message in red color
func Red(msg ...interface{}) {
	Print(red, msg, resetColor)
}

// ColorPrint prints a message on the specified color
func ColorPrint(color string, msg interface{}) {
	Print(color, msg, resetColor)
}

// Info
func Info(msg ...interface{}) {
	Print(msg...)
}

// Warn
func Warn(msg ...interface{}) {
	Purple(msg...)
}

// Error
func Error(msg ...interface{}) {
	Red(msg...)
}

// Debug
func Debug(msg ...interface{}) {
	Blue(msg...)
}

// PrintTable
func PrintTable(content [][]string) error {

	if ui == nil {
		Init(os.Stdout)
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

	fmt.Fprintf(ui.write, "%s\n", columnize.Format(table, config))

	return nil
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
