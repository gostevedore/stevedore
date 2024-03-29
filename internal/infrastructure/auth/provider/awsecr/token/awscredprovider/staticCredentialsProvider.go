package awscredprovider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type StaticCredentialsProvider struct{}

func NewStaticCredentialsProvider() *StaticCredentialsProvider {
	return &StaticCredentialsProvider{}
}

func (p *StaticCredentialsProvider) CredentialsProvider(key, secret, session string, options ...func(*config.LoadOptions) error) (aws.CredentialsProvider, error) {
	provider := credentials.NewStaticCredentialsProvider(key, secret, session)
	return provider, nil
}
