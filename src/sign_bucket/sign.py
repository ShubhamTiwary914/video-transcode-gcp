from typing import Optional, Union
from datetime import timedelta

from google import auth
from google.auth.transport import requests
from google.cloud.storage import Client
from fastapi import FastAPI
from pydantic import BaseModel


class SignParams(BaseModel):
    bucket: str
    filename: str
   


app = FastAPI()



@app.get("/health")
def health_check():
    return {"status": "Healthy (GET /health received)"}


@app.get("/")
def fetch_signURL(bucket: str, filename: str, content_type: str = "video/mp4"):
    return { "GCS-URI": make_signed_upload_url(bucket, filename, content_type) }


def make_signed_upload_url(bucket: str, blob: str,*, content_type="video/mp4",
            exp: Optional[timedelta] = None, min_size=1, max_size=int(1e8)): 
    """
    fetch GCS signed URL without private key (with GCP-SA) 
    ----------
    bucket : str
        name of the GCS bucket the signed URL will reference.
    blob : str
        Name of the GCS blob (in `bucket`) the signed URL will reference.
    exp : timedelta, optional
        Time from now when the signed url will expire.
    content_type : str, optional
        The required mime type of the data that is uploaded to the generated
        signed url.
    min_size : int, optional
        The minimum size the uploaded file can be, in bytes (inclusive).
        If the file is smaller than this, GCS will return a 400 code on upload.
    max_size : int, optional (100 mb default)
        The maximum size the uploaded file can be, in bytes (inclusive).
        If the file is larger than this, GCS will return a 400 code on upload.
    """
    if exp is None:
        exp = timedelta(hours=1)
    credentials, project_id = auth.default()
    if credentials.token is None: 
        credentials.refresh(requests.Request())
    client = Client()
    bucket = client.get_bucket(bucket)
    blob = bucket.blob(blob)
    return blob.generate_signed_url(
        version="v4",
        expiration=exp,
        method="PUT",
        content_type=content_type,
        headers={"X-Goog-Content-Length-Range": f"{min_size},{max_size}"},
        credentials=credentials 
    )