resource "google_cloudbuild_trigger" "ci-build-trigger" {
    project = var.project_id
    name   = var.trigger_name
    location = "europe-west1"
    filename = var.config_path
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
        _TARGET_IMAGE_URI  = var.artifact_packagename
    }
}

