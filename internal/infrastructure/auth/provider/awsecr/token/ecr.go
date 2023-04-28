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

func defaultCfgFunc(ctx context.Context, options ...func(*config.LoadOptions) error) (aws.Config, error) {

	errContext := "(token::AWSECRToken::defaultCfgFunc)"

	cfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		errors.New(errContext, "", err)
	}

	return cfg, nil
}

// Get return the authorization token
func (t *AWSECRToken) Get(ctx context.Context, cfgFunc func(context.Context, ...func(*config.LoadOptions) error) (aws.Config, error), credential *credentials.Credential) (*ecr.GetAuthorizationTokenOutput, error) {

	errContext := "(token::AWSECRToken::Token)"

	if cfgFunc == nil {
		cfgFunc = defaultCfgFunc
	}

	if credential == nil {
		return nil, errors.New(errContext, "To get an ECR authorization token, you must provide a credential")
	}

	options, err := t.awsConfigLoadOptions(credential)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	cfg, err := cfgFunc(ctx, options...)
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	credentialsProvider, err := t.resolveCredentialsProvider(cfg, credential, options...)
	if err != nil {
		errors.New(errContext, "", err)
	}

	if credentialsProvider != nil {
		// when resolve credentials provider returns an empty aws.CredentialsCache it means that no credentials provider was found
		// if reflect.TypeOf(credentialsProvider) == reflect.TypeOf(&aws.CredentialsCache{}) {
		if reflect.DeepEqual(credentialsProvider, &aws.CredentialsCache{}) {
			return nil, nil
		}

		cfg.Credentials = aws.NewCredentialsCache(credentialsProvider)
	}

	client := t.ecrClientFactory.Client(cfg)
	auth, err := client.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		errors.New(errContext, "", err)
	}

	return auth, nil
}

// awsConfigLoadOptions returns a list of load options to use when loading the default config.
func (t *AWSECRToken) awsConfigLoadOptions(credential *credentials.Credential) ([]func(*config.LoadOptions) error, error) {

	errContext := "(token::AWSECRToken::awsConfigLoadOptions)"

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
func (t *AWSECRToken) resolveCredentialsProvider(cfg aws.Config, credential *credentials.Credential, options ...func(*config.LoadOptions) error) (aws.CredentialsProvider, error) {
	var provider aws.CredentialsProvider
	var err error
	errContext := "(token::AWSECRToken::resolveCredentialsProvider)"

	if credential == nil {
		return nil, errors.New(errContext, "To get an ECR authorization token, you must provide a credential")
	}

	if credential.AWSRoleARN != "" {
		if t.assumeRoleARNProvider != nil {

			provider, err = t.assumeRoleARNProvider.CredentialsProvider(cfg, credential.AWSRoleARN, credential.AWSAccessKeyID, credential.AWSSecretAccessKey, "", options...)
			if err != nil {
				return nil, errors.New(errContext, "", err)
			}

			return provider, nil
		}
	}

	if credential.AWSAccessKeyID != "" && credential.AWSSecretAccessKey != "" {
		if t.staticCredentialsProvider != nil {
			provider, err = t.staticCredentialsProvider.CredentialsProvider(credential.AWSAccessKeyID, credential.AWSSecretAccessKey, "", options...)
			if err != nil {
				return nil, errors.New(errContext, "", err)
			}

			return provider, nil
		}
	}

	if credential.AWSUseDefaultCredentialsChain {
		return nil, nil
	}

	return &aws.CredentialsCache{}, nil
}
