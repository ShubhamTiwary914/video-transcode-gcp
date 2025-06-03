
module "sign_url_trigger" {
  source = "./modules/ci"
  project_id = var.project_id
  user_SA = var.user_SA
  trigger_name = "sign-url-trigger"
  config_path = "src/sign_bucket/cloudbuild.yaml"
  artifact_reponame = var.artifact_reponame
  artifact_packagename = var.artifact_packagename_signbucket
}

module "processjob_trigger" {
  source = "./modules/ci"
  project_id = var.project_id
  user_SA = var.user_SA
  trigger_name = "processjob-trigger"
  config_path = "src/process_job/cloudbuild.yaml"
  artifact_reponame = var.artifact_reponame
  artifact_packagename = var.artifact_packagename_processjob
}


module "jobrunner_trigger" {
  source = "./modules/ci"
  project_id = var.project_id
  user_SA = var.user_SA
  artifact_reponame = var.artifact_reponame
  artifact_packagename = var.artifact_packagename_jobrunner
  trigger_name = "jobrunner-trigger"
  config_path = "src/job_runner/cloudbuild.yaml"
}




