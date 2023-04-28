package ecr

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
)

// AWSECRTokenProvider is the interface for the ECR client that generates the authorization token
type AWSECRTokenProvider interface {
	Get(ctx context.Context, cfgFunc func(context.Context, ...func(*config.LoadOptions) error) (aws.Config, error), badge *credentials.Credential) (*ecr.GetAuthorizationTokenOutput, error)
}
