package plan

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// CascadePlan is the plan used to cascade build
type CascadePlan struct {
	BasePlan
	depth int
}

// NewCascadePlan creates a new CascadePlan
func NewCascadePlan(imagesStorer repository.ImagesStorerReader, depth int) *CascadePlan {
	return &CascadePlan{
		BasePlan{
			imagesStorer,
		},
		depth,
	}
}

// Plan return a list of images to build
func (p *CascadePlan) Plan(name string, versions []string) ([]*Step, error) {

	var images []*image.Image
	var err error

	errContext := "(plan::Cascade::Plan)"
	steps := []*Step{}

	images, err = p.findImages(name, versions)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	for _, image := range images {
		steps = append(steps, p.plan(image, nil, p.depth)...)
	}

	return steps, nil
}

func (p *CascadePlan) plan(image *image.Image, parent *Step, depth int) []*Step {
	steps := []*Step{}
	var sync chan struct{}

	// root images does not require to sync
	if parent != nil {
		sync = make(chan struct{})
		parent.Subscribe(sync)
	}

	// not tested
	if p.images.IsWildcard(image) {
		return steps
	}

	step := NewStep(image, image.Name, sync)
	steps = append(steps, step)

	if depth == 0 {
		return steps
	}

	for _, child := range image.Children {
		steps = append(steps, p.plan(child, step, depth-1)...)
	}

	return steps
}
