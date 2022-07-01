package token

import "github.com/aws/aws-sdk-go-v2/aws"

type ECRClientFactory struct {
	factoryFunc func(aws.Config) ECRClienter
}

func NewECRClientFactory(factoryFunc func(aws.Config) ECRClienter) *ECRClientFactory {
	return &ECRClientFactory{
		factoryFunc: factoryFunc,
	}
}

func (f *ECRClientFactory) Client(cfg aws.Config) ECRClienter {
	return f.factoryFunc(cfg)
}
