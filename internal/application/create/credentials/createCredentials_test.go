package credentials

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/mock"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	errContext := "(application::create::credentials::Run)"

	tests := []struct {
		desc              string
		app               *CreateCredentialsApplication
		id                string
		badge             *credentials.Badge
		prepareAssertFunc func(CredentialsStorer)
		err               error
	}{
		{
			desc: "Testing run create credentials application without store",
			app:  NewCreateCredentialsApplication(),
			err:  errors.New(errContext, "To run the create credentials application, a credentials storer must be provided"),
		},
		{
			desc: "Testing run create credentials application without badge",
			app:  NewCreateCredentialsApplication(WithCredentialsStore(mock.NewMockStore())),
			err:  errors.New(errContext, "To run the create credentials application, a badge must be provided"),
		},
		{
			desc:  "Testing run create credentials application without badge id",
			app:   NewCreateCredentialsApplication(WithCredentialsStore(mock.NewMockStore())),
			badge: &credentials.Badge{},
			err:   errors.New(errContext, "To run the create credentials application, a id for credentials must be provided"),
		},
		{
			desc: "Testing run create credentials application",
			app:  NewCreateCredentialsApplication(WithCredentialsStore(mock.NewMockStore())),
			badge: &credentials.Badge{
				Username: "username",
				Password: "password",
			},
			id: "id",
			prepareAssertFunc: func(store CredentialsStorer) {
				store.(*mock.MockStore).On("Store", "id", &credentials.Badge{
					Username: "username",
					Password: "password",
				}).Return(nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil && test.app.store != nil {
				test.prepareAssertFunc(test.app.store)
			}

			err := test.app.Run(context.TODO(), test.id, test.badge)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				test.app.store.(*mock.MockStore).AssertExpectations(t)
			}
		})
	}
}
