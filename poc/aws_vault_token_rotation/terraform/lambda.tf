# # 4. Lambda Function
data "archive_file" "lambda_zip" {
  type        = "zip"
  source_dir  = "../lambda"
  output_path = "../lambda.zip"
}

resource "aws_lambda_function" "vault_lambda" {
  function_name = "vault-token-example"
  runtime       = "python3.11"
  handler       = "handler.lambda_handler"
  role          = aws_iam_role.lambda_role.arn
  filename      = "../lambda.zip"
  timeout       = 10
}