#Cloud run service  ========================================

resource "google_cloud_run_v2_service" "default" {
  name     = "sign-url-service"
  location = var.region
  deletion_protection = false
  ingress = "INGRESS_TRAFFIC_ALL"
  
  
  template {
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_reponame}/${var.artifact_packagename_signbucket}:latest"
      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }  
  }
}

resource "google_cloud_run_v2_service_iam_member" "public_invoker" {
  name = google_cloud_run_v2_service.default.name
  location        = google_cloud_run_v2_service.default.location   
  project         = google_cloud_run_v2_service.default.project
  role            = "roles/run.invoker"
  member          = "allUsers"
}
