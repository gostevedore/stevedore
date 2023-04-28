package local

import (
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	credentialscompatibility "github.com/gostevedore/stevedore/internal/infrastructure/compatibility/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/format/credentials/mock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestStore_LocalStoreWithSafeStore(t *testing.T) {
	var err error

	errContext := "(store::credentials::local::LocalStoreWithSafeStore)"

	credentialsPath := filepath.Join("credentials")
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(credentialsPath, 0755)

	err = afero.WriteFile(testFs, filepath.Join("credentials", "52a3dd11c26f43983739cec4b383af28"), []byte(`
{
	  "username": "username",
	  "password": "password"
}
`), 0666)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc              string
		store             *LocalStoreWithSafeStore
		id                string
		credential        *credentials.Credential
		prepareAssertFunc func(*LocalStoreWithSafeStore)
		res               string
		err               error
	}{
		{
			desc: "Testing error persisting a credential into local store that already exist with safe store",
			id:   "existing_id",
			store: NewLocalStoreWithSafeStore(
				WithFilesystem(testFs),
				WithPath(credentialsPath),
				WithFormater(mock.NewMockFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			err: errors.New(errContext, "Credentials 'existing_id' already exist"),
		},
		{
			desc: "Testing persist a credential into local store with safe store",
			store: NewLocalStoreWithSafeStore(
				WithFilesystem(testFs),
				WithPath(credentialsPath),
				WithFormater(mock.NewMockFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			id: "id",
			credential: &credentials.Credential{
				Username: "username",
				Password: "password",
			},
			prepareAssertFunc: func(s *LocalStoreWithSafeStore) {
				s.store.formater.(*mock.MockFormater).On("Marshal",
					&credentials.Credential{
						ID:       "id",
						Username: "username",
						Password: "password",
					}).Return("formated", nil)
			},
			res: "formated",
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.store)
			}

			err := test.store.Store(test.id, test.credential)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				testFs := afero.Afero{Fs: test.store.store.fs}
				content, err := testFs.ReadFile(filepath.Join(credentialsPath, "b80bb7740288fda1f201890375a60c8f"))
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.res, string(content))
			}
		})
	}
}
