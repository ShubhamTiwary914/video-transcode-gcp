#!/usr/bin/env zsh
source ~/.zshrc

kconda
conda activate gcp
uvicorn runner:app --reload --host 0.0.0.0 --port 9027