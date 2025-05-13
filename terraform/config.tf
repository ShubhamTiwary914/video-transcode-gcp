
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