package awscredprovider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type StaticCredentialsProvider struct{}

func NewStaticCredentialsProvider() *StaticCredentialsProvider {
	return &StaticCredentialsProvider{}
}

func (p *StaticCredentialsProvider) Credentials(key, secret, session string) (aws.CredentialsProvider, error) {

	provider := credentials.NewStaticCredentialsProvider(key, secret, session)

	return provider, nil
}
