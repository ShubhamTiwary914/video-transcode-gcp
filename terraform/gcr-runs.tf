

module "gcr-signer" {
  source = "./modules/run-service"
  project_id = var.project_id
  region = var.region
  user_SA = var.user_SA
  run-service-name = "sign-url-service"
  artifact_reponame = var.artifact_reponame
  artifact_packagename = var.artifact_packagename_signbucket
}

module "job-runner" {
  source = "./modules/run-service"
  project_id = var.project_id
  region = var.region
  user_SA = var.user_SA
  run-service-name = "job-runner-service"
  artifact_reponame = var.artifact_reponame
  artifact_packagename = var.artifact_packagename_jobrunner

  container_envs = [
    {name = "MODE", value = "prod"},
    {name= "REGION", value = var.region },
    {name= "PROJECT_ID", value = var.project_id },
    {name= "JOB_NAME", value = var.job_name },
  ]
}

