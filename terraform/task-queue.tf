

# resource "google_cloud_tasks_queue" "simple" {
#   name = "buffer"
#   location = var.region_queue

#   http_target {
#     uri_override {
#         scheme = "HTTPS"
#         host = "587f-2401-4900-5f75-4ee3-54dd-9ede-e8f6-c381.ngrok-free.app" 
#         path_override {
#           path = "/queue-hook"
#         }
#     }
#   }
# }


# resource "google_project_iam_member" "enqueuer" {
#     project = var.project_id
#     role   = "roles/cloudtasks.enqueuer"
#     member = "serviceAccount:${var.user_SA}"
# }
