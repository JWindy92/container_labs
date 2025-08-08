resource "vault_auth_backend" "aws" {
  type = "aws"
  path = "aws"
}

resource "vault_aws_auth_backend_client" "example" {
  backend         = vault_auth_backend.aws.path
  access_key      = "test"
  secret_key      = "test"
  iam_server_id_header_value = "vault.example.com"  #? optional
  sts_region = "us-east-1"  #? optional, for LocalStack
  sts_endpoint = "http://localstack:4566"  #? optional, for LocalStack
#   rotation_schedule = "0 * * * SAT"
#   rotation_window   = 3600
}


resource "vault_aws_auth_backend_role" "lambda_role" {
  depends_on          = [aws_iam_role.lambda_role]
  backend             = vault_auth_backend.aws.path
  role                = aws_iam_role.lambda_role.name
  bound_iam_principal_arns = ["arn:aws:iam::000000000000:role/vaultLambdaRole"]
  auth_type           = "iam"
  resolve_aws_unique_ids = false
  token_policies      = ["default", vault_policy.myapp_read_policy.name]
}