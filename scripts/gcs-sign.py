import os
import sys
import requests
from dotenv import load_dotenv


load_dotenv()
bucket= os.getenv("BUCKET_NAME")


def sign_url_upload(filepath, filename, signer_API):
    if(not os.path.isfile(filepath)):
        print("file doesn't exist")
        exit(1)
    data = {
        "bucket": bucket,
        "filename": filename
    }
    print(data)
    signed = requests.post(signer_API, json=data)
    signed.raise_for_status()
    return signed.text.strip('"')  


def uploadFile(filepath, filename, sign):
    with open(filepath, "rb") as f:
        response = requests.put(
            sign, data=f,
            headers={
                "Content-Type": "video/mp4"
            }
        )
    if not response.raise_for_status():
        print(f"[200]file uploaded, at: gs://{bucket}/{filename}")

def sign_url_download(filename: str, signer_API: str) -> str:
    signed = requests.get(signer_API, params={
        "bucket": bucket,
        "filename": filename
    }, headers={
        "Content-Type": "video/mp4"
    }) 
    signed.raise_for_status()
    return signed.text.strip('"')  



#args:  (1 - filename) (2 - local/gcp : signer) (3 - upload/download : mode)
def main():
    filepath = ""
    filename = ""
    mode = None
    signer = None

    try:
        if(len(sys.argv) < 4):
            raise Exception("[arg-1]: filename, [arg-2]: signer(local/gcp), [arg-3]: mode(upload/download)")
        filepath = sys.argv[1]
        filename = os.path.basename(filepath)
        signer = sys.argv[2]
        mode = sys.argv[3]
    except:
        print("[arg-1]: filename, [arg-2]: signer(local/gcp), [arg-3]: mode(upload/download)")
        exit(1)
    signer_API = os.getenv("SIGNER_API_LOCAL") if signer == 'local' else os.getenv("SIGNER_API_GCR")
    if mode == 'upload': 
        sign=sign_url_upload(filepath, filename, signer_API)
        print(sign)
        uploadFile(filepath, filename, sign)
    else:
        print(sign_url_download(filename, signer_API))



if __name__ == "__main__":
    main()