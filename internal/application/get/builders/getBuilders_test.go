package builders

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/builder"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	filter "github.com/gostevedore/stevedore/internal/infrastructure/filters/builders"
	operationfilter "github.com/gostevedore/stevedore/internal/infrastructure/filters/operation"
	output "github.com/gostevedore/stevedore/internal/infrastructure/output/builders"
	buildersstore "github.com/gostevedore/stevedore/internal/infrastructure/store/builders"
	imagesstore "github.com/gostevedore/stevedore/internal/infrastructure/store/images"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	errContext := "(application::get::builders::Run)"

	tests := []struct {
		desc            string
		app             *GetBuildersApplication
		options         *Options
		prepareMockFunc func(*GetBuildersApplication)
		err             error
	}{
		{
			desc:            "Testing error getting builders without providing build store",
			app:             NewGetBuildersApplication(),
			options:         &Options{},
			prepareMockFunc: func(a *GetBuildersApplication) {},
			err:             errors.New(errContext, "On get builders application, builders store must be provided"),
		},
		{
			desc: "Testing error getting builders without providing images store",
			app: NewGetBuildersApplication(
				WithBuildersStore(
					buildersstore.NewMockStore(),
				),
			),
			options:         &Options{},
			prepareMockFunc: func(a *GetBuildersApplication) {},
			err:             errors.New(errContext, "On get builders application, images store must be provided"),
		},
		{
			desc: "Testing error getting builders without providing selectors",
			app: NewGetBuildersApplication(
				WithBuildersStore(
					buildersstore.NewMockStore(),
				),
				WithImagesStore(
					imagesstore.NewMockStore(),
				),
			),
			options:         &Options{},
			prepareMockFunc: func(a *GetBuildersApplication) {},
			err:             errors.New(errContext, "On get builders application, selectors must be provided"),
		},
		{
			desc: "Testing error getting builders without providing filter factory",
			app: NewGetBuildersApplication(
				WithBuildersStore(
					buildersstore.NewMockStore(),
				),
				WithImagesStore(
					imagesstore.NewMockStore(),
				),
				WithSelector(
					map[string]repository.BuildersSelector{},
				),
			),
			options:         &Options{},
			prepareMockFunc: func(a *GetBuildersApplication) {},
			err:             errors.New(errContext, "On get builders application, filter factory must be provided"),
		},
		{
			desc: "Testing error getting builders without providing builders output",
			app: NewGetBuildersApplication(
				WithBuildersStore(
					buildersstore.NewMockStore(),
				),
				WithImagesStore(
					imagesstore.NewMockStore(),
				),
				WithSelector(
					map[string]repository.BuildersSelector{},
				),
				WithFilterFactory(operationfilter.NewFilterOperationFactory()),
			),
			options:         &Options{},
			prepareMockFunc: func(a *GetBuildersApplication) {},
			err:             errors.New(errContext, "On get builders application, builders output must be provided"),
		},
		{
			desc: "Testing get builders",
			app: NewGetBuildersApplication(
				WithBuildersStore(
					buildersstore.NewMockStore(),
				),
				WithImagesStore(
					imagesstore.NewMockStore(),
				),
				WithSelector(
					map[string]repository.BuildersSelector{
						"name": filter.NewBuilderNameFilter(),
					},
				),
				WithOutput(
					output.NewMockOutput(),
				),
				WithFilterFactory(operationfilter.NewFilterOperationFactory()),
			),
			options: &Options{},
			prepareMockFunc: func(a *GetBuildersApplication) {
				a.imagesStore.(*imagesstore.MockStore).On("List").Return([]*image.Image{
					{
						Name: "image1",
						Builder: &builder.Builder{
							Name:   "image1",
							Driver: "driver",
						},
					},
				}, nil)
				a.buildersStore.(*buildersstore.MockStore).On("List").Return([]*builder.Builder{
					{
						Name:   "builder1",
						Driver: "docker",
					},
					{
						Name:   "builder2",
						Driver: "docker",
					},
				}, nil)
				a.output.(*output.MockOutput).On("Output", []*builder.Builder{
					{
						Name:   "builder1",
						Driver: "docker",
					},
					{
						Name:   "builder2",
						Driver: "docker",
					},
					{
						Name:   "image1",
						Driver: "driver",
					},
				}).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing get builders filtering by name",
			app: NewGetBuildersApplication(
				WithBuildersStore(
					buildersstore.NewMockStore(),
				),
				WithImagesStore(
					imagesstore.NewMockStore(),
				),
				WithSelector(
					map[string]repository.BuildersSelector{
						"name": filter.NewBuilderNameFilter(),
					},
				),
				WithOutput(
					output.NewMockOutput(),
				),
				WithFilterFactory(operationfilter.NewFilterOperationFactory()),
			),
			options: &Options{
				Filter: []string{
					"name=builder1",
				},
			},
			prepareMockFunc: func(a *GetBuildersApplication) {
				a.imagesStore.(*imagesstore.MockStore).On("List").Return([]*image.Image{}, nil)
				a.buildersStore.(*buildersstore.MockStore).On("List").Return([]*builder.Builder{
					{
						Name:   "builder1",
						Driver: "docker",
					},
					{
						Name:   "builder2",
						Driver: "docker",
					},
				}, nil)
				a.output.(*output.MockOutput).On("Output", []*builder.Builder{
					{
						Name:   "builder1",
						Driver: "docker",
					},
				}).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing get builders filtering by driver",
			app: NewGetBuildersApplication(
				WithBuildersStore(
					buildersstore.NewMockStore(),
				),
				WithImagesStore(
					imagesstore.NewMockStore(),
				),
				WithSelector(
					map[string]repository.BuildersSelector{
						"name":   filter.NewBuilderNameFilter(),
						"driver": filter.NewBuilderDriverFilter(),
					},
				),
				WithOutput(
					output.NewMockOutput(),
				),
				WithFilterFactory(operationfilter.NewFilterOperationFactory()),
			),
			options: &Options{
				Filter: []string{
					"driver=driver1",
				},
			},
			prepareMockFunc: func(a *GetBuildersApplication) {
				a.imagesStore.(*imagesstore.MockStore).On("List").Return([]*image.Image{}, nil)
				a.buildersStore.(*buildersstore.MockStore).On("List").Return([]*builder.Builder{
					{
						Name:   "builder1",
						Driver: "driver1",
					},
					{
						Name:   "builder2",
						Driver: "docker",
					},
				}, nil)
				a.output.(*output.MockOutput).On("Output", []*builder.Builder{
					{
						Name:   "builder1",
						Driver: "driver1",
					},
				}).Return(nil)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing get builders from an image",
			app: NewGetBuildersApplication(
				WithBuildersStore(
					buildersstore.NewMockStore(),
				),
				WithImagesStore(
					imagesstore.NewMockStore(),
				),
				WithSelector(
					map[string]repository.BuildersSelector{
						"name": filter.NewBuilderNameFilter(),
					},
				),
				WithOutput(
					output.NewMockOutput(),
				),
				WithFilterFactory(operationfilter.NewFilterOperationFactory()),
			),
			options: &Options{
				Filter: []string{
					"name=image1",
				},
			},
			prepareMockFunc: func(a *GetBuildersApplication) {
				a.imagesStore.(*imagesstore.MockStore).On("List").Return([]*image.Image{
					{
						Name: "image1",
						Builder: &builder.Builder{
							Name:   "image1",
							Driver: "docker",
						},
					},
					{
						Name: "image2",
						Builder: &builder.Builder{
							Name:   "image2",
							Driver: "docker",
						},
					},
				}, nil)
				a.buildersStore.(*buildersstore.MockStore).On("List").Return([]*builder.Builder{
					{
						Name:   "builder1",
						Driver: "docker",
					},
					{
						Name:   "builder2",
						Driver: "docker",
					},
				}, nil)
				a.output.(*output.MockOutput).On("Output", []*builder.Builder{
					{
						Name:   "image1",
						Driver: "docker",
					},
				}).Return(nil)
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
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.app.imagesStore.(*imagesstore.MockStore).AssertExpectations(t)
				test.app.buildersStore.(*buildersstore.MockStore).AssertExpectations(t)
				test.app.output.(*output.MockOutput).AssertExpectations(t)
			}
		})
	}
}
