package token

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/stretchr/testify/mock"
)

// MockAWSECRToken is a mock of AWSECRToken
type MockAWSECRToken struct {
	mock.Mock
}

// NewMockAWSECRToken creates a new mock
func NewMockAWSECRToken() *MockAWSECRToken {
	return &MockAWSECRToken{}
}

// Get is a mock function to get the authorization token
func (c *MockAWSECRToken) Get(ctx context.Context, cfgFunc func(context.Context, ...func(*config.LoadOptions) error) (aws.Config, error), badge *credentials.Credential) (*ecr.GetAuthorizationTokenOutput, error) {
	args := c.Called(ctx, cfgFunc, badge)
	if args.Error(1) != nil {
		return args.Get(0).(*ecr.GetAuthorizationTokenOutput), args.Error(1)
	}
	return args.Get(0).(*ecr.GetAuthorizationTokenOutput), nil
}
