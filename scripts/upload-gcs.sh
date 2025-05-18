#!/bin/bash


source "$(pwd)/.env"


filename=""
signer_API="$SIGNER_API"
bucket="$BUCKET_NAME"
signed=""


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


sign_url() { 
  curl -sL -G "$signer_API?bucket=$bucket&filename=$filename" | jq -r '.'
}


#main execution ========
handle_options "$@"

#check inputs
if [ -z "$filename" ]; then
  echo "Error: -f or --file (file-path) are required."
  usage
  exit 1
fi


{
  echo "Trying to sign URL..."
  signed=$(sign_url)
  echo "Signed successfully: $signed"

  echo -e "\nUploading..."
  curl -f -X PUT -H "Content-Type: video/mp4" --upload-file "$filename" "$signed"
  echo "Upload successful, file: gs://$bucket/$filename"
} || {
  echo "Error occurred during signing or upload." >&2
  exit 1
}
