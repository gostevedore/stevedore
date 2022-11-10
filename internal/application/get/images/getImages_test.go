package images

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	filter "github.com/gostevedore/stevedore/internal/infrastructure/filters/images"
	output "github.com/gostevedore/stevedore/internal/infrastructure/output/images"
	store "github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	errContext := "(application::get::images::Run)"

	tests := []struct {
		desc            string
		app             *GetImagesApplication
		options         *Options
		prepareMockFunc func(*GetImagesApplication)
		err             error
	}{
		{
			desc: "Testing error on get images application and store is not defined",
			app:  &GetImagesApplication{},
			err:  errors.New(errContext, "On get images application, images store must be provided"),
		},
		{
			desc: "Testing error on get images application and output is not defined",
			app: NewGetImagesApplication(
				WithStore(
					store.NewMockStore(),
				),
			),
			err: errors.New(errContext, "On get images application, images output must be provided"),
		},
		{
			desc: "Testing application get images",
			app: NewGetImagesApplication(
				WithStore(
					store.NewMockStore(),
				),
				WithOutput(
					output.NewMockOutput(),
				),
				WithSelector(
					map[string]repository.ImagesSelector{
						"name": filter.NewImageNameFilter(),
					},
				),
			),
			options: &Options{},
			prepareMockFunc: func(a *GetImagesApplication) {
				a.store.(*store.MockStore).On("List").Return([]*image.Image{
					{
						Name:    "image1",
						Version: "v1",
					},
					{
						Name:    "image2",
						Version: "v1",
					},
				}, nil)
				a.output.(*output.MockOutput).On("Output",
					[]*image.Image{
						{
							Name:    "image1",
							Version: "v1",
						},
						{
							Name:    "image2",
							Version: "v1",
						},
					},
				).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing application get images filtering by name",
			app: NewGetImagesApplication(
				WithStore(
					store.NewMockStore(),
				),
				WithOutput(
					output.NewMockOutput(),
				),
				WithSelector(
					map[string]repository.ImagesSelector{
						"name": filter.NewImageNameFilter(),
					},
				),
			),
			options: &Options{
				Filter: []string{
					"name=image1",
				},
			},
			prepareMockFunc: func(a *GetImagesApplication) {
				a.store.(*store.MockStore).On("List").Return([]*image.Image{
					{
						Name:    "image1",
						Version: "v1",
					},
					{
						Name:    "image2",
						Version: "v1",
					},
				}, nil)
				a.output.(*output.MockOutput).On("Output",
					[]*image.Image{
						{
							Name:    "image1",
							Version: "v1",
						},
					},
				).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing application get images with not valid filter",
			app: NewGetImagesApplication(
				WithStore(
					store.NewMockStore(),
				),
				WithOutput(
					output.NewMockOutput(),
				),
				WithSelector(
					map[string]repository.ImagesSelector{
						"name": filter.NewImageNameFilter(),
					},
				),
			),
			options: &Options{
				Filter: []string{
					"attr=image1",
				},
			},
			prepareMockFunc: func(a *GetImagesApplication) {
				a.store.(*store.MockStore).On("List").Return([]*image.Image{
					{
						Name:    "image1",
						Version: "v1",
					},
					{
						Name:    "image2",
						Version: "v1",
					},
				}, nil)
				a.output.(*output.MockOutput).On("Output",
					[]*image.Image{
						{
							Name:    "image1",
							Version: "v1",
						},
						{
							Name:    "image2",
							Version: "v1",
						},
					},
				).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMockFunc != nil && test.app != nil {
				test.prepareMockFunc(test.app)
			}

			err := test.app.Run(context.TODO(), test.options)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				test.app.store.(*store.MockStore).AssertExpectations(t)
				test.app.output.(*output.MockOutput).AssertExpectations(t)
			}
		})
	}
}
