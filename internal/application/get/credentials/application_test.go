package credentials

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	mockoutput "github.com/gostevedore/stevedore/internal/infrastructure/output/credentials"
	mockstore "github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/mock"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	tests := []struct {
		desc              string
		app               *Application
		prepareAssertFunc func(app *Application)
		err               error
	}{
		{
			desc: "Testing get credentials application",
			app: NewApplication(
				WithCredentials(
					mockstore.NewMockStore(),
				),
				WithOutput(
					mockoutput.NewMockOutput(),
				),
			),
			prepareAssertFunc: func(app *Application) {
				app.credentials.(*mockstore.MockStore).On("All").Return([]*credentials.Badge{
					{
						ID:       "id",
						Username: "username",
						Password: "password",
					},
				}, nil)
				app.output.(*mockoutput.MockOutput).On("Print", []*credentials.Badge{
					{
						ID:       "id",
						Username: "username",
						Password: "password",
					},
				}).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.app)
			}

			err := test.app.Run(context.TODO())
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				test.app.credentials.(*mockstore.MockStore).AssertExpectations(t)
				test.app.output.(*mockoutput.MockOutput).AssertExpectations(t)
			}

		})
	}
}
