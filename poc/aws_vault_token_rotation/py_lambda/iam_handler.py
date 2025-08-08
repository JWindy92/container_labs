import boto3
import base64
import requests
import json
import datetime
from botocore.auth import SigV4Auth
from botocore.awsrequest import AWSRequest
from botocore.session import get_session

def headers_to_go_style(headers):
    retval = {}
    for k, v in headers.items():
        if isinstance(v, bytes):
            retval[k] = [str(v, 'ascii')]
        else:
            retval[k] = [v]
    return retval

def generate_vault_request(role_name=""):
    session = boto3.session.Session()
    # if you have credentials from non-default sources, call
    # session.set_credentials here, before calling session.create_client
    client = session.client('sts')
    endpoint = client._endpoint
    operation_model = client._service_model.operation_model('GetCallerIdentity')
    request_dict = client._convert_to_request_dict({}, operation_model)

    awsIamServerId = 'vault.example.com'
    request_dict['headers']['X-Vault-AWS-IAM-Server-ID'] = awsIamServerId

    request = endpoint.create_request(request_dict, operation_model)
    # It's now signed...
    return {
        'iam_http_request_method': request.method,
        'iam_request_url': str(base64.b64encode(request.url.encode('ascii')), 'ascii'),
        'iam_request_body': str(base64.b64encode(request.body.encode('ascii')), 'ascii'),
        'iam_request_headers': str(base64.b64encode(bytes(json.dumps(headers_to_go_style(dict(request.headers))), 'ascii')), 'ascii'), # It's a CaseInsensitiveDict, which is not JSON-serializable
        'role': role_name,
    }
# def get_request_payload(role_name: str):
#     session = get_session()
#     credentials = session.get_credentials()
#     region = session.get_config_variable("region")

#     request = AWSRequest(
#         method="POST",
#         url="https://sts.amazonaws.com/",
#         data="Action=GetCallerIdentity&Version=2011-06-15",
#         headers={"Content-Type": "application/x-www-form-urlencoded","X-Vault-AWS-IAM-Server-ID": "vault.example.com"}
#     )

#     SigV4Auth(credentials, "sts", region).add_auth(request)

#     headers = dict(request.headers)

#     payload = {
#         "iam_http_request_method": "POST",
#         "iam_request_url": base64.b64encode(request.url.encode()).decode(),
#         "iam_request_body": base64.b64encode(request.data.encode()).decode(),
#         "iam_request_headers": base64.b64encode(
#             json.dumps({k: [v] for k, v in headers.items()}).encode()
#         ).decode(),
#         "role": "vaultLambdaRole"
#     }

#     return payload

def login_to_vault(vault_addr: str, payload: dict):
    url = f"{vault_addr}/v1/auth/aws/login"
    headers={
            "Content-Type": "application/x-www-form-urlencoded"
        }

    response = requests.post(url, headers=headers, json=payload)
    print("Vault response:", response.status_code, response.text)
    response.raise_for_status()
    print(response.json())
    return response.json()["auth"]["client_token"]

def get_vault_secret(vault_addr: str, token: str, secret_path: str):
    url = f"{vault_addr}/v1/{secret_path}"
    headers = {"X-Vault-Token": token}

    response = requests.get(url, headers=headers)
    response.raise_for_status()

    # For KV v2, secret data is nested inside ["data"]["data"]
    return response.json()["data"]["data"]


if __name__ == "__main__":
    vault_addr = "http://localhost:8200"
    role_name = "vaultLambdaRole"
    secret_path = "secret/data/someorg/someapp/db-creds"  # KV v2 path

    print(json.dumps(generate_vault_request(role_name)))
    # vault_token = login_to_vault(vault_addr, payload)
    # secret = get_vault_secret(vault_addr, vault_token, secret_path)

    # print("Vault Token:", vault_token)
    # print("Secret Data:", secret)