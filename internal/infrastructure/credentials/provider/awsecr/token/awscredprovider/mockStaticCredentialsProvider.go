package awscredprovider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/stretchr/testify/mock"
)

type MockStaticCredentialsProvider struct {
	mock.Mock
}

func NewMockStaticCredentialsProvider() *MockStaticCredentialsProvider {
	return &MockStaticCredentialsProvider{}
}

func (p *MockStaticCredentialsProvider) CredentialsProvider(key, secret, session string, options ...func(*config.LoadOptions) error) (aws.CredentialsProvider, error) {

	args := p.Called(key, secret, session, options)

	return args.Get(0).(credentials.StaticCredentialsProvider), nil
}
