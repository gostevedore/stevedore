package ecr

import (
	"context"
	"fmt"
	"testing"

	sdkcred "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gostevedore/stevedore/internal/core/domain/credentials"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/method/basic"
	"github.com/gostevedore/stevedore/internal/infrastructure/credentials/provider/awsecr/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	tests := []struct {
		desc              string
		credentials       *AWSECRCredentialsProvider
		prepareAssertFunc func(*AWSECRCredentialsProvider)
		badge             *credentials.Badge
		res               *basic.BasicAuthMethod
		err               error
	}{
		{
			desc:        "Testing get credentials from aws ecr credentials provisioner",
			credentials: NewAWSECRCredentialsProvider(token.NewMockAWSECRToken()),
			badge:       &credentials.Badge{},
			prepareAssertFunc: func(cred *AWSECRCredentialsProvider) {
				cred.tokenProvider.(*token.MockAWSECRToken).On("Get", context.TODO(), mock.Anything, &credentials.Badge{}).Return(&ecr.GetAuthorizationTokenOutput{
					AuthorizationData: []types.AuthorizationData{
						{
							AuthorizationToken: aws.String(`QVdTOmV3b2dJQ0FnY0dGNWJHOWhaRG9nWTBkR05XSkhPV2hhUVQwOUxBb2dJQ0FnWkdGMFlXdGxl
VG9nV2tkR01GbFhkR3hsVVQwOUxBb2cKSUNBZ2RtVnljMmx2YmpvZ01pd0tJQ0FnSUhSNWNHVTZJ
RVJCVkVGZlMwVlpMQW9nSUNBZ1pYaHdhWEpoZEdsdmJqb2dNVFkxTlRVeQpOVE0yTmdwOQ==`),
						},
					},
				}, nil)
			},
			res: &basic.BasicAuthMethod{
				Username: "AWS",
				Password: `ewogICAgcGF5bG9hZDogY0dGNWJHOWhaQT09LAogICAgZGF0YWtleTogWkdGMFlXdGxlUT09LAog
ICAgdmVyc2lvbjogMiwKICAgIHR5cGU6IERBVEFfS0VZLAogICAgZXhwaXJhdGlvbjogMTY1NTUy
NTM2Ngp9`,
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.credentials)
			}

			res, err := test.credentials.Get(test.badge)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}
}

func TestAuthMethod(t *testing.T) {

	errContext := "(credentials::provider::AWSECRCredentialsProvider::AuthMethod)"

	tests := []struct {
		desc        string
		credentials *AWSECRCredentialsProvider
		// prepareAssertFunc func(*AWSECRCredentialsProvider)
		token string
		res   *basic.BasicAuthMethod
		err   error
	}{
		{
			desc:        "Testing get auth method from authorization token",
			credentials: NewAWSECRCredentialsProvider(token.NewMockAWSECRToken()),
			token: `QVdTOmV3b2dJQ0FnY0dGNWJHOWhaRG9nWTBkR05XSkhPV2hhUVQwOUxBb2dJQ0FnWkdGMFlXdGxl
VG9nV2tkR01GbFhkR3hsVVQwOUxBb2cKSUNBZ2RtVnljMmx2YmpvZ01pd0tJQ0FnSUhSNWNHVTZJ
RVJCVkVGZlMwVlpMQW9nSUNBZ1pYaHdhWEpoZEdsdmJqb2dNVFkxTlRVeQpOVE0yTmdwOQ==`,
			res: &basic.BasicAuthMethod{
				Username: "AWS",
				Password: `ewogICAgcGF5bG9hZDogY0dGNWJHOWhaQT09LAogICAgZGF0YWtleTogWkdGMFlXdGxlUT09LAog
ICAgdmVyc2lvbjogMiwKICAgIHR5cGU6IERBVEFfS0VZLAogICAgZXhwaXJhdGlvbjogMTY1NTUy
NTM2Ngp9`,
			},
			err: &errors.Error{},
		},
		{
			desc:        "Testing error when authorization token is not valid",
			credentials: NewAWSECRCredentialsProvider(token.NewMockAWSECRToken()),
			token:       `bm90dmFpbGRhdXRob3JpemF0aW9udG9rZW4=`,
			err:         errors.New(errContext, "Credentials could not be extracted from AWS token"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			pass, err := test.credentials.AuthMethod(test.token)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.res, pass)
			}
		})
	}
}

// TestFunctional is kept as skiped to avoid running this test automatically
func TestFunctional(t *testing.T) {

	t.Skip()

	// credprov := NewAWSECRCredentialsProvider(token.NewAWSECRToken(
	// 	token.WithStaticCredentialsProvider(awscredprovider.NewStaticCredentialsProvider()),
	// 	token.WithAssumeRoleARNProvider(awscredprovider.NewAssumerRoleARNProvider()),
	// 	token.WithECRClientFactory(
	// 		token.NewECRClientFactory(
	// 			func(cfg aws.Config) token.ECRClienter {
	// 				c := ecr.NewFromConfig(cfg)

	// 				return c
	// 			}),
	// 	),
	// ))

	// badge := &credentials.Badge{
	// 	AWSUseDefaultCredentialsChain: true,
	// }

	// auth, err := credprov.Get(badge)
	// if err != nil {
	// 	t.Error(err)
	// } else {
	// 	t.Log(auth)
	// }

	assumecnf, err := config.LoadDefaultConfig(
		context.TODO(), config.WithRegion("eu-west-1"),
		config.WithCredentialsProvider(aws.NewCredentialsCache(
			sdkcred.NewStaticCredentialsProvider(
				"AKIASWM3KPCNAVF5AOM5",
				"FbktdeFJLd7sVtM9laVTw9QvCqcJhTuKydcMm/m1",
				"",
			)),
		),
	)
	if err != nil {
		t.Error(err)
	}

	stsclient := sts.NewFromConfig(assumecnf)
	cfg, err := config.LoadDefaultConfig(
		context.TODO(), config.WithRegion("eu-west-1"),
		config.WithCredentialsProvider(aws.NewCredentialsCache(
			stscreds.NewAssumeRoleProvider(
				stsclient,
				// "arn:aws:iam::486857178526:role/role-admin-admin-access",
				"arn:aws:iam::460972834675:role/role-lab01-admin-access",
			)),
		),
	)
	if err != nil {
		t.Error(err)
	}

	client := ecr.NewFromConfig(cfg)

	auth, err := client.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v", auth.AuthorizationData)

}
