package graph

// GraphTemplateFactory is a factory for the graph template
type GraphTemplateFactory struct {
	isMock bool
}

// NewGraphTemplateFactory creates a new graph template factory
func NewGraphTemplateFactory(isMock bool) *GraphTemplateFactory {
	return &GraphTemplateFactory{
		isMock: isMock,
	}
}

// NewGraphTemplate creates a new graph template node
func (f *GraphTemplateFactory) NewGraphTemplate() GraphTemplater {

	if f.isMock {
		return NewMockGraphTemplate()
	}

	return NewGraphTemplate()
}

// NewGraphTemplateNode creates a new graph template node
func (f *GraphTemplateFactory) NewGraphTemplateNode(name string) GraphTemplateNoder {
	return NewGraphTemplateNode(name)
}
