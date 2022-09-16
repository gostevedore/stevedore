package awscredprovider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AssumerRoleARNProvider struct{}

func NewAssumerRoleARNProvider() *AssumerRoleARNProvider {
	return &AssumerRoleARNProvider{}
}

func (p *AssumerRoleARNProvider) CredentialsProvider(cfg aws.Config, roleARN, awsAccessKeyID, awsSecretAccessKey, session string, options ...func(*config.LoadOptions) error) (aws.CredentialsProvider, error) {

	var stsclient *sts.Client

	if awsAccessKeyID != "" && awsSecretAccessKey != "" {
		options = append(options,
			config.WithCredentialsProvider(aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(
					awsAccessKeyID,
					awsSecretAccessKey,
					session,
				)),
			),
		)

		assumecnf, err := config.LoadDefaultConfig(
			context.Background(),
			options...,
		)

		if err != nil {
			return nil, err
		}

		stsclient = sts.NewFromConfig(assumecnf)
	} else {
		stsclient = sts.NewFromConfig(cfg)
	}

	provider := stscreds.NewAssumeRoleProvider(stsclient, roleARN)

	return provider, nil
}
