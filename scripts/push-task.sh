ACCESS_TOKEN=$(gcloud auth print-access-token)

curl -X POST \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  "https://cloudtasks.googleapis.com/v2/projects/$TF_VAR_project_id/locations/asia-south1/queues/buffer/tasks" \
  -d '{
    "task": {
      "httpRequest": {
        "httpMethod": "POST",
        "url": "$QUEUE_GCR_URL/queue-hook",
        "headers": {
          "Content-Type": "application/json"
        },
        "body": "'$(echo -n '{"filepath":"gs://$BUCKET/somefile.txt"}' | base64)'"
      }
    }
  }'
