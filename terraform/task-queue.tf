

resource "google_cloud_tasks_queue" "simple" {
  name = "build-queue"
  location = var.region_queue

    http_target {
        uri_override {
            scheme = "HTTPS"
            host = "job-runner-service-${var.project_number}.${var.region}.run.app" 
            path_override {
              path = "/job-run"
            }
        }
        http_method = "POST"
        header_overrides {
          header {
            key = "Content-Type"
            value = "application/json"
          }
        }

        oidc_token {
          service_account_email = var.user_SA
        }
    }

    rate_limits {
      max_concurrent_dispatches = 20
      max_dispatches_per_second = 5
    }
}


resource "google_project_iam_member" "enqueuer" {
    project = var.project_id
    role   = "roles/cloudtasks.enqueuer"
    member = "serviceAccount:${var.user_SA}"
}
