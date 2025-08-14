# Setup userpass authentication for demonstration purposes


resource "vault_auth_backend" "userpass" {
  type = "userpass"
  path = "userpass"
}

resource "vault_generic_endpoint" "myapp_user" {
  path = "auth/userpass/users/john"
  data_json = jsonencode({
    password = "password"
    policies = [vault_policy.secrets_admin_policy.name]
  })
  depends_on = [vault_auth_backend.userpass]
}

resource "vault_policy" "secrets_admin_policy" {
  name = "admin"

  policy = <<EOT
path "*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
EOT
}