package plan

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
)

// CascadePlan is the plan used to cascade build
type CascadePlan struct {
	BasePlan
	depth int
}

// NewCascadePlan creates a new CascadePlan
func NewCascadePlan(imagesStorer ImagesStorer, depth int) *CascadePlan {
	return &CascadePlan{
		BasePlan{
			imagesStorer,
		},
		depth,
	}
}

// Plan return a list of images to build
func (p *CascadePlan) Plan(name string, versions []string) ([]*Step, error) {

	var steps []*Step
	var images []*image.Image
	var err error

	errContext := "(plan::Cascade::Plan)"
	_ = errContext

	images, err = p.findImages(name, versions)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	for _, image := range images {
		steps = append(steps, plan(image, nil, p.depth)...)
	}

	return steps, nil
}

func plan(image *image.Image, parent *Step, depth int) []*Step {
	steps := []*Step{}
	var sync chan struct{}

	if depth == 0 {
		return steps
	}

	// root images does not require to sync
	if parent != nil {
		sync = make(chan struct{})
		parent.Subscribe(sync)
	}

	step := NewStep(image, image.Name, sync)
	steps = append(steps, step)

	for _, child := range image.Children {
		steps = append(steps, plan(child, step, depth-1)...)
	}

	return steps
}
