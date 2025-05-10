#Cloud run service  ========================================

# resource "google_cloud_run_service" "default" {
#   name     = "cloudrun-srv-go-web"
#   location = var.region

#   template {
#     spec {
#       containers {
#         image = "docker.io/sardinesszsz/go-hello:v2"
#         ports {
#           container_port = 8080
#         }
#       }
#     }
#   }

#   traffic {
#     percent         = 100
#     latest_revision = true
#   }
# }

# resource "google_cloud_run_service_iam_member" "public_invoker" {
#   location        = google_cloud_run_service.default.location
#   project         = google_cloud_run_service.default.project
#   service         = google_cloud_run_service.default.name
#   role            = "roles/run.invoker"
#   member          = "allUsers"
# }





