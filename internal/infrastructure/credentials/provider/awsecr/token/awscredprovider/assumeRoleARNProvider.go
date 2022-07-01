package awscredprovider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AssumerRoleARNProvider struct{}

func NewAssumerRoleARNProvider() *AssumerRoleARNProvider {
	return &AssumerRoleARNProvider{}
}

func (p *AssumerRoleARNProvider) Credentials(cfg aws.Config, roleARN string) (aws.CredentialsProvider, error) {

	stsclient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsclient, roleARN)

	return provider, nil
}
