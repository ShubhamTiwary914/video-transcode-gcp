resource "google_cloud_run_v2_job" "processjob" {
  name     = var.job_name 
  location = var.region

  template {
    template {
      containers {
        # image = "docker.io/sardinesszsz/processjob:v2"
        image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_reponame}/${var.artifact_packagename_processjob}:latest"

        resources {
          limits = {
            cpu    = "4"
            memory = "8Gi"
          }
        }

        env {
          name  = "HLS_BUCKETNAME"
          value = var.hls_bucketname
        }
        env {
          name = "MODE"
          value = "prod"
        }
      }

      service_account = var.user_SA
    }
  }

  deletion_protection = false
}
