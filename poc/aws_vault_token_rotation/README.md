

### CLI Commands

List all secretsmanager secrets
`aws secretsmanager list-secrets --region us-east-1`

Vault uname/pass login
`vault login -method=userpass username=john password=password`

Get secret
`vault kv get secret/someorg/someapp/db-creds`

vault write auth/aws/role/lambda-role \
    auth_type=iam \
    bound_iam_principal_arn="arn:aws:iam::000000000000:role/vaultLambdaRole" \
    policies="lambda-access-policy" \
    max_ttl="1h"