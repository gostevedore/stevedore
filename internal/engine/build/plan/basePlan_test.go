package plan

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/image/store"
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
				images: store.NewMockImageStore(),
			},
			name:     "image",
			versions: []string{"version1", "version2"},
			err:      &errors.Error{},
			prepareAssertFunc: func(p *BasePlan) {
				p.images.(*store.MockImageStore).On("Find", "image", "version1").Return(&image.Image{
					Name:    "image",
					Version: "version1",
				}, nil)
				p.images.(*store.MockImageStore).On("Find", "image", "version2").Return(&image.Image{
					Name:    "image",
					Version: "version2",
				}, nil)
			},
			assertFunc: func(p *BasePlan) bool {
				return p.images.(*store.MockImageStore).AssertExpectations(t)
			},
		},
		{
			desc: "Testing generate plan when no version is provided",
			plan: &BasePlan{
				images: store.NewMockImageStore(),
			},
			name:     "image",
			versions: []string{},
			err:      &errors.Error{},
			prepareAssertFunc: func(p *BasePlan) {
				p.images.(*store.MockImageStore).On("All", "image").Return([]*image.Image{
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
				return p.images.(*store.MockImageStore).AssertExpectations(t)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			//t.Log(test.desc)
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
