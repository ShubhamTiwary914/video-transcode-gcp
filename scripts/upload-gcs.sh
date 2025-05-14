#!/bin/bash

sign=""
filename=""

usage() {
  echo "Usage: $0 [OPTIONS]"
  echo "Options:"
  echo " -h, --help         Display this help message"
  echo " -s, --sign <url>   GCS signed URL"
  echo " -f, --file <path>  File path to upload"
}

handle_options() {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h|--help)
        usage
        exit 0
        ;;
      -s|--sign)
        shift
        sign="$1"
        ;;
      -f|--file)
        shift
        filename="$1"
        ;;
      *)
        echo "Invalid option: $1"
        usage
        exit 1
        ;;
    esac
    shift
  done
}

handle_options "$@"

# Validate inputs
if [ -z "$sign" ] || [ -z "$filename" ]; then
  echo "Error: Both --sign and --file are required."
  usage
  exit 1
fi

curl -X PUT -H "Content-Type: application/octet-stream" \
     --upload-file "$filename" \
     "$sign"
