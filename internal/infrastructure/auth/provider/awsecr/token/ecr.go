package token

import (
	"context"
	"reflect"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

type OptionsFunc func(*AWSECRToken)

// AWSECRToken is the interface for the ECR client that generates the authorization token
type AWSECRToken struct {
	staticCredentialsProvider StaticCredentialsProviderer
	assumeRoleARNProvider     AssumerRoleARNProviderer
	ecrClientFactory          *ECRClientFactory
}

// NewAWSECRToken creates a new ECR client
func NewAWSECRToken(options ...OptionsFunc) *AWSECRToken {
	token := &AWSECRToken{}

	for _, option := range options {
		option(token)
	}

	return token
}

// WithStaticCredentialsProvider is a function that sets the static credentials provider
func WithStaticCredentialsProvider(staticCredentialsProviderer StaticCredentialsProviderer) OptionsFunc {
	return func(token *AWSECRToken) {
		token.staticCredentialsProvider = staticCredentialsProviderer
	}
}

// WithAssumeRoleARNProvider is a function that sets the assume role ARN provider
func WithAssumeRoleARNProvider(assumeRoleARNProvider AssumerRoleARNProviderer) OptionsFunc {
	return func(token *AWSECRToken) {
		token.assumeRoleARNProvider = assumeRoleARNProvider
	}
}

// WithECRClientFactory is a function that sets the ECR client factory
func WithECRClientFactory(ecrClientFactory *ECRClientFactory) OptionsFunc {
	return func(token *AWSECRToken) {
		token.ecrClientFactory = ecrClientFactory
	}
}

func defaultAWSConfigFunc(ctx context.Context, options ...func(*config.LoadOptions) error) (aws.Config, error) {

	// if LoadDefaultConfig returns an error the process continues with out returning an error. It is prefered to return nil and let the caller to decide what to do. In worst case scenario, docker API will return an error when trying to pull or push the image because the authorization token is not found
	// TODO: log the error
	cfg, _ := config.LoadDefaultConfig(ctx, options...)

	return cfg, nil
}

// Get return the authorization token
func (token *AWSECRToken) Get(ctx context.Context, AWSConfigFunc func(context.Context, ...func(*config.LoadOptions) error) (aws.Config, error), credential *credentials.Credential) (*ecr.GetAuthorizationTokenOutput, error) {

	errContext := "(token::AWSECRToken::Token)"

	if credential == nil {
		return nil, errors.New(errContext, "To get an ECR authorization token, you must provide a credential")
	}

	if AWSConfigFunc == nil {
		AWSConfigFunc = defaultAWSConfigFunc
	}

	options, err := token.getAWSConfigLoadOptionFuncs(credential)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	AWSConfig, err := AWSConfigFunc(ctx, options...)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	// The error is not being controled because is prefered to return nil and let the caller to decide what to do. AWS would be required for very specific cases and controlling the error would cause anoying errors to the user
	// In worst case scenario, docker API will return an error when trying to pull or push the image because the authorization token is not found
	// TODO: log the error
	credentialsProvider, _ := token.resolveCredentialsProvider(AWSConfig, credential, options...)

	if credentialsProvider != nil {
		// when resolve credentials provider returns an empty aws.CredentialsCache it means that no credentials provider was found
		// if reflect.TypeOf(credentialsProvider) == reflect.TypeOf(&aws.CredentialsCache{}) {
		if reflect.DeepEqual(credentialsProvider, &aws.CredentialsCache{}) {
			return nil, nil
		}

		// configure the AWSConfig with the credentials provider
		AWSConfig.Credentials = aws.NewCredentialsCache(credentialsProvider)
	}

	client := token.ecrClientFactory.Client(AWSConfig)
	// The error is not being controled because is prefered to return nil and let the caller to decide what to do. AWS would be required for very specific cases and controlling the error would cause anoying errors to the user
	// In worst case scenario, docker API will return an error when trying to pull or push the image because the authorization token is not found
	// TODO: log the error
	auth, _ := client.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})

	return auth, nil
}

// getAWSConfigLoadOptionFuncs returns a list of functions to configure the default config
func (token *AWSECRToken) getAWSConfigLoadOptionFuncs(credential *credentials.Credential) ([]func(*config.LoadOptions) error, error) {

	errContext := "(token::AWSECRToken::getAWSConfigLoadOptionFuncs)"

	if credential == nil {
		return nil, errors.New(errContext, "To get an ECR authorization token, you must provide a credential")
	}

	optFuncs := []func(*config.LoadOptions) error{}

	if credential.AWSRegion != "" {
		optFuncs = append(optFuncs, config.WithRegion(credential.AWSRegion))
	}

	if credential.AWSProfile != "" {
		optFuncs = append(optFuncs, config.WithSharedConfigProfile(credential.AWSProfile))
	}

	if len(credential.AWSSharedCredentialsFiles) > 0 {
		optFuncs = append(optFuncs, config.WithSharedCredentialsFiles(credential.AWSSharedCredentialsFiles))
	}

	if len(credential.AWSSharedConfigFiles) > 0 {
		optFuncs = append(optFuncs, config.WithSharedConfigFiles(credential.AWSSharedConfigFiles))
	}

	return optFuncs, nil
}

// resolveCredentialsProvider returns the credentials provider to use for the given config. To use the default aws configuration, is returned a nil provider. When no provider is found it is returned an empty CredentialsCache.
func (token *AWSECRToken) resolveCredentialsProvider(cfg aws.Config, credential *credentials.Credential, options ...func(*config.LoadOptions) error) (aws.CredentialsProvider, error) {
	var provider aws.CredentialsProvider
	var err error
	errContext := "(token::AWSECRToken::resolveCredentialsProvider)"

	if credential == nil {
		return nil, errors.New(errContext, "To get an ECR authorization token, you must provide a credential")
	}

	if credential.AWSRoleARN != "" {
		if token.assumeRoleARNProvider != nil {

			provider, err = token.assumeRoleARNProvider.CredentialsProvider(cfg, credential.AWSRoleARN, credential.AWSAccessKeyID, credential.AWSSecretAccessKey, "", options...)
			if err != nil {
				return nil, errors.New(errContext, "", err)
			}

			return provider, nil
		}
	}

	if credential.AWSAccessKeyID != "" && credential.AWSSecretAccessKey != "" {
		if token.staticCredentialsProvider != nil {
			provider, err = token.staticCredentialsProvider.CredentialsProvider(credential.AWSAccessKeyID, credential.AWSSecretAccessKey, "", options...)
			if err != nil {
				return nil, errors.New(errContext, "", err)
			}

			return provider, nil
		}
	}

	if credential.AWSUseDefaultCredentialsChain {
		// return a nil provider to use the default aws configuration
		return nil, nil
	}

	return &aws.CredentialsCache{}, nil
}
