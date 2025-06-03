FROM python:3.10-alpine 


EXPOSE 8080 
WORKDIR /app

COPY . .

ENV MODE=prod
ENV REGION=""
ENV JOB_NAME=processjob
ENV PROJECT_ID=""

RUN pip install --upgrade pip
RUN pip install --no-cache-dir -r requirements.txt
CMD ["uvicorn", "runner:app", "--reload", "--host", "0.0.0.0", "--port", "8080"]

