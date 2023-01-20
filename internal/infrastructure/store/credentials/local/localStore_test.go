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
	"github.com/gostevedore/stevedore/internal/infrastructure/store/credentials/encryption"
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
		{
			desc: "Testing persist a badge into local store using encryption",
			store: NewLocalStore(
				WithFilesystem(afero.NewMemMapFs()),
				WithPath(credentialsPath),
				WithFormater(mock.NewMockFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
				WithEncryption(encryption.NewMockEncryption()),
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
					}).Return(`{
  "ID": "id",
  "aws_access_key_id": "",
  "aws_region": "",
  "aws_role_arn": "",
  "aws_secret_access_key": "",
  "aws_profile": "",
  "aws_shared_credentials_files": null,
  "aws_shared_config_files": null,
  "aws_use_default_credentials_chain": false,
  "docker_login_password": "",
  "docker_login_username": "",
  "password": "password",
  "username": "username",
  "private_key_file": "",
  "private_key_password": "",
  "git_ssh_user": "",
  "use_ssh_agent": false
}`, nil)

				s.encryption.(*encryption.MockEncription).On("Encrypt", `{
  "ID": "id",
  "aws_access_key_id": "",
  "aws_region": "",
  "aws_role_arn": "",
  "aws_secret_access_key": "",
  "aws_profile": "",
  "aws_shared_credentials_files": null,
  "aws_shared_config_files": null,
  "aws_use_default_credentials_chain": false,
  "docker_login_password": "",
  "docker_login_username": "",
  "password": "password",
  "username": "username",
  "private_key_file": "",
  "private_key_password": "",
  "git_ssh_user": "",
  "use_ssh_agent": false
}`).Return("encrypted-text", nil)
			},
			res: "encrypted-text",
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
	encryptedCredentialsPath := filepath.Join("credentials.encrypted")
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
	err = afero.WriteFile(testFs, filepath.Join("credentials.encrypted", "b80bb7740288fda1f201890375a60c8f"), []byte(`eea3a34338172045310bba615d4669f8806dc36ac814c8db152d5c5cb980b2b2da2fa6c5962d17f2032332af099fe6d412d8385b0fd0dd23b34bb7a459254f9ba306530bee03fb1a51ecd42b0caaf0b2120ca2b6f7d4f4a79b46105997e92c933f469e322f002dd9bed8f1a28820769bf40f1fa2cd8abaa8ab0f9b3e52dd45d4ef6220ad87bed9a36c6af8bc6683616a4f7f337f446a837091aab040f53d1c2de1da2a1dd6ccff02ce81e30de58278a09351c64f003b1d71a26d3eabbbe127df39932061c9ca14400570c3aa5215edb1b22997a5226d497bc6a305ffa7b73e39d7f53fcf1119630fd9cf167634e845dc0e285bb4d7999c04c70688fb8cd711a0c302f1df71ead2917ceabc43740ef2c6a9f3bebbb60b219b4eac63ce89c5a437e824bea8ddc3d8cc87b8e1bfda3efcf56b9b63ac96bf35cee4e88d87fbea6fe2a4e4d67ca145a4b98b915d6cd464a329dbd156070c1559b18a8d16d8135ed43a1b13c4a922dc5574a58263f584faa77be00db8cc3bb63df2fb1796d35d0f21bc6365ef17470f67919dc7583fa728481c3e488042037b4ed7fa0ab7772b933f6d176e89a6c89158f8c989fbd83d198113b913059cb9e11333dfa236987d70b8348cfef9b701a3a5d5559fb9e70c372f9ad7f9b2ef7399c873fd18fb144b578d6fbde7577b75113fe512504262bc530a64387b74eb8dd50707eabe0b1d64dd5cd91faa9ec609464a8f`), 0666)
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
		{
			desc: "Testing get encrypted credentials badge from local store",
			store: NewLocalStore(
				WithFilesystem(testFs),
				WithPath(encryptedCredentialsPath),
				WithFormater(json.NewJSONFormater()),
				WithCompatibility(
					credentialscompatibility.NewCredentialsCompatibility(
						compatibility.NewMockCompatibility(),
					),
				),
				WithEncryption(encryption.NewEncryption(
					encryption.WithKey("encryption-key"),
				)),
			),
			id: "id",
			res: &credentials.Badge{
				ID:                        "id",
				Username:                  "username",
				Password:                  "password",
				AWSSharedConfigFiles:      []string{""},
				AWSSharedCredentialsFiles: []string{""},
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

// func TestHashID(t *testing.T) {

// 	errContext := "(store::credentials::local::hashID)"
// 	tests := []struct {
// 		desc string
// 		id   string
// 		res  string
// 		err  error
// 	}{
// 		{
// 			desc: "Testing error when hashing an id with providing the id",
// 			id:   "",
// 			err:  errors.New(errContext, "Hash method requires an id"),
// 		},
// 		{
// 			desc: "Testing hashing an id",
// 			id:   "id",
// 			res:  "b80bb7740288fda1f201890375a60c8f",
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			res, err := hashID(test.id)
// 			if err != nil {
// 				assert.Equal(t, test.err.Error(), err.Error())
// 			} else {
// 				assert.Equal(t, test.res, res)
// 			}
// 		})
// 	}
// }
