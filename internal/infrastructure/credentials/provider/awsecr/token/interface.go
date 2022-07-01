package token

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

// StaticCredentialsProvider is an interface that provides static credentials credentials provider.
type StaticCredentialsProviderer interface {
	Credentials(key, secret, session string) (aws.CredentialsProvider, error)
}

// AssumerRoleARNProviderer is an interface that provides AssumerRoleARN credendials provider.
type AssumerRoleARNProviderer interface {
	Credentials(cfg aws.Config, roleARN string) (aws.CredentialsProvider, error)
}

// ECRClienter is an interface that provides ECR client.
type ECRClienter interface {
	GetAuthorizationToken(ctx context.Context, input *ecr.GetAuthorizationTokenInput, options ...func(*ecr.Options)) (*ecr.GetAuthorizationTokenOutput, error)
}
