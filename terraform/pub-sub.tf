

resource "google_pubsub_topic" "job_complete_topic" {
  name = "process-status"
}

resource "google_pubsub_topic_iam_member" "job_complete_publisher" {
  topic = google_pubsub_topic.job_complete_topic.name
  role  = "roles/pubsub.publisher"
  member = "serviceAccount:${var.user_SA}"
}


