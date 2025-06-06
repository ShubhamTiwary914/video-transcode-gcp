#Storage Bucket =========================


resource "google_storage_bucket" "bucket_module" {
    name = var.bucket_name
    location = var.region_bucket
    force_destroy = true
    uniform_bucket_level_access = true 

    cors {
      origin = ["*"]
      method = ["*"]
      response_header = [
        "Content-Type",
        "Access-Control-Allow-Origin",
        "X-Goog-Content-Length-Range"  #restrict max upload file size 
      ]
      max_age_seconds = 3600 #5 min
    }

    #object 3 days TTL
    lifecycle_rule {
      condition {
        age = 3
      }

      action {
        type = "Delete"
      }
    } 
}


#GCS admin role for 2nd user SA (temp bucket)
resource "google_storage_bucket_iam_member" "bucket_admin" {
  bucket = google_storage_bucket.bucket_module.name
  role   = "roles/storage.admin"  
  member = "serviceAccount:${var.user_SA}"
}


output "bucket_name" {
  value = google_storage_bucket.bucket_module.name
}