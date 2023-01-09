package plan

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/stretchr/testify/assert"
)

func TestCascadePlanPlan(t *testing.T) {
	errContext := "(plan::Simple::Plan)"

	tests := []struct {
		desc              string
		plan              *CascadePlan
		name              string
		versions          []string
		res               int
		prepareAssertFunc func(*CascadePlan)
		assertFunc        func(*CascadePlan) bool
		err               error
	}{
		{
			desc: "Testing error when images storer is nil",
			plan: &CascadePlan{},
			err:  errors.New(errContext, "Images storer is nil"),
		},
		{
			desc: "Testing generate cascade plan with three images",
			plan: &CascadePlan{
				BasePlan{
					images: images.NewMockStore(),
				},
				// Depth
				-1,
			},
			name:     "image",
			versions: []string{"version1"},
			err:      &errors.Error{},
			res:      3,
			prepareAssertFunc: func(p *CascadePlan) {
				p.images.(*images.MockStore).On("FindGuaranteed", "image", "version1").Return([]*image.Image{
					{
						Name:    "image",
						Version: "version1",
						Children: []*image.Image{
							{
								Name:    "image2",
								Version: "version2",
								Children: []*image.Image{
									{
										Name:    "image3",
										Version: "version3",
									},
								},
							},
						},
					},
				}, nil)

				p.images.(*images.MockStore).On("IsWildcard",
					&image.Image{
						Name:    "image",
						Version: "version1",
						Children: []*image.Image{
							{
								Name:    "image2",
								Version: "version2",
								Children: []*image.Image{
									{
										Name:    "image3",
										Version: "version3",
									},
								},
							},
						},
					}).Return(false)

				p.images.(*images.MockStore).On("IsWildcard",
					&image.Image{
						Name:    "image2",
						Version: "version2",
						Children: []*image.Image{
							{
								Name:    "image3",
								Version: "version3",
							},
						},
					}).Return(false)

				p.images.(*images.MockStore).On("IsWildcard",
					&image.Image{

						Name:    "image3",
						Version: "version3",
					}).Return(false)

			},
			assertFunc: func(p *CascadePlan) bool {
				return p.images.(*images.MockStore).AssertExpectations(t)
			},
		},
		{
			desc: "Testing generate cascade plan with depth defined",
			plan: &CascadePlan{
				BasePlan{
					images: images.NewMockStore(),
				},
				// Depth
				1,
			},
			name:     "image",
			versions: []string{"version1"},
			err:      &errors.Error{},
			res:      2,
			prepareAssertFunc: func(p *CascadePlan) {
				p.images.(*images.MockStore).On("FindGuaranteed", "image", "version1").Return([]*image.Image{
					{
						Name:    "image",
						Version: "version1",
						Children: []*image.Image{
							{
								Name:    "image2",
								Version: "version2",
								Children: []*image.Image{
									{
										Name:    "image3",
										Version: "version3",
									},
								},
							},
						},
					},
				}, nil)

				p.images.(*images.MockStore).On("IsWildcard",
					&image.Image{
						Name:    "image",
						Version: "version1",
						Children: []*image.Image{
							{
								Name:    "image2",
								Version: "version2",
								Children: []*image.Image{
									{
										Name:    "image3",
										Version: "version3",
									},
								},
							},
						},
					}).Return(false)

				p.images.(*images.MockStore).On("IsWildcard",
					&image.Image{
						Name:    "image2",
						Version: "version2",
						Children: []*image.Image{
							{
								Name:    "image3",
								Version: "version3",
							},
						},
					}).Return(false)
			},
			assertFunc: func(p *CascadePlan) bool {
				return p.images.(*images.MockStore).AssertExpectations(t)
			},
		},
		{
			desc: "Testing generate cascade plan when no version is provided",
			plan: &CascadePlan{
				BasePlan{
					images: images.NewMockStore(),
				},
				// Depth
				-1,
			},
			name:     "image",
			versions: []string{},
			err:      &errors.Error{},
			prepareAssertFunc: func(p *CascadePlan) {
				p.images.(*images.MockStore).On("FindByName", "image").Return([]*image.Image{
					{
						Name:    "image",
						Version: "version1",
						Children: []*image.Image{
							{
								Name:    "image2",
								Version: "version2",
								Children: []*image.Image{
									{
										Name:    "image3",
										Version: "version3",
									},
								},
							},
						},
					},
					{
						Name:    "image",
						Version: "version4",
					},
				}, nil)

				p.images.(*images.MockStore).On("IsWildcard",
					&image.Image{
						Name:    "image",
						Version: "version1",
						Children: []*image.Image{
							{
								Name:    "image2",
								Version: "version2",
								Children: []*image.Image{
									{
										Name:    "image3",
										Version: "version3",
									},
								},
							},
						},
					}).Return(false)

				p.images.(*images.MockStore).On("IsWildcard",
					&image.Image{
						Name:    "image2",
						Version: "version2",
						Children: []*image.Image{
							{
								Name:    "image3",
								Version: "version3",
							},
						},
					}).Return(false)

				p.images.(*images.MockStore).On("IsWildcard",
					&image.Image{

						Name:    "image3",
						Version: "version3",
					}).Return(false)

				p.images.(*images.MockStore).On("IsWildcard",
					&image.Image{

						Name:    "image",
						Version: "version4",
					}).Return(false)
			},
			res: 4,
			assertFunc: func(p *CascadePlan) bool {
				return p.images.(*images.MockStore).AssertExpectations(t)
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
