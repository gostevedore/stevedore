package local

import (
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/formater/json"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/formater/mock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	errContext := "(store::credentials::local::Store)"

	credentialsPath := "/credentials"

	tests := []struct {
		desc              string
		store             *LocalStore
		id                string
		badge             *credentials.Badge
		prepareAssertFunc func(*LocalStore)
		res               string
		err               error
	}{
		{
			desc:  "Testing error when storing to local store a badge without local path",
			store: NewLocalStore(afero.NewMemMapFs(), "", mock.NewMockFormater()),
			err:   errors.New(errContext, "To store a badge into local store, local store path must be provided"),
		},
		{
			desc:  "Testing error when storing to local store a badge without an id",
			store: NewLocalStore(afero.NewMemMapFs(), credentialsPath, mock.NewMockFormater()),
			err:   errors.New(errContext, "To store a badge into local store, id must be provided"),
		},
		{
			desc:  "Testing error when storing to local store a badge with a nil badge",
			store: NewLocalStore(afero.NewMemMapFs(), credentialsPath, mock.NewMockFormater()),
			id:    "id",
			err:   errors.New(errContext, "To store a badge for 'id' into local store, credentials badge must be provided"),
		},
		{
			desc:  "Testing persist a badge into local store",
			store: NewLocalStore(afero.NewMemMapFs(), credentialsPath, mock.NewMockFormater()),
			id:    "id",
			badge: &credentials.Badge{
				Username: "username",
				Password: "password",
			},
			prepareAssertFunc: func(s *LocalStore) {
				s.formater.(*mock.MockFormater).On("Marshal",
					&credentials.Badge{
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

			err := test.store.Store(test.id, test.badge)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				testFs := afero.Afero{Fs: test.store.fs}
				content, err := testFs.ReadFile(filepath.Join(credentialsPath, "b80bb7740288fda1f201890375a60c8f"))
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.res, string(content))
			}
		})
	}
}

func TestGet(t *testing.T) {
	var err error

	credentialsPath := filepath.Join("credentials")
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(credentialsPath, 0755)

	err = afero.WriteFile(testFs, filepath.Join("credentials", "b80bb7740288fda1f201890375a60c8f"), []byte(`
{
	  "username": "username",
	  "password": "password"
}
`), 0666)
	if err != nil {
		t.Log(err)
	}

	errContext := "(store::credentials::local::Get)"
	tests := []struct {
		desc  string
		store *LocalStore
		id    string
		res   *credentials.Badge
		err   error
	}{
		{
			desc:  "Testing error when getting a badge from local storewithout giving an id",
			store: NewLocalStore(afero.NewMemMapFs(), credentialsPath, json.NewJSONFormater()),
			err:   errors.New(errContext, "To get a badge from the store, id must be provided"),
		},
		{
			desc:  "Testing get credentials badge from local store",
			store: NewLocalStore(testFs, credentialsPath, json.NewJSONFormater()),
			id:    "id",
			res: &credentials.Badge{
				Username: "username",
				Password: "password",
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			badge, err := test.store.Get(test.id)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, badge)
			}
		})
	}
}

func TestHashID(t *testing.T) {

	errContext := "(store::credentials::local::hashID)"
	tests := []struct {
		desc string
		id   string
		res  string
		err  error
	}{
		{
			desc: "Testing error when hashing an id with providing the id",
			id:   "",
			err:  errors.New(errContext, "Hash method requires an id"),
		},
		{
			desc: "Testing hashing an id",
			id:   "id",
			res:  "b80bb7740288fda1f201890375a60c8f",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			res, err := hashID(test.id)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}
