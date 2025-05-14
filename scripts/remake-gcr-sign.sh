cd ../terraform
terraform destroy -target=google_cloud_run_v2_service.default -auto-approve
terraform apply -auto-approve