package graph

// GraphTemplateNoder is the interface for the graph template node
type GraphTemplateNoder interface {
	AddChild(GraphTemplateNoder) error
	AddParent(GraphTemplateNoder) error
	AddItem(interface{})
	Name() string
	Item() interface{}
	Children() []GraphTemplateNoder
	Parents() []GraphTemplateNoder
}

// GraphTemplater is the interface for the graph template for images
type GraphTemplater interface {
	GetNode(string) GraphTemplateNoder
	AddNode(GraphTemplateNoder) error
	AddRelationship(GraphTemplateNoder, GraphTemplateNoder) error
	Exists(string) bool
	HasCycles() bool
	Iterate() <-chan GraphTemplateNoder
}
