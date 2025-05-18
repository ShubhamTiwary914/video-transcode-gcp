import os
from typing import Optional
from datetime import timedelta

from google import auth
from google.auth.transport import requests
from google.cloud.storage import Client
from fastapi import FastAPI, Request
from google.oauth2 import service_account

from dotenv import load_dotenv


load_dotenv()
app = FastAPI()


@app.get("/health")
def health_check():
    return {"status": "Healthy (GET /health received)"}


@app.get("/sign")
def fetch_signURL(bucket: str, filename: str):
    return GCS_signedURL_SA(bucket, filename, method="GET") 
@app.post("/sign")
def fetch_signURL(bucket: str, filename: str):
    return GCS_signedURL_SA(bucket, filename, method="PUT") 


@app.get("/sign/key")
def fetch_signURL(bucket: str, filename: str):
    return GCS_signedURL_keyfile(bucket, filename, method="GET")
@app.post("/sign/key")
def fetch_signURL(bucket: str, filename: str):
    return GCS_signedURL_keyfile(bucket, filename, method="PUT")



def GCS_signedURL_SA(bucket: str, blob: str,*, content_type="video/mp4",
            exp: Optional[timedelta] = None, min_size=1, max_size=int(1e8), method="PUT"):  
    """
        Generate GCS (PUT) signed URL (without key file) - with SA
    """
    if exp is None:
        exp = timedelta(hours=1)
    credentials, _ = auth.default()
    if credentials.token is None: 
        credentials.refresh(requests.Request())
    client = Client()
    bucket = client.get_bucket(bucket)
    blob = bucket.blob(blob)
    return blob.generate_signed_url(
        version="v4",
        expiration=exp,
        service_account_email=credentials.service_account_email,
        access_token=credentials.token,
        method=method,
        content_type=content_type,
        headers={"Content-Type": content_type}
    )


def GCS_signedURL_keyfile(bucket: str, blob: str,*, 
        content_type="video/mp4", exp: Optional[timedelta] = None, min_size=1, max_size=int(1e8), method="PUT"): 
    """
        Generate GCS (PUT) signed URL with SA key file 
    """
    if exp is None:
        exp = timedelta(hours=1)

    sa_path =  os.getenv('SA_FILE_PATH') 
    credentials = service_account.Credentials.from_service_account_file(sa_path)
    client = Client(credentials=credentials, project=credentials.project_id)
    bucket = client.get_bucket(bucket)
    blob = bucket.blob(blob)

    return blob.generate_signed_url(
        version="v4",
        expiration=exp,
        credentials=credentials,
        method=method,
        content_type=content_type,
        headers={"Content-Type": content_type},
    )