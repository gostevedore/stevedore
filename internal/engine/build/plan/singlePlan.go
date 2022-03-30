package plan

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
)

// SinglePlan is a build plan
type SinglePlan struct {
	BasePlan
}

// NewSinglePlan returns a new SinglePlan
func NewSinglePlan(imagesStorer ImagesStorer) *SinglePlan {
	return &SinglePlan{
		BasePlan{
			imagesStorer,
		},
	}
}

// Plan return a list of images to build
func (p *SinglePlan) Plan(name string, versions []string) ([]*Step, error) {
	var steps []*Step
	var images []*image.Image
	var err error

	errContext := "(plan::Simple::Plan)"

	images, err = p.findImages(name, versions)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	for _, image := range images {
		steps = append(steps, NewStep(image, image.Name, nil))
	}

	return steps, nil
}
