package local

import (
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/compatibility"
	credentialscompatibility "github.com/gostevedore/stevedore/internal/infrastructure/credentials/compatibility"
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
			desc: "Testing error when storing to local store a badge without local path",
			store: NewLocalStore(
				WithFilesystem(afero.NewMemMapFs()),
				WithPath(""),
				WithFormater(mock.NewMockFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			err: errors.New(errContext, "To store a badge into local store, local store path must be provided"),
		},
		{
			desc: "Testing error when storing to local store a badge without an id",
			store: NewLocalStore(
				WithFilesystem(afero.NewMemMapFs()),
				WithPath(credentialsPath),
				WithFormater(mock.NewMockFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			err: errors.New(errContext, "To store a badge into local store, id must be provided"),
		},
		{
			desc: "Testing error when storing to local store a badge with a nil badge",
			store: NewLocalStore(
				WithFilesystem(afero.NewMemMapFs()),
				WithPath(credentialsPath),
				WithFormater(mock.NewMockFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			id:  "id",
			err: errors.New(errContext, "To store a badge for 'id' into local store, credentials badge must be provided"),
		},
		{
			desc: "Testing persist a badge into local store",
			store: NewLocalStore(
				WithFilesystem(afero.NewMemMapFs()),
				WithPath(credentialsPath),
				WithFormater(mock.NewMockFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			id: "id",
			badge: &credentials.Badge{
				Username: "username",
				Password: "password",
			},
			prepareAssertFunc: func(s *LocalStore) {
				s.formater.(*mock.MockFormater).On("Marshal",
					&credentials.Badge{
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

func TestSafeStore(t *testing.T) {
	var err error

	errContext := "(store::credentials::local::StoreSafe)"

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
		store             *LocalStore
		id                string
		badge             *credentials.Badge
		prepareAssertFunc func(*LocalStore)
		res               string
		err               error
	}{
		{
			desc: "Testing error persisting a badge into local store that already exist",
			id:   "existing_id",
			store: NewLocalStore(
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
			desc: "Testing persist a badge into local store",
			store: NewLocalStore(
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
			badge: &credentials.Badge{
				Username: "username",
				Password: "password",
			},
			prepareAssertFunc: func(s *LocalStore) {
				s.formater.(*mock.MockFormater).On("Marshal",
					&credentials.Badge{
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

			err := test.store.SafeStore(test.id, test.badge)
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
			desc: "Testing error when getting a badge from local store without giving an id",
			store: NewLocalStore(
				WithFilesystem(afero.NewMemMapFs()),
				WithPath(credentialsPath),
				WithFormater(json.NewJSONFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			err: errors.New(errContext, "To get a badge from the store, id must be provided"),
		},
		{
			desc: "Testing get credentials badge from local store",
			store: NewLocalStore(
				WithFilesystem(testFs),
				WithPath(credentialsPath),
				WithFormater(json.NewJSONFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			id: "id",
			res: &credentials.Badge{
				ID:       "b80bb7740288fda1f201890375a60c8f",
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

func TestAll(t *testing.T) {

	var err error

	credentialsPath := filepath.Join("credentials")
	testFs := afero.NewMemMapFs()
	testFs.MkdirAll(credentialsPath, 0755)
	emptyPath := filepath.Join("empty")
	testFs.MkdirAll(emptyPath, 0755)

	err = afero.WriteFile(testFs, filepath.Join("credentials", "b80bb7740288fda1f201890375a60c8f"), []byte(`
{
	  "username": "username",
	  "password": "password"
}
`), 0666)
	if err != nil {
		t.Log(err)
	}

	tests := []struct {
		desc  string
		store *LocalStore
		res   []*credentials.Badge
		err   error
	}{
		{
			desc: "Testing get all credentials badges from local store",
			store: NewLocalStore(
				WithFilesystem(testFs),
				WithPath(credentialsPath),
				WithFormater(json.NewJSONFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			res: []*credentials.Badge{
				{
					ID:       "b80bb7740288fda1f201890375a60c8f",
					Username: "username",
					Password: "password",
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing get all credentials badges from an empty local store",
			store: NewLocalStore(
				WithFilesystem(testFs),
				WithPath(emptyPath),
				WithFormater(json.NewJSONFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
			),
			res: []*credentials.Badge{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			badges, err := test.store.All()
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, badges)
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
