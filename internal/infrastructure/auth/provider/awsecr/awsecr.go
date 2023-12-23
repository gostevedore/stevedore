package ecr

import (
	"context"
	"encoding/base64"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
	"github.com/gostevedore/stevedore/internal/infrastructure/auth/method/basic"
)

const (
	AWSECRUserName = "AWS"
)

// AWSECRAuthProvider return auth method from credential
type AWSECRAuthProvider struct {
	tokenProvider AWSECRTokenProvider
}

// NewAWSECRAuthProvider return new instance of AWSECRAuthProvider
func NewAWSECRAuthProvider(provider AWSECRTokenProvider) *AWSECRAuthProvider {
	return &AWSECRAuthProvider{
		tokenProvider: provider,
	}
}

// Get returns the most appropiate AuthMethodReader for the credential received
func (p *AWSECRAuthProvider) Get(credential *credentials.Credential) (repository.AuthMethodReader, error) {

	// The error is not being controled because is prefered to return nil and let the caller to decide what to do. AWS would be required for very specific cases and controlling the error would cause anoying errors to the user
	// In worst case scenario, docker API will return an error when trying to pull or push the image because the authorization token is not found
	// TODO: log the error
	token, _ := p.tokenProvider.Get(context.TODO(),
		// That funcion is used to load the AWS Configuration that will be used to create the ECR client to get the authorization token
		func(ctx context.Context, loadOptionsFuncs ...func(*config.LoadOptions) error) (aws.Config, error) {
			// The error is not being controled because is prefered to return nil and let the caller to decide what to do. AWS would be required for very specific cases and controlling the error would cause anoying errors to the user
			// In worst case scenario, docker API will return an error when trying to pull or push the image because the authorization token is not found
			// TODO: log the error
			cfg, _ := config.LoadDefaultConfig(ctx, loadOptionsFuncs...)

			return cfg, nil
		}, credential)

	if token != nil {
		for _, a := range token.AuthorizationData {
			// if error is returned the process continues with out returning an error. It is prefered to return nil and continues exploring the other credentials.
			// TODO: log the error
			auth, _ := p.AuthMethod(aws.ToString(a.AuthorizationToken))

			if auth != nil {
				return auth, nil
			}
		}
	}

	return nil, nil
}

// AuthMethod returns a BasicAuthMethod having the username and password from authorization token
func (p *AWSECRAuthProvider) AuthMethod(authorizationToken string) (repository.AuthMethodReader, error) {

	errContext := "(credentials::provider::AWSECRAuthProvider::AuthMethod)"

	decodedToken, err := base64.StdEncoding.DecodeString(authorizationToken)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	parts := strings.SplitN(string(decodedToken), ":", 2)
	if len(parts) < 2 {
		return nil, errors.New(errContext, "Credentials could not be extracted from AWS token")
	}

	auth := &basic.BasicAuthMethod{
		Username: parts[0],
		Password: parts[1],
	}

	return auth, nil
}
