#Cloud run service  ========================================

resource "google_cloud_run_service" "default" {
  name     = "sign-url-service"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_reponame}/${var.artifact_packagename_signbucket}:latest"
        ports {
          container_port = 8080
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "public_invoker" {
  location        = google_cloud_run_service.default.location
  project         = google_cloud_run_service.default.project
  service         = google_cloud_run_service.default.name
  role            = "roles/run.invoker"
  member          = "allUsers"
}





