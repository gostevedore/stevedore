package local

import (
	"path/filepath"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/formater/mock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {

	errContext := "(store::credentials::local::Store)"

	tests := []struct {
		desc  string
		store *LocalStore
		id    string
		badge *credentials.Badge
		res   map[string]*credentials.Badge
		err   error
	}{
		{
			desc:  "Testing error when storing to local store a badge without an id",
			store: NewLocalStore(afero.NewMemMapFs(), mock.NewMockFormater()),
			err:   errors.New(errContext, "To store a badge into local store, id must be provided"),
		},
		{
			desc:  "Testing error when storing to local store a badge with a nil badge",
			store: NewLocalStore(afero.NewMemMapFs(), mock.NewMockFormater()),
			id:    "id",
			err:   errors.New(errContext, "To store a badge for 'id' into local store, credentials badge must be provided"),
		},
		{
			desc:  "Testing store a badge into local store",
			store: NewLocalStore(afero.NewMemMapFs(), mock.NewMockFormater()),
			id:    "id",
			badge: &credentials.Badge{
				Username: "username",
				Password: "password",
			},

			res: map[string]*credentials.Badge{
				"b80bb7740288fda1f201890375a60c8f": {
					Username: "username",
					Password: "password",
				},
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing error when storing an existing id already found on local store",
			store: &LocalStore{
				store: map[string]*credentials.Badge{
					"b80bb7740288fda1f201890375a60c8f": {
						Username: "username",
						Password: "password",
					},
				},
			},
			id:    "id",
			badge: &credentials.Badge{},
			res: map[string]*credentials.Badge{
				"b80bb7740288fda1f201890375a60c8f": {
					Username: "username",
					Password: "password",
				},
			},
			err: errors.New(errContext, "Badge with id 'id' already exists"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.store.Store(test.id, test.badge)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, test.store.store)
			}
		})
	}
}

func TestPersist(t *testing.T) {
	errContext := "(store::credentials::local::Persist)"

	tests := []struct {
		desc              string
		store             *LocalStore
		path              string
		id                string
		badge             *credentials.Badge
		prepareAssertFunc func(*LocalStore)
		res               string
		err               error
	}{
		{
			desc:  "Testing error when persisting to local store a badge without local path",
			store: NewLocalStore(afero.NewMemMapFs(), mock.NewMockFormater()),
			err:   errors.New(errContext, "To persist a badge into local store, local store path must be provided"),
		},
		{
			desc:  "Testing error when persisting to local store a badge without an id",
			store: NewLocalStore(afero.NewMemMapFs(), mock.NewMockFormater()),
			path:  "/credentials",
			err:   errors.New(errContext, "To persist a badge into local store, id must be provided"),
		},
		{
			desc:  "Testing error when persisting to local store a badge with a nil badge",
			store: NewLocalStore(afero.NewMemMapFs(), mock.NewMockFormater()),
			id:    "id",
			path:  "/credentials",
			err:   errors.New(errContext, "To persist a badge for 'id' into local store, credentials badge must be provided"),
		},
		{
			desc:  "Testing persist a badge into local store",
			store: NewLocalStore(afero.NewMemMapFs(), mock.NewMockFormater()),
			id:    "id",
			path:  "/credentials",
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

			err := test.store.Persist(test.path, test.id, test.badge)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				testFs := afero.Afero{Fs: test.store.fs}
				content, err := testFs.ReadFile(filepath.Join(test.path, "b80bb7740288fda1f201890375a60c8f"))
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.res, string(content))
			}
		})
	}
}

func TestGet(t *testing.T) {

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
			store: NewLocalStore(afero.NewMemMapFs(), mock.NewMockFormater()),
			err:   errors.New(errContext, "To get a badge from the store, id must be provided"),
		},
		{
			desc: "Testing get credentials badge from local store",
			store: &LocalStore{
				store: map[string]*credentials.Badge{
					"b80bb7740288fda1f201890375a60c8f": {
						Username: "username",
						Password: "password",
					},
				},
			},
			id: "id",
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

func TestLoadCredentialsFromFile(t *testing.T) {
	tests := []struct {
		desc  string
		store *LocalStore
		id    string
		badge *credentials.Badge
		err   error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
		})
	}
}

func TestLoadCredentialsFromDir(t *testing.T) {
	tests := []struct {
		desc  string
		store *LocalStore
		id    string
		badge *credentials.Badge
		err   error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
		})
	}
}

func TestLoadCredentials(t *testing.T) {
	tests := []struct {
		desc  string
		store *LocalStore
		id    string
		badge *credentials.Badge
		err   error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
		})
	}
}
