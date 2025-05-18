
provider "google" {
  project = var.project_id
  region = var.region
}



#variables
variable "project_id" {
  description = "GCP Project ID"
}

variable "project_number" {
  description = "GCP Project Number"
}


variable "region" {
  description = "GCP Project region"
}

variable "region_bucket" {
  description = "GCP Project region (for storage buckets)"
}

variable "region_queue" {
  description = "GCP Project region supported (for task queue)"
}

variable "artifact_reponame" {
  description = "Repository name for Artifact Registry"
}

variable "artifact_packagename_signbucket" {
  description = "Package name under the Artifact Repository"
}


variable "default_SA" {
  description = "Email of GCP default service account for the project"
}

variable "user_SA" {
  description = "Email of GCP 2nd user service account for the project"
}