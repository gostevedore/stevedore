package credentials

// func TestAddCredential(t *testing.T) {

// 	type args struct {
// 		id   string
// 		auth *UserPasswordAuth
// 	}
// 	tests := []struct {
// 		desc        string
// 		credentials map[string]*UserPasswordAuth
// 		args        *args
// 		res         map[string]*UserPasswordAuth
// 		err         error
// 	}{
// 		{
// 			desc:        "Test add credential no nil Credentials",
// 			credentials: nil,
// 			args: &args{
// 				id: "123",
// 				auth: &UserPasswordAuth{
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 			res: map[string]*UserPasswordAuth{
// 				"123": {
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 		},
// 		{
// 			desc: "Test add credential",
// 			credentials: map[string]*UserPasswordAuth{
// 				"123": {
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 			args: &args{
// 				id: "456",
// 				auth: &UserPasswordAuth{
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 			res: map[string]*UserPasswordAuth{
// 				"123": {
// 					Username: "username",
// 					Password: "password",
// 				},
// 				"456": {
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 		},
// 		{
// 			desc: "Test add existing credential",
// 			credentials: map[string]*UserPasswordAuth{
// 				"123": {
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 			args: &args{
// 				id: "123",
// 				auth: &UserPasswordAuth{
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 			res: nil,
// 			err: errors.New("(credentials::AddCredential)", "Auth method with id '123' already exist"),
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)
// 			ClearCredentials()
// 			credentials = test.credentials

// 			err := AddCredential(test.args.id, test.args.auth)
// 			if err != nil {
// 				assert.Equal(t, test.err, err)
// 			} else {
// 				for id, credential := range credentials {
// 					expectedCredential, exist := test.res[id]
// 					assert.True(t, exist, "Credential '"+id+"' does not exist")
// 					assert.Equal(t, expectedCredential.Username, credential.Username, "Username not expected")
// 					assert.Equal(t, expectedCredential.Password, credential.Password, "Password not expected")
// 				}
// 			}
// 		})
// 	}
// }

// func TestAchieveCredential(t *testing.T) {

// 	type args struct {
// 		registry string
// 	}
// 	tests := []struct {
// 		desc        string
// 		credentials map[string]*UserPasswordAuth
// 		args        *args
// 		res         *UserPasswordAuth
// 		err         error
// 	}{
// 		{
// 			desc:        "Test achieve credentials from a nil Credentials",
// 			credentials: nil,
// 			args: &args{
// 				registry: "registry",
// 			},
// 			res: nil,
// 			err: errors.New("(credentials::AchieveCredential)", "Credentials has not been initialized"),
// 		},
// 		{
// 			desc: "Test achieve credentials",
// 			credentials: map[string]*UserPasswordAuth{
// 				"a9205dcfd4a6f7c2cbe8be01566ff84a": {
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 			args: &args{
// 				registry: "registry",
// 			},
// 			res: &UserPasswordAuth{
// 				Username: "username",
// 				Password: "password",
// 			},
// 		},
// 		{
// 			desc:        "Test achieve unexisting credential",
// 			credentials: map[string]*UserPasswordAuth{},
// 			args: &args{
// 				registry: "registry",
// 			},
// 			res: nil,
// 			err: errors.New("(credentials::AchieveCredential)", "No credential found for 'registry'"),
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			ClearCredentials()
// 			credentials = test.credentials

// 			credential, err := AchieveCredential(test.args.registry)
// 			if err != nil {
// 				assert.Equal(t, test.err, err)
// 			} else {
// 				assert.Equal(t, test.res.Username, credential.Username, "Username not expected")
// 				assert.Equal(t, test.res.Password, credential.Password, "Password not expected")

// 			}
// 		})
// 	}
// }

// func TestListRegistryCredentials(t *testing.T) {
// 	type args struct {
// 		wideList bool
// 	}
// 	tests := []struct {
// 		desc        string
// 		credentials map[string]*UserPasswordAuth
// 		args        *args
// 		res         [][]string
// 		err         error
// 	}{
// 		{
// 			desc:        "Test list credentials from nil Credentials",
// 			credentials: nil,
// 			args: &args{
// 				wideList: false,
// 			},
// 			res: nil,
// 			err: errors.New("(credentials::ListRegistryCredentials)", "Credentials has not been initialized"),
// 		},
// 		{
// 			desc: "Test list credentials",
// 			credentials: map[string]*UserPasswordAuth{
// 				"a9205dcfd4a6f7c2cbe8be01566ff84a": {
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 			args: &args{
// 				wideList: false,
// 			},
// 			res: [][]string{
// 				{"a9205dcfd4a6f7c2cbe8be01566ff84a", "username"},
// 			},
// 			err: nil,
// 		},
// 		{
// 			desc: "Test list credentials returning a wide output",
// 			credentials: map[string]*UserPasswordAuth{
// 				"a9205dcfd4a6f7c2cbe8be01566ff84a": {
// 					Username: "username",
// 					Password: "password",
// 				},
// 			},
// 			args: &args{
// 				wideList: true,
// 			},
// 			res: [][]string{
// 				{"a9205dcfd4a6f7c2cbe8be01566ff84a", "username", "password"},
// 			},
// 			err: nil,
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			ClearCredentials()
// 			credentials = test.credentials

// 			credentialList, err := ListRegistryCredentials(test.args.wideList)
// 			if err != nil {
// 				assert.Equal(t, test.err, err)
// 			} else {
// 				equal := reflect.DeepEqual(credentialList, test.res)
// 				assert.True(t, equal, "Credential list not equal")

// 			}
// 		})
// 	}
// }

// func TestListRegistryCredentialsHeader(t *testing.T) {
// 	type args struct {
// 		wideList bool
// 	}
// 	tests := []struct {
// 		desc string
// 		args *args
// 		res  []string
// 	}{
// 		{
// 			desc: "Test list registry credentials header",
// 			args: &args{
// 				wideList: false,
// 			},
// 			res: []string{
// 				"CREDENTIAL ID",
// 				"USERNAME",
// 			},
// 		},
// 		{
// 			desc: "Test list registry credentials header returning a wide output",
// 			args: &args{
// 				wideList: true,
// 			},
// 			res: []string{
// 				"CREDENTIAL ID",
// 				"USERNAME",
// 				"PASSWORD",
// 			},
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.desc, func(t *testing.T) {
// 			t.Log(test.desc)

// 			header := ListRegistryCredentialsHeader(test.args.wideList)
// 			equal := reflect.DeepEqual(header, test.res)
// 			assert.True(t, equal, "Header list not equal")
// 		})
// 	}
// }

// func TestHashRegistryName(t *testing.T) {
// 	t.Log("Test hash registry name")
// 	res := "a9205dcfd4a6f7c2cbe8be01566ff84a"
// 	hash := hashRegistryName("registry")
// 	assert.Equal(t, res, hash, "Return an unexpected hash")

// }

// // func TestLoadCredentials(t *testing.T) {
// // 	type args struct {
// // 		dir string
// // 	}
// // 	tests := []struct {
// // 		name    string
// // 		args    args
// // 		want    *RegistryCredentials
// // 		wantErr bool
// // 	}{
// // 		// TODO: Add test cases.
// // 	}
// // 	for _, test := range tests {
// // 		t.Run(test.name, func(t *testing.T) {
// // 			got, err := LoadCredentials(test.args.dir)
// // 			if (err != nil) != test.wantErr {
// // 				t.Errorf("LoadCredentials() error = %v, wantErr %v", err, test.wantErr)
// // 				return
// // 			}
// // 			if !reflect.DeepEqual(got, test.want) {
// // 				t.Errorf("LoadCredentials() = %v, want %v", got, test.want)
// // 			}
// // 		})
// // 	}
// // }

// // func TestCreateCredential(t *testing.T) {
// // 	type args struct {
// // 		dir      string
// // 		username string
// // 		password string
// // 		registry string
// // 	}
// // 	tests := []struct {
// // 		name    string
// // 		args    args
// // 		wantErr bool
// // 	}{
// // 		// TODO: Add test cases.
// // 	}
// // 	for _, test := range tests {
// // 		t.Run(test.name, func(t *testing.T) {
// // 			if err := CreateCredential(test.args.dir, test.args.username, test.args.password, test.args.registry); (err != nil) != test.wantErr {
// // 				t.Errorf("CreateCredential() error = %v, wantErr %v", err, test.wantErr)
// // 			}
// // 		})
// // 	}
// // }
