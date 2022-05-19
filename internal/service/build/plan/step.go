package plan

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// Step is a plan step
type Step struct {
	// image is the image to build
	image *image.Image
	// description is the description of the step
	description string
	// sync
	sync chan struct{}
	// subscriptions is a list of channels to sync to children images
	subscriptions []chan struct{}
}

// NewStep returns a new instance of the Step
func NewStep(image *image.Image, desc string, sync chan struct{}) *Step {
	return &Step{
		image:         image,
		description:   desc,
		sync:          sync,
		subscriptions: []chan struct{}{},
	}
}

// Image returns the image to build
func (p *Step) Image() *image.Image {
	return p.image
}

// Subscribe adds a channel to the list of channels to notify
func (p *Step) Subscribe(sync chan struct{}) error {

	errContext := "(Step::Subscribe)"
	if p.subscriptions == nil {
		return errors.New(errContext, "Subscribtions list is nil")
	}

	if sync == nil {
		return errors.New(errContext, "Sync channel is nil")
	}

	p.subscriptions = append(p.subscriptions, sync)

	return nil
}

// Wait blocks Step until it is notified
func (p *Step) Wait() {
	if p.sync != nil {
		<-p.sync
	}
}

// Notify notifies the notify channels
func (p *Step) Notify() {
	for _, subscrition := range p.subscriptions {
		close(subscrition)
	}
}
