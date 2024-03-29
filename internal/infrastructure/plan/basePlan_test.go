package plan

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/stretchr/testify/assert"
)

func TestFindImages(t *testing.T) {

	errContext := "(plan::BasePlan::images)"
	_ = errContext

	tests := []struct {
		desc              string
		plan              *BasePlan
		name              string
		versions          []string
		prepareAssertFunc func(*BasePlan)
		assertFunc        func(*BasePlan) bool
		err               error
	}{
		{
			desc: "Testing error when images storer is nil",
			plan: &BasePlan{},
			err:  errors.New(errContext, "Images storer is nil"),
		},
		{
			desc: "Testing generate plan with an image name and versions",
			plan: &BasePlan{
				images: images.NewMockStore(),
			},
			name:     "image",
			versions: []string{"version1", "version2"},
			err:      &errors.Error{},
			prepareAssertFunc: func(p *BasePlan) {
				p.images.(*images.MockStore).On("FindGuaranteed", "image", "version1").Return([]*image.Image{
					{
						Name:    "image",
						Version: "version1",
					},
				}, nil)
				p.images.(*images.MockStore).On("FindGuaranteed", "image", "version2").Return([]*image.Image{
					{
						Name:    "image",
						Version: "version2",
					},
				}, nil)
			},
			assertFunc: func(p *BasePlan) bool {
				return p.images.(*images.MockStore).AssertExpectations(t)
			},
		},
		{
			desc: "Testing generate plan when no version is provided",
			plan: &BasePlan{
				images: images.NewMockStore(),
			},
			name:     "image",
			versions: []string{},
			err:      &errors.Error{},
			prepareAssertFunc: func(p *BasePlan) {
				p.images.(*images.MockStore).On("FindByName", "image").Return([]*image.Image{
					{
						Name:    "image",
						Version: "version1",
					},
					{
						Name:    "image",
						Version: "version2",
					},
				}, nil)
			},
			assertFunc: func(p *BasePlan) bool {
				return p.images.(*images.MockStore).AssertExpectations(t)
			},
		},
		{
			desc: "Testing error generating plan when image is not defined",
			plan: &BasePlan{
				images: images.NewMockStore(),
			},
			name:     "image",
			versions: []string{},
			prepareAssertFunc: func(p *BasePlan) {
				p.images.(*images.MockStore).On("FindByName", "image").Return([]*image.Image{}, nil)
			},
			assertFunc: func(p *BasePlan) bool {
				return p.images.(*images.MockStore).AssertExpectations(t)
			},
			err: errors.New(errContext, "The image 'image' seems not to be defined"),
		},
		{
			desc: "Testing error generating plan when image is not defined when version is provided",
			plan: &BasePlan{
				images: images.NewMockStore(),
			},
			name:     "image",
			versions: []string{"version"},
			prepareAssertFunc: func(p *BasePlan) {
				p.images.(*images.MockStore).On("FindGuaranteed", "image", "version").Return([]*image.Image{}, nil)
			},
			assertFunc: func(p *BasePlan) bool {
				return p.images.(*images.MockStore).AssertExpectations(t)
			},
			err: errors.New(errContext, "The version(s) [version] for the image 'image' seems not to be defined"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.plan)
			}

			_, err := test.plan.findImages(test.name, test.versions)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.True(t, test.assertFunc(test.plan))
			}

		})
	}
}
