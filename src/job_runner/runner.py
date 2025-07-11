import os
from fastapi import FastAPI, Request
from pydantic import BaseModel
from google.cloud import run_v2
from google.cloud.run_v2.types import RunJobRequest
from google.cloud.run_v2.types.k8s_min import EnvVar
from google.auth import default
from dotenv import load_dotenv

app = FastAPI()
load_dotenv()

JOB_NAME = os.getenv("JOB_NAME")
REGION = os.getenv("REGION")
MODE = os.getenv("MODE")
PROJECT_ID = os.getenv("PROJECT_ID")


class JobRunRequest(BaseModel):
    file_id: str
    hls_bucketname: str
    input_path: str
    pub_topic : str #pub-topic after job done (completion status)
    creds_file: bool = False


@app.get("/healthcheck")
def checkEnvs():
    return {
        "JOB_NAME": JOB_NAME,
        "REGION": REGION,
        "MODE": MODE,
        "PROJECT_ID": PROJECT_ID
    }


@app.post("/job-run/")
@app.post("/job-run")
async def run_job(req: JobRunRequest):
    # GOOGLE_APPLICATION_CREDENTIALS for credentials path
    if req.creds_file:
        client = run_v2.JobsClient()
    else:
        creds, _ = default()
        client = run_v2.JobsClient(credentials=creds)

    job_name = client.job_path(project=PROJECT_ID, location=REGION, job=JOB_NAME)

    container_override = RunJobRequest.Overrides.ContainerOverride(
        env=[
            EnvVar(name="PROJECT_ID", value=PROJECT_ID),
            EnvVar(name="FILE_ID", value=req.file_id),
            EnvVar(name="HLS_BUCKETNAME", value=req.hls_bucketname),
            EnvVar(name="MODE", value=MODE),
            EnvVar(name="INPUT_PATH", value=req.input_path),
            EnvVar(name="PUB_TOPIC", value=req.pub_topic)
        ]
    )
    overrides = RunJobRequest.Overrides(container_overrides=[container_override])

    request = RunJobRequest(
        name=job_name,
        overrides=overrides
    )

    operation = client.run_job(request=request)
    return {"status": "Job started", "operation_name": operation.operation.name}
