import os
from fastapi import FastAPI
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

class JobRunRequest(BaseModel):
    file_id: str
    hls_bucketname: str
    input_path: str
    creds_file: bool = False


@app.post("/job-run/")
@app.post("/job-run")
async def run_job(req: JobRunRequest):
    # GOOGLE_APPLICATION_CREDENTIALS for credentials path
    if req.creds_file:
        client = run_v2.JobsClient()
    else:
        creds, _ = default()
        client = run_v2.JobsClient(credentials=creds)

    job_name = client.job_path(project=os.getenv("PROJECT_ID"), location=REGION, job=JOB_NAME)

    container_override = RunJobRequest.Overrides.ContainerOverride(
        env=[
            EnvVar(name="FILE_ID", value=req.file_id),
            EnvVar(name="HLS_BUCKETNAME", value=req.hls_bucketname),
            EnvVar(name="MODE", value=MODE),
            EnvVar(name="INPUT_PATH", value=req.input_path),
        ]
    )
    overrides = RunJobRequest.Overrides(container_overrides=[container_override])

    request = RunJobRequest(
        name=job_name,
        overrides=overrides
    )

    operation = client.run_job(request=request)
    return {"status": "Job started", "operation_name": operation.operation.name}
