# Configure IAM Policies to allow Lambda to authenticate with Vault using AWS IAM Auth Method



# IAM Role for Lambda
resource "aws_iam_role" "lambda_role" {
  name = "vaultLambdaRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}



# IAM Policy for accessing Secrets Manager
resource "aws_iam_policy" "vault_iam_auth_policy" {
  name = "VaultSecretsAccess"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        "Effect": "Allow",
        "Action": ["sts:GetCallerIdentity"],
        "Resource": "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "attach_lambda_vault_auth_policy" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.vault_iam_auth_policy.arn
}

resource "aws_iam_role_policy_attachment" "attach_lambda_basic_execution" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

output "lambda_role_arn" {
  value = aws_iam_role.lambda_role.arn
}