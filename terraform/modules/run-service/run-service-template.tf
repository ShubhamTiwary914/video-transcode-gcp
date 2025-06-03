# Cloud run service  ========================================

resource "google_cloud_run_v2_service" "gcr-service" {
  name     = var.run-service-name 
  location = var.region
  deletion_protection = false
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = var.user_SA
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_reponame}/${var.artifact_packagename}:latest"
      ports {
        container_port = 8080
      }

      dynamic "env" {
        for_each = var.container_envs 
        content {
          name = env.value.name
          value = env.value.value
        }
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

# public access (0.0.0.0/0) -> no auth 
resource "google_cloud_run_v2_service_iam_member" "public_invoker" {
  name            = google_cloud_run_v2_service.gcr-service.name
  location        = google_cloud_run_v2_service.gcr-service.location   
  project         = google_cloud_run_v2_service.gcr-service.project
  role            = "roles/run.invoker"
  member          = "allUsers"
}


#provide GCR with signing role (sign bucket URLs)
resource "google_project_iam_member" "token_creator" {
  project = var.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:${var.user_SA}"
}