resource "google_cloudbuild_trigger" "sign_url_trigger" {
    project = var.project_id
    name   = "sign-url-trigger" 
    location = "europe-west1"
    filename = "src/sign_bucket/cloudbuild.yaml"
    service_account = "projects/${var.project_id}/serviceAccounts/${var.user_SA}"


    github {
        owner = "ShubhamTiwary914"
        name  = "video-transcode-gcp"
        push {
            branch = "main" 
        }
    }

    substitutions = {
        _ARTIFACT_REPO     = var.artifact_reponame
        _PROJECT_ID        = var.project_id
        _TARGET_IMAGE_URI  = var.artifact_packagename_signbucket
    }
}


resource "google_cloudbuild_trigger" "processjob_trigger" {
    project = var.project_id
    name   = "processjob-trigger" 
    location = "europe-west1"
    filename = "src/process_job/cloudbuild.yaml"
    service_account = "projects/${var.project_id}/serviceAccounts/${var.user_SA}"


    github {
        owner = "ShubhamTiwary914"
        name  = "video-transcode-gcp"
        push {
            branch = "main" 
        }
    }

    substitutions = {
        _ARTIFACT_REPO     = var.artifact_reponame
        _PROJECT_ID        = var.project_id
        _TARGET_IMAGE_URI  = var.artifact_packagename_processjob
    }
}

