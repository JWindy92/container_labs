# Lambda Function
data "archive_file" "lambda_zip" {
  type        = "zip"
#   source_dir  = "../lambda"
  source_file = "bootstrap"
  output_path = "../lambda.zip"
}

resource "aws_lambda_function" "vault_lambda" {
  function_name = "vault-token-example"
  runtime       = "provided.al2"
  handler       = "bootstrap"
  role          = aws_iam_role.lambda_role.arn
  filename      = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  timeout       = 10

  environment {
        variables = {
            VAULT_ADDR = "http://vault:8200"
        }
    }
}

resource "aws_cloudwatch_log_group" "vault_lambda_log_group" {
  name              = "/aws/lambda/vault-token-example"
  retention_in_days = 7  # Set retention as needed
}