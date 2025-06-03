
variable "project_id" {
  description = "GCP Project ID"
}

variable "region" {
  description = "GCP Project region"
}

variable "run-service-name" {
  description = "Name for the GCR service"
}

variable "user_SA" {
  description = "Email of GCP 2nd user service account for the project"
}

variable "artifact_reponame" {
  description = "Repository name for Artifact Registry"
}

variable "artifact_packagename" {
  description = "Package name under the Artifact Repository [container to be used for GCR service]"
}


variable "container_envs" {
  type = list(object({ name=string, value=string }))
  description = "Optional list of environment variables for the container"
  default = []
}