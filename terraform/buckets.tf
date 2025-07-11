
resource "random_id" "bucket_suffix" {
    byte_length = 4
}


module "temp_bucket" {
  source = "./modules/bucket/"
  bucket_name = "temp-${random_id.bucket_suffix.hex}"
  region_bucket = var.region_bucket
  user_SA = var.user_SA
}



//for performing mock-tests on
module "mock_bucket" {
  source = "./modules/bucket/"
  bucket_name = "mock-${random_id.bucket_suffix.hex}"
  region_bucket = var.region_bucket
  user_SA = var.user_SA  
}


module "hls_bucket" {
  source = "./modules/bucket/"
  bucket_name = "hls-${random_id.bucket_suffix.hex}"
  region_bucket = var.region_bucket
  user_SA = var.user_SA
}

//GET public from HLS bucket (for streaming from)
resource "google_storage_bucket_iam_member" "public_read" {
  bucket = module.hls_bucket.bucket_name
  role = "roles/storage.objectViewer"
  member = "allUsers"
}