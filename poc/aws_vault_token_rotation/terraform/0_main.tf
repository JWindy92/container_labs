provider "aws" {
    region = "us-east-1"
    access_key = "test"
    secret_key = "test"
  # LocalStack configuration
    endpoints {
        secretsmanager = "http://localhost:4566"
        lambda         = "http://localhost:4566"
        iam            = "http://localhost:4566"
        logs           = "http://localhost:4566"
        sts            = "http://localhost:4566"
    }
}

# provider "vault" {
#   address = "http://localhost:8200"
#   token   = var.vault_token
#   version = ">= 3.15.0" 
# }

terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = ">= 5.1.0" 
    }
  }
}