import os
import sys
import requests
from dotenv import load_dotenv




filepath = ""
filename = ""
mode = None

load_dotenv()
signer_API = os.getenv("SIGNER_API")
bucket= os.getenv("BUCKET_NAME")



def sign_url_upload():
    if(not os.path.isfile(filepath)):
        print("file doesn't exist")
        exit(1)
    signed = requests.post(signer_API, json={
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

def sign_url_download(filename):
    signed = requests.get(signer_API, params={
        "bucket": bucket,
        "filename": filename
    }, headers={
        "Content-Type": "video/mp4"
    }) 
    signed.raise_for_status()
    return signed.text.strip('"')  



def main():
    try:
        filepath = sys.argv[1]
        filename = os.path.basename(filepath)
        mode = sys.argv[2]
    except IndexError:
        print("[arg-1]: filename, [arg-2]: mode(upload/download)")
        exit(1)
    if mode == 'upload':
        sign=sign_url_upload()
        uploadFile(sign)
    else:
        print(sign_url_download(filename))



if __name__ == "__main__":
    main()