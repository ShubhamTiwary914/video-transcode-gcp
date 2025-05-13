

resource "google_artifact_registry_repository" "video_transcode" {
  repository_id  = var.artifact_reponame 
  format         = "DOCKER"
  location       = var.region 
  description    = "Docker images for project (video-transcode)"
  mode           = "STANDARD_REPOSITORY"
}
