#!/usr/bin/env zsh
source ~/.zshrc

conda activate gcp
cd ../src/task_consumer

uvicorn main:app --reload --host 0.0.0.0 --port 9090