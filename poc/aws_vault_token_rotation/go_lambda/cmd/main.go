package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JWindy92/golang_vault_iam/internal/awscreds"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

// type CredentialsInterface interface {
// 	GetCreds(ctx context.Context) (aws.Credentials, error)
// }

func main() {
	// lambda.Start(handler)
	run_v1()
}

func handler(ctx context.Context) (map[string]interface{}, error) {
	resp, err := getSecretWithAWSAuthIAM()
	if err != nil {
		return nil, fmt.Errorf("error getting secret: %w", err)
	}
	data := map[string]interface{}{
		"secret": resp,
	}
	return data, nil
}

func formatHeaders(h http.Header) string {
	var b bytes.Buffer
	for k, vs := range h {
		for _, v := range vs {
			b.WriteString(fmt.Sprintf("%s:%s\n", k, v))
		}
	}
	return b.String()
}

func PrettyPrint(v interface{}) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("PrettyPrint error:", err)
		return
	}
	fmt.Println(string(bytes))
}

func run_v1() {
	ctx := context.Background()

	vaultAddr := "http://127.0.0.1:8200"
	vaultRole := "vaultLambdaRole"
	secretPath := "secret/data/someorg/someapp/db-creds"
	awsRegion := "us-east-1"

	credsProvider := awscreds.NewCredentialsProvider()

	assumedCreds, err := credsProvider.GetCreds(ctx)
	if err != nil {
		log.Fatalf("failed to get AWS credentials: %v", err)
	}

	// Create HTTP request for STS GetCallerIdentity
	body := "Action=GetCallerIdentity&Version=2011-06-15"
	req, err := http.NewRequest("POST", "https://sts.amazonaws.com/", bytes.NewBufferString(body))
	if err != nil {
		log.Fatalf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("X-Vault-AWS-IAM-Server-ID", "vault.example.com")

	// Sign the request with AWS SigV4 signer v2
	signer := v4.NewSigner()
	err = signer.SignHTTP(ctx, assumedCreds, req, body, "sts", awsRegion, time.Now())
	if err != nil {
		log.Fatalf("failed to sign request: %v", err)
	}

	// Prepare Vault login payload with base64-encoded request components
	identityRequest := map[string]interface{}{
		"iam_http_request_method": req.Method,
		"iam_request_url":         base64.StdEncoding.EncodeToString([]byte(req.URL.String())),
		"iam_request_body":        base64.StdEncoding.EncodeToString([]byte(body)),
		"iam_request_headers":     req.Header,
		"role":                    vaultRole,
	}

	PrettyPrint(identityRequest)
	// Create Vault client
	vaultClient, err := vault.NewClient(&vault.Config{Address: vaultAddr})
	if err != nil {
		log.Fatalf("failed to create Vault client: %v", err)
	}

	// Login to Vault AWS auth backend
	secret, err := vaultClient.Logical().Write("auth/aws/login", identityRequest)
	if err != nil {
		log.Fatalf("vault AWS login failed: %v", err)
	}
	if secret == nil || secret.Auth == nil {
		log.Fatalf("vault login response missing auth data")
	}

	token := secret.Auth.ClientToken
	fmt.Printf("Vault token: %s\n", token)

	// Use Vault token to read secret
	vaultClient.SetToken(token)
	secretData, err := vaultClient.Logical().Read(secretPath)
	if err != nil {
		log.Fatalf("failed to read secret at %s: %v", secretPath, err)
	}
	if secretData == nil || secretData.Data == nil {
		log.Fatalf("secret data not found at %s", secretPath)
	}

	PrettyPrint(secretData.Data)
}

func getSecretWithAWSAuthIAM() (string, error) {
	vaultAddr := "http://127.0.0.1:8200"
	vaultRole := "vaultLambdaRole"
	secretPath := "secret/data/someorg/someapp/db-creds"
	awsRegion := "us-east-1"

	// config := vault.DefaultConfig()             // modify for more granular configuration
	config := &vault.Config{Address: vaultAddr} // modify for more granular configuration

	client, err := vault.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	awsAuth, err := auth.NewAWSAuth(
		auth.WithRole(vaultRole), // if not provided, Vault will fall back on looking for a role with the IAM role name if you're using the iam auth type, or the EC2 instance's AMI id if using the ec2 auth type
		auth.WithRegion(awsRegion),
		auth.WithIAMServerIDHeader("vault.example.com"),
	)
	if err != nil {
		return "", fmt.Errorf("unable to initialize AWS auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.Background(), awsAuth)

	if err != nil {
		return "", fmt.Errorf("unable to login to AWS auth method: %w", err)
	}
	if authInfo == nil {
		return "", fmt.Errorf("no auth info was returned after login")
	}
	PrettyPrint(authInfo)
	// get secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := client.KVv2("secret").Get(context.Background(), secretPath)
	if err != nil {
		return "", fmt.Errorf("unable to read secret: %w", err)
	}

	// data map can contain more than one key-value pair,
	// in this case we're just grabbing one of them
	value, ok := secret.Data["password"].(string)
	if !ok {
		return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data["password"], secret.Data["password"])
	}

	return value, nil
}
