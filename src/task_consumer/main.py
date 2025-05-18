from fastapi import FastAPI, Request, HTTPException
from pydantic import BaseModel
import logging


app = FastAPI(debug=True)

class TaskPayload(BaseModel):
    filepath: str  



@app.route("/health", methods=["GET", "POST"])
def healthCheck(req: Request) -> dict:
    return {"status", f"successful health check at {req.method} /"}


@app.post("/queue-hook")
def queue_worker(payload: TaskPayload, req: Request) -> bool:
    """
        Webhook for GCP task queue to hit after GCS-bucket upload
        Args:
            payload (TaskPayload): GCS bucket path [format: gs://<bucket>/<filename.ext>]
    """
    logging.info(f"Payload: {payload.filepath}\n")
    if not payload.filepath.startswith("gs://"):
        raise HTTPException(status_code=400, detail="Invalid GCS filepath")
    return True
