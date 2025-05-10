#Storage Bucket =========================

resource "random_id" "suffix" {
  byte_length = 4
}


resource "google_storage_bucket" "static" {
    name = "temp-bucket-${random_id.suffix.hex}"
    location = "ASIA-SOUTH2"
    force_destroy = true
    uniform_bucket_level_access = true 

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

resource "google_storage_bucket_object" "demo" {
  name = "samplefile.txt"
  source = "/home/dev/work/projects/video-transcode/tests/sample.txt"
  content_type = "text/plain"
  bucket = google_storage_bucket.static.id
}