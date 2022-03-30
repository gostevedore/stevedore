package plan

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
	"github.com/gostevedore/stevedore/internal/images/store"
	"github.com/stretchr/testify/assert"
)

func TestSinglePlanPlan(t *testing.T) {

	errContext := "(plan::Simple::Plan)"

	tests := []struct {
		desc              string
		plan              *SinglePlan
		name              string
		versions          []string
		res               int
		prepareAssertFunc func(*SinglePlan)
		assertFunc        func(*SinglePlan) bool
		err               error
	}{
		{
			desc: "Testing error when images storer is nil",
			plan: &SinglePlan{},
			err:  errors.New(errContext, "Images storer is nil"),
		},
		{
			desc: "Testing generate plan with an image name and versions",
			plan: &SinglePlan{
				BasePlan{
					images: store.NewMockImageStore(),
				},
			},
			name:     "image",
			versions: []string{"version1", "version2"},
			err:      &errors.Error{},
			res:      2,
			prepareAssertFunc: func(p *SinglePlan) {
				p.images.(*store.MockImageStore).On("Find", "image", "version1").Return(&image.Image{
					Name:    "image",
					Version: "version1",
				}, nil)
				p.images.(*store.MockImageStore).On("Find", "image", "version2").Return(&image.Image{
					Name:    "image",
					Version: "version2",
				}, nil)
			},
			assertFunc: func(p *SinglePlan) bool {
				return p.images.(*store.MockImageStore).AssertExpectations(t)
			},
		},
		{
			desc: "Testing generate plan when no version is provided",
			plan: &SinglePlan{
				BasePlan{
					images: store.NewMockImageStore(),
				},
			},
			name:     "image",
			versions: []string{},
			err:      &errors.Error{},
			prepareAssertFunc: func(p *SinglePlan) {
				p.images.(*store.MockImageStore).On("FindByName", "image").Return([]*image.Image{
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
			res: 2,
			assertFunc: func(p *SinglePlan) bool {
				return p.images.(*store.MockImageStore).AssertExpectations(t)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.plan)
			}

			res, err := test.plan.Plan(test.name, test.versions)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.True(t, test.assertFunc(test.plan))
				assert.Equal(t, test.res, len(res))
			}

		})
	}

}
