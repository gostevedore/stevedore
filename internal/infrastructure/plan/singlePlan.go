package plan

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// SinglePlan is a build plan
type SinglePlan struct {
	BasePlan
}

// NewSinglePlan returns a new SinglePlan
func NewSinglePlan(imagesStorer repository.ImagesStorerReader) *SinglePlan {
	return &SinglePlan{
		BasePlan{
			imagesStorer,
		},
	}
}

// Plan return a list of images to build
func (p *SinglePlan) Plan(name string, versions []string) ([]*Step, error) {
	var images []*image.Image
	var err error

	errContext := "(plan::Simple::Plan)"
	steps := []*Step{}

	images, err = p.findImages(name, versions)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	for _, image := range images {
		steps = append(steps, NewStep(image, image.Name, nil))
	}

	return steps, nil
}
