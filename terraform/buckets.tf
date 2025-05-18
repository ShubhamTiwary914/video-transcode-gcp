#Storage Bucket =========================

resource "random_id" "suffix" {
  byte_length = 4
}


resource "google_storage_bucket" "temp" {
    name = "temp-bucket-${random_id.suffix.hex}"
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

#GCS admin role for 2nd user SA 
resource "google_storage_bucket_iam_member" "bucket_admin" {
  bucket = google_storage_bucket.temp.name
  role   = "roles/storage.admin"  
  member = "serviceAccount:${var.user_SA}"
}
