package graph

// GraphTemplateNodeFactory is a factory for the graph template
type GraphTemplateNodeFactory struct{}

// NewGraphTemplateNodeFactory creates a new graph template factory
func NewGraphTemplateNodeFactory() *GraphTemplateNodeFactory {
	return &GraphTemplateNodeFactory{}
}

// NewGraphTemplate creates a new graph template node
func (f *GraphTemplateNodeFactory) NewGraphTemplateNode(isMock bool) GraphTemplater {

	if isMock {
		return NewMockGraphTemplate()
	}

	return NewGraphTemplate()
}
