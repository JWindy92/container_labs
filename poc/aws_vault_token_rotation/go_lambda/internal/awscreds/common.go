package awscreds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type CredentialsInterface interface {
	GetCreds(ctx context.Context) (aws.Credentials, error)
}
