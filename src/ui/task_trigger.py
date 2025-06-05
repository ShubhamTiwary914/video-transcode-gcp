import os
import json
import random
import string
from subprocess import run, PIPE
from dotenv import load_dotenv

load_dotenv()

REGION = os.getenv("TASK_REGION")
QUEUE = os.getenv("TASK_QUEUE")
SERVICE_ACCOUNT = os.getenv("TASK_SERVICE_ACCOUNT")
TARGET_URL = os.getenv("TASK_TARGET_URL")
HLS_BUCKET = os.getenv("HLS_BUCKET")


def random_id(length=12):
    chars = string.ascii_lowercase + string.digits + "_"
    return ''.join(random.choices(chars, k=length))

def trigger_task(download_url: str):
    file_id = random_id()
    body = {
        "file_id": file_id,
        "hls_bucketname": HLS_BUCKET,
        "input_path": download_url,
        "creds_file": True
    }

    cmd = [
        "gcloud", "tasks", "create-http-task",
        f"--queue={QUEUE}",
        f"--location={REGION}",
        f"--url={TARGET_URL}",
        "--method=POST",
        f"--body-content={json.dumps(body)}",
        f"--oidc-service-account-email={SERVICE_ACCOUNT}"
    ]

    result = run(cmd, stdout=PIPE, stderr=PIPE, text=True)
    if result.returncode != 0:
        print("[Task Error]", result.stderr.strip())
        return None
    print("[Task Triggered]", result.stdout.strip())
    return file_id
