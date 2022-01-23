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

// // graphTemplateNoder is the interface for the graph template node
// type graphTemplateNoder interface {
// 	GraphTemplateNoder
// 	getNode() *gdsexttree.Node
// }

// GraphTemplater is the interface for the graph template for images
type GraphTemplater interface {
	GetNode(string) GraphTemplateNoder
	AddNode(GraphTemplateNoder) error
	AddRelationship(GraphTemplateNoder, GraphTemplateNoder) error
	Exists(string) bool
	HasCycles() bool
	Iterate() <-chan GraphTemplateNoder
}

// // Grapher is a graph template for images
// type Grapher interface {
// 	AddNode(*gdsexttree.Node) error
// 	AddRelationship(*gdsexttree.Node, *gdsexttree.Node) error
// 	GetNode(string) *gdsexttree.Node
// }

// // Noder is a node for the graph template
// type Noder interface {
// 	AddChild(*gdsexttree.Node) error
// 	AddParent(*gdsexttree.Node) error
// }
