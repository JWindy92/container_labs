import boto3
import requests
import os

aws_secrets = boto3.client('secretsmanager')

VAULT_ADDR = os.getenv("VAULT_ADDR", "http://localhost:8200")

def get_vault_token():
    response = aws_secrets.get_secret_value(SecretId='vault-access-token')
    print(f"Using Vault token: {response['SecretString']}")
    return response['SecretString']

def lambda_handler(event, context):
    token = get_vault_token()
    headers = {"X-Vault-Token": token}
    vault_url = f"https://{VAULT_ADDR}/v1/secret/data/myapp"

    response = requests.get(vault_url, headers=headers)
    return {
        "statusCode": response.status_code,
        "body": response.json()
    }

if __name__ == "__main__":
    # get_vault_token()
    lambda_handler({}, {})