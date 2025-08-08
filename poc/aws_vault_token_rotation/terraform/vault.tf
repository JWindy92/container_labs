resource "vault_kv_secret_v2" "myapp_db_creds" {
#   mount = vault_mount.kv.path
  mount = "secret"
  name  = "someorg/someapp/db-creds"

  data_json = jsonencode({
    dbuser = "dbusername"
    dbpass = "dbpassword"
  })
}

resource "vault_policy" "secrets_admin_policy" {
  name = "admin"

  policy = <<EOT
path "*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
EOT
}

resource "vault_policy" "myapp_read_policy" {
  name = "myapp-read"

  policy = <<EOT
path "secret/data/someorg/someapp/*" {
  capabilities = ["read"]
}

path "secret/metadata/someorg/someapp/*" {
  capabilities = ["list"]
}
EOT
}

