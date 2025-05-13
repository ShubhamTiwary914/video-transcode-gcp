gcloud storage sign-url gs://temp-bucket-01936b24/assembly-line.mp4 \
    --private-key-file=/home/dev/.keys/user-o1-gcp.json \
    --http-verb=PUT --duration=1h \
    --headers=Content-Type=video/mp4