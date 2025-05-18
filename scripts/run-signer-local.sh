#!/usr/bin/env zsh
source ~/.zshrc

conda activate gcp
cd ../src/sign_bucket/
uvicorn sign:app --reload --host 0.0.0.0 --port 8090