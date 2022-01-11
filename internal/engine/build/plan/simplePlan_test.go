package plan

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/image"
	"github.com/gostevedore/stevedore/internal/image/store"
	"github.com/stretchr/testify/assert"
)

func TestSimplePlanPlan(t *testing.T) {

	errContext := "(plan::Simple::Plan)"

	tests := []struct {
		desc              string
		plan              *SimplePlan
		name              string
		versions          []string
		res               int
		prepareAssertFunc func(*SimplePlan)
		assertFunc        func(*SimplePlan) bool
		err               error
	}{
		{
			desc: "Testing error when images storer is nil",
			plan: &SimplePlan{},
			err:  errors.New(errContext, "Images storer is nil"),
		},
		{
			desc: "Testing generate plan with an image name and versions",
			plan: &SimplePlan{
				BasePlan{
					images: store.NewMockImageStore(),
				},
			},
			name:     "image",
			versions: []string{"version1", "version2"},
			err:      &errors.Error{},
			res:      2,
			prepareAssertFunc: func(p *SimplePlan) {
				p.images.(*store.MockImageStore).On("Find", "image", "version1").Return(&image.Image{
					Name:    "image",
					Version: "version1",
				}, nil)
				p.images.(*store.MockImageStore).On("Find", "image", "version2").Return(&image.Image{
					Name:    "image",
					Version: "version2",
				}, nil)
			},
			assertFunc: func(p *SimplePlan) bool {
				return p.images.(*store.MockImageStore).AssertExpectations(t)
			},
		},
		{
			desc: "Testing generate plan when no version is provided",
			plan: &SimplePlan{
				BasePlan{
					images: store.NewMockImageStore(),
				},
			},
			name:     "image",
			versions: []string{},
			err:      &errors.Error{},
			prepareAssertFunc: func(p *SimplePlan) {
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
			res: 2,
			assertFunc: func(p *SimplePlan) bool {
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

// func TestGenerateImagesList(t *testing.T) {
// 	errContext := "(build::generateImagesList)"

// 	tests := []struct {
// 		desc              string
// 		service           *Service
// 		name              string
// 		versions          []string
// 		res               []*image.Image
// 		prepareAssertFunc func(*Service)
// 		assertFunc        func(*Service) bool
// 		err               error
// 	}{
// 		{
// 			desc: "Testing error when no image name is provided",
// 			name: "",
// 			err:  errors.New(errContext, "Image name is required to build an image"),
// 		},
// 		{
// 			desc: "Testing generate images list",
// 			service: &Service{
// 				images: imagestore.NewMockImageStore(),
// 			},
// 			name:     "image",
// 			versions: []string{"version1", "version2"},
// 			res: []*image.Image{
// 				{
// 					Name:    "image",
// 					Version: "version1",
// 				},
// 				{
// 					Name:    "image",
// 					Version: "version2",
// 				},
// 			},
// 			err: &errors.Error{},
// 			prepareAssertFunc: func(s *Service) {
// 				s.images.(*imagestore.MockImageStore).On("Find", "image", "version1").Return(&image.Image{
// 					Name:    "image",
// 					Version: "version1",
// 				}, nil)
// 				s.images.(*imagestore.MockImageStore).On("Find", "image", "version2").Return(&image.Image{
// 					Name:    "image",
// 					Version: "version2",
// 				}, nil)
// 			},
// 			assertFunc: func(s *Service) bool {
// 				return s.images.(*imagestore.MockImageStore).AssertExpectations(t)
// 			},
// 		},

// 		{
// 			desc: "Testing generate images list when no version is provided",
// 			service: &Service{
// 				images: imagestore.NewMockImageStore(),
// 			},
// 			name:     "image",
// 			versions: []string{},
// 			res: []*image.Image{
// 				{
// 					Name:    "image",
// 					Version: "version1",
// 				},
// 				{
// 					Name:    "image",
// 					Version: "version2",
// 				},
// 			},
// 			err: &errors.Error{},
// 			prepareAssertFunc: func(s *Service) {
// 				s.images.(*imagestore.MockImageStore).On("All", "image").Return([]*image.Image{
// 					{
// 						Name:    "image",
// 						Version: "version1",
// 					},
// 					{
// 						Name:    "image",
// 						Version: "version2",
// 					},
// 				}, nil)
// 			},
// 			assertFunc: func(s *Service) bool {
// 				return s.images.(*imagestore.MockImageStore).AssertExpectations(t)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			if test.prepareAssertFunc != nil {
// 				test.prepareAssertFunc(test.service)
// 			}

// 			res, err := test.service.generateImagesList(test.name, test.versions)

// 			if err != nil {
// 				assert.Equal(t, test.err.Error(), err.Error())
// 			} else {
// 				assert.True(t, test.assertFunc(test.service))
// 				assert.Equal(t, test.res, res)
// 			}
// 		})
// 	}

// }
