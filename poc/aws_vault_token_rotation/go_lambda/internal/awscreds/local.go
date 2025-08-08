//go:build local
// +build local

package awscreds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AssumeRoleCreds struct {
	RoleArn string
}

func (a *AssumeRoleCreds) GetCreds(ctx context.Context) (aws.Credentials, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}
	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, a.RoleArn)
	return provider.Retrieve(ctx)
}

func NewCredentialsProvider() CredentialsInterface {
	return &AssumeRoleCreds{
		RoleArn: "arn:aws:iam::000000000000:role/vaultLambdaRole",
	}
}
