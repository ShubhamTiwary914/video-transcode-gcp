FROM python:3.10-alpine 

WORKDIR /app
COPY . .

RUN pip install --upgrade pip
RUN pip install --no-cache-dir -r requirements.txt
CMD ["uvicorn", "sign:app", "--reload", "--host", "0.0.0.0", "--port", "8080"]