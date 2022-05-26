package plan

// Planner interfaces defines the storage of images
type Planner interface {
	Plan(name string, versions []string) ([]*Step, error)
}
