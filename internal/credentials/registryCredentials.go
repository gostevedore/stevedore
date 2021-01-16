package credentials

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	common "github.com/apenella/go-common-utils/data"
	errors "github.com/apenella/go-common-utils/error"
)

var credentials map[string]*RegistryUserPassAuth

// RegistryCredentials
type RegistryCredentials map[string]*RegistryUserPassAuth

// Init initializes the credentials store
func Init() {
	if credentials == nil {
		credentials = make(map[string]*RegistryUserPassAuth)
	}
}

// LoadCredentials loads credentials defined on dir path to credentials store
func LoadCredentials(dir string) error {

	var err error

	if credentials == nil {
		credentials = make(map[string]*RegistryUserPassAuth)
	}

	_, err = os.Stat(dir)
	if err == nil {

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return errors.New("(credentials::LoadCredentials)", fmt.Sprintf("Error reading directory '%s'", dir), err)
		}

		for _, file := range files {
			userpass := &RegistryUserPassAuth{}
			if file.Mode().IsRegular() {
				filename := file.Name()
				err := common.LoadJSONFile(strings.Join([]string{dir, filename}, string(os.PathSeparator)), userpass)
				if err == nil {
					AddCredential(filename, userpass)
				}
			}
		}
	}

	return nil
}

// ClearCredentials
func ClearCredentials() {
	credentials = make(map[string]*RegistryUserPassAuth)
}

// AddCredential
func AddCredential(id string, auth *RegistryUserPassAuth) error {

	if credentials == nil {
		credentials = make(map[string]*RegistryUserPassAuth)
	}

	_, exists := credentials[id]
	if exists {
		return errors.New("(credentials::AddCredential)", fmt.Sprintf("Auth method with id '%s' already exist", id))
	}

	credentials[id] = auth

	return nil
}

// AchieveCredential
func AchieveCredential(registry string) (*RegistryUserPassAuth, error) {

	if credentials == nil {
		return nil, errors.New("(credentials::AchieveCredential)", "Credentials has not been initialized")
	}

	hashedRegisty := hashRegistryName(registry)

	credential, exists := credentials[hashedRegisty]
	if !exists {
		return nil, errors.New("(credentials::AchieveCredential)", "No credential found for '"+registry+"'")
	}

	return credential, nil
}

// ListRegistryCredentials
func ListRegistryCredentials(wideList bool) ([][]string, error) {
	list := [][]string{}

	if credentials == nil {
		return nil, errors.New("(credentials::ListRegistryCredentials)", "Credentials has not been initialized")
	}

	for id, credential := range credentials {

		credentialArray := []string{id, credential.Username}
		if wideList {
			credentialArray = append(credentialArray, credential.Password)
		}

		list = append(list, credentialArray)
	}

	return list, nil
}

// ListRegistryCredentialsHeader
func ListRegistryCredentialsHeader(wideList bool) []string {
	h := []string{
		"CREDENTIAL ID",
		"USERNAME",
	}

	if wideList {
		h = append(h, "PASSWORD")
	}

	return h
}

func CreateCredential(dir, username, password, registry string) error {

	var err error
	var userPassJSON string
	var credentialFile *os.File

	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	registryHashed := hashRegistryName(registry)

	credentialFileName := strings.Join([]string{dir, registryHashed}, string(os.PathSeparator))
	credentialFile, err = os.OpenFile(credentialFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("(credentials::CreateCredentials)", "File '"+credentialFileName+"' could not be created", err)
	}
	defer credentialFile.Close()

	userPass := &RegistryUserPassAuth{
		Username: username,
		Password: password,
	}

	userPassJSON, err = common.ObjectToJSONStringPretty(userPass)
	if err != nil {
		return errors.New("(credentials::CreateCredential)", "Error converting user-pass auth to []byte. ", err)
	}
	credentialFile.WriteString(userPassJSON)

	return nil
}

// hashRegistryName
func hashRegistryName(registry string) string {
	hasher := md5.New()
	hasher.Write([]byte(registry))
	registryHashed := hex.EncodeToString(hasher.Sum(nil))

	return registryHashed
}
