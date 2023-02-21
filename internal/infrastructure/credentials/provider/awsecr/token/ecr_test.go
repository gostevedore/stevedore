package token

import (
	"context"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awscredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/provider/awsecr/token/awscredprovider"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/provider/awsecr/token/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	errContext := "(token::AWSECRToken::Token)"
	tests := []struct {
		desc              string
		ecr               *AWSECRToken
		cfgFunc           func(context.Context, ...func(*config.LoadOptions) error) (aws.Config, error)
		badge             *credentials.Badge
		res               *ecr.GetAuthorizationTokenOutput
		prepareAssertFunc func(*AWSECRToken)
		err               error
	}{
		{
			desc: "Testing error when getting an authorization token with a nil badge",
			ecr:  NewAWSECRToken(),
			cfgFunc: func(context.Context, ...func(*config.LoadOptions) error) (aws.Config, error) {
				return aws.Config{}, nil
			},
			badge: nil,
			err:   errors.New(errContext, "To get an ECR authorization token, you must provide a badge"),
		},
		{
			desc: "Testing get ecr authorization token",
			ecr: NewAWSECRToken(
				WithAssumeRoleARNProvider(awscredprovider.NewMockAssumerRoleARNProvider()),
				WithStaticCredentialsProvider(awscredprovider.NewMockStaticCredentialsProvider()),
				WithECRClientFactory(
					NewECRClientFactory(
						func(cfg aws.Config) ECRClienter {
							c := client.NewMockECRClient()
							c.On("GetAuthorizationToken", context.TODO(), &ecr.GetAuthorizationTokenInput{}, mock.Anything).Return(
								&ecr.GetAuthorizationTokenOutput{},
								nil)

							return c
						})),
			),
			badge: &credentials.Badge{
				AWSAccessKeyID:     "accessKey",
				AWSSecretAccessKey: "secretKey",
			},
			prepareAssertFunc: func(ecr *AWSECRToken) {
				ecr.staticCredentialsProvider.(*awscredprovider.MockStaticCredentialsProvider).On(
					"CredentialsProvider",
					"accessKey",
					"secretKey",
					"",
					[]func(*config.LoadOptions) error{},
				).Return(
					awscredentials.StaticCredentialsProvider{},
					nil)
			},
			res: &ecr.GetAuthorizationTokenOutput{},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.ecr)
			}

			res, err := test.ecr.Get(context.TODO(), test.cfgFunc, test.badge)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestResolveCredentialsProvider(t *testing.T) {
	tests := []struct {
		desc              string
		ecr               *AWSECRToken
		config            aws.Config
		badge             *credentials.Badge
		options           []func(*config.LoadOptions) error
		res               aws.CredentialsProvider
		prepareAssertFunc func(*AWSECRToken)
		err               error
	}{
		{
			desc: "Testing resolve static credentials",
			ecr: NewAWSECRToken(
				WithAssumeRoleARNProvider(awscredprovider.NewMockAssumerRoleARNProvider()),
				WithStaticCredentialsProvider(awscredprovider.NewMockStaticCredentialsProvider()),
			),
			config: aws.Config{},
			badge: &credentials.Badge{
				AWSAccessKeyID:     "accessKey",
				AWSSecretAccessKey: "secretKey",
			},
			options: []func(*config.LoadOptions) error{},
			prepareAssertFunc: func(ecr *AWSECRToken) {
				ecr.staticCredentialsProvider.(*awscredprovider.MockStaticCredentialsProvider).On(
					"CredentialsProvider",
					"accessKey",
					"secretKey",
					"",
					[]func(*config.LoadOptions) error{},
				).Return(
					awscredentials.StaticCredentialsProvider{}, nil)
			},
			res: awscredentials.StaticCredentialsProvider{},
		},
		{
			desc: "Testing resolve assume role arn credentials",
			ecr: NewAWSECRToken(
				WithAssumeRoleARNProvider(awscredprovider.NewMockAssumerRoleARNProvider()),
				WithStaticCredentialsProvider(awscredprovider.NewMockStaticCredentialsProvider()),
			),
			config: aws.Config{},
			badge: &credentials.Badge{
				AWSRoleARN:         "arn:aws:iam::1234567890:role/testing-role",
				AWSAccessKeyID:     "accessKey",
				AWSSecretAccessKey: "secretKey",
			},
			options: []func(*config.LoadOptions) error{},
			prepareAssertFunc: func(ecr *AWSECRToken) {
				ecr.assumeRoleARNProvider.(*awscredprovider.MockAssumerRoleARNProvider).On(
					"CredentialsProvider",
					aws.Config{},
					"arn:aws:iam::1234567890:role/testing-role",
					"accessKey",
					"secretKey",
					"",
					[]func(*config.LoadOptions) error{},
				).Return(
					&stscreds.AssumeRoleProvider{}, nil)
			},
			res: &stscreds.AssumeRoleProvider{},
		},
		{
			desc: "Testing use default credentials chain",
			ecr: NewAWSECRToken(
				WithAssumeRoleARNProvider(awscredprovider.NewMockAssumerRoleARNProvider()),
				WithStaticCredentialsProvider(awscredprovider.NewMockStaticCredentialsProvider()),
			),
			config: aws.Config{},
			badge: &credentials.Badge{
				AWSUseDefaultCredentialsChain: true,
			},
			res: nil,
		},
		{
			desc: "Testing no credentials resolved",
			ecr: NewAWSECRToken(
				WithAssumeRoleARNProvider(awscredprovider.NewMockAssumerRoleARNProvider()),
				WithStaticCredentialsProvider(awscredprovider.NewMockStaticCredentialsProvider()),
			),
			config: aws.Config{},
			badge:  &credentials.Badge{},
			res:    &aws.CredentialsCache{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.ecr)
			}

			res, err := test.ecr.resolveCredentialsProvider(test.config, test.badge, test.options...)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.IsType(t, test.res, res)
			}
		})
	}
}

func TestAWSConfigLoadOption(t *testing.T) {
	tests := []struct {
		desc  string
		ecr   *AWSECRToken
		badge *credentials.Badge
		res   []func(*config.LoadOptions) error
		err   error
	}{
		{
			desc: "Testing load options",
			ecr:  &AWSECRToken{},
			badge: &credentials.Badge{
				AWSRegion:  "us-east-1",
				AWSProfile: "test-profile",
				AWSRoleARN: "arn:aws:iam::1234567890:role/testing-role",
				AWSSharedConfigFiles: []string{
					"/aws/config",
				},
				AWSSharedCredentialsFiles: []string{
					"/aws/credentials",
				},
				AWSUseDefaultCredentialsChain: true,
			},
			res: []func(*config.LoadOptions) error{
				config.WithSharedConfigProfile("test-profile"),
				config.WithRegion("us-east-1"),
				config.WithSharedConfigFiles([]string{
					"/aws/config",
				}),
				config.WithSharedCredentialsFiles([]string{
					"/aws/credentials",
				}),
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			options, err := test.ecr.awsConfigLoadOptions(test.badge)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, len(test.res), len(options))
			}
		})
	}
}
