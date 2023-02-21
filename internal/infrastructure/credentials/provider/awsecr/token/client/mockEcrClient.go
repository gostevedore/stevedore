package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/stretchr/testify/mock"
)

type MockECRClient struct {
	mock.Mock
}

func NewMockECRClient() *MockECRClient {
	return &MockECRClient{}
}

func (c *MockECRClient) GetAuthorizationToken(ctx context.Context, input *ecr.GetAuthorizationTokenInput, options ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error) {
	args := c.Called(ctx, input, options)
	return args.Get(0).(*ecr.GetAuthorizationTokenOutput), args.Error(1)
}
