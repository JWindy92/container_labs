## Vault configuration


# Create a KV v2 secrets engine at path "secret/"
resource "vault_kv_secret_v2" "myapp_db_creds" {
  mount = "secret"
  name  = "someorg/someapp/db-creds"

  data_json = jsonencode({
    dbuser = "dbusername"
    dbpass = "dbpassword"
  })
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

