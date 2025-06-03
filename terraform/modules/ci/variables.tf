variable "project_id" {
    description = "GCP Project ID"
}

variable "trigger_name" {
    description = "Name for the cloud build trigger service"
}

variable "config_path" {
    description = "Path to cloudbuild.yaml in project"
}

variable "user_SA" {
  description = "Email of GCP 2nd user service account for the project"
}

variable "artifact_reponame" {
  description = "Repository name for Artifact Registry"
}


variable "artifact_packagename" {
  description = "Package name for the project's Artifact Registry"
}