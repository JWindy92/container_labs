//go:build deploy
// +build deploy

package awscreds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type DefaultCreds struct{}

func (d *DefaultCreds) GetCreds(ctx context.Context) (aws.Credentials, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}
	return cfg.Credentials.Retrieve(ctx)
}

func NewCredentialsProvider() CredentialsInterface {
	return &DefaultCreds{}
}
