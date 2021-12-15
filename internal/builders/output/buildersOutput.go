package output

// BuildersOutput is an output for the builders
type BuildersOutput struct {
	BuildersFilterer
	Write BuildersPrinter
}

// NewBuildersOutput creates a new BuildersOutput
func NewBuildersOutput(write BuildersPrinter, builders BuildersFilterer) *BuildersOutput {
	return &BuildersOutput{
		builders, write,
	}
}

// builderHeader returns the header for the builders
func builderHeader() []string {
	return []string{"NAME", "DRIVER"}
}

// PrintAll prints all the builders
func (o *BuildersOutput) PrintAll() {

	content := [][]string{}
	content = append(content, builderHeader())

	for _, builder := range o.All() {
		builderContent := []string{builder.Name, builder.Driver}
		content = append(content, builderContent)
	}

	o.Write.PrintTable(content)
}

// func (o *BuildersOutput) Filter() *builder.Builder {
// 	return nil
// }
