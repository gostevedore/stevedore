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

// Get return user password auth for docker registry
func (p *AWSECRAuthProvider) Get(credential *credentials.Credential) (repository.AuthMethodReader, error) {
	errContext := "(credentials::provider::AWSECRAuthProvider::Get)"

	token, err := p.tokenProvider.Get(context.TODO(),
		func(ctx context.Context, loadOptionsFuncs ...func(*config.LoadOptions) error) (aws.Config, error) {
			cfg, err := config.LoadDefaultConfig(ctx, loadOptionsFuncs...)
			if err != nil {
				errors.New(errContext, "", err)
			}

			return cfg, nil
		}, credential)
	if err != nil {
		errors.New(errContext, "", err)
	}

	if token != nil {
		for _, a := range token.AuthorizationData {
			auth, err := p.AuthMethod(aws.ToString(a.AuthorizationToken))
			if err != nil {
				errors.New(errContext, "", err)
			}

			if auth != nil {
				return auth, nil
			}
		}
	}

	return nil, nil
}

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
