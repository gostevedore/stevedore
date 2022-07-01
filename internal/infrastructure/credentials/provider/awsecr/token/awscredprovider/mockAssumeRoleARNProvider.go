package awscredprovider

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/stretchr/testify/mock"
)

type MockAssumerRoleARNProvider struct {
	mock.Mock
}

func NewMockAssumerRoleARNProvider() *MockAssumerRoleARNProvider {
	return &MockAssumerRoleARNProvider{}
}

func (p *MockAssumerRoleARNProvider) Credentials(cfg aws.Config, roleARN string) (aws.CredentialsProvider, error) {

	args := p.Called(cfg, roleARN)

	return args.Get(0).(*stscreds.AssumeRoleProvider), nil
}
