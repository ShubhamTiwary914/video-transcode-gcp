
provider "google" {
  project = var.project_id
  region = var.region
}



#variables
variable "project_id" {
  description = "GCP Project ID"
}

variable "region" {
  description = "GCP Project region"
}

variable "region_bucket" {
  description = "GCP Project region (for storage buckets)"
}

variable "artifact_reponame" {
  description = "Repository name for Artifact Registry"
}