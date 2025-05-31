import os
import sys
import requests
from dotenv import load_dotenv



load_dotenv()

filepath = ""
filename = ""

try:
    filepath = sys.argv[1]
    filename = os.path.basename(filepath)
except IndexError:
    print("[arg-1] pass file argument")
    exit(1)
signer_API = os.getenv("SIGNER_API")
bucket= os.getenv("BUCKET_NAME")


if(not os.path.isfile(filepath)):
    print("file doesn't exist")
    exit(1)


def sign_url_upload():
    signed = requests.post(signer_API, json={
        "bucket": bucket,
        "filename": filename
    })
    signed.raise_for_status()
    return signed.text.strip('"')  


def sign_url_download():
    signed = requests.get(signer_API, params={
        "bucket": bucket,
        "filename": filename
    }) 
    signed.raise_for_status()
    return signed.text.strip('"')  


def uploadFile(sign):
    with open(filepath, "rb") as f:
        response = requests.put(
            sign, data=f,
            headers={
                "Content-Type": "video/mp4"
            }
        )
    if not response.raise_for_status():
        print(f"[200]file uploaded, at: gs://{bucket}/{filename}")
 

# Upload
# sign=sign_url_upload()
# uploadFile(sign)


# Download
print(sign_url_download())