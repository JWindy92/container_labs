# # 1. Store Vault token in Secrets Manager
resource "aws_secretsmanager_secret" "vault_token" {
  name = "vault-access-token"
}

resource "aws_secretsmanager_secret_version" "vault_token_value" {
  secret_id     = aws_secretsmanager_secret.vault_token.id
  secret_string = "your-initial-vault-token"
}