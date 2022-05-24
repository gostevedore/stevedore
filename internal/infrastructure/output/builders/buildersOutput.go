package builders

// Output is an output for the builders
type Output struct {
	BuildersFilterer
	Write BuildersPrinter
}

// NewOutput creates a new Output
func NewOutput(write BuildersPrinter, builders BuildersFilterer) *Output {
	return &Output{
		builders, write,
	}
}

// builderHeader returns the header for the builders
func builderHeader() []string {
	return []string{"NAME", "DRIVER"}
}

// PrintAll prints all the builders
func (o *Output) PrintAll() {

	content := [][]string{}
	content = append(content, builderHeader())

	for _, builder := range o.All() {
		builderContent := []string{builder.Name, builder.Driver}
		content = append(content, builderContent)
	}

	o.Write.PrintTable(content)
}

// func (o *Output) Filter() *builder.Builder {
// 	return nil
// }
